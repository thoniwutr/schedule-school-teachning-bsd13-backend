package service

import (
	"fmt"

	"github.com/Beam-Data-Company/merchant-config-svc/db"
	"github.com/Beam-Data-Company/merchant-config-svc/messaging"
	"github.com/Beam-Data-Company/merchant-config-svc/model"
	"github.com/Beam-Data-Company/merchant-config-svc/server/dto"
)

type MerchantService struct {
	db  db.MerchantDB
	pub messaging.Publisher
}

func NewMerchantService(db db.MerchantDB, pub messaging.Publisher) *MerchantService {
	return &MerchantService{db: db, pub: pub}
}

func (s *MerchantService) AddMerchant(nm *dto.NewMerchant, kf *dto.KymField) error {
	merchant := nm.ToModel()
	if err := s.db.AddMerchant(merchant); err != nil {
		return err
	}

	// add merchant successful, publish merchant to topic
	if err := s.pub.PublishMerchant(dto.ToMerchantPublish(merchant, kf)); err != nil {
		return err
	}
	return nil
}

func (s *MerchantService) GetMerchant(id string) (m *model.Merchant, err error) {
	m, err = s.db.GetMerchant(id)
	if err != nil {
		return nil, fmt.Errorf("unable to find merchant with id %v: %w", id, err)
	}
	return m, nil
}

func (s *MerchantService) UpdateMerchant(m *dto.Merchant) error {
	return s.db.UpdateMerchant(m.ToModel())
}

func (s *MerchantService) UpsertPayOutConfig(merchantID string, poc *dto.PayOutConfig) error {
	return s.db.UpsertPayOutConfig(merchantID, poc.ToModel())
}
