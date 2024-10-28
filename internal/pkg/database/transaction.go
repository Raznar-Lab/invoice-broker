package database

import (
	"errors"
	"sync"

	"raznar.id/invoice-broker/internal/models"
)

type transactionContainer struct {
	m      sync.Mutex
	models []*models.TransactionModel
}

func (d *Database) GetTransactionByTrID(trId string) (transaction *models.TransactionModel) {
	d.tc.m.Lock()
	defer d.tc.m.Unlock()
	for _, t := range d.tc.models {
		if t.TransactionID == trId {
			transaction = t
			return
		}
	}

	return
}

func (d *Database) GetTransaction(invoiceId string) (transaction *models.TransactionModel) {
	d.tc.m.Lock()
	defer d.tc.m.Unlock()
	for _, t := range d.tc.models {
		if t.Id == invoiceId {
			transaction = t
			return
		}
	}

	return
}

func (d *Database) AddTransaction(transaction *models.TransactionModel) (err error) {
	if d.GetTransactionByTrID(transaction.TransactionID) != nil {
		err = errors.New("short url with that id is already exists")
		return
	}

	d.tc.m.Lock()
	defer d.tc.m.Unlock()

	d.tc.models = append(d.tc.models, transaction)

	err = d.Save()
	return
}
