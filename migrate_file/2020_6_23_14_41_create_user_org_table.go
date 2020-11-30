package migrate_file

import (
	"fmt"
	"uims/internal/model"
	"uims/pkg/db"
)

type CreateUserOrgTableMigrate struct {
}

func (CreateUserOrgTableMigrate) Key() string {
	return "2020_6_23_14_41_create_user_org_table"
}

// migrate
func (CreateUserOrgTableMigrate) Up() (err error) {
	if db.Def().HasTable(model.UserOrg{}.TableName()) {
		err = fmt.Errorf("organization table alreay exist")
		return
	}
	err = db.Def().
		Set("gorm:table_options", "CHARSET=utf8mb4,COMMENT='用户与组织关系表'").
		CreateTable(&model.UserOrg{}).Error
	return
}

// rollback
func (CreateUserOrgTableMigrate) Down() (err error) {
	err = db.Def().DropTableIfExists(&model.UserOrg{}).Error
	return
}
