package thirdparty

// OrganisationClient defines internal api for organisation
type OrganisationClient interface {
	CreateOrganisation(req *OrganisationRequest, userInfo string) (*OrganisationResponse, error)
}

// RecipientClient defines internal api for recipient
type RecipientClient interface {
	CreateRecipient(crReq *CreateRecipientRequest) error
	FollowEntity(source string, organisationId string) error
}

// KymServiceClient defines only client type use for kym service
type KymServiceClient interface {
	OrganisationClient
	RecipientClient
}
