package database

import (
	"encoding/json"
	"fmt"
	"os"
	"sync"

	"raznar.id/invoice-broker/pkg/internal/database/models"
)

type Database struct {
	filePath     string
	transactions []*models.TransactionModel
	m sync.Mutex 
}

func (d *Database) Get(transactionId string) (transaction *models.TransactionModel) {
	d.m.Lock()
	defer d.m.Unlock()
	for _, t := range d.transactions {
		if t.TransactionID == transactionId {
			transaction = t
			return
		}
	}

	return
}

func (d *Database) Add(transaction *models.TransactionModel) (err error) {
	d.m.Lock()
	defer d.m.Unlock()
	if d.Get(transaction.TransactionID) != nil {
		return
	}

	fmt.Printf("added %s\n", transaction.TransactionID)
	d.transactions = append(d.transactions, transaction)
	return d.Save()
}

func (d *Database) Load() (err error) {
	res, err := os.ReadFile(d.filePath)
	if os.IsNotExist(err) {
		err = d.Save()
		if err != nil {
			return
		}

		return d.Load()
	}

	if err != nil {
		return
	}

	return json.Unmarshal(res, &d.transactions)
}

func (d *Database) Save() (err error) {
	res, err := json.Marshal(d.transactions)
	if err != nil {
		return
	}

	fmt.Printf("saved with %d data(s)\n", len(d.transactions))
	return os.WriteFile(d.filePath, res, 0644)
}

func New(filePath string) *Database {
	return &Database{filePath: filePath, transactions: []*models.TransactionModel{}}
}
