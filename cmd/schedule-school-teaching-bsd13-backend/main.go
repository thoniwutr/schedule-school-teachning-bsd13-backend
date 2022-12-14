package main

import (
	"context"
	"fmt"
	"github.com/rs/cors"
	"github.com/thoniwutr/schedule-school-teachning-bsd13-backend/api"
	"github.com/thoniwutr/schedule-school-teachning-bsd13-backend/config"
	"github.com/thoniwutr/schedule-school-teachning-bsd13-backend/db"
	"github.com/thoniwutr/schedule-school-teachning-bsd13-backend/messaging"
	"github.com/thoniwutr/schedule-school-teachning-bsd13-backend/server/handlers"
	"github.com/thoniwutr/schedule-school-teachning-bsd13-backend/server/router"
	"github.com/thoniwutr/schedule-school-teachning-bsd13-backend/server/security"
	"github.com/thoniwutr/schedule-school-teachning-bsd13-backend/service"
	"github.com/thoniwutr/schedule-school-teachning-bsd13-backend/thirdparty"
	"github.com/thoniwutr/schedule-school-teachning-bsd13-backend/util"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"
)

// @title Merchant Config Service API
// @version 0.1
// @description Merchant Config Service manages the Merchant entities.
// @description Merchants on our platform would be able to manage configuration for their payments such as
// @description payment methods they allow, and how they would like to receive their payouts.
// @description During customer purchases, this configuration is used to determine how users are able to pay.

// @contact.name Beam Developers
// @contact.url http://www.beamcheckout.com
// @contact.email developers@beamcheckout.com
func main() {

	cfg, err := config.AppConfig()
	if err != nil {
		log.Fatalf("failed to load app config %v", err)
	}

	l := util.NewLogger(cfg.Debug)

	// Setup Cloud Storage Client
	cs, err := util.NewCloudStorage(cfg.KymBucketName)
	if err != nil {
		log.Fatalf("failed to connect cloud storage: %v", err)
	}

	appDb, err := db.NewAppDatastore(cfg.ProjectID)
	if err != nil {
		log.Fatalf("failed to init appDb connection: %v", err)
	}
	defer appDb.Close()

	publisher, err := messaging.NewPubsubPublisher(cfg.ProjectID, cfg.MerchantCreatePublisher.TopicID, l)
	if err != nil {
		log.Fatalf("failed to init publisher connection: %v", err)
	}

	hu := util.NewHandlerUtil(l)
	roleAuth := security.NewRoleAuth(appDb)

	apiClient := thirdparty.NewApiClient(l, cfg.OrganisationServiceURL, cfg.RecipientServiceURL)
	merchantService := service.NewMerchantService(appDb, publisher)
	newKymService := service.NewKymService(appDb, cs, apiClient, merchantService)
	teacherService := service.NewTeacherService(appDb)
	mainSubjectService := service.NewMainSubjectService(appDb)
	subjectService := service.NewSubjectService(appDb)
	confirmationService := service.NewConfirmationService(appDb)

	mh := handlers.NewMerchantsHandler(hu, roleAuth, publisher, merchantService)
	kym := handlers.NewKymHandler(hu, roleAuth, newKymService)
	th := handlers.NewTeacherHandler(hu,teacherService)
	msh := handlers.NewMainSubjectHandler(hu,mainSubjectService)
	sh := handlers.NewSubjectHandler(hu,subjectService)
	ch := handlers.NewConfirmationHandler(hu,confirmationService)

	r := router.NewRouter(l, mh, kym,th,msh, sh,ch)
	cor := cors.New(cors.Options{
		AllowedOrigins: []string{"http://localhost:3000"},
		AllowCredentials: true,
	})

	handler := cor.Handler(r)

	srv := &http.Server{
		Addr:         fmt.Sprintf("0.0.0.0:%s", cfg.Server.Port),
		Handler:      handler,
		WriteTimeout: cfg.Server.TimeoutWrite,
		ReadTimeout:  cfg.Server.TimeoutRead,
		IdleTimeout:  cfg.Server.TimeoutIdle,
	}
	api.SwaggerInfo.Host = "localhost:" + cfg.Server.Port

	// Run our server in a goroutine so that it doesn't block.
	go func() {
		l.Info().Msgf("Starting up server on %v", srv.Addr)
		if err := srv.ListenAndServe(); err != nil {
			l.Error().Err(err).Msg("Server shutting down...")
		}
	}()

	// below code allows for graceful shut down
	c := make(chan os.Signal, 1)
	// We'll accept graceful shutdowns when quit via SIGINT (Ctrl+C)
	// SIGKILL, SIGQUIT or SIGTERM (Ctrl+/) will not be caught.
	signal.Notify(c, os.Interrupt)

	<-c

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	srv.Shutdown(ctx)
	l.Info().Msg("Graceful shutdown")
	os.Exit(0)
}
