package model

// PayInConfig defines model for PayInConfig.
type PayInConfig struct {

	// the channel through which the purchase is being made
	Channel string

	// the payment methods that the merchant is willing to accept through this channel
	PaymentMethods []string
}

// PayOutConfig defines model for PayOutConfig.
type PayOutConfig struct {
	BankAccount BankAccount

	// Three-letter ISO currency code in caps
	CurrencyCode string

	// how often the payout should occur
	Schedule PayOutConfigSchedule
}

// BankAccount defines model for BankAccount.
type BankAccount struct {

	// Name of the account
	AccountName string

	// Account number
	AccountNumber string

	// Name of the bank
	BankName string
}

// PayOutConfigSchedule how often the payout should occur
type PayOutConfigSchedule string

// Defines values for PayOutConfigSchedule.
const (
	PayOutConfigScheduleMonthly PayOutConfigSchedule = "monthly"

	PayOutConfigScheduleWeekly PayOutConfigSchedule = "weekly"
)
