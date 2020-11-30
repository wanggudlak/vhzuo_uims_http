package migrate_file

import (
	"fmt"
	"uims/internal/model"
	"uims/pkg/db"
)

type CreateOrganizationTableMigrate struct {
}

func (CreateOrganizationTableMigrate) Key() string {
	return "2020_6_23_14_19_create_organization_table"
}

// migrate
func (CreateOrganizationTableMigrate) Up() (err error) {
	if db.Def().HasTable(model.Org{}.TableName()) {
		err = fmt.Errorf("organization table alreay exist")
		return
	}
	err = db.Def().
		Set("gorm:table_options", "CHARSET=utf8mb4,COMMENT='组织表'").
		CreateTable(&model.Org{}).Error
	return
}

// rollback
func (CreateOrganizationTableMigrate) Down() (err error) {
	err = db.Def().DropTableIfExists(&model.Org{}).Error
	return
}
