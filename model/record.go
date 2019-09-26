package model

import (
	"github.com/jinzhu/gorm"
)

type Record struct {
	db        *gorm.DB
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
type RecordModel interface {
	FindBy(map[string]interface{}) (Records, error)
	UpdateByID(string, *Record) (bool, error)
	DeleteByID(string) (bool, error)
	Create(*Record) error
}

func NewRecordModel(db *gorm.DB) *Record {
	return &Record{
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

func (d *Record) FindBy(params map[string]interface{}) (Records, error) {
	query := d.db.New().Preload("Domain")
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
func (d *Record) UpdateByID(id string, newRecord *Record) (bool, error) {
	record := &Record{}
	r := d.db.New().Where("id = ?", id).Take(&record)
	if r.Error != nil {
		if r.RecordNotFound() {
			return false, nil
		} else {
			return false, r.Error
		}
	}

	newRecord.DomainID = d.DomainID
	r = d.db.Model(&d).Updates(&newRecord)

	if r.Error != nil {
		return false, r.Error
	}
	return true, nil
}
func (d *Record) DeleteByID(id string) (bool, error) {
	record := &Record{}
	r := d.db.New().Where("id = ?", id).Take(&record)
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

func (d *Record) Create(newRecord *Record) error {
	if err := d.db.New().Create(newRecord).Error; err != nil {
		return err
	}
	return nil
}
