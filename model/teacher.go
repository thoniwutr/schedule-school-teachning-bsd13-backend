package model

// Teacher defines model for Teacher.
type Teacher struct {
	ID string
	// Full URL to the merchant's logo to display
	ImgURL string
	// The merchant's company name in full
	FistName string
	// The merchant's company name in full
	LastName string
	// The merchant's main contact number
	ContactNumber string
	// Only THB is supported currently
	Capacity string
	// Only THB is supported currently
	MainSubjectID string
}

// NewTeacher is a constructor for Merchant which populates the MerchantID and the timestamps
func NewTeacher(id string, imgURL string, fistName string, lastName string, contactNumber string, capacity string, mainSubjectID string) *Teacher {
	return &Teacher{
		ID:            id,
		ImgURL:        imgURL,
		FistName:      fistName,
		LastName:      lastName,
		ContactNumber: contactNumber,
		Capacity:      capacity,
		MainSubjectID: mainSubjectID,
	}
}
