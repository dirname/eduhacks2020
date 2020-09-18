package database

import (
	"gorm.io/gorm"
	"testing"
)

func TestOrm_Init(t *testing.T) {
	type fields struct {
		DB *gorm.DB
	}
	var orm ORM
	DefaultSetting()
	orm.Init()
	tests := []struct {
		name   string
		fields fields
	}{
		{"ping", fields{DB: orm.DB}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			o := &ORM{
				DB: tt.fields.DB,
			}
			sqlDB, err := o.DB.DB()
			if err != nil {
				t.Error(err.Error())
			}
			err = sqlDB.Ping()
			if err != nil {
				t.Error(err.Error())
			} else {
				t.Log("success")
			}
			defer o.Close()
		})
	}
}
