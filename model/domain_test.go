package model

import (
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jinzhu/gorm"
)

func TestDomain_FindBy(t *testing.T) {
	type args struct {
		params map[string]interface{}
	}
	tests := []struct {
		name       string
		domainRows *sqlmock.Rows
		recordRows *sqlmock.Rows
		args       args
		retErr     error
		want       Domains
		wantErr    bool
	}{
		{
			name: "ok",
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
			want: Domains{
				Domain{
					ID:             1,
					Name:           "test.com",
					LastCheck:      1,
					NotifiedSerial: 1,
					Account:        "test",
					Records: Records{
						Record{
							ID:        1,
							DomainID:  1,
							Name:      "test.com",
							Type:      "A",
							Content:   "1.1.1.1",
							TTL:       100,
							Prio:      1,
							Disabled:  newBool(false),
							OrderName: "",
							Auth:      newBool(false),
						},
					},
				},
			},
		},
		{
			name: "notfound",
			args: args{
				params: map[string]interface{}{
					"id": 1,
				},
			},
			retErr: gorm.ErrRecordNotFound,
			want:   nil,
		},
		{
			name: "other error",
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
				mock.ExpectQuery("SELECT \\* FROM `domains` WHERE \\(id in\\(\\?\\)\\)").
					WithArgs(1).
					WillReturnRows(tt.domainRows)

				mock.ExpectQuery("SELECT \\* FROM `records`  WHERE \\(`domain_id` IN \\(\\?\\)\\)").
					WithArgs(1).
					WillReturnRows(tt.recordRows)
			} else {
				mock.ExpectQuery("SELECT \\* FROM `domains` WHERE \\(id in\\(\\?\\)\\)").
					WillReturnError(tt.retErr)

				mock.ExpectQuery("SELECT \\* FROM `records`  WHERE \\(`domain_id` IN \\(\\?\\)\\)").
					WillReturnError(tt.retErr)
			}

			gdb, _ := gorm.Open("mysql", db)
			d := &DomainModel{
				db: gdb,
			}

			req := httptest.NewRequest("GET", "/", nil)
			got, _, err := d.FindBy(req, tt.args.params)
			if (err != nil) != tt.wantErr {
				t.Errorf("Domain.FindBy() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Domain.FindBy() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDomain_UpdateByName(t *testing.T) {
	type args struct {
		name      string
		newDomain *Domain
	}
	tests := []struct {
		name       string
		domainRows *sqlmock.Rows
		args       args
		retErr     error
		want       bool
		wantErr    bool
	}{
		{
			name: "ok",
			args: args{
				name: "test.com",
				newDomain: &Domain{
					ID:      1,
					Name:    "test.com",
					Account: "update",
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
			want: true,
		},
		{
			name: "notfound",
			args: args{
				name: "test.com",
			},
			retErr: gorm.ErrRecordNotFound,
			want:   false,
		},
		{
			name: "other error",
			args: args{
				name: "test.com",
				newDomain: &Domain{
					ID:      1,
					Name:    "test.com",
					Account: "update",
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
				mock.ExpectQuery("SELECT \\* FROM `domains` WHERE \\(name = \\?\\) LIMIT 1").
					WithArgs("test.com").
					WillReturnRows(tt.domainRows)
				mock.ExpectBegin()
				mock.ExpectExec("UPDATE `domains` SET `account` = \\?, `id` = \\?, `name` = \\?  WHERE `domains`.`id` = \\?").
					WithArgs(`update`, 1, "test.com", 1).WillReturnResult(
					sqlmock.NewResult(
						1,
						1,
					),
				)
				mock.ExpectCommit()
			} else {
				mock.ExpectQuery("SELECT \\* FROM `domains` WHERE \\(name = \\?\\) LIMIT 1").
					WillReturnError(tt.retErr)
			}

			gdb, _ := gorm.Open("mysql", db)
			d := &DomainModel{
				db: gdb,
			}

			got, err := d.UpdateByName(tt.args.name, tt.args.newDomain)
			if (err != nil) != tt.wantErr {
				t.Errorf("Domain.UpdateByName() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Domain.UpdateByName() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDomain_DeleteByName(t *testing.T) {
	type args struct {
		name string
	}
	tests := []struct {
		name       string
		domainRows *sqlmock.Rows
		args       args
		retErr     error
		want       bool
		wantErr    bool
	}{
		{
			name: "ok",
			args: args{
				name: "test.com",
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
			want: true,
		},
		{
			name: "notfound",
			args: args{
				name: "test.com",
			},
			retErr: gorm.ErrRecordNotFound,
			want:   false,
		},
		{
			name: "other error",
			args: args{
				name: "test.com",
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
				mock.ExpectQuery("SELECT \\* FROM `domains` WHERE \\(name = \\?\\) LIMIT 1").
					WithArgs("test.com").
					WillReturnRows(tt.domainRows)
				mock.ExpectBegin()

				mock.ExpectExec("DELETE FROM `records` WHERE \\(domain_id = \\?\\)").
					WithArgs(1).WillReturnResult(
					sqlmock.NewResult(
						1,
						1,
					),
				)

				mock.ExpectExec("DELETE FROM `domains` WHERE `domains`.`id` = \\?").
					WithArgs(1).WillReturnResult(
					sqlmock.NewResult(
						1,
						1,
					),
				)
				mock.ExpectCommit()
			} else {
				mock.ExpectQuery("SELECT \\* FROM `domains` WHERE \\(name = \\?\\) LIMIT 1").
					WillReturnError(tt.retErr)
			}

			gdb, _ := gorm.Open("mysql", db)
			d := &DomainModel{
				db: gdb,
			}

			got, err := d.DeleteByName(tt.args.name)
			if (err != nil) != tt.wantErr {
				t.Errorf("Domain.DeleteByName() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Domain.DeleteByName() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDomain_Create(t *testing.T) {
	type args struct {
		newDomain *Domain
	}
	tests := []struct {
		name    string
		args    args
		retErr  error
		wantErr bool
	}{
		{
			name: "ok",
			args: args{
				newDomain: &Domain{
					ID:   1,
					Name: "test.com",
				},
			},
		},
		{
			name: "other error",
			args: args{
				newDomain: &Domain{},
			},
			retErr:  gorm.ErrInvalidSQL,
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
				mock.ExpectExec("INSERT INTO `domains` \\(`id`,`name`,`master`,`last_check`,`type`,`notified_serial`,`account`\\) VALUES \\(\\?,\\?,\\?,\\?,\\?,\\?,\\?\\)").
					WithArgs(1, "test.com", "", 0, "", 0, "").WillReturnResult(
					sqlmock.NewResult(
						1,
						1,
					),
				)
				mock.ExpectCommit()
			} else {
				mock.ExpectExec("INSERT INTO `domains` \\(`id`,`name`,`master`,`last_check`,`type`,`notified_serial`,`account`\\) VALUES \\(\\?,\\?,\\?,\\?,\\?,\\?,\\?\\)").
					WillReturnError(tt.retErr)
			}

			gdb, _ := gorm.Open("mysql", db)
			d := &DomainModel{
				db: gdb,
			}

			err = d.Create(tt.args.newDomain)
			if (err != nil) != tt.wantErr {
				t.Errorf("Domain.Create() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func newBool(b bool) *bool {
	return &b
}

func TestDomain_UpdateByID(t *testing.T) {
	type args struct {
		id        string
		newDomain *Domain
	}
	tests := []struct {
		name       string
		domainRows *sqlmock.Rows
		args       args
		retErr     error
		want       bool
		wantErr    bool
	}{
		{
			name: "ok",
			args: args{
				id: "1",
				newDomain: &Domain{
					ID:      1,
					Name:    "test.com",
					Account: "update",
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
			want: true,
		},
		{
			name: "notfound",
			args: args{
				id: "2",
			},
			retErr: gorm.ErrRecordNotFound,
			want:   false,
		},
		{
			name: "other error",
			args: args{
				id: "3",
				newDomain: &Domain{
					ID:      1,
					Name:    "test.com",
					Account: "update",
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
				mock.ExpectQuery("SELECT \\* FROM `domains` WHERE \\(id = \\?\\) LIMIT 1").
					WithArgs("1").
					WillReturnRows(tt.domainRows)
				mock.ExpectBegin()
				mock.ExpectExec("UPDATE `domains` SET `account` = \\?, `id` = \\?, `name` = \\?  WHERE `domains`.`id` = \\?").
					WithArgs(`update`, 1, "test.com", 1).WillReturnResult(
					sqlmock.NewResult(
						1,
						1,
					),
				)
				mock.ExpectCommit()
			} else {
				mock.ExpectQuery("SELECT \\* FROM `domains` WHERE \\(id = \\?\\) LIMIT 1").
					WillReturnError(tt.retErr)
			}

			gdb, _ := gorm.Open("mysql", db)
			d := &DomainModel{
				db: gdb,
			}

			got, err := d.UpdateByID(tt.args.id, tt.args.newDomain)
			if (err != nil) != tt.wantErr {
				t.Errorf("Domain.UpdateByID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Domain.UpdateByID() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDomain_DeleteByID(t *testing.T) {
	type args struct {
		id string
	}
	tests := []struct {
		name       string
		domainRows *sqlmock.Rows
		args       args
		retErr     error
		want       bool
		wantErr    bool
	}{
		{
			name: "ok",
			args: args{
				id: "1",
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
			want: true,
		},
		{
			name: "notfound",
			args: args{
				id: "2",
			},
			retErr: gorm.ErrRecordNotFound,
			want:   false,
		},
		{
			name: "other error",
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
				mock.ExpectQuery("SELECT \\* FROM `domains` WHERE \\(id = \\?\\) LIMIT 1").
					WithArgs("1").
					WillReturnRows(tt.domainRows)
				mock.ExpectBegin()

				mock.ExpectExec("DELETE FROM `records` WHERE \\(domain_id = \\?\\)").
					WithArgs(1).WillReturnResult(
					sqlmock.NewResult(
						1,
						1,
					),
				)

				mock.ExpectExec("DELETE FROM `domains` WHERE `domains`.`id` = \\?").
					WithArgs(1).WillReturnResult(
					sqlmock.NewResult(
						1,
						1,
					),
				)
				mock.ExpectCommit()
			} else {
				mock.ExpectQuery("SELECT \\* FROM `domains` WHERE \\(id = \\?\\) LIMIT 1").
					WillReturnError(tt.retErr)
			}

			gdb, _ := gorm.Open("mysql", db)
			d := &DomainModel{
				db: gdb,
			}

			got, err := d.DeleteByID(tt.args.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("Domain.DeleteByID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Domain.DeleteByID() = %v, want %v", got, tt.want)
			}
		})
	}
}
