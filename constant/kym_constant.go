package constant

// internal app specific errors to be used
const (
	OrganisationBeamDataCompany  = "@beamdatacompany"
	OrganisationZort             = "@zort"
	SourceLighthouse             = "@lighthouse"
	SourceZort                   = "@zort"
	KymStatusApproved            = "approved"
	KymStatusPending             = "pending"
	KymStatusRejected            = "rejected"
	PaymentMethodCreditCard      = "creditCard"
	PaymentMethodEWallet         = "eWallet"
	PaymentMethodInternetBanking = "internetBanking"
	RecipientTypeMerchant        = "MERCHANT"
	CurrencyCodeTHB              = "THB"
)

func IsValidKymStatus(status string) bool {
	arr := []string{KymStatusPending, KymStatusApproved, KymStatusRejected}
	var result = false
	for _, x := range arr {
		if x == status {
			result = true
			break
		}
	}
	return result
}
