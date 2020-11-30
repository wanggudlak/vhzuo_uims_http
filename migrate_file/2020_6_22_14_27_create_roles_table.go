package migrate_file

import (
	"fmt"
	"uims/internal/model"
	"uims/pkg/db"
)

type CreateRolesTableMigrate struct {
}

func (CreateRolesTableMigrate) Key() string {
	return "2020_6_22_14_27_create_roles_table"
}

// migrate
func (CreateRolesTableMigrate) Up() (err error) {
	if db.Def().HasTable(model.Role{}.TableName()) {
		err = fmt.Errorf("role table alreay exist")
		return
	}
	err = db.Def().
		Set("gorm:table_options", "CHARSET=utf8mb4,COMMENT='角色表'").
		CreateTable(&model.Role{}).Error
	return
}

// rollback
func (CreateRolesTableMigrate) Down() (err error) {
	err = db.Def().DropTableIfExists(&model.Role{}).Error
	return
}
