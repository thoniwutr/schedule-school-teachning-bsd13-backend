package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/Beam-Data-Company/merchant-config-svc/api"
	"github.com/Beam-Data-Company/merchant-config-svc/config"
	"github.com/Beam-Data-Company/merchant-config-svc/db"
	"github.com/Beam-Data-Company/merchant-config-svc/messaging"
	"github.com/Beam-Data-Company/merchant-config-svc/server/handlers"
	"github.com/Beam-Data-Company/merchant-config-svc/server/router"
	"github.com/Beam-Data-Company/merchant-config-svc/server/security"
	"github.com/Beam-Data-Company/merchant-config-svc/service"
	"github.com/Beam-Data-Company/merchant-config-svc/thirdparty"
	"github.com/Beam-Data-Company/merchant-config-svc/util"
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

	mh := handlers.NewMerchantsHandler(hu, roleAuth, publisher, merchantService)
	kym := handlers.NewKymHandler(hu, roleAuth, newKymService)
	th := handlers.NewTeacherHandler(hu,teacherService)

	r := router.NewRouter(l, mh, kym,th)

	srv := &http.Server{
		Addr:         fmt.Sprintf("0.0.0.0:%s", cfg.Server.Port),
		Handler:      r,
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
