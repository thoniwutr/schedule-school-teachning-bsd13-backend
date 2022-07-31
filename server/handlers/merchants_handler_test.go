package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	c "github.com/Beam-Data-Company/merchant-config-svc/constant"
	"github.com/Beam-Data-Company/merchant-config-svc/messaging"
	"github.com/Beam-Data-Company/merchant-config-svc/server/security"
	"github.com/Beam-Data-Company/merchant-config-svc/service"
	"github.com/Beam-Data-Company/merchant-config-svc/util"
	"github.com/gorilla/mux"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/Beam-Data-Company/merchant-config-svc/model"
	"github.com/Beam-Data-Company/merchant-config-svc/server/dto"
)

func TestMerchantsHandler_GetMerchant(t *testing.T) {
	type fields struct {
		util *util.HandlerUtil
		ra   security.RoleAuthenticator
		pub  messaging.Publisher
		ms   service.MerchantServiceInterface
	}
	type args struct {
		rw http.ResponseWriter
		r  *http.Request
	}
	l := util.NewLogger(true)
	tests := []struct {
		name       string
		fields     fields
		merchantID string
		expStatus  int
	}{
		{
			name: "successfulGet",
			fields: fields{
				util: util.NewHandlerUtil(l),
				ra:   &RoleAuthenticatorStub{},
				pub:  testPublisher{},
				ms: MerchantServiceStub{
					getMerchant: func(id string) (m *model.Merchant, err error) {
						return &model.Merchant{
							Address: model.Address{
								City:        "bkk",
								Country:     "thailand",
								District:    "district",
								HouseNumber: "123",
								Province:    "province",
								Street:      "street",
								Subdistrict: "subdistrict",
								Zipcode:     "12345",
							},
							ContactNumber:  "0123456789",
							Email:          "test@example.com",
							FullName:       "test-organisation",
							OrganisationID: "m1",
							MerchantID:     "m1",
							ShortName:      "wahey",
							Created:        time.Now(),
							Updated:        time.Now(),
						}, nil
					},
				},
			},
			merchantID: "m1",
			expStatus:  http.StatusOK,
		},
		{
			name: "merchantNotFound",
			fields: fields{
				util: util.NewHandlerUtil(l),
				ra:   &RoleAuthenticatorStub{},
				pub:  testPublisher{},
				ms: MerchantServiceStub{
					getMerchant: func(id string) (m *model.Merchant, err error) {
						return nil, c.ErrDBNoSuchEntity
					},
				},
			},
			merchantID: "notfound",
			expStatus:  http.StatusNotFound,
		},
		{
			name: "internalServerError",
			fields: fields{
				util: util.NewHandlerUtil(l),
				ra:   &RoleAuthenticatorStub{},
				pub:  testPublisher{},
				ms: MerchantServiceStub{
					getMerchant: func(id string) (m *model.Merchant, err error) {
						return nil, errors.New("some internal db failure")
					},
				},
			},
			merchantID: "503",
			expStatus:  http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mh := &MerchantsHandler{
				util: tt.fields.util,
				ra:   tt.fields.ra,
				pub:  tt.fields.pub,
				ms:   tt.fields.ms,
			}

			req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/merchants/%s", tt.merchantID), nil)
			rr := httptest.NewRecorder()

			router := mux.NewRouter()
			router.Methods(http.MethodGet).Path("/merchants/{merchantId:[a-zA-Z0-9]+}").HandlerFunc(mh.GetMerchant)
			router.ServeHTTP(rr, req)
			resp := rr.Result()

			if tt.expStatus != resp.StatusCode {
				t.Errorf("unexpected status code: got %v want %v", resp.StatusCode, tt.expStatus)
			}

			if isHTTPSuccess(resp.StatusCode) {
				mget := &dto.MerchantGetResponse{}
				if err := json.NewDecoder(rr.Body).Decode(mget); err != nil {
					t.Errorf("could not unmarshal response to dto type: %v", err)
				}
			}
		})
	}
}

func TestMerchantsHandler_AddMerchant(t *testing.T) {
	type fields struct {
		util *util.HandlerUtil
		ra   security.RoleAuthenticator
		pub  messaging.Publisher
		ms   service.MerchantServiceInterface
	}
	l := util.NewLogger(true)
	tests := []struct {
		name      string
		fields    fields
		reqBody   string
		expStatus int
	}{
		{
			name: "successfulAdd",
			fields: fields{
				util: util.NewHandlerUtil(l),
				ra:   &RoleAuthenticatorStub{},
				pub:  testPublisher{},
				ms: MerchantServiceStub{
					addMerchant: func(nm *dto.NewMerchant, kf *dto.KymField) error {
						return nil
					},
				},
			},
			reqBody: `{
				"organisationId": "org123",
				"shortName": "wahey",
				"fullName": "test-organisation",
				"email": "user@example.com",
				"contactNumber": "123456789",
				"address": {
				  "houseNumber": "string",
				  "street": "string",
				  "district": "string",
				  "subdistrict": "string",
				  "city": "bkk",
				  "province": "string",
				  "country": "thailand",
				  "zipcode": "12345"
				},
				"logoUrl" : "abc",
				"currencyCode" : "THB",
				"availablePaymentMethods" : ["creditCard"]
			}`,
			expStatus: http.StatusCreated,
		},
		{
			name: "validationError",
			fields: fields{
				util: util.NewHandlerUtil(l),
				ra:   &RoleAuthenticatorStub{},
				pub:  testPublisher{},
				ms: MerchantServiceStub{
					addMerchant: func(nm *dto.NewMerchant, kf *dto.KymField) error {
						return nil
					},
				},
			},
			reqBody: `{
				"hello":"world"
			}`,
			expStatus: http.StatusBadRequest,
		},
		{
			name: "unableToDecode",
			fields: fields{
				util: util.NewHandlerUtil(l),
				ra:   &RoleAuthenticatorStub{},
				pub:  testPublisher{},
				ms: MerchantServiceStub{
					addMerchant: func(nm *dto.NewMerchant, kf *dto.KymField) error {
						return nil
					},
				},
			},
			reqBody:   `{`,
			expStatus: http.StatusBadRequest,
		},
		{
			name: "addDuplicate",
			fields: fields{
				util: util.NewHandlerUtil(l),
				ra:   &RoleAuthenticatorStub{},
				pub:  testPublisher{},
				ms: MerchantServiceStub{
					addMerchant: func(nm *dto.NewMerchant, kf *dto.KymField) error {
						return c.ErrDBEntityAlreadyExists
					},
				},
			},
			reqBody: `{
				"organisationId": "org123",
				"shortName": "wahey",
				"fullName": "test-organisation",
				"email": "user@example.com",
				"contactNumber": "123456789",
				"address": {
				  "houseNumber": "string",
				  "street": "string",
				  "district": "string",
				  "subdistrict": "string",
				  "city": "bkk",
				  "province": "string",
				  "country": "thailand",
				  "zipcode": "12345"
				},
				"logoUrl" : "abc",
				"currencyCode" : "THB",
				"availablePaymentMethods" : ["creditCard"]
			}`,
			expStatus: http.StatusConflict,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mh := &MerchantsHandler{
				util: tt.fields.util,
				ra:   tt.fields.ra,
				pub:  tt.fields.pub,
				ms:   tt.fields.ms,
			}

			req := httptest.NewRequest(http.MethodPost, "/merchants", strings.NewReader(tt.reqBody))
			rr := httptest.NewRecorder()

			router := mux.NewRouter()
			router.Methods(http.MethodPost).Path("/merchants").HandlerFunc(mh.AddMerchant)
			router.ServeHTTP(rr, req)
			resp := rr.Result()

			if tt.expStatus != resp.StatusCode {
				t.Errorf("unexpected status code: got %v want %v", resp.StatusCode, tt.expStatus)
			}

			if !isHTTPSuccess(resp.StatusCode) {
				apiresp := &util.APIResponse{}
				if err := json.NewDecoder(rr.Body).Decode(apiresp); err != nil {
					t.Errorf("error unmarshelling response: %v", err)
				}
				t.Logf("[debug] api response: %v", apiresp)
			}

		})
	}
}

func isHTTPSuccess(code int) bool {
	return code >= 200 && code < 300
}

type MerchantServiceStub struct {
	addMerchant        func(nm *dto.NewMerchant, kf *dto.KymField) error
	getMerchant        func(id string) (m *model.Merchant, err error)
	updateMerchant     func(m *dto.Merchant) error
	upsertPayOutConfig func(merchantID string, poc *dto.PayOutConfig) error
}

func (stub MerchantServiceStub) AddMerchant(nm *dto.NewMerchant, kf *dto.KymField) error {
	return stub.addMerchant(nm, kf)
}

func (stub MerchantServiceStub) GetMerchant(id string) (m *model.Merchant, err error) {
	return stub.getMerchant(id)
}

func (stub MerchantServiceStub) UpdateMerchant(m *dto.Merchant) error {
	return stub.updateMerchant(m)
}

func (stub MerchantServiceStub) UpsertPayOutConfig(merchantID string, poc *dto.PayOutConfig) error {
	return stub.upsertPayOutConfig(merchantID, poc)
}

type RoleAuthenticatorStub struct{}

func (r RoleAuthenticatorStub) CheckPermission(_ *http.Request, _ string, _ ...model.RoleType) error {
	return nil
}

type testPublisher struct {
}

func (p testPublisher) PublishMerchant(_ *dto.MerchantPublish) error {
	return nil
}
