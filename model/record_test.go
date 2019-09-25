package model

import (
	"reflect"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jinzhu/gorm"
)

func TestRecord_FindBy(t *testing.T) {
	type fields struct {
		db        *gorm.DB
		ID        int
		DomainID  int
		Name      string
		Type      string
		Content   string
		TTL       int
		Prio      int
		Disabled  bool
		OrderName string
		Auth      bool
		Domain    Domain
	}
	type args struct {
		params map[string]interface{}
	}
	tests := []struct {
		name       string
		fields     fields
		domainRows *sqlmock.Rows
		recordRows *sqlmock.Rows
		retErr     error
		args       args
		want       Records
		wantErr    bool
	}{
		{
			name:   "ok",
			fields: fields{},
			args: args{
				params: map[string]interface{}{
					"id": 1,
				},
			},
			domainRows: sqlmock.NewRows([]string{
				"id",
				"name",
				"master",
				"last_check",
				"type",
				"notified_serial",
				"account",
			}).
				AddRow(
					1,
					"test.com",
					"",
					1,
					"",
					1,
					"test",
				),
			recordRows: sqlmock.NewRows([]string{
				"id",
				"domain_id",
				"name",
				"type",
				"content",
				"ttl",
				"prio",
				"disabled",
				"ordername",
				"auth",
			}).
				AddRow(
					1,
					1,
					"test.com",
					"A",
					"1.1.1.1",
					"100",
					1,
					false,
					"",
					false,
				),
			want: Records{
				Record{
					ID:       1,
					DomainID: 1,
					Name:     "test.com",
					Type:     "A",
					Content:  "1.1.1.1",
					TTL:      100,
					Prio:     1,
					Domain: Domain{
						ID:             1,
						Name:           "test.com",
						LastCheck:      1,
						NotifiedSerial: 1,
						Account:        "test",
					},
				},
			},
		},
		{
			name:   "notfound",
			fields: fields{},
			args: args{
				params: map[string]interface{}{
					"id": 1,
				},
			},
			retErr: gorm.ErrRecordNotFound,
			want:   nil,
		},
		{
			name:   "other error",
			fields: fields{},
			args: args{
				params: map[string]interface{}{
					"id": 1,
				},
			},
			retErr:  gorm.ErrInvalidSQL,
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, mock, err := sqlmock.New()
			if err != nil {
				t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
			}
			defer db.Close()

			if tt.retErr == nil {
				mock.ExpectQuery("SELECT \\* FROM `records` WHERE \\(id in\\(\\?\\)\\)").
					WithArgs(1).
					WillReturnRows(tt.recordRows)

				mock.ExpectQuery("SELECT \\* FROM `domains` WHERE \\(`id` IN \\(\\?\\)\\)").
					WithArgs(1).
					WillReturnRows(tt.domainRows)

			} else {
				mock.ExpectQuery("SELECT \\* FROM `records` WHERE \\(id in\\(\\?\\)\\)").
					WillReturnError(tt.retErr)
				mock.ExpectQuery("SELECT \\* FROM `domains` WHERE \\(`id` IN \\(\\?\\)\\)").
					WillReturnError(tt.retErr)
			}

			gdb, _ := gorm.Open("mysql", db)

			d := &Record{
				db:        gdb,
				ID:        tt.fields.ID,
				DomainID:  tt.fields.DomainID,
				Name:      tt.fields.Name,
				Type:      tt.fields.Type,
				Content:   tt.fields.Content,
				TTL:       tt.fields.TTL,
				Prio:      tt.fields.Prio,
				Disabled:  tt.fields.Disabled,
				OrderName: tt.fields.OrderName,
				Auth:      tt.fields.Auth,
				Domain:    tt.fields.Domain,
			}
			got, err := d.FindBy(tt.args.params)
			if (err != nil) != tt.wantErr {
				t.Errorf("Record.FindBy() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Record.FindBy() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRecord_UpdateByID(t *testing.T) {
	type fields struct {
		db        *gorm.DB
		ID        int
		DomainID  int
		Name      string
		Type      string
		Content   string
		TTL       int
		Prio      int
		Disabled  bool
		OrderName string
		Auth      bool
		Domain    Domain
	}
	type args struct {
		id        string
		newRecord *Record
	}
	tests := []struct {
		name       string
		fields     fields
		args       args
		recordRows *sqlmock.Rows
		retErr     error
		want       bool
		wantErr    bool
	}{
		{
			name:   "ok",
			fields: fields{},
			args: args{
				id: "1",
				newRecord: &Record{
					Content: "2.2.2.2",
					Type:    "CNAME",
				},
			},
			recordRows: sqlmock.NewRows([]string{
				"id",
				"domain_id",
				"name",
				"type",
				"content",
				"ttl",
				"prio",
				"disabled",
				"ordername",
				"auth",
			}).
				AddRow(
					1,
					1,
					"test.com",
					"A",
					"1.1.1.1",
					"100",
					1,
					false,
					"",
					false,
				),
			want: true,
		},
		{
			name:   "notfound",
			fields: fields{},
			args: args{
				id: "2",
			},
			retErr: gorm.ErrRecordNotFound,
			want:   false,
		},
		{
			name:   "other error",
			fields: fields{},
			args: args{
				id: "3",
			},
			retErr:  gorm.ErrInvalidSQL,
			want:    false,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, mock, err := sqlmock.New()
			if err != nil {
				t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
			}
			defer db.Close()

			if tt.retErr == nil {
				mock.ExpectQuery("SELECT \\* FROM `records` WHERE \\(id = \\?\\) LIMIT 1").
					WithArgs("1").
					WillReturnRows(tt.recordRows)
				mock.ExpectBegin()
				mock.ExpectExec("UPDATE `records` SET `content` = \\?, `domain_id` = \\?, `type` = \\? WHERE `records`.`id` = \\?").
					WithArgs(`2.2.2.2`, 1, "CNAME", 1).WillReturnResult(
					sqlmock.NewResult(
						1,
						1,
					),
				)
				mock.ExpectCommit()
			} else {
				mock.ExpectQuery("SELECT \\* FROM `records` WHERE \\(id = \\?\\) LIMIT 1").
					WillReturnError(tt.retErr)
			}

			gdb, _ := gorm.Open("mysql", db)
			d := &Record{
				db:        gdb,
				ID:        tt.fields.ID,
				DomainID:  tt.fields.DomainID,
				Name:      tt.fields.Name,
				Type:      tt.fields.Type,
				Content:   tt.fields.Content,
				TTL:       tt.fields.TTL,
				Prio:      tt.fields.Prio,
				Disabled:  tt.fields.Disabled,
				OrderName: tt.fields.OrderName,
				Auth:      tt.fields.Auth,
				Domain:    tt.fields.Domain,
			}
			got, err := d.UpdateByID(tt.args.id, tt.args.newRecord)
			if (err != nil) != tt.wantErr {
				t.Errorf("Record.UpdateByID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Record.UpdateByID() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRecord_DeleteByID(t *testing.T) {
	type fields struct {
		db        *gorm.DB
		ID        int
		DomainID  int
		Name      string
		Type      string
		Content   string
		TTL       int
		Prio      int
		Disabled  bool
		OrderName string
		Auth      bool
		Domain    Domain
	}
	type args struct {
		id string
	}
	tests := []struct {
		name       string
		fields     fields
		args       args
		recordRows *sqlmock.Rows
		retErr     error
		want       bool
		wantErr    bool
	}{
		{
			name:   "ok",
			fields: fields{},
			args: args{
				id: "1",
			},
			recordRows: sqlmock.NewRows([]string{
				"id",
				"domain_id",
				"name",
				"type",
				"content",
				"ttl",
				"prio",
				"disabled",
				"ordername",
				"auth",
			}).
				AddRow(
					1,
					1,
					"test.com",
					"A",
					"1.1.1.1",
					"100",
					1,
					false,
					"",
					false,
				),
			want: true,
		},
		{
			name:   "notfound",
			fields: fields{},
			args: args{
				id: "2",
			},
			retErr: gorm.ErrRecordNotFound,
			want:   false,
		},
		{
			name:   "other error",
			fields: fields{},
			args: args{
				id: "3",
			},
			retErr:  gorm.ErrInvalidSQL,
			want:    false,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, mock, err := sqlmock.New()
			if err != nil {
				t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
			}
			defer db.Close()

			if tt.retErr == nil {
				mock.ExpectQuery("SELECT \\* FROM `records` WHERE \\(id = \\?\\) LIMIT 1").
					WithArgs("1").
					WillReturnRows(tt.recordRows)
				mock.ExpectBegin()
				mock.ExpectExec("DELETE FROM `records` WHERE `records`.`id` = \\?").
					WithArgs(1).WillReturnResult(
					sqlmock.NewResult(
						1,
						1,
					),
				)
				mock.ExpectCommit()
			} else {
				mock.ExpectQuery("SELECT \\* FROM `records` WHERE \\(id = \\?\\) LIMIT 1").
					WillReturnError(tt.retErr)
			}

			gdb, _ := gorm.Open("mysql", db)
			d := &Record{
				db:        gdb,
				ID:        tt.fields.ID,
				DomainID:  tt.fields.DomainID,
				Name:      tt.fields.Name,
				Type:      tt.fields.Type,
				Content:   tt.fields.Content,
				TTL:       tt.fields.TTL,
				Prio:      tt.fields.Prio,
				Disabled:  tt.fields.Disabled,
				OrderName: tt.fields.OrderName,
				Auth:      tt.fields.Auth,
				Domain:    tt.fields.Domain,
			}
			got, err := d.DeleteByID(tt.args.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("Record.DeleteByID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Record.DeleteByID() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRecord_Create(t *testing.T) {
	type fields struct {
		db        *gorm.DB
		ID        int
		DomainID  int
		Name      string
		Type      string
		Content   string
		TTL       int
		Prio      int
		Disabled  bool
		OrderName string
		Auth      bool
		Domain    Domain
	}
	type args struct {
		newRecord *Record
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		retErr  error
		want    bool
		wantErr bool
	}{
		{
			name:   "ok",
			fields: fields{},
			args: args{
				newRecord: &Record{
					DomainID: 1,
					Name:     "test.com",
					Type:     "A",
					Content:  "1.1.1.1",
					TTL:      100,
					Prio:     200,
					Disabled: false,
					Auth:     false,
				},
			},
			want: true,
		},
		{
			name:   "other error",
			fields: fields{},
			args: args{
				newRecord: &Record{
					Content: "2.2.2.2",
					Type:    "CNAME",
				},
			},
			retErr:  gorm.ErrInvalidSQL,
			want:    false,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, mock, err := sqlmock.New()
			if err != nil {
				t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
			}
			defer db.Close()

			if tt.retErr == nil {
				mock.ExpectBegin()
				mock.ExpectExec("INSERT INTO `records` \\(`domain_id`,`name`,`type`,`content`,`ttl`,`prio`,`disabled`,`ordername`,`auth`\\) VALUES \\(\\?,\\?,\\?,\\?,\\?,\\?,\\?,\\?,\\?\\)").
					WithArgs(
						1,
						"test.com",
						"A",
						"1.1.1.1",
						100,
						200,
						false,
						"",
						false,
					).WillReturnResult(
					sqlmock.NewResult(
						1,
						1,
					),
				)
				mock.ExpectCommit()

			} else {
				mock.ExpectExec("INSERT INTO `records` \\(`domain_id`,`name`,`type`,`content`,`ttl`,`prio`,`disabled`,`ordername`,`auth`\\) VALUES \\(\\?,\\?,\\?,\\?,\\?,\\?,\\?,\\?,\\?\\)").
					WillReturnError(tt.retErr)
			}

			gdb, _ := gorm.Open("mysql", db)
			d := &Record{
				db:        gdb,
				ID:        tt.fields.ID,
				DomainID:  tt.fields.DomainID,
				Name:      tt.fields.Name,
				Type:      tt.fields.Type,
				Content:   tt.fields.Content,
				TTL:       tt.fields.TTL,
				Prio:      tt.fields.Prio,
				Disabled:  tt.fields.Disabled,
				OrderName: tt.fields.OrderName,
				Auth:      tt.fields.Auth,
				Domain:    tt.fields.Domain,
			}
			err = d.Create(tt.args.newRecord)
			if (err != nil) != tt.wantErr {
				t.Errorf("Record.Create() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}
