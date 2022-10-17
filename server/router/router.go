package router

import (
	"github.com/gorilla/mux"
	"github.com/schedule-school-teachning-bsd13-backend/server/handlers"
	"github.com/schedule-school-teachning-bsd13-backend/server/middleware"
	"github.com/schedule-school-teachning-bsd13-backend/util"
	"net/http"
)

// NewRouter returns a http.Handler which handles different routes for the server
func NewRouter(
	logger *util.Logger,
	th *handlers.TeacherHandler,
	msh *handlers.MainSubjectHandler,
	sh *handlers.SubjectHandler,
	ch *handlers.ConfirmationHandler,
	sch *handlers.ScheduleHandler,
) http.Handler {

	r := mux.NewRouter()

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

	cr.Methods(http.MethodPost).Path("/{id}").HandlerFunc(ch.AddConfirmationDetail)
	cr.Methods(http.MethodGet).Path("/{id}").HandlerFunc(ch.GetAllConfirmationDetail)

	scr := r.PathPrefix("/create-schedule").Subrouter()
	scr.Methods(http.MethodPost).Path("").HandlerFunc(sch.ScheduleClass)

	return middleware.RemoveTrailingSlash(r)
}
