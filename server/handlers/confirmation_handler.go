package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"net/http"

	c "github.com/thoniwutr/schedule-school-teachning-bsd13-backend/constant"
	"github.com/thoniwutr/schedule-school-teachning-bsd13-backend/server/dto"
	"github.com/thoniwutr/schedule-school-teachning-bsd13-backend/service"
	"github.com/thoniwutr/schedule-school-teachning-bsd13-backend/util"
)


type ConfirmationHandler struct {
	util       *util.HandlerUtil
	confirmationService service.ConfirmationServiceInterface
}

func NewConfirmationHandler(util *util.HandlerUtil, subjectService service.ConfirmationServiceInterface) *ConfirmationHandler {
	return &ConfirmationHandler{util: util, confirmationService: subjectService}
}


func (m *ConfirmationHandler) GetAllConfirmation(rw http.ResponseWriter, r *http.Request) {

	response, err := m.confirmationService.GetAllConfirmation()
	if err != nil {
		m.util.WrappedError(rw, err)
		return
	}

	resp, err := json.Marshal(response)
	if err != nil {
		m.util.WrappedError(rw, fmt.Errorf("unable to marshal json: %w", err))
		return
	}

	rw.WriteHeader(http.StatusOK)
	rw.Write(resp)
}

func (m *ConfirmationHandler) GetAllConfirmationDetail(rw http.ResponseWriter, r *http.Request) {

	id, ok := mux.Vars(r)["id"]
	if !ok {
		m.util.HTTPError(rw, fmt.Errorf("invalid id in path"), http.StatusBadRequest)
		return
	}

	response, err := m.confirmationService.GetAllConfirmationDetail(id)
	if err != nil {
		m.util.WrappedError(rw, err)
		return
	}

	resp, err := json.Marshal(response)
	if err != nil {
		m.util.WrappedError(rw, fmt.Errorf("unable to marshal json: %w", err))
		return
	}

	rw.WriteHeader(http.StatusOK)
	rw.Write(resp)
}


func (m *ConfirmationHandler) AddConfirmation(rw http.ResponseWriter, r *http.Request) {
	req := &dto.NewConfirmation{}
	if err := json.NewDecoder(r.Body).Decode(req); err != nil {
		m.util.HTTPError(rw, fmt.Errorf("error deserializing main subject %w", err), http.StatusBadRequest)
		return
	}

	if err := req.Validate(); err != nil {
		m.util.WrappedError(rw, &c.ErrValidation{Violations: err})
		return
	}

	if err := m.confirmationService.AddConfirmation(req); err != nil {
		m.util.WrappedError(rw, err)
		return
	}


	m.util.HTTPSuccess(rw, "success", http.StatusCreated)
}


func (m *ConfirmationHandler) AddConfirmationDetail(rw http.ResponseWriter, r *http.Request) {

	req := &dto.NewConfirmationDetail{}
	if err := json.NewDecoder(r.Body).Decode(req); err != nil {
		m.util.HTTPError(rw, fmt.Errorf("error deserializing main subject %w", err), http.StatusBadRequest)
		return
	}

	if err := req.Validate(); err != nil {
		m.util.WrappedError(rw, &c.ErrValidation{Violations: err})
		return
	}

	if err := m.confirmationService.AddConfirmationDetail(req); err != nil {
		m.util.WrappedError(rw, err)
		return
	}


	m.util.HTTPSuccess(rw, "success", http.StatusCreated)
}
