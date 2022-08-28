package dto

import (
	"github.com/go-playground/validator/v10"
	"github.com/thoniwutr/schedule-school-teachning-bsd13-backend/model"
)


type NewSubject struct {
	SubjectName string `json:"subjectName" validate:"required"`
	MainSubjectId string `json:"mainSubjectId" validate:"required"`
	MinOfStudent int `json:"minOfStudent" validate:"required"`
}


type Subject struct {
	NewSubject
	ID string `json:"id" validate:"required"`
}

func ToSubjectDTO(subjectList []*model.Subject) []*Subject {

	var subjectRes []*Subject

	for _, item := range subjectList {
		var mainSubjectResItem = &Subject{
			ID:                  item.ID,
			NewSubject: NewSubject{
				SubjectName: item.SubjectName,
				MainSubjectId: item.MainSubjectId,
				MinOfStudent: item.MinOfStudent,
			},
		}
		subjectRes = append(subjectRes, mainSubjectResItem)

	}

	return subjectRes
}


// Validate does some simple validation on the NewMerchant object per annotations
func (ns *NewSubject) Validate() error {
	return validator.New().Struct(ns)
}


// ToModel converts dto.NewMerchant to model.Merchant
func (ns *NewSubject) ToModel(id string) *model.Subject {
	return model.NewSubject(
		id,
		ns.SubjectName,
		ns.MainSubjectId,
		ns.MinOfStudent,
	)
}

