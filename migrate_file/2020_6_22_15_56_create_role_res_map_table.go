package migrate_file

import (
	"fmt"
	"uims/internal/model"
	"uims/pkg/db"
)

type CreateRoleResMapTableMigrate struct {
}

func (CreateRoleResMapTableMigrate) Key() string {
	return "2020_6_22_15_56_create_role_res_map_table"
}

// migrate
func (CreateRoleResMapTableMigrate) Up() (err error) {
	if db.Def().HasTable(model.RoleResMap{}.TableName()) {
		err = fmt.Errorf("role_res table alreay exist")
		return
	}
	err = db.Def().
		Set("gorm:table_options", "CHARSET=utf8mb4,COMMENT='角色与资源组关系'").
		CreateTable(&model.RoleResMap{}).Error
	return
}

// rollback
func (CreateRoleResMapTableMigrate) Down() (err error) {
	err = db.Def().DropTableIfExists(&model.RoleResMap{}).Error
	return
}
