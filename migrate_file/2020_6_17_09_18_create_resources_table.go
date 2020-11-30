package migrate_file

import (
	"fmt"
	"uims/internal/model"
	"uims/pkg/db"
)

type CreateResourcesTableMigrate struct {
}

func (CreateResourcesTableMigrate) Key() string {
	return "2020_6_17_09_18_create_resources_table"
}

// migrate
func (CreateResourcesTableMigrate) Up() (err error) {
	if db.Def().HasTable(model.Resource{}.TableName()) {
		err = fmt.Errorf("users table alreay exist")
		return
	}
	err = db.Def().
		Set("gorm:table_options", "CHARSET=utf8mb4,COMMENT='权限资源表'").
		CreateTable(&model.Resource{}).Error
	return
}

// rollback
func (CreateResourcesTableMigrate) Down() (err error) {
	err = db.Def().DropTableIfExists(&model.Resource{}).Error
	return
}
