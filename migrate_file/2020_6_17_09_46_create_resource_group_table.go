package migrate_file

import (
	"fmt"
	"uims/internal/model"
	"uims/pkg/db"
)

type CreateResourceGroupTableMigrate struct {
}

func (CreateResourceGroupTableMigrate) Key() string {
	return "2020_6_17_09_46_create_resource_group_table"
}

// migrate
func (CreateResourceGroupTableMigrate) Up() (err error) {
	if db.Def().HasTable(model.ResourceGroup{}.TableName()) {
		err = fmt.Errorf("users table alreay exist")
		return
	}
	err = db.Def().
		Set("gorm:table_options", "CHARSET=utf8mb4,COMMENT='资源策略组'").
		CreateTable(&model.ResourceGroup{}).Error
	return
}

// rollback
func (CreateResourceGroupTableMigrate) Down() (err error) {
	err = db.Def().DropTableIfExists(&model.ResourceGroup{}).Error
	return
}
