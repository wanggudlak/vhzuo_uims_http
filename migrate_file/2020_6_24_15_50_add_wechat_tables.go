package migrate_file

import (
	"uims/internal/model"
	"uims/pkg/db"
)

type AddWeChatTableMigrate struct {
}

func (AddWeChatTableMigrate) Key() string {
	return "2020_6_24_15_50_add_wechat_tables"
}

// migrate
func (AddWeChatTableMigrate) Up() error {
	return db.Def().
		Set("gorm:table_options", "CHARSET=utf8mb4").
		CreateTable(&model.UserWeChat{}, &model.WeChat{}, &model.ClientWeChat{}).Error
}

// rollback
func (AddWeChatTableMigrate) Down() error {
	return db.Def().DropTableIfExists(&model.UserWeChat{}, &model.WeChat{}, &model.ClientWeChat{}).Error
}
