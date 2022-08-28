package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	c "github.com/thoniwutr/schedule-school-teachning-bsd13-backend/constant"
	"github.com/thoniwutr/schedule-school-teachning-bsd13-backend/server/dto"
	"github.com/thoniwutr/schedule-school-teachning-bsd13-backend/service"
	"github.com/thoniwutr/schedule-school-teachning-bsd13-backend/util"
)


type SubjectHandler struct {
	util       *util.HandlerUtil
	subjectService service.SubjectServiceInterface
}

func NewSubjectHandler(util *util.HandlerUtil, subjectService service.SubjectServiceInterface) *SubjectHandler {
	return &SubjectHandler{util: util, subjectService: subjectService}
}


func (m *SubjectHandler) GetAllSubject(rw http.ResponseWriter, r *http.Request) {

	subjectResponse, err := m.subjectService.GetAllSubject()
	if err != nil {
		m.util.WrappedError(rw, err)
		return
	}

	resp, err := json.Marshal(subjectResponse)
	if err != nil {
		m.util.WrappedError(rw, fmt.Errorf("unable to marshal json: %w", err))
		return
	}

	rw.WriteHeader(http.StatusOK)
	rw.Write(resp)
}


func (m *SubjectHandler) AddSubject(rw http.ResponseWriter, r *http.Request) {
	nSubject := &dto.NewSubject{}
	if err := json.NewDecoder(r.Body).Decode(nSubject); err != nil {
		m.util.HTTPError(rw, fmt.Errorf("error deserializing main subject %w", err), http.StatusBadRequest)
		return
	}

	if err := nSubject.Validate(); err != nil {
		m.util.WrappedError(rw, &c.ErrValidation{Violations: err})
		return
	}

	if err := m.subjectService.AddSubject(nSubject); err != nil {
		m.util.WrappedError(rw, err)
		return
	}


	m.util.HTTPSuccess(rw, "success", http.StatusCreated)
}
