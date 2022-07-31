package thirdparty

import (
	"github.com/go-playground/validator/v10"

	"github.com/thoniwutr/schedule-school-teachning-bsd13-backend/model"
)

// OrganisationRequest struct used for creating new organisation
type OrganisationRequest struct {
	// Full address of organisation
	Address string `json:"address" validate:"required" example:"012345678"`
	// Billing account id
	BillingAccountId string `json:"billingAccountId"`
	// Organisation's email address
	ContactEmail string `json:"contactEmail" validate:"required"`
	// Organisation description
	Description string `json:"description" validate:"required"`
	// Organisation name
	DisplayName string `json:"displayName" validate:"required"`
	// Organisation's id
	ID string `json:"id" validate:"required"`
	// Organisation's logo
	ImageURL string `json:"imageUrl"`
	// Organisation's package
	Package string `json:"package"`
	// Organisation phone number
	Phone string `json:"phone" validate:"required"`
	// Organisation's website
	Website string `json:"website"`
	// Partner refer from source attribute, It'll determine the value only when they're our partner
	Partner string `json:"partner" validate:"required"`
	// api key
	Apikey model.ApiKey `json:"apikey" validate:"required"`
}

func (o *OrganisationRequest) Validate() error {
	return validator.New().Struct(o)
}

// OrganisationResponse struct for received orgasition id and apikey
type OrganisationResponse struct {
	// organisation id
	ID string `json:"id"`
	// api key
	ApiKey string `json:"apiKey"`
}

// CreateRecipientRequest struct used for creating new organisation response
type CreateRecipientRequest struct {
	// Recipient id
	RecipientID string `json:"recipientId"`
	// Recipient type
	RecipientType string `json:"recipientType"`
	// The subscription source
	Subscriptions []string `json:"subscriptions"`
}
