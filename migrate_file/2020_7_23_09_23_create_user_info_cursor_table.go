package migrate_file

import (
	"fmt"
	"uims/internal/model"
	"uims/pkg/db"
)

type CreateUserInfoCursorTableMigrate struct {
}

func (CreateUserInfoCursorTableMigrate) Key() string {
	return "2020_7_23_09_23_create_user_info_cursor_table"
}

// migrate
func (CreateUserInfoCursorTableMigrate) Up() (err error) {
	if db.Def().HasTable(model.UserInfoCursor{}.TableName()) {
		err = fmt.Errorf("user_info_cursor table alreay exist")
		return
	}
	err = db.Def().
		Set("gorm:table_options", "CHARSET=utf8mb4,COMMENT='用户信息临时表'").
		CreateTable(&model.UserInfoCursor{}).Error
	return
}

// rollback
func (CreateUserInfoCursorTableMigrate) Down() (err error) {
	err = db.Def().DropTableIfExists(&model.UserInfoCursor{}).Error
	return
}
