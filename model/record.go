package model

import (
	"github.com/jinzhu/gorm"
)

type RecordModel struct {
	db *gorm.DB
}
type Record struct {
	ID        int    `json:"id"`
	DomainID  int    `json:"domain_id"`
	Name      string `json:"name"`
	Type      string `json:"type"`
	Content   string `json:"content"`
	TTL       int    `json:"ttl"`
	Prio      int    `json:"prio"`
	Disabled  *bool  `json:"disabled"`
	OrderName string `json:"ordername" gorm:"column:ordername"`
	Auth      *bool  `json:"auth"`
	Domain    Domain `json:"-"`
}

type Records []Record
type RecordModeler interface {
	FindBy(map[string]interface{}) (Records, error)
	UpdateByID(string, *Record) (bool, error)
	DeleteByID(string) (bool, error)
	Create(*Record) error
}

func NewRecordModeler(db *gorm.DB) *RecordModel {
	return &RecordModel{
		db: db,
	}
}
func (rs *Records) ToIntreface() []interface{} {
	ret := []interface{}{}
	if rs != nil {
		for _, d := range *rs {
			ret = append(ret, d)
		}
	}
	return ret
}

func (d *RecordModel) FindBy(params map[string]interface{}) (Records, error) {
	query := d.db.Preload("Domain")
	for k, v := range params {
		query = query.Where(k+" in(?)", v)
	}

	ds := Records{}
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
func (d *RecordModel) UpdateByID(id string, newRecord *Record) (bool, error) {
	record := &Record{}
	r := d.db.Where("id = ?", id).Take(&record)
	if r.Error != nil {
		if r.RecordNotFound() {
			return false, nil
		} else {
			return false, r.Error
		}
	}

	newRecord.DomainID = record.DomainID
	r = d.db.Model(&record).Updates(&newRecord)

	if r.Error != nil {
		return false, r.Error
	}
	return true, nil
}
func (d *RecordModel) DeleteByID(id string) (bool, error) {
	record := &Record{}
	r := d.db.Where("id = ?", id).Take(&record)
	if r.Error != nil {
		if r.RecordNotFound() {
			return false, nil
		} else {
			return false, r.Error
		}
	}

	r = d.db.Delete(record)
	if r.Error != nil {
		return false, r.Error
	}
	return true, nil
}

func (d *RecordModel) Create(newRecord *Record) error {
	f := false
	if newRecord.Disabled == nil {
		newRecord.Disabled = &f
	}
	if newRecord.Auth == nil {
		newRecord.Auth = &f
	}
	if err := d.db.Create(newRecord).Error; err != nil {
		return err
	}
	return nil
}
