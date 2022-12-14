package db

import (
	"github.com/thoniwutr/schedule-school-teachning-bsd13-backend/model"
)

type AppDB interface {
	MerchantDB

	RoleDB
}

// MerchantDB defines an interface for our Application's data access methods
type MerchantDB interface {
	// GetMerchant gets the Merchant from the given id
	GetMerchant(merchantID string) (*model.Merchant, error)

	// AddMerchant creates the Merchant to the db
	AddMerchant(m *model.Merchant) error

	// UpdateMerchant updates an existing Merchant model
	UpdateMerchant(m *model.Merchant) error

	// UpsertPayOutConfig Merchant's PayOutConfig to be updated or inserted
	UpsertPayOutConfig(merchantID string, poc *model.PayOutConfig) error
}

// RoleDB defines an interface for our Application's data access methods
type RoleDB interface {
	// GetRole retrieves an organisation role by ID
	GetRole(organisationID, userID string) (*model.Role, error)
}

// KymDB defines an interface for our Application's data access methods
type KymDB interface {
	// GetAllKym gets all kym detail from db
	GetAllKym(status string) ([]*model.Kym, error)

	// GetKym gets all kym detail from db
	GetKym(id string) (*model.Kym, error)

	// AddKym creates the klm detail to db
	AddKym(kym *model.Kym) error

	// UpdateKymStatus update kym status
	UpdateKymStatus(kym *model.Kym, status string, notes string) error
}



// TeacherDB defines an interface for our Application's data access methods
type TeacherDB interface {
	// GetAllTeacher gets all kym detail from db
	GetAllTeacher() ([]*model.Teacher, error)


	// AddTeacher creates the klm detail to db
	AddTeacher(kym *model.Teacher) error

}

// MainSubjectDB defines an interface for our Application's data access methods
type MainSubjectDB interface {
	// GetAllMainSubject gets all kym detail from db
	GetAllMainSubject() ([]*model.MainSubject, error)


	// AddMainSubject creates the klm detail to db
	AddMainSubject(kym *model.MainSubject) error

}

type SubjectDB interface {
	GetAllSubject() ([]*model.Subject, error)
	AddSubject(kym *model.Subject) error

}


type ConfirmationDB interface {
	AddConfirmation(m *model.Confirmation) error
	AddConfirmationDetail(m *model.ConfirmationDetail) error
	GetAllConfirmation() ([]*model.Confirmation, error)
	GetAllConfirmationDetail(confirmationId string) ([]*model.ConfirmationDetail, error)
}
