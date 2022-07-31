package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"

	c "github.com/Beam-Data-Company/merchant-config-svc/constant"
	"github.com/Beam-Data-Company/merchant-config-svc/model"
	"github.com/Beam-Data-Company/merchant-config-svc/server/dto"
	"github.com/Beam-Data-Company/merchant-config-svc/server/security"
	"github.com/Beam-Data-Company/merchant-config-svc/service"
	"github.com/Beam-Data-Company/merchant-config-svc/util"
)

// KymHandler is a handler for /kym path
type KymHandler struct {
	util       *util.HandlerUtil
	ra         security.RoleAuthenticator
	kymService service.KymServiceInterface
}

func NewKymHandler(util *util.HandlerUtil, ra security.RoleAuthenticator, kymService service.KymServiceInterface) *KymHandler {
	return &KymHandler{util: util, ra: ra, kymService: kymService}
}

// GetAllKym godoc
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
func (h *KymHandler) GetAllKym(rw http.ResponseWriter, r *http.Request) {

	if err := h.ra.CheckPermission(r, c.OrganisationBeamDataCompany, model.RoleTypeOwner, model.RoleTypeEditor); err != nil {
		h.util.HTTPError(rw, err, http.StatusForbidden)
		return
	}

	status := r.URL.Query().Get("status")

	kymResponse, err := h.kymService.GetAllKym(status)
	if err != nil {
		h.util.WrappedError(rw, err)
		return
	}

	resp, err := json.Marshal(kymResponse)
	if err != nil {
		h.util.WrappedError(rw, fmt.Errorf("unable to marshal json: %w", err))
		return
	}

	rw.WriteHeader(http.StatusOK)
	rw.Write(resp)
}

// GetKym godoc
// @Id GetKym
// @Summary Get Kym Full Detail
// @Description Returns all Kym full detail including document uploaded via downloadURL
// @Tags kym
// @Produce json
// @Accept json
// @Param id path string true "id"
// @Success 200 {object} dto.KymFullDetailResponse "success"
// @Failure default {object} util.APIResponse "fail"
// @Router /kym/{id} [get]
func (h *KymHandler) GetKym(rw http.ResponseWriter, r *http.Request) {

	id, ok := mux.Vars(r)["id"]
	if !ok {
		h.util.HTTPError(rw, errors.New("invalid id in path"), http.StatusBadRequest)
		return
	}

	if err := h.ra.CheckPermission(r, c.OrganisationBeamDataCompany, model.RoleTypeOwner, model.RoleTypeEditor); err != nil {
		h.util.HTTPError(rw, err, http.StatusForbidden)
		return
	}

	kymFullDetailResponse, err := h.kymService.GetKym(id)
	if err != nil {
		h.util.WrappedError(rw, err)
		return
	}

	resp, err := json.Marshal(kymFullDetailResponse)
	if err != nil {
		h.util.WrappedError(rw, fmt.Errorf("unable to marshal json: %w", err))
		return
	}

	rw.WriteHeader(http.StatusOK)
	rw.Write(resp)
}

// AddKym godoc
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
func (h *KymHandler) AddKym(rw http.ResponseWriter, r *http.Request) {
	nKym := &dto.NewKym{}
	if err := json.NewDecoder(r.Body).Decode(nKym); err != nil {
		h.util.HTTPError(rw, fmt.Errorf("error deserializing kym : %w", err), http.StatusBadRequest)
		return
	}

	if err := nKym.Validate(); err != nil {
		h.util.WrappedError(rw, &c.ErrValidation{Violations: err})
		return
	}

	// Check User Authentication
	switch nKym.Source {
	case c.SourceLighthouse:
		if err := h.ra.CheckPermission(r, nKym.OrganisationID, model.RoleTypeOwner); err != nil {
			h.util.HTTPError(rw, err, http.StatusForbidden)
			return
		}
	case c.SourceZort:
		if err := h.ra.CheckPermission(r, c.OrganisationZort, model.RoleTypeOwner); err != nil {
			h.util.HTTPError(rw, err, http.StatusForbidden)
			return
		}
	default:
		h.util.HTTPError(rw, fmt.Errorf("unsupported source %v", nKym.Source), http.StatusBadRequest)
		return
	}

	if err := h.kymService.AddKym(nKym); err != nil {
		h.util.WrappedError(rw, err)
		return
	}

	h.util.HTTPSuccess(rw, "success", http.StatusCreated)
}

// UpdateKymStatus godoc
// @Id UpdateKymStatus
// @Summary Update Kym status and take note if the documents are not completed.
// @Description Admin change status of KYM registration to "approved" or "rejected" and make additional notes
// @Tags kym
// @Produce json
// @Accept json
// @Param id path string true "id"
// @Param requestBody body dto.UpdateKymStatusRequest true "UpdateKymStatusRequest entity"
// @Success 200 {object} util.APIResponse "success"
// @Failure default {object} util.APIResponse "fail"
// @Router /kym/{id}/status [put]
func (h *KymHandler) UpdateKymStatus(rw http.ResponseWriter, r *http.Request) {

	id, ok := mux.Vars(r)["id"]
	if !ok {
		h.util.HTTPError(rw, errors.New("invalid id in path"), http.StatusBadRequest)
		return
	}

	if err := h.ra.CheckPermission(r, c.OrganisationBeamDataCompany, model.RoleTypeOwner, model.RoleTypeEditor); err != nil {
		h.util.HTTPError(rw, err, http.StatusForbidden)
		return
	}

	kymStatusReq := &dto.UpdateKymStatusRequest{}
	if err := json.NewDecoder(r.Body).Decode(kymStatusReq); err != nil {
		h.util.HTTPError(rw, fmt.Errorf("error decode kym status request with error : %w", err), http.StatusBadRequest)
		return
	}

	if err := kymStatusReq.Validate(); err != nil {
		h.util.WrappedError(rw, &c.ErrValidation{Violations: err})
		return
	}

	userInfo := r.Header.Get("X-Endpoint-API-UserInfo")

	if err := h.kymService.UpdateKymStatus(id, kymStatusReq, userInfo); err != nil {
		h.util.WrappedError(rw, err)
		return
	}

	h.util.HTTPSuccess(rw, "success", http.StatusOK)
}
