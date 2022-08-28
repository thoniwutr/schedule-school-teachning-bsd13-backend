package service

import (
	"github.com/google/uuid"
	"github.com/thoniwutr/schedule-school-teachning-bsd13-backend/db"
	"github.com/thoniwutr/schedule-school-teachning-bsd13-backend/server/dto"
)

type MainSubjectService struct {
	db  db.MainSubjectDB
}

func NewMainSubjectService(db db.MainSubjectDB) *MainSubjectService {
	return &MainSubjectService{db: db}
}


func (m *MainSubjectService) AddMainSubject(mainSubjectRequest *dto.NewMainSubject) error {

	id := uuid.New()
	mainSubjectDetail := mainSubjectRequest.ToModel(id.String())

	return m.db.AddMainSubject(mainSubjectDetail)
}



func (m *MainSubjectService) GetAllMainSubject() ([]*dto.MainSubject, error) {

	mainSubjectList, err := m.db.GetAllMainSubject()
	if err != nil {
		return nil, err
	}

	mainSubjectResponse := dto.ToMainSubjectDTO(mainSubjectList)

	return mainSubjectResponse, nil
}

