package dto

import (
	"errors"
	"fmt"
	"time"

	"github.com/go-playground/validator/v10"

	c "github.com/thoniwutr/-schedule-school-teachning-bsd13-backend/constant"
	"github.com/thoniwutr/-schedule-school-teachning-bsd13-backend/model"
)

// NewKym struct used for creating new KYM requests
type NewKym struct {
	// OrganisationID represent id of organisation
	OrganisationID string `json:"organisationId"`
	// PartnerRefID defines partner reference id
	PartnerRefID string `json:"partnerRefId"`
	// ImageURL represent company logo
	ImageURL string `json:"imageUrl"`
	// Source indicate where kym is submitted from
	Source string `json:"source" validate:"required" example:"@lighthouse"`
	// DocumentContent represent documentation uploaded convert to base64
	DocumentContent string `json:"documentContent" validate:"required"`
	// BusinessDetail represent primary information of company
	BusinessDetail BusinessDetail `json:"businessDetail" validate:"required"`
	// BusinessContact represent business of contact
	BusinessContact PointOfContact `json:"businessContact"`
	// TechnicalContact represent technical of contact
	TechnicalContact PointOfContact `json:"technicalContact"`
	// AccountingContact represent accounting of contact
	AccountingContact PointOfContact `json:"accountingContact"`
	// BankTransferDetail represent detail of bank detail
	BankTransferDetail BankTransferDetail `json:"bankTransferDetail" validate:"required"`
	// ApiKey defines model for api key.
	ApiKey ApiKey `json:"apiKey"`
}

type KymResponse struct {
	ID                  string    `json:"id" validate:"required"`
	BusinessName        string    `json:"businessName" validate:"required"`
	ContactPerson       string    `json:"contactPerson" validate:"required"`
	PhoneNumber         string    `json:"phoneNumber" validate:"required"`
	DatetimeCreated     time.Time `json:"datetimeCreated"`
	DocumentDownloadURL string    `json:"documentDownloadUrl"`
	Status              string    `json:"status"`
}

type KymFullDetailResponse struct {
	ID                  string             `json:"id" validate:"required"`
	OrganisationID      string             `json:"organisationId" validate:"required"`
	PartnerRefID        string             `json:"partnerRefId" validate:"required"`
	BusinessDetail      BusinessDetail     `json:"businessDetail"`
	BusinessContact     PointOfContact     `json:"businessContact"`
	TechnicalContact    PointOfContact     `json:"technicalContact"`
	AccountingContact   PointOfContact     `json:"accountingContact"`
	BankTransferDetail  BankTransferDetail `json:"bankTransferDetail"`
	ApiKey              ApiKey             `json:"apiKey"`
	DocumentDownloadURL string             `json:"documentDownloadUrl"`
	Source              string             `json:"source"`
	ImageURL            string             `json:"imageUrl"`
	DatetimeCreated     time.Time          `json:"datetimeCreated"`
	Status              string             `json:"status"`
	Notes               string             `json:"notes"`
}

// ApiKey defines the model for API key request for the business
type ApiKey struct {
	Required bool     `json:"required"`
	Scope    []string `json:"scope,omitempty"`
}

// BusinessDetail defines model for overview of business.
type BusinessDetail struct {
	RegisteredEntityName string          `json:"registeredEntityName" validate:"required"`
	BusinessName         string          `json:"businessName" validate:"required"`
	BusinessIndustry     string          `json:"businessIndustry,omitempty"`
	DomainName           string          `json:"domainName,omitempty"`
	IDNumber             string          `json:"idNumber" validate:"required"`
	PhoneNumber          string          `json:"phoneNumber" validate:"required"`
	Address              MerchantAddress `json:"address" validate:"required"`
}

// PointOfContact defines model for represent point of contact.
type PointOfContact struct {
	FullName    string `json:"fullName" validate:"required"`
	Role        string `json:"role,omitempty"`
	PhoneNumber string `json:"phoneNumber" validate:"required"`
	Email       string `json:"email,omitempty"`
}

// BankTransferDetail defines model for bank transfer detail.
type BankTransferDetail struct {
	AccountName   string `json:"accountName" validate:"required"`
	BankName      string `json:"bankName" validate:"required"`
	Branch        string `json:"branch,omitempty"`
	AccountNumber string `json:"accountNumber" validate:"required"`
	AccountType   string `json:"accountType,omitempty"`
}

// Validate does some simple validation on the Kym object per annotations
func (nKym *NewKym) Validate() error {
	return validator.New().Struct(nKym)
}

// GenerateID create unique id for new KYM object
func (nKym *NewKym) GenerateID() (string, error) {
	var suffix string
	if nKym.Source == c.SourceLighthouse {
		if len(nKym.OrganisationID) == 0 {
			return "", &c.ErrValidation{Violations: errors.New("failed to generate id with empty organisation id")}
		} else {
			suffix = nKym.OrganisationID
		}
	} else {
		if len(nKym.PartnerRefID) == 0 {
			return "", &c.ErrValidation{Violations: errors.New("failed to generate id with empty partner ref id")}
		} else {
			suffix = nKym.PartnerRefID
		}
	}
	return fmt.Sprintf("%v-%v", nKym.Source, suffix), nil
}

// ToModel converts kym request to kym model for save to datastore
func (nKym *NewKym) ToModel(id string, downloadURL string) *model.Kym {
	return model.NewKym(
		id,
		nKym.OrganisationID,
		nKym.PartnerRefID,
		model.BusinessDetail{
			RegisteredEntityName: nKym.BusinessDetail.RegisteredEntityName,
			BusinessName:         nKym.BusinessDetail.BusinessName,
			DomainName:           nKym.BusinessDetail.DomainName,
			BusinessIndustry:     nKym.BusinessDetail.BusinessIndustry,
			IDNumber:             nKym.BusinessDetail.IDNumber,
			PhoneNumber:          nKym.BusinessDetail.PhoneNumber,
			Address: model.Address{
				City:        nKym.BusinessDetail.Address.City,
				Country:     nKym.BusinessDetail.Address.Country,
				District:    nKym.BusinessDetail.Address.District,
				HouseNumber: nKym.BusinessDetail.Address.HouseNumber,
				Province:    nKym.BusinessDetail.Address.Province,
				Street:      nKym.BusinessDetail.Address.Street,
				Subdistrict: nKym.BusinessDetail.Address.Subdistrict,
				Zipcode:     nKym.BusinessDetail.Address.Zipcode,
			},
		},
		model.PointOfContact{
			FullName:    nKym.BusinessContact.FullName,
			Role:        nKym.BusinessContact.Role,
			PhoneNumber: nKym.BusinessContact.PhoneNumber,
			Email:       nKym.BusinessContact.Email,
		},
		model.PointOfContact{
			FullName:    nKym.TechnicalContact.FullName,
			Role:        nKym.TechnicalContact.Role,
			PhoneNumber: nKym.TechnicalContact.PhoneNumber,
			Email:       nKym.TechnicalContact.Email,
		},
		model.PointOfContact{
			FullName:    nKym.AccountingContact.FullName,
			Role:        nKym.AccountingContact.Role,
			PhoneNumber: nKym.AccountingContact.PhoneNumber,
			Email:       nKym.AccountingContact.Email,
		},
		model.BankTransferDetail{
			AccountName:   nKym.BankTransferDetail.AccountName,
			BankName:      nKym.BankTransferDetail.BankName,
			Branch:        nKym.BankTransferDetail.Branch,
			AccountNumber: nKym.BankTransferDetail.AccountNumber,
			AccountType:   nKym.BankTransferDetail.AccountType,
		},
		nKym.ImageURL,
		model.ApiKey(nKym.ApiKey),
		nKym.Source,
		downloadURL,
	)
}

// UpdateKymStatusRequest defines model for PayOutConfig.
type UpdateKymStatusRequest struct {
	OrganisationID string `json:"organisationId"`
	Status         string `json:"status" validate:"required"`
	Notes          string `json:"notes"`
}

// Validate does some simple validation on the UpdateKymStatusRequest object per annotations
func (kyms *UpdateKymStatusRequest) Validate() error {
	return validator.New().Struct(kyms)
}

func ToKymDTO(kymList []*model.Kym) []*KymResponse {

	var kymRes []*KymResponse

	for _, item := range kymList {
		var kymResItem = &KymResponse{
			ID:                  item.ID,
			BusinessName:        item.BusinessDetail.BusinessName,
			ContactPerson:       item.BusinessDetail.RegisteredEntityName,
			PhoneNumber:         item.BusinessDetail.PhoneNumber,
			DatetimeCreated:     item.DatetimeCreated,
			DocumentDownloadURL: item.DocumentDownloadURL,
			Status:              item.Status,
		}
		kymRes = append(kymRes, kymResItem)

	}

	return kymRes
}

func ToKymFullDetailDTO(kym *model.Kym) *KymFullDetailResponse {
	return &KymFullDetailResponse{
		ID:             kym.ID,
		OrganisationID: kym.OrganisationID,
		PartnerRefID:   kym.PartnerRefID,
		BusinessDetail: BusinessDetail{
			kym.BusinessDetail.RegisteredEntityName,
			kym.BusinessDetail.BusinessName,
			kym.BusinessDetail.BusinessIndustry,
			kym.BusinessDetail.DomainName,
			kym.BusinessDetail.IDNumber,
			kym.BusinessDetail.PhoneNumber,
			MerchantAddress{
				kym.BusinessDetail.Address.City,
				kym.BusinessDetail.Address.Country,
				kym.BusinessDetail.Address.District,
				kym.BusinessDetail.Address.HouseNumber,
				kym.BusinessDetail.Address.Province,
				kym.BusinessDetail.Address.City,
				kym.BusinessDetail.Address.City,
				kym.BusinessDetail.Address.City,
			},
		},
		BusinessContact: PointOfContact{
			kym.BusinessContact.FullName,
			kym.BusinessContact.Role,
			kym.BusinessContact.PhoneNumber,
			kym.BusinessContact.Email,
		},
		TechnicalContact: PointOfContact{
			kym.TechnicalContact.FullName,
			kym.TechnicalContact.Role,
			kym.TechnicalContact.PhoneNumber,
			kym.TechnicalContact.Email,
		},
		AccountingContact: PointOfContact{
			kym.AccountingContact.FullName,
			kym.AccountingContact.Role,
			kym.AccountingContact.PhoneNumber,
			kym.AccountingContact.Email,
		},
		BankTransferDetail: BankTransferDetail{
			kym.BankTransferDetail.AccountName,
			kym.BankTransferDetail.BankName,
			kym.BankTransferDetail.Branch,
			kym.BankTransferDetail.AccountNumber,
			kym.BankTransferDetail.AccountType,
		},
		DocumentDownloadURL: kym.DocumentDownloadURL,
		ImageURL:            kym.ImageURL,
		ApiKey: ApiKey{
			kym.ApiKey.Required,
			kym.ApiKey.Scope,
		},
		Source:          kym.Source,
		DatetimeCreated: kym.DatetimeCreated,
		Status:          kym.Status,
		Notes:           kym.Notes,
	}

}
