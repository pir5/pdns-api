package model

import "github.com/jinzhu/gorm"

type Record struct {
	db        *gorm.DB
	ID        int    `json:"id"`
	DomainID  int    `json:"domain_id"`
	Name      string `json:"name"`
	Type      string `json:"type"`
	Content   string `json:"content"`
	TTL       int    `json:"ttl"`
	Prio      int    `json:"prio"`
	Disabled  bool   `json:"disabled"`
	OlderName string `json:"older_name"`
	Auth      bool   `json:"auth"`
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
func (d *Record) FindBy(params map[string]interface{}) (Records, error) {
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
func (d *Record) UpdateByID(id string, newRecord *Record) (bool, error) {
	r := d.db.Where("id = ?", id).Take(&d)
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
	r := d.db.Where("id = ?", id).Take(&d)
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

func (d *Record) Create(newRecord *Record) error {
	if err := d.db.Create(newRecord).Error; err != nil {
		return err
	}
	return nil
}
