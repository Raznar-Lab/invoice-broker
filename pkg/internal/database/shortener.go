package database

import (
	"errors"
	"sync"

	"raznar.id/invoice-broker/pkg/internal/database/models"
)

type shortenerContainer struct {
	m      sync.Mutex
	models []*models.ShortenerModel
}

func (d *Database) GetShortener(id string) (model *models.ShortenerModel) {
	d.sc.m.Lock()
	defer d.sc.m.Unlock()
	for _, s := range d.sc.models {
		if s.ID == id {
			model = s
			return
		}
	}

	return
}

func (d *Database) AddShortener(model *models.ShortenerModel) (err error) {
	if d.GetShortener(model.ID) != nil {
		err = errors.New("short url with that id is already exists")
		return
	}

	d.sc.m.Lock()
	defer d.sc.m.Unlock()
	d.sc.models = append(d.sc.models, model)

	d.SilentSave()
	return
}
