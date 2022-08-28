package service

import (
	"github.com/google/uuid"
	"github.com/thoniwutr/schedule-school-teachning-bsd13-backend/db"
	"github.com/thoniwutr/schedule-school-teachning-bsd13-backend/server/dto"
)

type SubjectService struct {
	db  db.SubjectDB
}

func NewSubjectService(db db.SubjectDB) *SubjectService {
	return &SubjectService{db: db}
}


func (m *SubjectService) AddSubject(subjectRequest *dto.NewSubject) error {

	id := uuid.New()
	mainSubjectDetail := subjectRequest.ToModel(id.String())

	return m.db.AddSubject(mainSubjectDetail)
}



func (m *SubjectService) GetAllSubject() ([]*dto.Subject, error) {

	subjectList, err := m.db.GetAllSubject()
	if err != nil {
		return nil, err
	}

	subjectResponse := dto.ToSubjectDTO(subjectList)

	return subjectResponse, nil
}

