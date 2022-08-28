package model


type Subject struct {
	ID string
	SubjectName string
	MainSubjectId string
	MinOfStudent int
}

func NewSubject(id string, subjectName string, mainSubjectId string, minOfStudent int) *Subject {
	return &Subject{
		ID:            id,
		SubjectName: subjectName,
		MainSubjectId: mainSubjectId,
		MinOfStudent: minOfStudent,
	}
}
