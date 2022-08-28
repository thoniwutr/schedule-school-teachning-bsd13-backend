package model


type TeacherResponsibility struct {
	ID string
	TeacherId string
	SubjectId string
}

func NewTeacherResponsibility(id string, teacherId string, subjectId string) *TeacherResponsibility {
	return &TeacherResponsibility{
		ID:            id,
		TeacherId: teacherId,
		SubjectId: subjectId,
	}
}
