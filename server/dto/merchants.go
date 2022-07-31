package dto

import (
	"time"

	"github.com/go-playground/validator/v10"

	"github.com/Beam-Data-Company/merchant-config-svc/model"
)

// NewMerchant struct used for creating new merchant requests
type NewMerchant struct {
	Address MerchantAddress `json:"address"`
	// The merchant's main contact number
	ContactNumber string `json:"contactNumber" validate:"required" example:"012345678"`
	// The merchant's email address
	Email string `json:"email" validate:"required" example:"email@merchant.com"`
	// The merchant's company name in full
	FullName string `json:"fullName" validate:"required" example:"merchant company"`
	// Beam's organisation ID of the merchant
	OrganisationID string `json:"organisationId" validate:"required" example:"merchant123"`
	// Full URL to the merchant's logo to display
	LogoURL string `json:"logoUrl,omitempty" example:"https://url.to.logo/image.png"`
	// acceptable payment methods
	AvailablePaymentMethods []string `json:"availablePaymentMethods" validate:"required,unique" example:"creditCard,internetBanking"`
	// Only THB is supported currently
	CurrencyCode string `json:"currencyCode" validate:"required,eq=THB"`
}

// Merchant struct used for updating Merchant
type Merchant struct {
	NewMerchant
	MerchantID string `json:"merchantId" validate:"required"`
}

// MerchantGetResponse represents the GET /merchants response
type MerchantGetResponse struct {
	Merchant

	// timestamp the merchant entity was created
	Created time.Time `json:"created"`

	// timestamp the merchant entity was last updated
	Updated time.Time `json:"updated"`

	PayOutConfig *PayOutConfig `json:"payOutConfig"`
}

// MerchantPublish represents the Merchant struct to be published downstream
type MerchantPublish struct {
	Merchant

	Source       string `json:"source"`
	ApiKey       string `json:"apiKey"`
	PartnerRefID string `json:"partnerRefID"`

	// timestamp the merchant entity was created
	Created time.Time `json:"created"`

	// timestamp the merchant entity was last updated
	Updated time.Time `json:"updated"`

	PayOutConfig *PayOutConfig `json:"payOutConfig"`
}

// MerchantAddress defines model for MerchantAddress.
type MerchantAddress struct {
	City        string `json:"city" validate:"required"`
	Country     string `json:"country" validate:"required"`
	District    string `json:"district,omitempty"`
	HouseNumber string `json:"houseNumber,omitempty"`
	Province    string `json:"province,omitempty"`
	Street      string `json:"street,omitempty"`
	Subdistrict string `json:"subdistrict,omitempty"`
	Zipcode     string `json:"zipcode" validate:"required"`
}

// BankAccount defines model for BankAccount.
type BankAccount struct {
	AccountName   string `json:"accountName" validate:"required"`
	AccountNumber string `json:"accountNumber" validate:"required"`
	BankName      string `json:"bankName" validate:"required"`
}

// PayOutConfig defines model for PayOutConfig.
type PayOutConfig struct {
	BankAccount  BankAccount `json:"bankAccount" validate:"required"`
	CurrencyCode string      `json:"currencyCode" validate:"required" example:"THB"`
	// how often the payout should occur
	Schedule PayOutConfigSchedule `json:"schedule" validate:"required" example:"weekly,monthly"`
}

// PayOutConfigSchedule defines how often the payout should occur
type PayOutConfigSchedule string

const (
	// PayOutConfigScheduleMonthly monthly
	PayOutConfigScheduleMonthly PayOutConfigSchedule = "monthly"

	// PayOutConfigScheduleWeekly weekly
	PayOutConfigScheduleWeekly PayOutConfigSchedule = "weekly"
)

// Validate does some simple validation on the Merchant object per annotations
func (nm *NewMerchant) Validate() error {
	return validator.New().Struct(nm)
}

// ToModel converts dto.NewMerchant to model.Merchant
func (nm *NewMerchant) ToModel() *model.Merchant {
	return model.NewMerchant(
		model.Address{
			City:        nm.Address.City,
			Country:     nm.Address.Country,
			District:    nm.Address.District,
			HouseNumber: nm.Address.HouseNumber,
			Province:    nm.Address.Province,
			Street:      nm.Address.Street,
			Subdistrict: nm.Address.Subdistrict,
			Zipcode:     nm.Address.Zipcode,
		},
		nm.ContactNumber,
		nm.Email,
		nm.FullName,
		nm.OrganisationID,
		nm.LogoURL,
		nm.CurrencyCode,
		nm.AvailablePaymentMethods,
	)
}

// Validate does some simple validation on the NewMerchant object per annotations
func (m *Merchant) Validate() error {
	return validator.New().Struct(m)
}

// ToModel converts dto.Merchant to model.Merchant
func (m *Merchant) ToModel() *model.Merchant {
	return &model.Merchant{
		Address: model.Address{
			City:        m.Address.City,
			Country:     m.Address.Country,
			District:    m.Address.District,
			HouseNumber: m.Address.HouseNumber,
			Province:    m.Address.Province,
			Street:      m.Address.Street,
			Subdistrict: m.Address.Subdistrict,
			Zipcode:     m.Address.Zipcode,
		},
		CurrencyCode:   m.CurrencyCode,
		ContactNumber:  m.ContactNumber,
		Email:          m.Email,
		FullName:       m.FullName,
		OrganisationID: m.OrganisationID,
		MerchantID:     m.MerchantID,
		LogoURL:        m.LogoURL,
		PaymentMethods: m.AvailablePaymentMethods,
	}
}

// Validate does some simple validation on the Merchant object per annotations
func (poc *PayOutConfig) Validate() error {
	return validator.New().Struct(poc)
}

// ToModel converts dto.PayOutConfig to model.PayOutConfig
func (poc *PayOutConfig) ToModel() *model.PayOutConfig {
	return &model.PayOutConfig{
		BankAccount: model.BankAccount{
			AccountName:   poc.BankAccount.AccountName,
			AccountNumber: poc.BankAccount.AccountNumber,
			BankName:      poc.BankAccount.BankName,
		},
		CurrencyCode: poc.CurrencyCode,
		Schedule:     model.PayOutConfigSchedule(poc.Schedule),
	}
}

// ToMerchantDTO converts model.Merchant to MerchantGetResponse DTO
func ToMerchantDTO(m *model.Merchant) *MerchantGetResponse {
	resp := &MerchantGetResponse{
		Merchant: Merchant{
			NewMerchant: NewMerchant{
				Address: MerchantAddress{
					City:        m.Address.City,
					Country:     m.Address.Country,
					District:    m.Address.District,
					HouseNumber: m.Address.HouseNumber,
					Province:    m.Address.Province,
					Street:      m.Address.Street,
					Subdistrict: m.Address.Subdistrict,
					Zipcode:     m.Address.Zipcode,
				},
				ContactNumber:           m.ContactNumber,
				Email:                   m.Email,
				FullName:                m.FullName,
				OrganisationID:          m.OrganisationID,
				LogoURL:                 m.LogoURL,
				CurrencyCode:            m.CurrencyCode,
				AvailablePaymentMethods: m.PaymentMethods,
			},
			MerchantID: m.MerchantID,
		},
		Created: m.Created,
		Updated: m.Updated,
	}

	if m.PayOutConfig != nil {
		resp.PayOutConfig = &PayOutConfig{
			BankAccount: BankAccount{
				AccountName:   m.PayOutConfig.BankAccount.AccountName,
				AccountNumber: m.PayOutConfig.BankAccount.AccountNumber,
				BankName:      m.PayOutConfig.BankAccount.BankName,
			},
			CurrencyCode: m.PayOutConfig.CurrencyCode,
			Schedule:     PayOutConfigSchedule(m.PayOutConfig.Schedule),
		}
	}

	return resp
}

// KymField converts model.kym to MerchantPublish
type KymField struct {
	ApiKey       string `json:"apiKey"`
	PartnerRefId string `json:"partnerRefId"`
	Source       string `json:"source"`
}

// ToMerchantPublish converts model.Merchant to MerchantPublish
func ToMerchantPublish(m *model.Merchant, np *KymField) *MerchantPublish {
	resp := &MerchantPublish{
		Merchant: Merchant{
			NewMerchant: NewMerchant{
				Address: MerchantAddress{
					City:        m.Address.City,
					Country:     m.Address.Country,
					District:    m.Address.District,
					HouseNumber: m.Address.HouseNumber,
					Province:    m.Address.Province,
					Street:      m.Address.Street,
					Subdistrict: m.Address.Subdistrict,
					Zipcode:     m.Address.Zipcode,
				},
				ContactNumber:           m.ContactNumber,
				Email:                   m.Email,
				FullName:                m.FullName,
				OrganisationID:          m.OrganisationID,
				LogoURL:                 m.LogoURL,
				CurrencyCode:            m.CurrencyCode,
				AvailablePaymentMethods: m.PaymentMethods,
			},
			MerchantID: m.MerchantID,
		},
		Created: m.Created,
		Updated: m.Updated,
	}

	if np != nil {
		resp.ApiKey = np.ApiKey
		resp.Source = np.Source
		resp.PartnerRefID = np.PartnerRefId
	}

	if m.PayOutConfig != nil {
		resp.PayOutConfig = &PayOutConfig{
			BankAccount: BankAccount{
				AccountName:   m.PayOutConfig.BankAccount.AccountName,
				AccountNumber: m.PayOutConfig.BankAccount.AccountNumber,
				BankName:      m.PayOutConfig.BankAccount.BankName,
			},
			CurrencyCode: m.PayOutConfig.CurrencyCode,
			Schedule:     PayOutConfigSchedule(m.PayOutConfig.Schedule),
		}
	}

	return resp
}
