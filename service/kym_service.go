package service

import (
	"encoding/base64"
	"errors"
	"fmt"

	c "github.com/thoniwutr/-schedule-school-teachning-bsd13-backend/constant"
	"github.com/thoniwutr/-schedule-school-teachning-bsd13-backend/db"
	"github.com/thoniwutr/-schedule-school-teachning-bsd13-backend/server/dto"
	"github.com/thoniwutr/-schedule-school-teachning-bsd13-backend/thirdparty"
	"github.com/thoniwutr/-schedule-school-teachning-bsd13-backend/util"
)

type KymService struct {
	db db.KymDB
	cs util.CloudStorageManager
	ks thirdparty.KymServiceClient
	ms MerchantServiceInterface
}

func NewKymService(db db.KymDB, cs util.CloudStorageManager, ks thirdparty.KymServiceClient, ms MerchantServiceInterface) *KymService {
	return &KymService{db: db, cs: cs, ks: ks, ms: ms}
}

func (s *KymService) GetAllKym(status string) ([]*dto.KymResponse, error) {

	kymList, err := s.db.GetAllKym(status)
	if err != nil {
		return nil, err
	}

	klmResponse := dto.ToKymDTO(kymList)

	return klmResponse, nil
}

func (s *KymService) GetKym(id string) (*dto.KymFullDetailResponse, error) {

	kym, err := s.db.GetKym(id)
	if err != nil {
		return nil, err
	}

	klmFullDetailResponse := dto.ToKymFullDetailDTO(kym)

	return klmFullDetailResponse, nil
}

func (s *KymService) AddKym(kymRequest *dto.NewKym) error {

	documentData, err := base64.StdEncoding.DecodeString(kymRequest.DocumentContent)
	if err != nil {
		return fmt.Errorf("failed to decode base64 : %w", err)
	}

	id, err := kymRequest.GenerateID()
	if err != nil {
		return err
	}

	path, err := s.cs.UploadFile(id, documentData)
	if err != nil {
		return fmt.Errorf("failed to upload file to cloud storage : %w", err)
	}

	kymDetail := kymRequest.ToModel(id, path)

	return s.db.AddKym(kymDetail)

}

func (s *KymService) UpdateKymStatus(id string, kymStatusReq *dto.UpdateKymStatusRequest, userInfo string) error {

	kym, err := s.db.GetKym(id)
	if err != nil {
		return err
	}

	// check current status in database
	if kym.Status == c.KymStatusApproved {
		return &c.ErrValidation{Violations: errors.New("status is already approved")}
	}

	// check request status
	if !c.IsValidKymStatus(kymStatusReq.Status) {
		return &c.ErrValidation{Violations: errors.New(fmt.Sprintf("kym status is not allowed with : %v", kymStatusReq.Status))}
	}

	err = s.db.UpdateKymStatus(kym, kymStatusReq.Status, kymStatusReq.Notes)
	if err != nil {
		return err
	}

	nm := &dto.NewMerchant{
		Address: dto.MerchantAddress{
			City:        kym.BusinessDetail.Address.City,
			Country:     kym.BusinessDetail.Address.Country,
			District:    kym.BusinessDetail.Address.District,
			HouseNumber: kym.BusinessDetail.Address.HouseNumber,
			Province:    kym.BusinessDetail.Address.Province,
			Street:      kym.BusinessDetail.Address.Street,
			Subdistrict: kym.BusinessDetail.Address.Subdistrict,
			Zipcode:     kym.BusinessDetail.Address.Zipcode,
		},
		FullName:                kym.BusinessDetail.BusinessName,
		Email:                   kym.BusinessContact.Email,
		ContactNumber:           kym.BusinessDetail.PhoneNumber,
		AvailablePaymentMethods: []string{c.PaymentMethodCreditCard, c.PaymentMethodEWallet, c.PaymentMethodInternetBanking},
		CurrencyCode:            c.CurrencyCodeTHB,
	}

	kymField := &dto.KymField{
		PartnerRefId: kym.PartnerRefID,
		Source:       kym.Source,
	}

	if kym.Source != c.SourceLighthouse {

		if len(kymStatusReq.OrganisationID) == 0 {
			return &c.ErrValidation{Violations: errors.New("organisation is empty")}
		}

		// step a : create organisation
		coReq := &thirdparty.OrganisationRequest{
			Address:      kym.GetFullAddress(),
			ContactEmail: kym.BusinessContact.Email,
			Description:  kym.BusinessDetail.BusinessIndustry,
			DisplayName:  kym.BusinessDetail.BusinessName,
			ID:           kymStatusReq.OrganisationID,
			Phone:        kym.BusinessDetail.PhoneNumber,
			Website:      kym.BusinessDetail.DomainName,
			Partner:      kym.Source,
			Apikey:       kym.ApiKey,
		}
		organisationRes, err := s.ks.CreateOrganisation(coReq, userInfo)
		if err != nil {
			return err
		}

		// step b : Create Recipient
		crReq := &thirdparty.CreateRecipientRequest{
			RecipientID:   organisationRes.ID,
			RecipientType: c.RecipientTypeMerchant,
			Subscriptions: []string{},
		}
		if err := s.ks.CreateRecipient(crReq); err != nil {
			return err
		}

		// step c : Follow Entity
		if err := s.ks.FollowEntity(kym.Source, organisationRes.ID); err != nil {
			return err
		}

		nm.OrganisationID = organisationRes.ID
		kymField.ApiKey = organisationRes.ApiKey

	} else {
		nm.OrganisationID = kym.OrganisationID
	}

	// d + f : Create New Merchant + Kym Field + Add Merchant to DB
	if err = s.ms.AddMerchant(nm, kymField); err != nil {
		return err
	}

	return nil
}
