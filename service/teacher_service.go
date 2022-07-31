package service

import (
	"github.com/Beam-Data-Company/merchant-config-svc/db"
	"github.com/Beam-Data-Company/merchant-config-svc/server/dto"
	"github.com/google/uuid"
)

type TeacherService struct {
	db  db.TeacherDB
}

func NewTeacherService(db db.TeacherDB) *TeacherService {
	return &TeacherService{db: db}
}



func (s *TeacherService) AddTeacher(teacherRequest *dto.NewTeacher) error {

	id := uuid.New()
	kymDetail := teacherRequest.ToModel(id.String())

	return s.db.AddTeacher(kymDetail)
}



func (s *TeacherService) GetAllTeacher() ([]*dto.Teacher, error) {

	teacherList, err := s.db.GetAllTeacher()
	if err != nil {
		return nil, err
	}

	klmResponse := dto.ToTeacherDTO(teacherList)

	return klmResponse, nil
}

