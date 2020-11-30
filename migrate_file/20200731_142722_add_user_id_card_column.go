package migrate_file

import (
	"database/sql"
	"fmt"
	"github.com/pkg/errors"
	"uims/internal/model"
	"uims/pkg/db"
)

type AddUserIdCardColumn struct {
}

func (AddUserIdCardColumn) Key() string {
	return "20200731_142722_add_user_id_card_column.go"
}

func (AddUserIdCardColumn) Up() (err error) {
	db.Def().AutoMigrate(&model.User{}, &model.UserInfo{})
	//修改身份证号默认值，初始化数据，添加索引
	err = db.Def().Exec("ALTER TABLE `uims_user_info` MODIFY COLUMN `identity_card_no` varchar(20) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT NULL COMMENT '用户的身份证号' AFTER `name_full_py`").Error
	if err != nil {
		fmt.Println(err)
	}
	err = db.Def().Exec("UPDATE `uims_user_info` SET `identity_card_no` = ? WHERE `identity_card_no` = ?", sql.NullString{}, "").Error
	if err != nil {
		fmt.Println(err)
	}
	err = db.Def().Model(&model.UserInfo{}).AddUniqueIndex("unx_identity_card_no", "identity_card_no").Error
	if err != nil {
		fmt.Println(err)
	}
	// 迁移 user_info 数据到 user_auth
	var userInfos []model.UserInfo
	err = db.Def().Where("length(identity_card_no) > ?", 15).Find(&userInfos).Error
	if err != nil {
		err = errors.Wrap(err, "查询 user_info 失败")
		return
	}
	for _, userInfo := range userInfos {
		db.Def().Model(&model.User{}).Where("id = ?", userInfo.UserID).Updates(&model.User{IdentityCardNo: userInfo.IdentityCardNo})
	}
	return
}

func (AddUserIdCardColumn) Down() (err error) {
	return
}
