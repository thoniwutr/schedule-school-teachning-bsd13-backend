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

// MainSubjectHandler is a handler for /teacher path
type MainSubjectHandler struct {
	util       *util.HandlerUtil
	mainSubjectService service.MainSubjectServiceInterface
}

func NewMainSubjectHandler(util *util.HandlerUtil, mainSubjectService service.MainSubjectServiceInterface) *MainSubjectHandler {
	return &MainSubjectHandler{util: util, mainSubjectService: mainSubjectService}
}

// GetAllMainSubject godoc
// @Id GetAllKym
// @Summary Get All Kym
// @Description Returns all Kym detail only necessary field
// @Tags kym
// @Produce json
// @Accept json
// @Param status query string false "string enums" Enums("approved", "rejected", "pending")
// @Success 200 {object} []dto.KymResponse "success"
// @Failure default {object} util.APIResponse "fail"
// @Router /kym [get]
func (m *MainSubjectHandler) GetAllMainSubject(rw http.ResponseWriter, r *http.Request) {

	mainSubjectResponse, err := m.mainSubjectService.GetAllMainSubject()
	if err != nil {
		m.util.WrappedError(rw, err)
		return
	}

	resp, err := json.Marshal(mainSubjectResponse)
	if err != nil {
		m.util.WrappedError(rw, fmt.Errorf("unable to marshal json: %w", err))
		return
	}

	rw.WriteHeader(http.StatusOK)
	rw.Write(resp)
}

// AddMainSubject godoc
// @Id AddKym
// @Summary Kym Registration
// @Description Submitted Kym detail and documents to register with Beam
// @Tags kym
// @Produce json
// @Accept json
// @Param requestBody body dto.NewKym true "NewKym entity"
// @Success 201 {object} util.APIResponse "success"
// @Failure default {object} util.APIResponse "fail"
// @Router /kym [post]
func (m *MainSubjectHandler) AddMainSubject(rw http.ResponseWriter, r *http.Request) {
	nMainSubject := &dto.NewMainSubject{}
	if err := json.NewDecoder(r.Body).Decode(nMainSubject); err != nil {
		m.util.HTTPError(rw, fmt.Errorf("error deserializing main subject %w", err), http.StatusBadRequest)
		return
	}

	if err := nMainSubject.Validate(); err != nil {
		m.util.WrappedError(rw, &c.ErrValidation{Violations: err})
		return
	}

	if err := m.mainSubjectService.AddMainSubject(nMainSubject); err != nil {
		m.util.WrappedError(rw, err)
		return
	}


	m.util.HTTPSuccess(rw, "success", http.StatusCreated)
}
