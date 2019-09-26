package model

import (
	"github.com/jinzhu/gorm"
)

type DomainModel interface {
	FindBy(map[string]interface{}) (Domains, error)
	UpdateByName(name string, newDoamin *Domain) (bool, error)
	DeleteByName(name string) (bool, error)
	UpdateByID(id string, newDoamin *Domain) (bool, error)
	DeleteByID(id string) (bool, error)
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

func (ds *Domains) ToIntreface() []interface{} {
	ret := []interface{}{}
	if ds != nil {
		for _, d := range *ds {
			ret = append(ret, d)
		}
	}
	return ret
}
func (d *Domain) FindBy(params map[string]interface{}) (Domains, error) {
	query := d.db.New().Preload("Records")
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

func (d *Domain) updateBy(db *gorm.DB, newDomain *Domain) (bool, error) {
	r := db.Take(&d)
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

func (d *Domain) UpdateByName(name string, newDomain *Domain) (bool, error) {
	return d.updateBy(d.db.New().Where("name = ?", name), newDomain)
}

func (d *Domain) UpdateByID(id string, newDomain *Domain) (bool, error) {
	return d.updateBy(d.db.New().Where("id = ?", id), newDomain)
}

func (d *Domain) DeleteBy(db *gorm.DB) (bool, error) {
	r := db.Take(&d)
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

	r = tx.Where("domain_id = ?", d.ID).Delete(&Record{})
	if r.Error != nil {
		if !r.RecordNotFound() {
			tx.Rollback()
			return false, r.Error
		}
	}

	r = tx.Delete(d)
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

func (d *Domain) DeleteByID(id string) (bool, error) {
	return d.DeleteBy(d.db.New().Where("id = ?", id))
}

func (d *Domain) DeleteByName(name string) (bool, error) {
	return d.DeleteBy(d.db.New().Where("name  = ?", name))
}

func (d *Domain) Create(newDomain *Domain) error {
	if err := d.db.New().Create(newDomain).Error; err != nil {
		return err
	}
	return nil
}
