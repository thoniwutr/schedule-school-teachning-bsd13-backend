package service

import (
	"github.com/google/uuid"
	"github.com/thoniwutr/schedule-school-teachning-bsd13-backend/db"
	"github.com/thoniwutr/schedule-school-teachning-bsd13-backend/server/dto"
)

type ConfirmationService struct {
	db  db.ConfirmationDB
}

func NewConfirmationService(db db.ConfirmationDB) *ConfirmationService {
	return &ConfirmationService{db: db}
}


func (m *ConfirmationService) AddConfirmation(confirmRequest *dto.NewConfirmation) error {

	id := uuid.New()
	confirmation := confirmRequest.ToModel(id.String())

	return m.db.AddConfirmation(confirmation)
}



func (m *ConfirmationService) GetAllConfirmation() ([]*dto.Confirmation, error) {

	list, err := m.db.GetAllConfirmation()
	if err != nil {
		return nil, err
	}

	response := dto.ToConfirmationDTO(list)

	return response, nil
}

func (m *ConfirmationService) AddConfirmationDetail(confirmRequest *dto.NewConfirmationDetail) error {

	id := uuid.New()
	confirmationDetail :=  confirmRequest.ToModel(id.String())

	return m.db.AddConfirmationDetail(confirmationDetail)
}



func (m *ConfirmationService) GetAllConfirmationDetail(id string) ([]*dto.ConfirmationDetail, error) {

	list, err := m.db.GetAllConfirmationDetail(id)
	if err != nil {
		return nil, err
	}

	response := dto.ToConfirmationDetailDTO(list)

	return response, nil
}

