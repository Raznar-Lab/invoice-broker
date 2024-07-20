package database

import (
	"errors"
	"sync"

	"raznar.id/invoice-broker/pkg/internal/database/models"
)

type transactionContainer struct {
	m      sync.Mutex
	models []*models.TransactionModel
}

func (d *Database) GetTransaction(transactionId string) (transaction *models.TransactionModel) {
	d.tc.m.Lock()
	defer d.tc.m.Unlock()
	for _, t := range d.tc.models {
		if t.TransactionID == transactionId {
			transaction = t
			return
		}
	}

	return
}

func (d *Database) AddTransaction(transaction *models.TransactionModel) (err error) {
	if d.GetTransaction(transaction.TransactionID) != nil {
		err = errors.New("short url with that id is already exists")
		return
	}

	d.tc.m.Lock()
	defer d.tc.m.Unlock()

	d.tc.models = append(d.tc.models, transaction)

	d.SilentSave()
	return
}
