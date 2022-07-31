package model_test

import (
	"testing"

	"github.com/thoniwutr/schedule-school-teachning-bsd13-backend/model"
)

func TestCreateMerchant(t *testing.T) {
	m := model.NewMerchant(
		model.Address{
			City:        "Bangkok",
			Country:     "Thailand",
			District:    "district",
			HouseNumber: "123",
			Province:    "province",
			Street:      "street",
			Subdistrict: "subdistrict",
			Zipcode:     "12345",
		},
		"123456",
		"test@user.com",
		"company",
		"org123",
		"org",
		"logo",
		[]string{},
	)

	// some tests
	if m.ContactNumber != "123456" {
		t.Errorf("contactNumber incorrect expected %s got %s", "123456", m.ContactNumber)
	}

	if m.MerchantID != "org123" {
		t.Errorf("incorrect merchantID expected %s got %s", "org123", m.MerchantID)
	}

	if m.Created.IsZero() {
		t.Error("property Created has not been set")
	}

	if m.Updated.IsZero() {
		t.Error("property Updated has not been set")
	}
}
