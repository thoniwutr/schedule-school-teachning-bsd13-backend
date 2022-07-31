package service

import (
	"reflect"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	c "github.com/thoniwutr/-schedule-school-teachning-bsd13-backend/constant"
	"github.com/thoniwutr/-schedule-school-teachning-bsd13-backend/db"
	"github.com/thoniwutr/-schedule-school-teachning-bsd13-backend/model"
	"github.com/thoniwutr/-schedule-school-teachning-bsd13-backend/server/dto"
	"github.com/thoniwutr/-schedule-school-teachning-bsd13-backend/thirdparty"
	"github.com/thoniwutr/-schedule-school-teachning-bsd13-backend/util"
)

func TestKymService_AddKym(t *testing.T) {
	type fields struct {
		db db.KymDB
		cs util.CloudStorageManager
		ks thirdparty.KymServiceClient
		ms MerchantServiceInterface
	}
	type args struct {
		kymRequest *dto.NewKym
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		asserts func(t *testing.T, gotErr error)
	}{
		{
			name: "addKymSuccessful",
			fields: fields{
				db: KymDBStub{
					addKym: func(kym *model.Kym) error {
						return nil
					},
				},
				cs: CloudStorageManagerStub{},
			},
			args: args{
				kymRequest: &dto.NewKym{
					PartnerRefID: "beam-test-001",
					BusinessDetail: dto.BusinessDetail{
						RegisteredEntityName: "Beam",
						BusinessName:         "Beam Data Company",
						BusinessIndustry:     "Payment Service",
						DomainName:           "https://www.beamcheckout.com/",
						IDNumber:             "0000000000000",
						PhoneNumber:          "000000000",
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
					DocumentContent: "UEsDBAoAAAAAABITMFMAAAAAAAAAAAAAAAAZAAAATmV3IFRleHQgRG9jdW1lbnQgKDcpLnR4dFBLAQIfAAoAAAAAABITMFMAAAAAAAAAAAAAAAAZACQAAAAAAAAAIAAAAAAAAABOZXcgVGV4dCBEb2N1bWVudCAoNykudHh0CgAgAAAAAAABABgAXm0xUmeq1wFebTFSZ6rXAV5tMVJnqtcBUEsFBgAAAAABAAEAawAAADcAAAAAAA==",
					ApiKey: dto.ApiKey{
						Required: true,
						Scope:    []string{},
					},
					Source: "lighthouse",
				},
			},
			asserts: func(t *testing.T, gotErr error) {
				assert.NoError(t, gotErr)
			},
		},
		{
			name: "addKymFailedWithDecodeDocumentContent",
			fields: fields{
				db: KymDBStub{
					addKym: func(kym *model.Kym) error {
						return nil
					},
				},
				cs: CloudStorageManagerStub{},
			},
			args: args{
				kymRequest: &dto.NewKym{
					PartnerRefID: "beam-test-001",
					BusinessDetail: dto.BusinessDetail{
						RegisteredEntityName: "Beam",
						BusinessName:         "Beam Data Company",
						BusinessIndustry:     "Payment Service",
						DomainName:           "https://www.beamcheckout.com/",
						IDNumber:             "0000000000000",
						PhoneNumber:          "000000000",
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
					DocumentContent: "12345",
					ApiKey: dto.ApiKey{
						Required: true,
						Scope:    []string{},
					},
					Source: "lighthouse",
				},
			},
			asserts: func(t *testing.T, gotErr error) {
				assert.EqualError(t, gotErr, "failed to decode base64 : illegal base64 data at input byte 4")
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &KymService{
				db: tt.fields.db,
				cs: tt.fields.cs,
				ks: tt.fields.ks,
				ms: tt.fields.ms,
			}
			err := s.AddKym(tt.args.kymRequest)
			tt.asserts(t, err)
		})
	}
}

func TestKymService_GetAllKym(t *testing.T) {
	type fields struct {
		db db.KymDB
		cs util.CloudStorageManager
		ks thirdparty.KymServiceClient
		ms MerchantServiceInterface
	}
	tests := []struct {
		name    string
		status  string
		fields  fields
		want    []*dto.KymResponse
		asserts func(t *testing.T, gotErr error)
	}{
		{
			name: "getAllKymSuccessful",
			fields: fields{
				db: KymDBStub{
					getAllKym: func(status string) ([]*model.Kym, error) {
						return []*model.Kym{
							{
								ID:             "GS-16SZ-oOu5iFpn9SMlj",
								OrganisationID: "test-organisation",
								PartnerRefID:   "beam-test-001",
								BusinessDetail: model.BusinessDetail{
									RegisteredEntityName: "Beam",
									BusinessName:         "Beam Data Company",
									BusinessIndustry:     "Payment Service",
									DomainName:           "https://www.beamcheckout.com/",
									IDNumber:             "0000000000000",
									PhoneNumber:          "000000000",
									Address: model.Address{
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
								BusinessContact: model.PointOfContact{
									FullName:    "test-organisation",
									Role:        "role",
									PhoneNumber: "000000000",
									Email:       "test@example.com",
								},
								TechnicalContact: model.PointOfContact{
									FullName:    "test-organisation",
									Role:        "role",
									PhoneNumber: "000000000",
									Email:       "test@example.com",
								},
								AccountingContact: model.PointOfContact{
									FullName:    "test-organisation",
									Role:        "role",
									PhoneNumber: "000000000",
									Email:       "test@example.com",
								},
								BankTransferDetail: model.BankTransferDetail{
									AccountName:   "accname",
									BankName:      "bankname",
									Branch:        "branch",
									AccountNumber: "0000000000000000",
									AccountType:   "type",
								},
								DocumentDownloadURL: "https://test.com",
								ApiKey: model.ApiKey{
									Required: true,
									Scope:    []string{},
								},
								Source: "lighthouse",
								Status: "approved",
								Notes:  "",
							},
						}, nil
					},
				},
			},
			want: []*dto.KymResponse{
				{
					ID:                  "GS-16SZ-oOu5iFpn9SMlj",
					BusinessName:        "Beam Data Company",
					ContactPerson:       "Beam",
					PhoneNumber:         "000000000",
					DocumentDownloadURL: "https://test.com",
					Status:              "approved",
				},
			},
			asserts: func(t *testing.T, gotErr error) {
				assert.NoError(t, gotErr)
			},
		},
		{
			name: "getAllKymFailedWithDBNoSuchEntity",
			fields: fields{
				db: KymDBStub{
					getAllKym: func(status string) ([]*model.Kym, error) {
						return nil, c.ErrDBNoSuchEntity
					},
				},
			},
			want: nil,
			asserts: func(t *testing.T, gotErr error) {
				assert.EqualError(t, gotErr, "no such entity")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &KymService{
				db: tt.fields.db,
				cs: tt.fields.cs,
				ks: tt.fields.ks,
				ms: tt.fields.ms,
			}
			kymList, err := s.GetAllKym(tt.status)
			if !reflect.DeepEqual(kymList, tt.want) {
				t.Errorf("GetAllKym() got = %v, want %v", kymList, tt.want)
			}
			tt.asserts(t, err)
		})
	}
}

func TestKymService_UpdateKymStatus(t *testing.T) {
	type fields struct {
		db db.KymDB
		cs util.CloudStorageManager
		ks thirdparty.KymServiceClient
		ms MerchantServiceInterface
	}
	type args struct {
		id           string
		kymStatusReq *dto.UpdateKymStatusRequest
		userInfo     string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		asserts func(t *testing.T, gotErr error)
	}{
		{
			name: "updateKymStatusSuccessful",
			fields: fields{
				db: KymDBStub{
					getKym: func(id string) (*model.Kym, error) {
						return &model.Kym{
							ID:           "GS-16SZ-oOu5iFpn9SMlj",
							PartnerRefID: "beam-test-001",
							BusinessDetail: model.BusinessDetail{
								RegisteredEntityName: "Beam",
								BusinessName:         "Beam Data Company",
								BusinessIndustry:     "Payment Service",
								DomainName:           "https://www.beamcheckout.com/",
								IDNumber:             "0000000000000",
								PhoneNumber:          "000000000",
								Address: model.Address{
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
							BusinessContact: model.PointOfContact{
								FullName:    "test-organisation",
								Role:        "role",
								PhoneNumber: "000000000",
								Email:       "test@example.com",
							},
							TechnicalContact: model.PointOfContact{
								FullName:    "test-organisation",
								Role:        "role",
								PhoneNumber: "000000000",
								Email:       "test@example.com",
							},
							AccountingContact: model.PointOfContact{
								FullName:    "test-organisation",
								Role:        "role",
								PhoneNumber: "000000000",
								Email:       "test@example.com",
							},
							BankTransferDetail: model.BankTransferDetail{
								AccountName:   "accname",
								BankName:      "bankname",
								Branch:        "branch",
								AccountNumber: "0000000000000000",
								AccountType:   "type",
							},
							DocumentDownloadURL: "https://test.com",
							ApiKey: model.ApiKey{
								Required: true,
								Scope:    []string{},
							},
							Source:          "lighthouse",
							DatetimeCreated: time.Now(),
							Status:          "",
							Notes:           "",
						}, nil
					},
					updateKymStatus: func(kym *model.Kym, status string, notes string) error {
						return nil
					},
				},
				ms: MerchantServiceStub{
					addMerchant: func(nm *dto.NewMerchant, np *dto.KymField) error {
						return nil
					},
				},
				ks: KymServiceClientStub{
					RecipientClientStub: RecipientClientStub{
						createRecipient: func(crReq *thirdparty.CreateRecipientRequest) error {
							return nil
						},
						followEntity: func(source string, organisationId string) error {
							return nil
						},
					},
					OrganisationClientStub: OrganisationClientStub{
						createOrganisation: func(req *thirdparty.OrganisationRequest, userInfo string) (*thirdparty.OrganisationResponse, error) {
							return &thirdparty.OrganisationResponse{
								ID:     "beamdatcompany",
								ApiKey: "api-key-test",
							}, nil
						},
					},
				},
			},
			args: args{
				id: "example-id",
				kymStatusReq: &dto.UpdateKymStatusRequest{
					OrganisationID: "beam-organisation",
					Status:         "approved",
				},
				userInfo: "",
			},
			asserts: func(t *testing.T, gotErr error) {
				assert.NoError(t, gotErr)
			},
		},
		{
			name: "updateKymStatusFailedWithInvalidStatus",
			fields: fields{
				db: KymDBStub{
					getKym: func(id string) (*model.Kym, error) {
						return &model.Kym{
							ID:           "GS-16SZ-oOu5iFpn9SMlj",
							PartnerRefID: "beam-test-001",
							BusinessDetail: model.BusinessDetail{
								RegisteredEntityName: "Beam",
								BusinessName:         "Beam Data Company",
								BusinessIndustry:     "Payment Service",
								DomainName:           "https://www.beamcheckout.com/",
								IDNumber:             "0000000000000",
								PhoneNumber:          "000000000",
								Address: model.Address{
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
							BusinessContact: model.PointOfContact{
								FullName:    "test-organisation",
								Role:        "role",
								PhoneNumber: "000000000",
								Email:       "test@example.com",
							},
							TechnicalContact: model.PointOfContact{
								FullName:    "test-organisation",
								Role:        "role",
								PhoneNumber: "000000000",
								Email:       "test@example.com",
							},
							AccountingContact: model.PointOfContact{
								FullName:    "test-organisation",
								Role:        "role",
								PhoneNumber: "000000000",
								Email:       "test@example.com",
							},
							BankTransferDetail: model.BankTransferDetail{
								AccountName:   "accname",
								BankName:      "bankname",
								Branch:        "branch",
								AccountNumber: "0000000000000000",
								AccountType:   "type",
							},
							DocumentDownloadURL: "https://test.com",
							ApiKey: model.ApiKey{
								Required: true,
								Scope:    []string{},
							},
							Source:          "lighthouse",
							DatetimeCreated: time.Now(),
							Status:          "",
							Notes:           "",
						}, nil
					},
					updateKymStatus: func(kym *model.Kym, status string, notes string) error {
						return nil
					},
				},
				ms: MerchantServiceStub{
					addMerchant: func(nm *dto.NewMerchant, np *dto.KymField) error {
						return nil
					},
				},
				ks: KymServiceClientStub{
					RecipientClientStub: RecipientClientStub{
						createRecipient: func(crReq *thirdparty.CreateRecipientRequest) error {
							return nil
						},
						followEntity: func(source string, organisationId string) error {
							return nil
						},
					},
					OrganisationClientStub: OrganisationClientStub{
						createOrganisation: func(req *thirdparty.OrganisationRequest, userInfo string) (*thirdparty.OrganisationResponse, error) {
							return &thirdparty.OrganisationResponse{
								ID:     "beamdatcompany",
								ApiKey: "api-key-test",
							}, nil
						},
					},
				},
			},
			args: args{
				id: "example-id",
				kymStatusReq: &dto.UpdateKymStatusRequest{
					OrganisationID: "beam-organisation",
					Status:         "12345",
				},
				userInfo: "",
			},
			asserts: func(t *testing.T, gotErr error) {
				assert.EqualError(t, gotErr, "there are validation errors: kym status is not allowed with : 12345")
			},
		},
		{
			name: "updateKymStatusFailedWithAlreadyApproved",
			fields: fields{
				db: KymDBStub{
					getKym: func(id string) (*model.Kym, error) {
						return &model.Kym{
							ID:           "GS-16SZ-oOu5iFpn9SMlj",
							PartnerRefID: "beam-test-001",
							BusinessDetail: model.BusinessDetail{
								RegisteredEntityName: "Beam",
								BusinessName:         "Beam Data Company",
								BusinessIndustry:     "Payment Service",
								DomainName:           "https://www.beamcheckout.com/",
								IDNumber:             "0000000000000",
								PhoneNumber:          "000000000",
								Address: model.Address{
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
							BusinessContact: model.PointOfContact{
								FullName:    "test-organisation",
								Role:        "role",
								PhoneNumber: "000000000",
								Email:       "test@example.com",
							},
							TechnicalContact: model.PointOfContact{
								FullName:    "test-organisation",
								Role:        "role",
								PhoneNumber: "000000000",
								Email:       "test@example.com",
							},
							AccountingContact: model.PointOfContact{
								FullName:    "test-organisation",
								Role:        "role",
								PhoneNumber: "000000000",
								Email:       "test@example.com",
							},
							BankTransferDetail: model.BankTransferDetail{
								AccountName:   "accname",
								BankName:      "bankname",
								Branch:        "branch",
								AccountNumber: "0000000000000000",
								AccountType:   "type",
							},
							DocumentDownloadURL: "https://test.com",
							ApiKey: model.ApiKey{
								Required: true,
								Scope:    []string{},
							},
							Source:          "lighthouse",
							DatetimeCreated: time.Now(),
							Status:          "approved",
							Notes:           "",
						}, nil
					},
					updateKymStatus: func(kym *model.Kym, status string, notes string) error {
						return nil
					},
				},
				ms: MerchantServiceStub{
					addMerchant: func(nm *dto.NewMerchant, np *dto.KymField) error {
						return nil
					},
				},
				ks: KymServiceClientStub{
					RecipientClientStub: RecipientClientStub{
						createRecipient: func(crReq *thirdparty.CreateRecipientRequest) error {
							return nil
						},
						followEntity: func(source string, organisationId string) error {
							return nil
						},
					},
					OrganisationClientStub: OrganisationClientStub{
						createOrganisation: func(req *thirdparty.OrganisationRequest, userInfo string) (*thirdparty.OrganisationResponse, error) {
							return &thirdparty.OrganisationResponse{
								ID:     "beamdatcompany",
								ApiKey: "api-key-test",
							}, nil
						},
					},
				},
			},
			args: args{
				id: "example-id",
				kymStatusReq: &dto.UpdateKymStatusRequest{
					OrganisationID: "beam-organisation",
					Status:         "approved",
				},
				userInfo: "",
			},
			asserts: func(t *testing.T, gotErr error) {
				assert.EqualError(t, gotErr, "there are validation errors: status is already approved")
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &KymService{
				db: tt.fields.db,
				cs: tt.fields.cs,
				ks: tt.fields.ks,
				ms: tt.fields.ms,
			}
			err := s.UpdateKymStatus(tt.args.id, tt.args.kymStatusReq, tt.args.userInfo)
			tt.asserts(t, err)
		})
	}
}

// KymDBStub is a stub struct that proxies method calls to function fields.
type KymDBStub struct {
	getAllKym func(status string) ([]*model.Kym, error)

	getKym func(id string) (*model.Kym, error)

	addKym func(kym *model.Kym) error

	updateKymStatus func(kym *model.Kym, status string, notes string) error
}

func (stub KymDBStub) AddKym(kym *model.Kym) error {
	return stub.addKym(kym)
}

func (stub KymDBStub) GetAllKym(status string) ([]*model.Kym, error) {
	return stub.getAllKym(status)
}

func (stub KymDBStub) GetKym(id string) (*model.Kym, error) {
	return stub.getKym(id)
}

func (stub KymDBStub) UpdateKymStatus(kym *model.Kym, status string, notes string) error {
	return stub.updateKymStatus(kym, status, notes)
}

// OrganisationClientStub is a stub struct that proxies method calls to function fields.
type OrganisationClientStub struct {
	createOrganisation func(req *thirdparty.OrganisationRequest, userInfo string) (*thirdparty.OrganisationResponse, error)
}

func (orgStub OrganisationClientStub) CreateOrganisation(req *thirdparty.OrganisationRequest, userInfo string) (*thirdparty.OrganisationResponse, error) {
	return orgStub.createOrganisation(req, userInfo)
}

// KymServiceClientStub is a stub struct that proxies method calls to function fields.
type KymServiceClientStub struct {
	RecipientClientStub
	OrganisationClientStub
}

// RecipientClientStub is a stub struct that proxies method calls to function fields.
type RecipientClientStub struct {
	followEntity func(source string, organisationId string) error

	createRecipient func(crReq *thirdparty.CreateRecipientRequest) error
}

func (orgStub RecipientClientStub) FollowEntity(source string, organisationId string) error {
	return orgStub.followEntity(source, organisationId)
}

func (orgStub RecipientClientStub) CreateRecipient(crReq *thirdparty.CreateRecipientRequest) error {
	return orgStub.createRecipient(crReq)
}

// CloudStorageManagerStub is a stub struct that proxies method calls to function fields.
type CloudStorageManagerStub struct{}

func (c CloudStorageManagerStub) UploadFile(_ string, _ []byte) (string, error) {
	return "test-download-url", nil
}

// MerchantServiceStub is a stub struct that proxies method calls to function fields.
type MerchantServiceStub struct {
	addMerchant func(nm *dto.NewMerchant, np *dto.KymField) error

	getMerchant func(id string) (m *model.Merchant, err error)

	updateMerchant func(m *dto.Merchant) error

	upsertPayOutConfig func(merchantID string, poc *dto.PayOutConfig) error
}

func (msStub MerchantServiceStub) UpdateMerchant(m *dto.Merchant) error {
	return msStub.updateMerchant(m)
}

func (msStub MerchantServiceStub) UpsertPayOutConfig(merchantID string, poc *dto.PayOutConfig) error {
	return msStub.upsertPayOutConfig(merchantID, poc)
}

func (msStub MerchantServiceStub) AddMerchant(nm *dto.NewMerchant, np *dto.KymField) error {
	return msStub.addMerchant(nm, np)
}

func (msStub MerchantServiceStub) GetMerchant(id string) (m *model.Merchant, err error) {
	return msStub.getMerchant(id)
}
