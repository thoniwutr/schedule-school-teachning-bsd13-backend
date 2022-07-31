package db

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"reflect"
	"testing"
	"time"

	"cloud.google.com/go/datastore"
	"github.com/google/go-cmp/cmp"
	"google.golang.org/api/iterator"

	c "github.com/thoniwutr/-schedule-school-teachning-bsd13-backend/constant"
	"github.com/thoniwutr/-schedule-school-teachning-bsd13-backend/model"
	"github.com/thoniwutr/-schedule-school-teachning-bsd13-backend/server/dto"
)

// merchantDB is the struct under test
var merchantDB *testDBEnv

// TestMain starts datastore emulator and waits for it to initialize then run unit tests

func TestMain(m *testing.M) {

	dse, ok := os.LookupEnv("DATASTORE_EMULATOR_HOST")
	if !ok {
		log.Fatal("this test depends on DATASTORE_EMULATOR_HOST being set")
	}

	if success := checkEmulatorHealth(dse); !success {
		log.Fatal("emulator health endpoint unhealthy")
	}

	merchantds, err := NewAppDatastore("12345")
	if err != nil {
		log.Fatalf("failed to init db connection: %v", err)
	}

	// initialize the struct under test here
	merchantDB = &testDBEnv{merchantds}

	// ensure clean db
	merchantDB.tearDown()

	// executes the unit tests within the file
	exitcode := m.Run()

	// teardown database
	merchantDB.tearDown()

	os.Exit(exitcode)
}

// TestAddMerchant tests adding a NewMerchant and then fetching it to ensure the data is as
// expected
func TestAddMerchant(t *testing.T) {

	defer merchantDB.tearDown()

	mdto := &dto.NewMerchant{
		Address: dto.MerchantAddress{
			City:        "bkk",
			Country:     "thailand",
			District:    "district",
			HouseNumber: "123",
			Province:    "province",
			Street:      "street",
			Subdistrict: "subdistrict",
			Zipcode:     "12345",
		},
		ContactNumber:           "0123456789",
		Email:                   "test@example.com",
		FullName:                "test-organisation",
		OrganisationID:          "org123",
		LogoURL:                 "logo",
		AvailablePaymentMethods: []string{"creditCard"},
	}

	added, err := merchantDB.setUpMerchant(mdto)
	if err != nil {
		t.Errorf("failed to add default merchant")
	}

	m, err := merchantDB.GetMerchant("org123")
	if err != nil {
		t.Errorf("merchantID org123 does not exist: %v", err)
	}

	// compare the two
	if diff := cmp.Diff(added, m, cmp.Comparer(merchantComparer)); diff != "" {
		t.Errorf("merchants unequal: %s", diff)
	}
}

func TestAddDuplicateMerchant(t *testing.T) {

	defer merchantDB.tearDown()

	mdto := &dto.NewMerchant{
		Address: dto.MerchantAddress{
			City:        "bkk",
			Country:     "thailand",
			District:    "district",
			HouseNumber: "123",
			Province:    "province",
			Street:      "street",
			Subdistrict: "subdistrict",
			Zipcode:     "12345",
		},
		ContactNumber:  "0123456789",
		Email:          "test@example.com",
		FullName:       "test-organisation",
		OrganisationID: "org123",
	}

	added, err := merchantDB.setUpMerchant(mdto)
	if err != nil {
		t.Errorf("failed to add default merchant")
	}

	m, err := merchantDB.GetMerchant("org123")
	if err != nil {
		t.Errorf("merchantID org123 does not exist: %v", err)
	}

	// compare the two
	if diff := cmp.Diff(added, m, cmp.Comparer(merchantComparer)); diff != "" {
		t.Errorf("merchants unequal: %s", diff)
	}

	// readd the Merchant
	if err := merchantDB.AddMerchant(added); err != c.ErrDBEntityAlreadyExists {
		t.Error("expecting add duplicate merchants to fail with c.ErrDBEntityAlreadyExists")
	}
}

func TestGetUnexistingMerchant(t *testing.T) {

	defer merchantDB.tearDown()

	_, err := merchantDB.GetMerchant("this_does_not_exist")
	if err == nil {
		t.Errorf("expecting merchantID 'this_does_not_exist' to throw an error")
	}
	if err != c.ErrDBNoSuchEntity {
		t.Errorf("expecting merchantID 'this_does_not_exist' to throw ErrDBNoSuchEntity error")
	}
}

// TestUpdateMerchant
func TestUpdateMerchant(t *testing.T) {

	defer merchantDB.tearDown()

	mdto := &dto.NewMerchant{
		Address: dto.MerchantAddress{
			City:        "bkk",
			Country:     "thailand",
			District:    "district",
			HouseNumber: "123",
			Province:    "province",
			Street:      "street",
			Subdistrict: "subdistrict",
			Zipcode:     "12345",
		},
		ContactNumber:  "0123456789",
		Email:          "test@example.com",
		FullName:       "test-organisation",
		OrganisationID: "org123",
	}

	added, err := merchantDB.setUpMerchant(mdto)
	if err != nil {
		t.Errorf("failed to add default merchant")
	}

	m, err := merchantDB.GetMerchant("org123")
	if err != nil {
		t.Errorf("merchantID org123 does not exist")
	}

	// compare the two
	if diff := cmp.Diff(added, m, cmp.Comparer(merchantComparer)); diff != "" {
		t.Errorf("merchants unequal: %s", diff)
	}

	added.Email = "changedEmail@test.com"
	added.ContactNumber = "987654321"

	if err := merchantDB.UpdateMerchant(added); err != nil {
		t.Errorf("failed to update merchant; %v", err)
	}

	// fetch again and expect no diff
	m, err = merchantDB.GetMerchant("org123")
	if err != nil {
		t.Errorf("merchantID org123 does not exist: %v", err)
	}

	if diff := cmp.Diff(added, m, cmp.Comparer(merchantComparer)); diff != "" {
		t.Errorf("merchants unequal: %s", diff)
	}
}

func TestUpdateMerchantFailedValidation(t *testing.T) {

	defer merchantDB.tearDown()

	mdto := &dto.NewMerchant{
		Address: dto.MerchantAddress{
			City:        "bkk",
			Country:     "thailand",
			District:    "district",
			HouseNumber: "123",
			Province:    "province",
			Street:      "street",
			Subdistrict: "subdistrict",
			Zipcode:     "12345",
		},
		ContactNumber:  "0123456789",
		Email:          "test@example.com",
		FullName:       "test-organisation",
		OrganisationID: "org123",
	}

	added, err := merchantDB.setUpMerchant(mdto)
	if err != nil {
		t.Errorf("failed to add default merchant")
	}

	m, err := merchantDB.GetMerchant("org123")
	if err != nil {
		t.Errorf("merchantID org123 does not exist")
	}

	// compare the two
	if diff := cmp.Diff(added, m, cmp.Comparer(merchantComparer)); diff != "" {
		t.Errorf("merchants unequal: %s", diff)
	}

	added.FullName = "thiscannotbechanged"
	added.OrganisationID = "thisalsocannotchange"

	err = merchantDB.UpdateMerchant(added)

	var validationErr *c.ErrValidation
	if !errors.As(err, &validationErr) {
		t.Error("expect merchant update to fail on validations")
	}

	// t.Logf("got expected error from UpdateMerchant (not supposed to fail): %v", err)
}

// TestUpsertPayOutConfig
func TestUpsertPayOutConfig(t *testing.T) {

	defer merchantDB.tearDown()

	mdto := &dto.NewMerchant{
		Address: dto.MerchantAddress{
			City:        "bkk",
			Country:     "thailand",
			District:    "district",
			HouseNumber: "123",
			Province:    "province",
			Street:      "street",
			Subdistrict: "subdistrict",
			Zipcode:     "12345",
		},
		ContactNumber:  "0123456789",
		Email:          "test@example.com",
		FullName:       "test-organisation",
		OrganisationID: "org123",
	}

	added, err := merchantDB.setUpMerchant(mdto)
	if err != nil {
		t.Errorf("failed to add default merchant")
	}

	m, err := merchantDB.GetMerchant("org123")
	if err != nil {
		t.Errorf("merchantID org123 does not exist: %v", err)
	}

	// compare the two
	if diff := cmp.Diff(added, m, cmp.Comparer(merchantComparer)); diff != "" {
		t.Errorf("merchants unequal: %s", diff)
	}

	poc := &model.PayOutConfig{
		BankAccount: model.BankAccount{
			AccountName:   "testaccount",
			AccountNumber: "12345",
			BankName:      "greenbank",
		},
		CurrencyCode: "THB",
		Schedule:     model.PayOutConfigScheduleWeekly,
	}

	// upsert PayOutConfig
	if err := merchantDB.UpsertPayOutConfig("org123", poc); err != nil {
		t.Errorf("unexpected error from pay out config: %v", err)
	}

	m, err = merchantDB.GetMerchant("org123")
	if err != nil {
		t.Errorf("merchantID org123 does not exist: %v", err)
	}

	if diff := cmp.Diff(m.PayOutConfig, poc); diff != "" {
		t.Errorf("found diffs between PayOutConfig %v", diff)
	}
}

// ----------- Helper functions ----------------

// merchantComparer leaves out comparing the Created and Updated fields
func merchantComparer(x, y model.Merchant) bool {
	return x.MerchantID == y.MerchantID &&
		x.Address == y.Address &&
		x.ContactNumber == y.ContactNumber &&
		x.Email == y.Email &&
		x.FullName == y.FullName &&
		x.OrganisationID == y.OrganisationID &&
		x.ShortName == y.ShortName &&
		x.LogoURL == y.LogoURL &&
		reflect.DeepEqual(x.PaymentMethods, y.PaymentMethods) &&
		reflect.DeepEqual(x.PayOutConfig, y.PayOutConfig)
}

// ------------ Merchant DB Test ----------------

// testDBEnv is a wrapper for db.MerchantDB but it also has a direct access to *datastore.Client.
// This allows us to run extra queries outside the scope of the app for testing purposes
type testDBEnv struct {
	*AppDatastore
}

func (tdb *testDBEnv) setUpMerchant(mdto *dto.NewMerchant) (*model.Merchant, error) {
	m := mdto.ToModel()
	err := tdb.AddMerchant(m)
	return m, err
}

// tearDown deletes all Merchant entity from the datastore emulator
func (tdb *testDBEnv) tearDown() error {

	// query all
	q := datastore.NewQuery(tdb.KindMerchant)

	keys := make([]*datastore.Key, 0)
	t := tdb.client.Run(context.Background(), q)
	for {
		var m model.Merchant
		key, err := t.Next(&m)
		if err == iterator.Done {
			break
		}
		if err != nil {
			return err
		}
		keys = append(keys, key)
	}

	// then we delete all
	return tdb.client.DeleteMulti(context.Background(), keys)
}

// ------------ Datastore Emulator ----------------

// checkEmulatorHealth checks health of the emulator running needed for this test
func checkEmulatorHealth(dshost string) bool {
	for retries := 5; retries > 0; retries-- {
		_, err := http.Get(fmt.Sprintf("http://%s/", dshost))
		if err != nil {
			log.Print("retrying connection...")
			time.Sleep(time.Second * 1)
		} else {
			log.Print("emulator is up")
			return true
		}
	}
	return false
}
