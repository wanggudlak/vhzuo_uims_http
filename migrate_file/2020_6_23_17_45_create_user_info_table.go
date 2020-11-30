package migrate_file

import (
	"fmt"
	"uims/internal/model"
	"uims/pkg/db"
)

type CreateUserInfoTableMigrate struct {
}

func (CreateUserInfoTableMigrate) Key() string {
	return "2020_6_23_17_45_create_user_info_table"
}

// migrate
func (CreateUserInfoTableMigrate) Up() (err error) {
	if db.Def().HasTable(model.UserInfo{}.TableName()) {
		err = fmt.Errorf("user_info table alreay exist")
		return
	}
	err = db.Def().
		Set("gorm:table_options", "CHARSET=utf8mb4,COMMENT='用户资料库'").
		CreateTable(&model.UserInfo{}).Error
	return
}

// rollback
func (CreateUserInfoTableMigrate) Down() (err error) {
	err = db.Def().DropTableIfExists(&model.UserInfo{}).Error
	return
}
