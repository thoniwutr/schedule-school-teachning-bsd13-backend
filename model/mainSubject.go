package model

// MainSubject defines model for MainSubject.
type MainSubject struct {
	ID string
	// Full URL to the merchant's logo to display
	MainSubjectName string
	// The merchant's company name in full
}

// NewMainSubject is a constructor for MainSubject which populates the MainSubject and the MainSubject
func NewMainSubject(id string, mainSubjectName string) *MainSubject {
	return &MainSubject{
		ID:            id,
		MainSubjectName: mainSubjectName,
	}
}
