package model

// Teacher defines model for Teacher.
type Teacher struct {
	ID string `json:"Id"`
	// The merchant's company name in full
	FirstName string `json:"firstName"`
	// The merchant's company name in full
	NickName string `json:"nickName"`
	// The merchant's company name in full
	LastName string `json:"lastName"`
	// The merchant's main contact number
	ContactNumber string `json:"contactNumber"`
	// Only THB is supported currently
	Capacity string `json:"capacity"`
	// Only THB is supported currently
	MainSubjectID string `json:"mainSubjectId"`
}

// NewTeacher is a constructor for Merchant which populates the MerchantID and the timestamps
func NewTeacher(id string, firstName string, nickName string, lastName string, contactNumber string, capacity string, mainSubjectID string) *Teacher {
	return &Teacher{
		ID:            id,
		NickName:      nickName,
		FirstName:     firstName,
		LastName:      lastName,
		ContactNumber: contactNumber,
		Capacity:      capacity,
		MainSubjectID: mainSubjectID,
	}
}
