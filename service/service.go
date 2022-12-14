package service

import (
	"github.com/thoniwutr/schedule-school-teachning-bsd13-backend/model"
	"github.com/thoniwutr/schedule-school-teachning-bsd13-backend/server/dto"
)

// KymServiceInterface defines business logic of kym api
type KymServiceInterface interface {
	AddKym(kymRequest *dto.NewKym) error

	GetKym(id string) (*dto.KymFullDetailResponse, error)

	GetAllKym(status string) ([]*dto.KymResponse, error)

	UpdateKymStatus(id string, kymStatusReq *dto.UpdateKymStatusRequest, userInfo string) error
}

// MerchantServiceInterface defines business logic of merchant api
type MerchantServiceInterface interface {
	AddMerchant(nm *dto.NewMerchant, kf *dto.KymField) error

	GetMerchant(id string) (m *model.Merchant, err error)

	UpdateMerchant(m *dto.Merchant) error

	UpsertPayOutConfig(merchantID string, poc *dto.PayOutConfig) error
}

// TeacherServiceInterface defines business logic of teacher api
type TeacherServiceInterface interface {
	AddTeacher(nt *dto.NewTeacher) error

	GetAllTeacher() ([]*dto.Teacher, error)
}

type MainSubjectServiceInterface interface {
	AddMainSubject(nt *dto.NewMainSubject) error

	GetAllMainSubject() ([]*dto.MainSubject, error)
}

type SubjectServiceInterface interface {
	AddSubject(nt *dto.NewSubject) error

	GetAllSubject() ([]*dto.Subject, error)
}

type ConfirmationServiceInterface interface {
	AddConfirmation(request *dto.NewConfirmation) error
	GetAllConfirmation() ([]*dto.Confirmation, error)
	AddConfirmationDetail(request *dto.NewConfirmationDetail) error
	GetAllConfirmationDetail(id string) ([]*dto.ConfirmationDetail, error)
}
