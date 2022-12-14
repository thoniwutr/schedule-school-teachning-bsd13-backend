package dto

import (
	"github.com/go-playground/validator/v10"
	"github.com/thoniwutr/schedule-school-teachning-bsd13-backend/model"
)

// NewTeacher struct used for creating new teacher requests
type NewTeacher struct {
	// Full URL to the merchant's logo to display
	NickName string `json:"nickName" validate:"required" example:"merchant company"`
	// The merchant's company name in full
	FirstName string `json:"firstName" validate:"required" example:"merchant company"`
	// The merchant's company name in full
	LastName string `json:"lastName" validate:"required" example:"merchant company"`
	// The merchant's main contact number
	ContactNumber string `json:"contactNumber" validate:"required" example:"012345678"`
	// Only THB is supported currently
	Capacity string `json:"capacity" validate:"required"`
	// Only THB is supported currently
	MainSubjectID string `json:"mainSubjectID" validate:"required"`
}

// Teacher struct used for updating Teacher
type Teacher struct {
	NewTeacher
	ID string `json:"id" validate:"required"`
}

func ToTeacherDTO(teacherList []*model.Teacher) []*Teacher {

	var teacherRes []*Teacher

	for _, item := range teacherList {
		var teacherResItem = &Teacher{
			ID:                  item.ID,
			NewTeacher: NewTeacher{
				NickName: item.NickName,
				FirstName: item.FirstName,
				LastName: item.ContactNumber,
				Capacity: item.Capacity,
				MainSubjectID: item.MainSubjectID,
			},
		}
		teacherRes = append(teacherRes, teacherResItem)

	}

	return teacherRes
}


// Validate does some simple validation on the NewMerchant object per annotations
func (nt *NewTeacher) Validate() error {
	return validator.New().Struct(nt)
}


// ToModel converts dto.NewMerchant to model.Merchant
func (nt *NewTeacher) ToModel(id string) *model.Teacher {
	return model.NewTeacher(
		id,
		nt.NickName,
		nt.FirstName,
		nt.LastName,
		nt.ContactNumber,
		nt.Capacity,
		nt.MainSubjectID,
	)
}

