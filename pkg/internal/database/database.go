package database

import (
	"encoding/json"
	"fmt"
	"os"
	"sync"

	"raznar.id/invoice-broker/pkg/internal/database/models"
)

type databaseData struct {
	Shortener    []*models.ShortenerModel   `json:"shortener"`
	Transactions []*models.TransactionModel `json:"transactions"`
}

type Database struct {
	filePath string
	tc       transactionContainer
	sc       shortenerContainer
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

	var data databaseData
	err = json.Unmarshal(res, &data)
	if err != nil {
		return
	}

	d.sc.models = data.Shortener
	d.tc.models = data.Transactions
	return
}

func (d *Database) Save() (err error) {
	res, err := json.Marshal(databaseData{
		Shortener:    d.sc.models,
		Transactions: d.tc.models,
	})
	if err != nil {
		return
	}

	return os.WriteFile(d.filePath, res, 0644)
}

// save in routine so the web can proceed faster.
func (d *Database) SilentSave()  {
	go func() {
		err := d.Save()
		if err != nil {
			fmt.Println("an error occured while saving the database: " + err.Error())
		}
	}()
}

func New(filePath string) *Database {
	return &Database{filePath: filePath, tc: transactionContainer{m: sync.Mutex{}, models: []*models.TransactionModel{}}, sc: shortenerContainer{m: sync.Mutex{}, models: []*models.ShortenerModel{}}}
}
