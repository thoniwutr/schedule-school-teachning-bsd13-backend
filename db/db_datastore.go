package db

import (
	"context"
	"errors"
	"fmt"
	"time"

	"cloud.google.com/go/datastore"

	c "github.com/thoniwutr/schedule-school-teachning-bsd13-backend/constant"
	"github.com/thoniwutr/schedule-school-teachning-bsd13-backend/model"
)

// AppDatastore is an internal type to be accessed through the database interface.
// Provide an interface to interact with Google Cloud Datastore
type AppDatastore struct {
	client       *datastore.Client
	KindMerchant string
	KindRole     string
	KindKym      string
	KindTeacher string
	KindMainSubject string
	KindSubject string
	KindConfirmation string
	KindConfirmationDetail string``
}

// NewAppDatastore create a datastore client to persist application data on Google Cloud Datastore
func NewAppDatastore(projectID string) (*AppDatastore, error) {

	ctx := context.Background()
	client, err := datastore.NewClient(ctx, projectID)
	if err != nil {
		return nil, err
	}
	// Verify that we can communicate and authenticate with the datastore security.
	// context with connection timeout
	ctx, cancel := context.WithTimeout(ctx, time.Second*3)
	defer cancel()
	t, err := client.NewTransaction(ctx)
	if err != nil {
		return nil, fmt.Errorf("unable to communicate to datastore: %w", err)
	}
	if err := t.Rollback(); err != nil {
		return nil, fmt.Errorf("unable to communicate to datastore: %w", err)
	}

	return &AppDatastore{client, "Merchant", "Role", "Kym","Teacher", "MainSubject", "Subject","Confirmation","ConfirmationDetail"}, nil
}

// GetMerchant returns the Merchant given the ID
func (db *AppDatastore) GetMerchant(merchantID string) (*model.Merchant, error) {
	key := db.merchantKey(merchantID)
	m := &model.Merchant{}
	err := db.client.Get(context.Background(), key, m)
	switch err {
	case nil:
		return m, nil
	case datastore.ErrNoSuchEntity:
		return nil, c.ErrDBNoSuchEntity
	}
	return nil, err
}

// AddMerchant attempts to add a NewMerchant to the datastore.
// Since merchantID == organisationID, we do not allow duplicate organisationIDs
func (db *AppDatastore) AddMerchant(m *model.Merchant) error {
	_, err := db.client.RunInTransaction(context.Background(), func(tx *datastore.Transaction) error {
		key := db.merchantKey(m.MerchantID)
		err := tx.Get(key, &model.Merchant{})

		switch err {
		case nil:
			return c.ErrDBEntityAlreadyExists
		case datastore.ErrNoSuchEntity:
			// no existing entity id, proceed with write
			_, err = tx.Put(key, m)
			return err
		}
		return err
	})

	return err
}

// UpdateMerchant updates the existing Merchant with the new properties
func (db *AppDatastore) UpdateMerchant(m *model.Merchant) error {

	_, err := db.client.RunInTransaction(context.Background(), func(tx *datastore.Transaction) error {
		key := db.merchantKey(m.MerchantID)
		old := &model.Merchant{}
		err := tx.Get(key, old)
		switch err {
		case nil:
			// entity exists, go ahead with update
			if err := validateMerchantUpdate(old, m); err != nil {
				return &c.ErrValidation{Violations: err}
			}

			// preserve other fields and update the last updated timestamp
			m.Updated = time.Now()
			m.Created = old.Created
			m.PayOutConfig = old.PayOutConfig

			// finally put
			_, err = tx.Put(key, m)
			return err
		case datastore.ErrNoSuchEntity:
			return c.ErrDBNoSuchEntity
		}
		return err
	})
	return err
}

// validateMerchantUpdate is a simple helper which checks where `to` can be updated from `from` or not
func validateMerchantUpdate(from *model.Merchant, to *model.Merchant) error {

	if from.OrganisationID != to.OrganisationID {
		return errors.New("changing organisationID is unallowed")
	}

	if from.MerchantID != to.MerchantID {
		return errors.New("changing organisationID is unallowed")
	}

	if from.FullName != to.FullName {
		return errors.New("changing fullName is unallowed")
	}

	if from.CurrencyCode != to.CurrencyCode {
		return errors.New("changing CurrencyCode is unallowed")
	}

	return nil
}

// UpsertPayOutConfig Merchant's PayOutConfig to be updated or inserted
func (db *AppDatastore) UpsertPayOutConfig(merchantID string, poc *model.PayOutConfig) error {
	_, err := db.client.RunInTransaction(context.Background(), func(tx *datastore.Transaction) error {
		key := db.merchantKey(merchantID)
		old := &model.Merchant{}
		err := tx.Get(key, old)
		switch err {
		case nil:
			old.PayOutConfig = poc
			old.Updated = time.Now()
			_, err = tx.Put(key, old)
			return err
		case datastore.ErrNoSuchEntity:
			// can't add config to non existent merchantID
			return c.ErrDBNoSuchEntity
		}
		return err
	})
	return err
}

func (db *AppDatastore) GetRole(organisationID, userID string) (*model.Role, error) {
	ctx := context.Background()
	k := db.roleDatastoreKey(organisationID, userID)
	role := &model.Role{}
	err := db.client.Get(ctx, k, role)

	switch err {
	case nil:
		return role, nil
	case datastore.ErrNoSuchEntity:
		return nil, c.ErrDBNoSuchEntity
	default:
		return nil, err
	}
}

// Close is needed to close the client connection
func (db *AppDatastore) Close() error {
	return db.client.Close()
}

// ----- Datastore keys -------------------------------------------------------

func (db *AppDatastore) merchantKey(id string) *datastore.Key {
	return datastore.NameKey(db.KindMerchant, id, nil)
}

func (db *AppDatastore) roleDatastoreKey(organisationID, userID string) *datastore.Key {
	return datastore.NameKey(db.KindRole, fmt.Sprintf("%v-%v", organisationID, userID), nil)
}

func (db *AppDatastore) kymKey(id string) *datastore.Key {
	return datastore.NameKey(db.KindKym, id, nil)
}

func (db *AppDatastore) teacherKey(id string) *datastore.Key {
	return datastore.NameKey(db.KindTeacher, id, nil)
}

func (db *AppDatastore) mainSubjectKey(id string) *datastore.Key {
	return datastore.NameKey(db.KindMainSubject, id, nil)
}

func (db *AppDatastore) subjectKey(id string) *datastore.Key {
	return datastore.NameKey(db.KindSubject, id, nil)
}

func (db *AppDatastore) confirmationKey(id string) *datastore.Key {
	return datastore.NameKey(db.KindConfirmation, id, nil)
}

func (db *AppDatastore) confirmationDetailKey(id string) *datastore.Key {
	return datastore.NameKey(db.KindConfirmationDetail, id, nil)
}

// AddKym attempts to add Kym to datastore.
func (db *AppDatastore) AddKym(kym *model.Kym) error {
	_, err := db.client.RunInTransaction(context.Background(), func(tx *datastore.Transaction) error {
		key := db.kymKey(kym.ID)
		err := tx.Get(key, &model.Kym{})

		switch err {
		case nil:
			return c.ErrDBEntityAlreadyExists
		case datastore.ErrNoSuchEntity:
			// no existing entity id, proceed with write
			_, err = tx.Put(key, kym)
			return err
		}
		return err
	})
	return err
}

// GetAllKym attempts to get all kym from datastore.
func (db *AppDatastore) GetAllKym(status string) ([]*model.Kym, error) {
	ctx := context.Background()

	var query *datastore.Query

	if len(status) > 0 {
		if c.IsValidKymStatus(status) {
			query = datastore.NewQuery(db.KindKym).Filter("Status=", status).Order("-DatetimeCreated").Limit(30)
		} else {
			return nil, &c.ErrValidation{Violations: errors.New(fmt.Sprintf("filter is not allowed with status : %v", status))}
		}
	} else {
		query = datastore.NewQuery(db.KindKym).Order("-DatetimeCreated").Limit(30)
	}

	klmList := make([]*model.Kym, 0)
	if _, err := db.client.GetAll(ctx, query, &klmList); err != nil {
		return nil, err
	}
	return klmList, nil
}

// GetKym attempts to get single kym from datastore by id.
func (db *AppDatastore) GetKym(id string) (*model.Kym, error) {
	key := db.kymKey(id)
	kym := &model.Kym{}
	err := db.client.Get(context.Background(), key, kym)
	switch err {
	case nil:
		return kym, nil
	case datastore.ErrNoSuchEntity:
		return nil, c.ErrDBNoSuchEntity
	}
	return nil, err
}

// UpdateKymStatus attempts to update kym status.
func (db *AppDatastore) UpdateKymStatus(kym *model.Kym, status string, notes string) error {
	_, err := db.client.RunInTransaction(context.Background(), func(tx *datastore.Transaction) error {
		key := db.kymKey(kym.ID)
		err := tx.Get(key, &model.Kym{})
		switch err {
		case nil:
			kym.Status = status
			kym.Notes = notes
			_, err = tx.Put(key, kym)
			return err
		case datastore.ErrNoSuchEntity:
			// can't add config to non-existent merchantID
			return c.ErrDBNoSuchEntity
		}
		return err
	})
	return err
}


// AddTeacher attempts to add a NewMerchant to the datastore.
// Since merchantID == organisationID, we do not allow duplicate organisationIDs
func (db *AppDatastore) AddTeacher(m *model.Teacher) error {
	_, err := db.client.RunInTransaction(context.Background(), func(tx *datastore.Transaction) error {
		key := db.teacherKey(m.ID)
		err := tx.Get(key, &model.Teacher{})

		switch err {
		case nil:
			return c.ErrDBEntityAlreadyExists
		case datastore.ErrNoSuchEntity:
			// no existing entity id, proceed with write
			_, err = tx.Put(key, m)
			return err
		}
		return err
	})

	return err
}

// GetAllTeacher attempts to get all kym from datastore.
func (db *AppDatastore) GetAllTeacher() ([]*model.Teacher, error) {
	ctx := context.Background()

	var query *datastore.Query

	query = datastore.NewQuery(db.KindTeacher).Limit(30)
	teacherList := make([]*model.Teacher, 0)
	if _, err := db.client.GetAll(ctx, query, &teacherList); err != nil {
		return nil, err
	}
	return teacherList, nil
}


// AddMainSubject attempts to add a NewMerchant to the datastore.
func (db *AppDatastore) AddMainSubject(m *model.MainSubject) error {
	_, err := db.client.RunInTransaction(context.Background(), func(tx *datastore.Transaction) error {
		key := db.mainSubjectKey(m.ID)
		err := tx.Get(key, &model.MainSubject{})

		switch err {
		case nil:
			return c.ErrDBEntityAlreadyExists
		case datastore.ErrNoSuchEntity:
			// no existing entity id, proceed with write
			_, err = tx.Put(key, m)
			return err
		}
		return err
	})

	return err
}

// GetAllMainSubject attempts to get all kym from datastore.
func (db *AppDatastore) GetAllMainSubject() ([]*model.MainSubject, error) {
	ctx := context.Background()

	var query *datastore.Query

	query = datastore.NewQuery(db.KindMainSubject).Limit(30)
	mainSubjectList := make([]*model.MainSubject, 0)
	if _, err := db.client.GetAll(ctx, query, &mainSubjectList); err != nil {
		return nil, err
	}
	return mainSubjectList, nil
}



func (db *AppDatastore) AddSubject(m *model.Subject) error {
	_, err := db.client.RunInTransaction(context.Background(), func(tx *datastore.Transaction) error {
		key := db.subjectKey(m.ID)
		err := tx.Get(key, &model.Subject{})

		switch err {
		case nil:
			return c.ErrDBEntityAlreadyExists
		case datastore.ErrNoSuchEntity:
			// no existing entity id, proceed with write
			_, err = tx.Put(key, m)
			return err
		}
		return err
	})

	return err
}


func (db *AppDatastore) GetAllSubject() ([]*model.Subject, error) {
	ctx := context.Background()

	var query *datastore.Query

	query = datastore.NewQuery(db.KindSubject).Limit(30)
	subjectList := make([]*model.Subject, 0)
	if _, err := db.client.GetAll(ctx, query, &subjectList); err != nil {
		return nil, err
	}
	return subjectList, nil
}


func (db *AppDatastore) AddConfirmation(m *model.Confirmation) error {
	_, err := db.client.RunInTransaction(context.Background(), func(tx *datastore.Transaction) error {
		key := db.confirmationKey(m.ID)
		err := tx.Get(key, &model.Confirmation{})

		switch err {
		case nil:
			return c.ErrDBEntityAlreadyExists
		case datastore.ErrNoSuchEntity:
			// no existing entity id, proceed with write
			_, err = tx.Put(key, m)
			return err
		}
		return err
	})

	return err
}


func (db *AppDatastore) AddConfirmationDetail(m *model.ConfirmationDetail) error {
	_, err := db.client.RunInTransaction(context.Background(), func(tx *datastore.Transaction) error {
		key := db.confirmationDetailKey(m.ID)
		err := tx.Get(key, &model.ConfirmationDetail{})

		switch err {
		case nil:
			return c.ErrDBEntityAlreadyExists
		case datastore.ErrNoSuchEntity:
			// no existing entity id, proceed with write
			_, err = tx.Put(key, m)
			return err
		}
		return err
	})

	return err
}



func (db *AppDatastore) GetAllConfirmation() ([]*model.Confirmation, error) {
	ctx := context.Background()

	var query *datastore.Query

	query = datastore.NewQuery(db.KindConfirmation).Limit(100)
	list := make([]*model.Confirmation, 0)
	if _, err := db.client.GetAll(ctx, query, &list); err != nil {
		return nil, err
	}
	return list, nil
}


func (db *AppDatastore) GetAllConfirmationDetail(confirmationId string) ([]*model.ConfirmationDetail, error) {
	ctx := context.Background()

	var query *datastore.Query

	query = datastore.NewQuery(db.KindConfirmationDetail).Filter("ConfirmationID =", confirmationId)
	list := make([]*model.ConfirmationDetail, 0)
	if _, err := db.client.GetAll(ctx, query, &list); err != nil {
		return nil, err
	}
	return list, nil
}


