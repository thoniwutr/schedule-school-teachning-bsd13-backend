package model

import (
	"time"
)

// Teacher defines model for Teacher.
type Merchant struct {

	// MerchantID is the OrganisationID
	MerchantID string

	Address Address

	// main contact number for the merchant
	ContactNumber string

	// main contact email for the merchant
	Email string

	// full name for the organisation/merchant
	FullName string

	// the original organisation ID
	OrganisationID string

	// deprecated - kept for datastore compatability
	ShortName string

	LogoURL string

	// the payment methods that the merchant is willing to accept
	PaymentMethods []string

	CurrencyCode string

	// timestamp the merchant entity was created
	Created time.Time

	// timestamp the merchant entity was last updated
	Updated time.Time

	PayOutConfig *PayOutConfig
}

// Address defines model for Address.
type Address struct {
	HouseNumber string
	Street      string
	District    string
	Subdistrict string
	City        string
	Province    string
	Zipcode     string
	Country     string
}

// NewMerchant is a constructor for Merchant which populates the MerchantID and the timestamps
func NewMerchant(address Address, contactNumber, email, fullName, orgID, logoURL, currencyCode string, paymentMethods []string) *Merchant {
	now := time.Now()
	return &Merchant{
		Address:        address,
		ContactNumber:  contactNumber,
		Email:          email,
		FullName:       fullName,
		OrganisationID: orgID,
		MerchantID:     orgID,
		LogoURL:        logoURL,
		CurrencyCode:   currencyCode,
		PaymentMethods: paymentMethods,
		Created:        now,
		Updated:        now,
	}
}
