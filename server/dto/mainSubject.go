package dto

import (
	"github.com/go-playground/validator/v10"
	"github.com/thoniwutr/schedule-school-teachning-bsd13-backend/model"
)

// NewMainSubject struct used for creating new teacher requests
type NewMainSubject struct {
	// Only THB is supported currently
	MainSubjectName string `json:"mainSubjectName" validate:"required"`
}

// MainSubject struct used for updating Teacher
type MainSubject struct {
	NewMainSubject
	ID string `json:"id" validate:"required"`
}

func ToMainSubjectDTO(mainSubjectList []*model.MainSubject) []*MainSubject {

	var mainSubjectRes []*MainSubject

	for _, item := range mainSubjectList {
		var mainSubjectResItem = &MainSubject{
			ID:                  item.ID,
			NewMainSubject: NewMainSubject{
				MainSubjectName: item.MainSubjectName,
			},
		}
		mainSubjectRes = append(mainSubjectRes, mainSubjectResItem)

	}

	return mainSubjectRes
}


// Validate does some simple validation on the NewMerchant object per annotations
func (nt *NewMainSubject) Validate() error {
	return validator.New().Struct(nt)
}


// ToModel converts dto.NewMerchant to model.Merchant
func (nt *NewMainSubject) ToModel(id string) *model.MainSubject {
	return model.NewMainSubject(
		id,
		nt.MainSubjectName,
	)
}

