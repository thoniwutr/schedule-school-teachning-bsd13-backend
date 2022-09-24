package router

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/swaggo/http-swagger"

	"github.com/thoniwutr/schedule-school-teachning-bsd13-backend/server/handlers"
	"github.com/thoniwutr/schedule-school-teachning-bsd13-backend/server/middleware"
	"github.com/thoniwutr/schedule-school-teachning-bsd13-backend/util"
)

// NewRouter returns a http.Handler which handles different routes for the server
func NewRouter(
	logger *util.Logger,
	mh *handlers.MerchantsHandler,
	kymh *handlers.KymHandler,
	th *handlers.TeacherHandler,
	msh *handlers.MainSubjectHandler,
	sh *handlers.SubjectHandler,
	ch *handlers.ConfirmationHandler,
) http.Handler {

	r := mux.NewRouter()

	r.Methods(http.MethodGet).Path("/merchants/healthz").HandlerFunc(mh.HealthCheck)
	// swagger docs
	r.Methods(http.MethodGet).PathPrefix("/merchants/swagger").Handler(httpSwagger.WrapHandler)


	// subrouter for /merchants
	msr := r.PathPrefix("/merchants").Subrouter()
	msr.Use(middleware.ContentTypeJSON)
	msr.Use(middleware.NewRequestLogger(logger).LogRequest)
	msr.Methods(http.MethodPost).Path("").HandlerFunc(mh.AddMerchant)
	msr.Methods(http.MethodGet).Path("/{merchantId}").HandlerFunc(mh.GetMerchant)
	msr.Methods(http.MethodPut).Path("/{merchantId}").HandlerFunc(mh.UpdateMerchant)

	// subrouter for /merchants/{merchantId}/pay-*
	mcsr := msr.PathPrefix("/{merchantId}").Subrouter()
	mcsr.Use(middleware.ContentTypeJSON)
	mcsr.Methods(http.MethodPost).Path("/pay-out-config").HandlerFunc(mh.UpsertPayOutConfig)
	mcsr.Methods(http.MethodGet).Path("/pay-out-config").HandlerFunc(mh.GetPayOutConfig)

	// subrouter for /kym
	kymr := r.PathPrefix("/kym").Subrouter()
	kymr.Use(middleware.ContentTypeJSON)
	kymr.Use(middleware.NewRequestLogger(logger).LogRequest)
	kymr.Methods(http.MethodPost).Path("").HandlerFunc(kymh.AddKym)
	kymr.Methods(http.MethodGet).Path("").HandlerFunc(kymh.GetAllKym)
	kymr.Methods(http.MethodGet).Path("/{id}").HandlerFunc(kymh.GetKym)
	kymr.Methods(http.MethodPut).Path("/{id}/status").HandlerFunc(kymh.UpdateKymStatus)

	// subrouter for /teacher
	str := r.PathPrefix("/teacher").Subrouter()
	str.Use(middleware.ContentTypeJSON)
	str.Use(middleware.NewRequestLogger(logger).LogRequest)
	str.Methods(http.MethodPost).Path("").HandlerFunc(th.AddTeacher)
	str.Methods(http.MethodGet).Path("").HandlerFunc(th.GetAllTeacher)


	// subrouter for /mainsubject
	msjr := r.PathPrefix("/mainsubject").Subrouter()
	msjr.Use(middleware.ContentTypeJSON)
	msjr.Use(middleware.NewRequestLogger(logger).LogRequest)
	msjr.Methods(http.MethodPost).Path("").HandlerFunc(msh.AddMainSubject)
	msjr.Methods(http.MethodGet).Path("").HandlerFunc(msh.GetAllMainSubject)


	// subrouter for /mainsubject
	sr := r.PathPrefix("/subject").Subrouter()
	sr.Use(middleware.ContentTypeJSON)
	sr.Use(middleware.NewRequestLogger(logger).LogRequest)
	sr.Methods(http.MethodPost).Path("").HandlerFunc(sh.AddSubject)
	sr.Methods(http.MethodGet).Path("").HandlerFunc(sh.GetAllSubject)


	// subrouter for /mainsubject
	cr := r.PathPrefix("/confirmation").Subrouter()
	cr.Use(middleware.ContentTypeJSON)
	cr.Use(middleware.NewRequestLogger(logger).LogRequest)
	cr.Methods(http.MethodPost).Path("").HandlerFunc(ch.AddConfirmation)
	cr.Methods(http.MethodGet).Path("").HandlerFunc(ch.GetAllConfirmation)

	cr.Methods(http.MethodPost).Path("/{id}").HandlerFunc(ch.AddConfirmation)
	cr.Methods(http.MethodGet).Path("/{id}").HandlerFunc(ch.GetAllConfirmationDetail)

	return middleware.RemoveTrailingSlash(r)
}
