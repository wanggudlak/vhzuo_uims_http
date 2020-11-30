package migrate_file

import (
	"fmt"
	"uims/internal/model"
	"uims/pkg/db"
)

type CreateUserRolesTableMigrate struct {
}

func (CreateUserRolesTableMigrate) Key() string {
	return "2020_6_22_15_43_create_user_role_table"
}

// migrate
func (CreateUserRolesTableMigrate) Up() (err error) {
	if db.Def().HasTable(model.UserRole{}.TableName()) {
		err = fmt.Errorf("user_role table alreay exist")
		return
	}
	err = db.Def().
		Set("gorm:table_options", "CHARSET=utf8mb4,COMMENT='用户与角色关系表'").
		CreateTable(&model.UserRole{}).Error
	return
}

// rollback
func (CreateUserRolesTableMigrate) Down() (err error) {
	err = db.Def().DropTableIfExists(&model.UserRole{}).Error
	return
}
