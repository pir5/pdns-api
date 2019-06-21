package model

import (
	"github.com/jinzhu/gorm"
)

type DomainModel interface {
	FindBy(map[string]interface{}) (Domains, error)
	UpdateByName(name string, newDoamin *Domain) (bool, error)
	DeleteByName(name string) (bool, error)
	Create(newDoamin *Domain) error
}

func NewDomainModel(db *gorm.DB) *Domain {
	return &Domain{
		db: db,
	}
}

type Domain struct {
	db             *gorm.DB
	ID             int     `json:"id"`
	Name           string  `json:"name"`
	Master         string  `json:"master"`
	LastCheck      int     `json:"last_check"`
	Type           string  `json:"type"`
	NotifiedSerial int32   `json:"notified_serial"`
	Account        string  `json:"account"`
	Records        Records `json:"records"`
}

type Domains []Domain

func (d *Domain) FindBy(params map[string]interface{}) (Domains, error) {
	query := d.db.Debug().Preload("Records")
	for k, v := range params {
		query = query.Where(k+" in(?)", v)
	}

	ds := Domains{}
	r := query.Find(&ds)
	if r.Error != nil {
		if r.RecordNotFound() {
			return nil, nil
		} else {
			return nil, r.Error
		}
	}

	return ds, nil
}
func (d *Domain) UpdateByName(name string, newDomain *Domain) (bool, error) {
	r := d.db.Where("name = ?", name).Take(&d)
	if r.Error != nil {
		if r.RecordNotFound() {
			return false, nil
		} else {
			return false, r.Error
		}
	}

	r = d.db.Model(&d).Updates(&newDomain)
	if r.Error != nil {
		return false, r.Error
	}
	return true, nil
}
func (d *Domain) DeleteByName(name string) (bool, error) {
	r := d.db.Where("name = ?", name).Take(&d)
	if r.Error != nil {
		if r.RecordNotFound() {
			return false, nil
		} else {
			return false, r.Error
		}
	}

	r = d.db.Delete(d)
	if r.Error != nil {
		return false, r.Error
	}
	return true, nil
}

func (d *Domain) Create(newDomain *Domain) error {
	if err := d.db.Create(newDomain).Error; err != nil {
		return err
	}
	return nil
}
