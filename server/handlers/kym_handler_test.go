package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/gorilla/mux"

	c "github.com/thoniwutr/schedule-school-teachning-bsd13-backend/constant"
	"github.com/thoniwutr/schedule-school-teachning-bsd13-backend/server/dto"
	"github.com/thoniwutr/schedule-school-teachning-bsd13-backend/server/security"
	"github.com/thoniwutr/schedule-school-teachning-bsd13-backend/service"
	"github.com/thoniwutr/schedule-school-teachning-bsd13-backend/util"
)

func TestKymHandler_AddKym(t *testing.T) {
	type fields struct {
		util       *util.HandlerUtil
		ra         security.RoleAuthenticator
		kymService service.KymServiceInterface
	}
	l := util.NewLogger(true)
	tests := []struct {
		name      string
		reqBody   string
		fields    fields
		expStatus int
	}{
		{
			name: "saveKymSuccessFull",
			reqBody: `{
				"partnerRefId": "partner-test-001",
 				"organisationId" : "beamdatacompany",
				"source": "@lighthouse",
				"documentContent": "UEsDBAoAAAAAABITMFMAAAAAAAAAAAAAAAAZAAAATmV3IFRleHQgRG9jdW1lbnQgKDcpLnR4dFBLAQIfAAoAAAAAABITMFMAAAAAAAAAAAAAAAAZACQAAAAAAAAAIAAAAAAAAABOZXcgVGV4dCBEb2N1bWVudCAoNykudHh0CgAgAAAAAAABABgAXm0xUmeq1wFebTFSZ6rXAV5tMVJnqtcBUEsFBgAAAAABAAEAawAAADcAAAAAAA==",
				"businessDetail": {
					"registeredEntityName": "Beam Data Company Limited",
					"businessName": "Beam Data Company",
					"idNumber": "0",
					"address": {
						"address": "string",
						"district": "string",
						"subdistrict": "string",
						"province": "string",
						"city": "string",
						"country": "string",
						"zipcode": "string"
					},
					"phoneNumber": "22134567",
					"businessIndustry": "software",
					"domainName": "https://www.beamcheckout.com/"
				},
				"businessContact": {
					"fullName": "Jenny Jones",
					"role": "accountant",
					"phoneNumber": "988899990",
					"email": "jenny.j@beamcheckout.com"
				},
				"technicalContact": {
					"fullName": "Jenny Jones",
					"role": "accountant",
					"phoneNumber": "988899990",
					"email": "jenny.j@beamcheckout.com"
				},
				"accountingContact": {
					"fullName": "Jenny Jones",
					"role": "accountant",
					"phoneNumber": "988899990",
					"email": "jenny.j@beamcheckout.com"
				},
				"bankTransferDetail": {
					"accountType": "Business",
					"accountName": "Beam Data Company",
					"accountNumber": "23490823904823410",
					"bankName": "Siam Commercial Bank",
					"branch": "Silom"
				},
				"apiKey": {
					"required": true,
					"scope": [
						"string"
					]
				}
			}`,
			fields: fields{
				util: util.NewHandlerUtil(l),
				ra:   &RoleAuthenticatorStub{},
				kymService: KymServiceStub{
					addKym: func(kymRequest *dto.NewKym) error {
						return nil
					},
				},
			},
			expStatus: http.StatusCreated,
		},
		{
			name: "saveKymSuccesFaildWithEmptySource",
			reqBody: `{
				"partnerRefId": "partner-test-001",
 				"organisationId" : "beamdatacompany",
				"documentContent": "UEsDBAoAAAAAABITMFMAAAAAAAAAAAAAAAAZAAAATmV3IFRleHQgRG9jdW1lbnQgKDcpLnR4dFBLAQIfAAoAAAAAABITMFMAAAAAAAAAAAAAAAAZACQAAAAAAAAAIAAAAAAAAABOZXcgVGV4dCBEb2N1bWVudCAoNykudHh0CgAgAAAAAAABABgAXm0xUmeq1wFebTFSZ6rXAV5tMVJnqtcBUEsFBgAAAAABAAEAawAAADcAAAAAAA==",
				"businessDetail": {
					"registeredEntityName": "Beam Data Company Limited",
					"businessName": "Beam Data Company",
					"idNumber": "0",
					"address": {
						"address": "string",
						"district": "string",
						"subdistrict": "string",
						"province": "string",
						"city": "string",
						"country": "string",
						"zipcode": "string"
					},
					"phoneNumber": "22134567",
					"businessIndustry": "software",
					"domainName": "https://www.beamcheckout.com/"
				},
				"businessContact": {
					"fullName": "Jenny Jones",
					"role": "accountant",
					"phoneNumber": "988899990",
					"email": "jenny.j@beamcheckout.com"
				},
				"technicalContact": {
					"fullName": "Jenny Jones",
					"role": "accountant",
					"phoneNumber": "988899990",
					"email": "jenny.j@beamcheckout.com"
				},
				"accountingContact": {
					"fullName": "Jenny Jones",
					"role": "accountant",
					"phoneNumber": "988899990",
					"email": "jenny.j@beamcheckout.com"
				},
				"bankTransferDetail": {
					"accountType": "Business",
					"accountName": "Beam Data Company",
					"accountNumber": "23490823904823410",
					"bankName": "Siam Commercial Bank",
					"branch": "Silom"
				},
				"apiKey": {
					"required": true,
					"scope": [
						"string"
					]
				}
			}`,
			fields: fields{
				util: util.NewHandlerUtil(l),
				ra:   &RoleAuthenticatorStub{},
				kymService: KymServiceStub{
					addKym: func(kymRequest *dto.NewKym) error {
						return nil
					},
				},
			},
			expStatus: http.StatusBadRequest,
		},
		{
			name: "saveKymSuccesFaildWithUnsupportedSource",
			reqBody: `{
				"partnerRefId": "partner-test-001",
 				"organisationId" : "beamdatacompany",
				"source" : "test-source"
				"documentContent": "UEsDBAoAAAAAABITMFMAAAAAAAAAAAAAAAAZAAAATmV3IFRleHQgRG9jdW1lbnQgKDcpLnR4dFBLAQIfAAoAAAAAABITMFMAAAAAAAAAAAAAAAAZACQAAAAAAAAAIAAAAAAAAABOZXcgVGV4dCBEb2N1bWVudCAoNykudHh0CgAgAAAAAAABABgAXm0xUmeq1wFebTFSZ6rXAV5tMVJnqtcBUEsFBgAAAAABAAEAawAAADcAAAAAAA==",
				"businessDetail": {
					"registeredEntityName": "Beam Data Company Limited",
					"businessName": "Beam Data Company",
					"idNumber": "0",
					"address": {
						"address": "string",
						"district": "string",
						"subdistrict": "string",
						"province": "string",
						"city": "string",
						"country": "string",
						"zipcode": "string"
					},
					"phoneNumber": "22134567",
					"businessIndustry": "software",
					"domainName": "https://www.beamcheckout.com/"
				},
				"businessContact": {
					"fullName": "Jenny Jones",
					"role": "accountant",
					"phoneNumber": "988899990",
					"email": "jenny.j@beamcheckout.com"
				},
				"technicalContact": {
					"fullName": "Jenny Jones",
					"role": "accountant",
					"phoneNumber": "988899990",
					"email": "jenny.j@beamcheckout.com"
				},
				"accountingContact": {
					"fullName": "Jenny Jones",
					"role": "accountant",
					"phoneNumber": "988899990",
					"email": "jenny.j@beamcheckout.com"
				},
				"bankTransferDetail": {
					"accountType": "Business",
					"accountName": "Beam Data Company",
					"accountNumber": "23490823904823410",
					"bankName": "Siam Commercial Bank",
					"branch": "Silom"
				},
				"apiKey": {
					"required": true,
					"scope": [
						"string"
					]
				}
			}`,
			fields: fields{
				util: util.NewHandlerUtil(l),
				ra:   &RoleAuthenticatorStub{},
				kymService: KymServiceStub{
					addKym: func(kymRequest *dto.NewKym) error {
						return nil
					},
				},
			},
			expStatus: http.StatusBadRequest,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &KymHandler{
				util:       tt.fields.util,
				ra:         tt.fields.ra,
				kymService: tt.fields.kymService,
			}

			req := httptest.NewRequest(http.MethodPost, "/kym", strings.NewReader(tt.reqBody))
			rr := httptest.NewRecorder()

			router := mux.NewRouter()
			router.Methods(http.MethodPost).Path("/kym").HandlerFunc(h.AddKym)
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

func TestKymHandler_GetAllKym(t *testing.T) {
	type fields struct {
		util       *util.HandlerUtil
		ra         security.RoleAuthenticator
		kymService service.KymServiceInterface
	}
	l := util.NewLogger(true)
	tests := []struct {
		name      string
		fields    fields
		expStatus int
	}{
		{
			name: "getAllKymSuccessFull",
			fields: fields{
				util: util.NewHandlerUtil(l),
				ra:   &RoleAuthenticatorStub{},
				kymService: KymServiceStub{
					getAllKym: func(status string) ([]*dto.KymResponse, error) {
						return []*dto.KymResponse{
							{
								ID:                  "GS-16SZ-oOu5iFpn9SMlj",
								BusinessName:        "business-name-example",
								ContactPerson:       "business-name-example",
								PhoneNumber:         "000000000",
								DatetimeCreated:     time.Now(),
								DocumentDownloadURL: "test-url",
								Status:              "approved",
							},
						}, nil
					},
				},
			},
			expStatus: http.StatusOK,
		},
		{
			name: "getAllKymFailedWithDBNoSuchEntity",
			fields: fields{
				util: util.NewHandlerUtil(l),
				ra:   &RoleAuthenticatorStub{},
				kymService: KymServiceStub{
					getAllKym: func(status string) ([]*dto.KymResponse, error) {
						return []*dto.KymResponse{}, c.ErrDBNoSuchEntity
					},
				},
			},
			expStatus: http.StatusNotFound,
		},
		{
			name: "getAllKymFailedWithErrorFormat",
			fields: fields{
				util: util.NewHandlerUtil(l),
				ra:   &RoleAuthenticatorStub{},
				kymService: KymServiceStub{
					getAllKym: func(status string) ([]*dto.KymResponse, error) {
						return nil, errors.New("some internal db failure")
					},
				},
			},
			expStatus: http.StatusInternalServerError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &KymHandler{
				util:       tt.fields.util,
				ra:         tt.fields.ra,
				kymService: tt.fields.kymService,
			}

			req := httptest.NewRequest(http.MethodGet, "/kym", nil)
			rr := httptest.NewRecorder()

			router := mux.NewRouter()
			router.Methods(http.MethodGet).Path("/kym").HandlerFunc(h.GetAllKym)
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

func TestKymHandler_GetKym(t *testing.T) {
	type fields struct {
		util       *util.HandlerUtil
		ra         security.RoleAuthenticator
		kymService service.KymServiceInterface
	}
	l := util.NewLogger(true)
	tests := []struct {
		name      string
		id        string
		fields    fields
		expStatus int
	}{
		{
			name: "getKymSuccessFull",
			id:   "test-id",
			fields: fields{
				util: util.NewHandlerUtil(l),
				ra:   &RoleAuthenticatorStub{},
				kymService: KymServiceStub{
					getKym: func(id string) (*dto.KymFullDetailResponse, error) {
						return &dto.KymFullDetailResponse{
							ID:             "GS-16SZ-oOu5iFpn9SMlj",
							OrganisationID: "testing-organisation",
							PartnerRefID:   "testing-partner-ref-id",
							BusinessDetail: dto.BusinessDetail{
								BusinessName:     "business-name-example",
								BusinessIndustry: "business-industry-example",
								DomainName:       "https://www.beamcheckout.com/",
								IDNumber:         "0000000000000",
								PhoneNumber:      "000000000",
								Address: dto.MerchantAddress{
									City:        "city",
									Country:     "country",
									District:    "district",
									HouseNumber: "123",
									Province:    "province",
									Street:      "street",
									Subdistrict: "subdistrict",
									Zipcode:     "12345",
								},
							},
							BusinessContact: dto.PointOfContact{
								FullName:    "test-organisation",
								Role:        "role",
								PhoneNumber: "000000000",
								Email:       "test@example.com",
							},
							TechnicalContact: dto.PointOfContact{
								FullName:    "test-organisation",
								Role:        "role",
								PhoneNumber: "000000000",
								Email:       "test@example.com",
							},
							AccountingContact: dto.PointOfContact{
								FullName:    "test-organisation",
								Role:        "role",
								PhoneNumber: "000000000",
								Email:       "test@example.com",
							},
							BankTransferDetail: dto.BankTransferDetail{
								AccountName:   "accname",
								BankName:      "bankname",
								Branch:        "branch",
								AccountNumber: "0000000000000000",
								AccountType:   "type",
							},
							DocumentDownloadURL: "https://test.com",
							ApiKey: dto.ApiKey{
								Required: true,
								Scope:    []string{},
							},
							Source:          "lighthouse",
							ImageURL:        "example-image-url",
							DatetimeCreated: time.Now(),
							Status:          "approved",
							Notes:           "",
						}, nil
					},
				},
			},
			expStatus: http.StatusOK,
		},
		{
			name: "getKymFailedWithNotFoundEnitity",
			id:   "test-id",
			fields: fields{
				util: util.NewHandlerUtil(l),
				ra:   &RoleAuthenticatorStub{},
				kymService: KymServiceStub{
					getKym: func(id string) (*dto.KymFullDetailResponse, error) {
						return &dto.KymFullDetailResponse{}, c.ErrDBNoSuchEntity
					},
				},
			},
			expStatus: http.StatusNotFound,
		},
		{
			name: "getKymFailedWithErrorFormat",
			id:   "test-id",
			fields: fields{
				util: util.NewHandlerUtil(l),
				ra:   &RoleAuthenticatorStub{},
				kymService: KymServiceStub{
					getKym: func(id string) (*dto.KymFullDetailResponse, error) {
						return nil, errors.New("some internal db failure")
					},
				},
			},
			expStatus: http.StatusInternalServerError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &KymHandler{
				util:       tt.fields.util,
				ra:         tt.fields.ra,
				kymService: tt.fields.kymService,
			}

			req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/kym/%s", tt.id), nil)
			rr := httptest.NewRecorder()

			router := mux.NewRouter()
			router.Methods(http.MethodGet).Path("/kym/{id}").HandlerFunc(h.GetKym)
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

func TestKymHandler_UpdateKymStatus(t *testing.T) {
	type fields struct {
		util       *util.HandlerUtil
		ra         security.RoleAuthenticator
		kymService service.KymServiceInterface
	}
	l := util.NewLogger(true)
	tests := []struct {
		name      string
		id        string
		reqBody   string
		fields    fields
		expStatus int
	}{
		{
			name: "getKymSuccessFull",
			id:   "test-id",
			reqBody: `{
		         "status": "approved",
 				 "notes": ""
			}`,
			fields: fields{
				util: util.NewHandlerUtil(l),
				ra:   &RoleAuthenticatorStub{},
				kymService: KymServiceStub{
					updateKymStatus: func(id string, kymStatusReq *dto.UpdateKymStatusRequest, userInfo string) error {
						return nil
					},
				},
			},
			expStatus: http.StatusOK,
		},
		{
			name: "getKymSuccessFailedWithInvalidRequest",
			id:   "test-id",
			reqBody: `{
		         abc
			}`,
			fields: fields{
				util: util.NewHandlerUtil(l),
				ra:   &RoleAuthenticatorStub{},
				kymService: KymServiceStub{
					updateKymStatus: func(id string, kymStatusReq *dto.UpdateKymStatusRequest, userInfo string) error {
						return nil
					},
				},
			},
			expStatus: http.StatusBadRequest,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &KymHandler{
				util:       tt.fields.util,
				ra:         tt.fields.ra,
				kymService: tt.fields.kymService,
			}

			req := httptest.NewRequest(http.MethodPut, fmt.Sprintf("/kym/%s/status", tt.id), strings.NewReader(tt.reqBody))
			rr := httptest.NewRecorder()

			router := mux.NewRouter()
			router.Methods(http.MethodPut).Path("/kym/{id}/status").HandlerFunc(h.UpdateKymStatus)
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

type KymServiceStub struct {
	addKym          func(kymRequest *dto.NewKym) error
	getKym          func(id string) (*dto.KymFullDetailResponse, error)
	getAllKym       func(status string) ([]*dto.KymResponse, error)
	updateKymStatus func(id string, kymStatusReq *dto.UpdateKymStatusRequest, userInfo string) error
}

func (stub KymServiceStub) AddKym(kymRequest *dto.NewKym) error {
	return stub.addKym(kymRequest)
}

func (stub KymServiceStub) GetKym(id string) (*dto.KymFullDetailResponse, error) {
	return stub.getKym(id)
}

func (stub KymServiceStub) GetAllKym(status string) ([]*dto.KymResponse, error) {
	return stub.getAllKym(status)
}

func (stub KymServiceStub) UpdateKymStatus(id string, kymStatusReq *dto.UpdateKymStatusRequest, userInfo string) error {
	return stub.updateKymStatus(id, kymStatusReq, userInfo)
}
