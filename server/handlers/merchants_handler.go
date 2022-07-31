package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"

	c "github.com/Beam-Data-Company/merchant-config-svc/constant"
	"github.com/Beam-Data-Company/merchant-config-svc/messaging"
	"github.com/Beam-Data-Company/merchant-config-svc/model"
	"github.com/Beam-Data-Company/merchant-config-svc/server/dto"
	"github.com/Beam-Data-Company/merchant-config-svc/server/security"
	"github.com/Beam-Data-Company/merchant-config-svc/service"
	"github.com/Beam-Data-Company/merchant-config-svc/util"
)

// MerchantsHandler is a handler for /merchants path
type MerchantsHandler struct {
	util *util.HandlerUtil
	ra   security.RoleAuthenticator
	pub  messaging.Publisher
	ms   service.MerchantServiceInterface
}

// NewMerchantsHandler a new Merchants handler instance
func NewMerchantsHandler(util *util.HandlerUtil, ra security.RoleAuthenticator, pub messaging.Publisher, ms service.MerchantServiceInterface) *MerchantsHandler {
	return &MerchantsHandler{util: util, ra: ra, pub: pub, ms: ms}
}

func (mh *MerchantsHandler) HealthCheck(w http.ResponseWriter, _ *http.Request) {
	mh.util.HTTPSuccess(w, "merchant healthy", http.StatusOK)
}

// GetMerchant godoc
// @Id GetMerchant
// @Summary retrieves the Merchant from the merchantId specified in the path
// @Description Returns a Merchant entity given the merchantId
// @Tags merchants
// @Produce json
// @Accept json
// @Param merchantId path string true "merchantId"
// @Success 200 {object} dto.MerchantGetResponse "success"
// @Failure default {object} util.APIResponse "fail"
// @Router /merchants/{merchantId} [get]
func (mh *MerchantsHandler) GetMerchant(rw http.ResponseWriter, r *http.Request) {
	id, ok := mux.Vars(r)["merchantId"]
	if !ok {
		mh.util.HTTPError(rw, errors.New("invalid merchantId in path"), http.StatusNotFound)
		return
	}

	if err := mh.ra.CheckPermission(r, id, model.RoleTypeOwner, model.RoleTypeEditor, model.RoleTypeViewer); err != nil {
		mh.util.HTTPError(rw, err, http.StatusForbidden)
		return
	}

	m, err := mh.ms.GetMerchant(id)
	if err != nil {
		mh.util.WrappedError(rw, err)
		return
	}

	resp, err := json.Marshal(dto.ToMerchantDTO(m))
	if err != nil {
		mh.util.WrappedError(rw, fmt.Errorf("unable to marshal json: %w", err))
		return
	}
	rw.WriteHeader(http.StatusOK)
	rw.Write(resp)
}

// AddMerchant godoc
// @Id AddMerchant
// @Summary Register Merchant
// @Description Creates a new Merchant entity to be able to accept payments through Beam
// @Tags merchants
// @Produce json
// @Accept json
// @Param requestBody body dto.NewMerchant true "NewMerchant entity"
// @Success 201 {object} util.APIResponse "success"
// @Failure default {object} util.APIResponse "fail"
// @Router /merchants [post]
func (mh *MerchantsHandler) AddMerchant(rw http.ResponseWriter, r *http.Request) {
	nm := &dto.NewMerchant{}
	if err := json.NewDecoder(r.Body).Decode(nm); err != nil {
		mh.util.HTTPError(rw, fmt.Errorf("error deserializing merchant : %w", err), http.StatusBadRequest)
		return
	}

	if err := nm.Validate(); err != nil {
		mh.util.WrappedError(rw, &c.ErrValidation{Violations: err})
		return
	}

	if err := mh.ra.CheckPermission(r, nm.OrganisationID, model.RoleTypeOwner, model.RoleTypeEditor); err != nil {
		mh.util.HTTPError(rw, err, http.StatusForbidden)
		return
	}

	if err := mh.ms.AddMerchant(nm, nil); err != nil {
		mh.util.WrappedError(rw, err)
		return
	}

	mh.util.HTTPSuccess(rw, "success", http.StatusCreated)
}

// UpdateMerchant godoc
// @Id UpdateMerchant
// @Summary Update Merchant
// @Description Updates the Merchant given an existing merchant ID. Note that some fields are not allowed to be updatable.
// @Tags merchants
// @Produce json
// @Accept json
// @Param merchantId path string true "merchantId"
// @Param requestBody body dto.Merchant true "Merchant entity"
// @Success 202 {object} util.APIResponse "success"
// @Failure default {object} util.APIResponse "fail"
// @Router /merchants/{merchantId} [put]
func (mh *MerchantsHandler) UpdateMerchant(rw http.ResponseWriter, r *http.Request) {
	id, ok := mux.Vars(r)["merchantId"]
	if !ok {
		mh.util.HTTPError(rw, errors.New("invalid merchantId in path"), http.StatusNotFound)
		return
	}

	m := &dto.Merchant{}
	if err := json.NewDecoder(r.Body).Decode(m); err != nil {
		mh.util.HTTPError(rw, fmt.Errorf("error deserializing merchant : %w", err), http.StatusBadRequest)
		return
	}

	if err := m.Validate(); err != nil {
		mh.util.WrappedError(rw, &c.ErrValidation{Violations: err})
		return
	}

	if id != m.MerchantID {
		mh.util.HTTPError(rw, errors.New("merchantId does not match"), http.StatusBadRequest)
		return
	}

	if err := mh.ra.CheckPermission(r, id, model.RoleTypeOwner, model.RoleTypeEditor); err != nil {
		mh.util.HTTPError(rw, err, http.StatusForbidden)
		return
	}

	if err := mh.ms.UpdateMerchant(m); err != nil {
		mh.util.WrappedError(rw, err)
		return
	}

	mh.util.HTTPSuccess(rw, "success", http.StatusAccepted)
}

// UpsertPayOutConfig godoc
// @Id UpsertPayOutConfig
// @Summary allows updating or inserting PayOutConfig
// @Description Update Merchant PayOutConfig
// @Tags merchants
// @Produce json
// @Accept json
// @Param merchantId path string true "merchantId"
// @Param requestBody body dto.PayOutConfig true "PayInConfig entity"
// @Success 202 {object} util.APIResponse "success"
// @Failure default {object} util.APIResponse "fail"
// @Router /merchants/{merchantId}/pay-out-config [post]
func (mh *MerchantsHandler) UpsertPayOutConfig(rw http.ResponseWriter, r *http.Request) {
	id, ok := mux.Vars(r)["merchantId"]
	if !ok {
		mh.util.HTTPError(rw, errors.New("invalid merchantId in path"), http.StatusNotFound)
		return
	}

	poc := &dto.PayOutConfig{}
	if err := json.NewDecoder(r.Body).Decode(poc); err != nil {
		mh.util.HTTPError(rw, fmt.Errorf("error deserializing PayInConfig : %w", err), http.StatusBadRequest)
		return
	}

	if err := poc.Validate(); err != nil {
		mh.util.WrappedError(rw, &c.ErrValidation{Violations: err})
		return
	}

	if err := mh.ra.CheckPermission(r, id, model.RoleTypeOwner, model.RoleTypeEditor); err != nil {
		mh.util.HTTPError(rw, err, http.StatusForbidden)
		return
	}

	if err := mh.ms.UpsertPayOutConfig(id, poc); err != nil {
		mh.util.WrappedError(rw, err)
		return
	}

	mh.util.HTTPSuccess(rw, "success", http.StatusAccepted)
}

// GetPayOutConfig godoc
// @Id GetPayOutConfig
// @Summary retrieve PayOutConfig
// @Description Returns a Merchant entity retrieves the Merchant's PayOutConfig
// @Tags merchants
// @Produce json
// @Accept json
// @Param merchantId path string true "merchantId"
// @Success 200 {array} dto.PayOutConfig "success"
// @Failure default {object} util.APIResponse "fail"
// @Router /merchants/{merchantId}/pay-out-config [get]
func (mh *MerchantsHandler) GetPayOutConfig(rw http.ResponseWriter, r *http.Request) {

	id, ok := mux.Vars(r)["merchantId"]
	if !ok {
		mh.util.HTTPError(rw, errors.New("invalid merchantId in path"), http.StatusNotFound)
		return
	}

	if err := mh.ra.CheckPermission(r, id, model.RoleTypeOwner, model.RoleTypeEditor, model.RoleTypeViewer); err != nil {
		mh.util.HTTPError(rw, err, http.StatusForbidden)
		return
	}

	m, err := mh.ms.GetMerchant(id)
	if err != nil {
		mh.util.WrappedError(rw, fmt.Errorf("unable to find merchant: %w", err))
		return
	}

	merchantDTO := dto.ToMerchantDTO(m)

	resp, err := json.Marshal(merchantDTO.PayOutConfig)
	if err != nil {
		mh.util.WrappedError(rw, fmt.Errorf("unable to marshal json: %w", err))
		return
	}
	rw.WriteHeader(http.StatusOK)
	rw.Write(resp)
}
