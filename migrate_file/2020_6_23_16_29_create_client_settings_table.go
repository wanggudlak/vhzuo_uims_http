package migrate_file

import (
	"fmt"
	"uims/internal/model"
	"uims/pkg/db"
)

type CreateClientSettingsTableMigrate struct {
}

func (CreateClientSettingsTableMigrate) Key() string {
	return "2020_6_23_14_41_create_user_org_table"
}

// migrate
func (CreateClientSettingsTableMigrate) Up() (err error) {
	if db.Def().HasTable(model.ClientSetting{}.TableName()) {
		err = fmt.Errorf("client_settings table alreay exist")
		return
	}
	err = db.Def().
		Set("gorm:table_options", "CHARSET=utf8mb4,COMMENT='客户端业务系统设置'").
		CreateTable(&model.ClientSetting{}).Error
	return
}

// rollback
func (CreateClientSettingsTableMigrate) Down() (err error) {
	err = db.Def().DropTableIfExists(&model.ClientSetting{}).Error
	return
}
