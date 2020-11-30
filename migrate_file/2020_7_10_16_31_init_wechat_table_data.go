package migrate_file

import (
	"github.com/jinzhu/gorm"
	"uims/internal/model"
	"uims/pkg/db"
)

type InitWeChatTableData struct {
}

func (InitWeChatTableData) Key() string {
	return "2020_7_10_16_31_init_wechat_table_data"
}

// migrate
func (InitWeChatTableData) Up() error {
	return db.Def().Transaction(func(tx *gorm.DB) error {
		var err error
		w1 := &model.WeChat{
			AppId:  "wxaa24bae3b89b7755",
			UUID:   "463c1a8bdc6c4f07a94290b524cd559c",
			Secret: "ccdd63e7b9dc18c0c09b17018631658c",
			Desc:   "代付系统微信扫码登录使用",
		}

		w2 := &model.WeChat{
			AppId:  "wxc7d9b02776155363",
			UUID:   "6a6b0e837a4f424191022f4fbe815514",
			Secret: "8e1e2cd9fa41895ae15433cfef1f40f4",
			Desc:   "微桌接活小程序",
		}

		w3 := &model.WeChat{
			AppId:  "wx3196cbb5d525f1de",
			UUID:   "c1d0940c0e56456b8b8f8fc8bd160632",
			Secret: "b1aec372d0cca3acc04a6250073fe830",
			Desc:   "微桌科技服务平台公众号(扫码登录用)",
		}

		err = db.Def().Create(&w1).Error
		if err != nil {
			return err
		}
		err = db.Def().Create(&w2).Error
		if err != nil {
			return err
		}
		err = db.Def().Create(&w3).Error
		if err != nil {
			return err
		}
		// w1 add wechat client bind
		var clients = []*model.Client{}
		err = db.Def().Where("client_flag_code in (?)", []string{"VDK_CASS_FRONT", "VDK_CASS_BACK"}).Find(&clients).Error
		if err != nil {
			return err
		}
		for _, client := range clients {
			cw := model.ClientWeChat{
				ClientId: client.ID,
				WeChatId: w1.ID,
			}
			err = db.Def().Create(&cw).Error
			if err != nil {
				return err
			}
		}

		// w2 小程序绑定到任务系统小程序中
		client := model.Client{}
		err = db.Def().Where("client_flag_code = ?", "VDK_ES_SAPP").First(&client).Error
		if err != nil {
			return err
		}
		cw := model.ClientWeChat{
			ClientId: client.ID,
			WeChatId: w2.ID,
		}
		err = db.Def().Create(&cw).Error
		if err != nil {
			return err
		}

		// w3 任务系统扫码登录绑定到任务系统前台
		clientTaskFront := model.Client{}
		err = db.Def().Where("client_flag_code = ?", "VDK_MP").First(&clientTaskFront).Error
		if err != nil {
			return err
		}
		cw = model.ClientWeChat{
			ClientId: clientTaskFront.ID,
			WeChatId: w3.ID,
		}
		err = db.Def().Create(&cw).Error
		if err != nil {
			return err
		}

		return nil
	})
}

// rollback
func (InitWeChatTableData) Down() error {
	return db.Def().DropTableIfExists(&model.UserWeChat{}, &model.WeChat{}, &model.ClientWeChat{}).Error
}
