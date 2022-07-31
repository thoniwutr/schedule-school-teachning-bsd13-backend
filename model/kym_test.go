package model

import (
	"testing"
	"time"
)

func TestKym_GetFullAddress(t *testing.T) {
	type fields struct {
		ID                  string
		OrganisationID      string
		ImageURL            string
		PartnerRefID        string
		BusinessDetail      BusinessDetail
		BusinessContact     PointOfContact
		TechnicalContact    PointOfContact
		AccountingContact   PointOfContact
		BankTransferDetail  BankTransferDetail
		DocumentDownloadURL string
		ApiKey              ApiKey
		Source              string
		DatetimeCreated     time.Time
		Status              string
		Notes               string
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name: "GetFullAddressSuccess",
			fields: fields{
				ID:             "id-test",
				OrganisationID: "org-test",
				BusinessDetail: BusinessDetail{
					Address: Address{
						City:        "city",
						Country:     "Thailand",
						District:    "district",
						HouseNumber: "9",
						Province:    "",
						Street:      "Phetkasem",
						Subdistrict: "Bang Khae",
						Zipcode:     "zipcode",
					},
				},
			},
			want: "9 Phetkasem district Bang Khae city zipcode Thailand",
		},
		{
			name: "GetFullAddressWithEmptyString",
			fields: fields{
				ID:             "id-test",
				OrganisationID: "org-test",
				BusinessDetail: BusinessDetail{
					Address: Address{
						City:        "",
						Country:     "",
						District:    "",
						HouseNumber: "",
						Province:    "",
						Street:      "",
						Subdistrict: "",
						Zipcode:     "",
					},
				},
			},
			want: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			k := &Kym{
				ID:                  tt.fields.ID,
				OrganisationID:      tt.fields.OrganisationID,
				ImageURL:            tt.fields.ImageURL,
				PartnerRefID:        tt.fields.PartnerRefID,
				BusinessDetail:      tt.fields.BusinessDetail,
				BusinessContact:     tt.fields.BusinessContact,
				TechnicalContact:    tt.fields.TechnicalContact,
				AccountingContact:   tt.fields.AccountingContact,
				BankTransferDetail:  tt.fields.BankTransferDetail,
				DocumentDownloadURL: tt.fields.DocumentDownloadURL,
				ApiKey:              tt.fields.ApiKey,
				Source:              tt.fields.Source,
				DatetimeCreated:     tt.fields.DatetimeCreated,
				Status:              tt.fields.Status,
				Notes:               tt.fields.Notes,
			}
			if got := k.GetFullAddress(); got != tt.want {
				t.Errorf("GetFullAddress() = %v, want %v", got, tt.want)
			}
		})
	}
}
