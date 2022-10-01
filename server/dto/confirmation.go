package dto

import (
	"github.com/go-playground/validator/v10"
	"github.com/thoniwutr/schedule-school-teachning-bsd13-backend/model"
)

type NewConfirmationDetail struct {
	ConfirmationID string `json:"confirmationId" validate:"required" example:"merchant company"`
	SubjectDetailID []string `json:"subjectDetailId" validate:"required" example:"merchant company"`
	StudentName string `json:"studentName" validate:"required" example:"merchant company"`
	Level string `json:"level" validate:"required" example:"merchant company"`
	Period string `json:"period" validate:"required" example:"merchant company"`
	Day string `json:"day" example:"merchant company"`
}


type ConfirmationDetail struct {
	NewConfirmationDetail
	ID string `json:"id" validate:"required"`
}

func ToConfirmationDetailDTO(confirmationList []*model.ConfirmationDetail) []*ConfirmationDetail {

	var res []*ConfirmationDetail

	for _, item := range confirmationList {
		var resItem = &ConfirmationDetail{
			ID:                  item.ID,
			NewConfirmationDetail: NewConfirmationDetail{
				ConfirmationID: item.ConfirmationID,
				SubjectDetailID: item.SubjectDetailID,
				StudentName: item.StudentName,
				Level: item.Level,
				Period: item.Period,
				Day: item.Day,
			},
		}
		res = append(res, resItem)
	}
	return res
}



func (cd *NewConfirmationDetail) Validate() error {
	return validator.New().Struct(cd)
}


func (cd *NewConfirmationDetail) ToModel(id string) *model.ConfirmationDetail {
	return model.NewConfirmationDetail(
		id,
		cd.ConfirmationID,
		cd.SubjectDetailID,
		cd.StudentName,
		cd.Level,
		cd.Period,
		cd.Day,
	)
}



// NewConfirmation struct used for creating new teacher requests
type NewConfirmation struct {
	// Full URL to the merchant's logo to display
	ConfirmationName string `json:"confirmationName" validate:"required" example:"merchant company"`
	// The merchant's company name in full
	CreateDate string `json:"createDate" validate:"required" example:"merchant company"`
}

// Confirmation struct used for updating Teacher
type Confirmation struct {
	NewConfirmation
	ID string `json:"id" validate:"required"`
}

func ToConfirmationDTO(confirmationList []*model.Confirmation) []*Confirmation {

	var res []*Confirmation

	for _, item := range confirmationList {
		var resItem = &Confirmation{
			ID:                  item.ID,
			NewConfirmation: NewConfirmation{
				ConfirmationName: item.ConfirmationName,
				CreateDate: item.CreateDate,
			},
		}
		res = append(res, resItem)
	}
	return res
}


// Validate does some simple validation on the NewMerchant object per annotations
func (nt *NewConfirmation) Validate() error {
	return validator.New().Struct(nt)
}


// ToModel converts dto.NewMerchant to model.Merchant
func (nt *NewConfirmation) ToModel(id string) *model.Confirmation {
	return model.NewConfirmation(
		id,
		nt.ConfirmationName,
		nt.CreateDate,
	)
}

