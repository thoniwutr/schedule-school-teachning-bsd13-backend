package model

import (
	"fmt"
	"time"

	c "github.com/Beam-Data-Company/merchant-config-svc/constant"
)

// Kym defines model for Kym.
type Kym struct {
	// ID defines partner ref ID
	ID string

	// OrganisationID defines id of organisation
	OrganisationID string

	// ImageURL defines company logo
	ImageURL string

	// PartnerRefID defines id of Partner Reference
	PartnerRefID string

	// BusinessDetail defines primary of information of company
	BusinessDetail BusinessDetail

	// BusinessContact defines contact detail of business department
	BusinessContact PointOfContact

	// TechnicalContact defines contact detail of technical department
	TechnicalContact PointOfContact

	// AccountingContact defines contact detail of accounting department
	AccountingContact PointOfContact

	// BankTransferDetail defines bank transfer detail
	BankTransferDetail BankTransferDetail

	// DocumentDownloadURL defines url for download the documentations
	DocumentDownloadURL string

	// ApiKey defines api key defines model for api key.
	ApiKey ApiKey

	// Source indicate where kym is submitted from
	Source string

	// DatetimeCreated defines timestamp for kym was created
	DatetimeCreated time.Time

	// Status defines the status of KYM document verification
	Status string

	// Notes defines the reason for KYM status update
	Notes string
}

// ApiKey defines model for api key.
type ApiKey struct {
	Required bool
	Scope    []string `datastore:",noindex"`
}

// BusinessDetail defines model for overview of business.
type BusinessDetail struct {
	RegisteredEntityName string
	BusinessName         string
	BusinessIndustry     string
	DomainName           string
	IDNumber             string
	PhoneNumber          string
	Address              Address
}

// PointOfContact defines model for represent point of contact.
type PointOfContact struct {
	FullName    string
	Role        string
	PhoneNumber string
	Email       string
}

// BankTransferDetail defines model for bank transfer detail.
type BankTransferDetail struct {
	AccountName   string
	BankName      string
	Branch        string
	AccountNumber string
	AccountType   string
}

// NewKym Kym is a constructor for Kym which populates the MerchantID and the timestamps
func NewKym(
	id string,
	organisationID string,
	partnerRefID string,
	businessDetail BusinessDetail,
	businessContact PointOfContact,
	technicalContact PointOfContact,
	accountingContact PointOfContact,
	bankTransfersDetail BankTransferDetail,
	imageUrl string,
	apiKey ApiKey,
	source string,
	documentDownloadURL string) *Kym {
	return &Kym{
		ID:                  id,
		OrganisationID:      organisationID,
		ImageURL:            imageUrl,
		PartnerRefID:        partnerRefID,
		BusinessDetail:      businessDetail,
		BusinessContact:     businessContact,
		TechnicalContact:    technicalContact,
		AccountingContact:   accountingContact,
		BankTransferDetail:  bankTransfersDetail,
		ApiKey:              apiKey,
		Source:              source,
		DocumentDownloadURL: documentDownloadURL,
		DatetimeCreated:     time.Now(),
		Status:              c.KymStatusPending,
	}
}

func (k *Kym) GetFullAddress() string {
	var fullAddress string
	fullAddress = appendAddress(fullAddress, k.BusinessDetail.Address.HouseNumber)
	fullAddress = appendAddress(fullAddress, k.BusinessDetail.Address.Street)
	fullAddress = appendAddress(fullAddress, k.BusinessDetail.Address.District)
	fullAddress = appendAddress(fullAddress, k.BusinessDetail.Address.Subdistrict)
	fullAddress = appendAddress(fullAddress, k.BusinessDetail.Address.City)
	fullAddress = appendAddress(fullAddress, k.BusinessDetail.Address.Province)
	fullAddress = appendAddress(fullAddress, k.BusinessDetail.Address.Zipcode)
	fullAddress = appendAddress(fullAddress, k.BusinessDetail.Address.Country)
	return fullAddress
}

func appendAddress(fullAddress, address string) string {
	switch {
	case len(fullAddress) > 0 && len(address) > 0:
		return fmt.Sprintf("%v %v", fullAddress, address)
	case len(fullAddress) > 0:
		return fullAddress
	case len(address) > 0:
		return address
	default:
		return ""
	}
}
