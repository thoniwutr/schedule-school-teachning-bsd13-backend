package model

type Confirmation struct {
	ID string
	ConfirmationName string
	CreateDate string
}

func NewConfirmation(id string,	confirmationName string, createDate string) *Confirmation {
	return &Confirmation{
		ID:   id,
		ConfirmationName: confirmationName,
		CreateDate: createDate,
	}
}


type ConfirmationDetail struct {
	ID string
	ConfirmationID string
	SubjectDetailID []string
	StudentName string
	Level string
	Period string
	Day string
}

func NewConfirmationDetail(id string,confirmationID string, subjectDetailID []string, studentName string, level string,period string,day string) *ConfirmationDetail {
	return &ConfirmationDetail{
		ID: id,
		ConfirmationID:   confirmationID,
		SubjectDetailID: subjectDetailID,
		StudentName: studentName,
		Level: level,
		Period: period,
		Day: day,
	}
}
