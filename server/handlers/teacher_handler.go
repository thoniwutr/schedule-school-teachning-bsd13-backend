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

// TeacherHandler is a handler for /teacher path
type TeacherHandler struct {
	util       *util.HandlerUtil
	teacherService service.TeacherServiceInterface
}

func NewTeacherHandler(util *util.HandlerUtil, teacherService service.TeacherServiceInterface) *TeacherHandler {
	return &TeacherHandler{util: util, teacherService: teacherService}
}

// GetAllTeacher godoc
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
func (t *TeacherHandler) GetAllTeacher(rw http.ResponseWriter, r *http.Request) {

	teacherResponse, err := t.teacherService.GetAllTeacher()
	if err != nil {
		t.util.WrappedError(rw, err)
		return
	}

	resp, err := json.Marshal(teacherResponse)
	if err != nil {
		t.util.WrappedError(rw, fmt.Errorf("unable to marshal json: %w", err))
		return
	}

	rw.WriteHeader(http.StatusOK)
	rw.Write(resp)
}

// AddTeacher godoc
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
func (t *TeacherHandler) AddTeacher(rw http.ResponseWriter, r *http.Request) {
	nTeacher := &dto.NewTeacher{}
	if err := json.NewDecoder(r.Body).Decode(nTeacher); err != nil {
		t.util.HTTPError(rw, fmt.Errorf("error deserializing kym : %w", err), http.StatusBadRequest)
		return
	}

	if err := nTeacher.Validate(); err != nil {
		t.util.WrappedError(rw, &c.ErrValidation{Violations: err})
		return
	}

	if err := t.teacherService.AddTeacher(nTeacher); err != nil {
		t.util.WrappedError(rw, err)
		return
	}


	t.util.HTTPSuccess(rw, "success", http.StatusCreated)
}
