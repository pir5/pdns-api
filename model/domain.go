package model

import (
	"net/http"

	"github.com/jinzhu/gorm"
	null "gopkg.in/guregu/null.v3"
)

type DomainModeler interface {
	FindBy(*http.Request, map[string]interface{}) (Domains, int64, error)
	UpdateByName(name string, newDoamin *Domain) (bool, error)
	DeleteByName(name string) (bool, error)
	UpdateByID(id string, newDoamin *Domain) (bool, error)
	DeleteByID(id string) (bool, error)
	Create(newDoamin *Domain) error
}

func NewDomainModeler(db *gorm.DB) *DomainModel {
	return &DomainModel{
		db: db,
	}
}

type DomainModel struct {
	db *gorm.DB
}

type Domain struct {
	ID             int         `json:"id" swaggertype:"integer" format:"int64" gorm:"primary_key"`
	Name           null.String `json:"name" validate:"required,fqdn" swaggertype:"string" format:"string" gorm:"column:name"`
	Master         null.String `json:"master" swaggertype:"string" format:"string" gorm:"column:master"`
	LastCheck      null.Int    `json:"last_check,number" swaggertype:"integer" format:"int64" gorm:"column:last_check"`
	Type           null.String `json:"type" validate:"oneof=NATIVE MASTER SLAVE" swaggertype:"string" format:"string" gorm:"column:type"`
	NotifiedSerial null.Int    `json:"notified_serial,number" swaggertype:"integer" format:"int64" gorm:"column:notified_serial"`
	Account        null.String `json:"account" swaggertype:"string" format:"string" gorm:"column:account"`
	Records        *Records    `json:"records"`
}

type Domains []Domain

func (d *DomainModel) TableName() string {
	return "domains"
}

func (d *DomainModel) FindBy(req *http.Request, params map[string]interface{}) (Domains, int64, error) {
	query := d.db
	for k, v := range params {
		query = query.Where(k+" in(?)", v)
	}

	ds := Domains{}
	r := query.Scopes(Paginate(req)).Find(&ds)
	if r.Error != nil {
		if r.RecordNotFound() {
			return nil, 0, nil
		} else {
			return nil, 0, r.Error
		}
	}

	totalCount := int64(0)
	r = query.Model(&Domains{}).Count(&totalCount)
	if r.Error != nil {
		return nil, 0, r.Error
	}

	return ds, totalCount, nil
}

func (d *DomainModel) updateBy(db *gorm.DB, newDomain *Domain) (bool, error) {
	r := db.Take(d)
	if r.Error != nil {
		if r.RecordNotFound() {
			return false, nil
		} else {
			return false, r.Error
		}
	}

	r = d.db.Model(&Domain{}).Updates(newDomain)
	if r.Error != nil {
		return false, r.Error
	}
	return true, nil
}

func (d *DomainModel) UpdateByName(name string, newDomain *Domain) (bool, error) {
	return d.updateBy(d.db.Where("name = ?", name), newDomain)
}

func (d *DomainModel) UpdateByID(id string, newDomain *Domain) (bool, error) {
	return d.updateBy(d.db.Where("id = ?", id), newDomain)
}

func (d *DomainModel) DeleteBy(db *gorm.DB) (bool, error) {
	domain := Domain{}
	r := db.Take(&domain)
	if r.Error != nil {
		if r.RecordNotFound() {
			return false, nil
		} else {
			return false, r.Error
		}
	}

	tx := d.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	r = tx.Where("domain_id = ?", domain.ID).Delete(&Record{})
	if r.Error != nil {
		if !r.RecordNotFound() {
			tx.Rollback()
			return false, r.Error
		}
	}

	r = tx.Delete(&domain)
	if r.Error != nil {
		tx.Rollback()
		return false, r.Error
	}

	r = tx.Commit()
	if r.Error != nil {
		return false, r.Error
	}
	return true, nil
}

func (d *DomainModel) DeleteByID(id string) (bool, error) {
	return d.DeleteBy(d.db.Where("id = ?", id))
}

func (d *DomainModel) DeleteByName(name string) (bool, error) {
	return d.DeleteBy(d.db.Where("name  = ?", name))
}

func (d *DomainModel) Create(newDomain *Domain) error {
	if err := d.db.Create(newDomain).Error; err != nil {
		return err
	}
	return nil
}
