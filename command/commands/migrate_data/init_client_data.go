package migrate_data

import (
	"fmt"
	"uims/command"
	"uims/command/commands/version"
	"uims/internal/model"
	"uims/pkg/db"
)

var CMDInitClientData = &command.Command{
	UsageLine: "init:client",
	Short:     "初始化客户端数据",
	Long:      `初始化uims客户端数据`,
	PreRun:    func(cmd *command.Command, args []string) { version.ShowShortVersionBanner() },
	Run:       initClientData,
}

func init() {
	command.CMD.Register(CMDInitClientData)
}

func initClientData(_ *command.Command, _ []string) int {
	needInitClientData := []model.Client{
		{
			AppId:          "YSyn5CeEqsVfEUqP",
			ClientSpm2Code: "YSyn5CeEqsVfEUqP",
			ClientType:     "CASS",
			ClientFlagCode: "VDK_CASS_FRONT",
			ClientSpm1Code: "1024",
			ClientName:     "结算系统前台",
			Status:         "Y",
			ClientHostIp:   "null",
			ClientHostUrl:  "null",
		},
		{
			AppId:          "DFASDF234FDAS231",
			ClientSpm2Code: "DFASDF234FDAS231",
			ClientType:     "CASS",
			ClientFlagCode: "VDK_CASS_BACK",
			ClientSpm1Code: "1024",
			ClientName:     "结算系统后台",
			Status:         "Y",
			ClientHostIp:   "null",
			ClientHostUrl:  "null",
		}, {
			AppId:          "ZBCDASDFGASDFASF",
			ClientSpm2Code: "ZBCDASDFGASDFASF",
			ClientType:     "VDK",
			ClientFlagCode: "VDK_MP",
			ClientSpm1Code: "1024",
			ClientName:     "任务系统前台",
			Status:         "Y",
			ClientHostIp:   "null",
			ClientHostUrl:  "null",
		},
		{
			AppId:          "WDYKUIYFJWPLJ6GS",
			ClientSpm2Code: "WDYKUIYFJWPLJ6GS",
			ClientType:     "VDK",
			ClientFlagCode: "VDK_ES_SAPP",
			ClientSpm1Code: "1024",
			ClientName:     "任务系统小程序",
			Status:         "Y",
			ClientHostIp:   "null",
			ClientHostUrl:  "null",
		},
		{
			AppId:          "H5TASKBCDASDFGAS",
			ClientSpm2Code: "H5TASKBCDASDFGAS",
			ClientType:     "VDK",
			ClientFlagCode: "VDK_H5",
			ClientSpm1Code: "1024",
			ClientName:     "任务系统H5",
			Status:         "Y",
			ClientHostIp:   "null",
			ClientHostUrl:  "null",
		},
	}

	for _, clientItem := range needInitClientData {
		var client model.Client
		db.Def().Table("uims_client").
			Where("client_type = ? AND client_flag_code = ? AND client_name = ?", clientItem.ClientType, clientItem.ClientFlagCode, clientItem.ClientName).
			First(&client)
		if client.ID != 0 {
			continue
		}
		err := db.Def().Create(&clientItem).Error
		if err != nil {
			fmt.Println("保存客户端基础数据失败", err)
			return 0
		}
	}

	err := db.Def().Exec(`
INSERT INTO uims_client_settings 
(id, client_id, type, bus_channel_id, page_id, spm_full_code, form_fields, page_template_file, isdel, created_at, updated_at)
VALUES
(1,1,'LGN','100','101','1024.YSyn5CeEqsVfEUqP.100.101',NULL,'{"a": "/resource/1024.YSyn5CeEqsVfEUqP.100.101/20200630163952/html_template/index.html"}','N','2020-06-11 00:00:00.000000','2020-06-11 00:00:00.000000'),
(2,2,'LGN','100','101','1024.DFASDF234FDAS231.100.101',NULL,'{"a": "/resource/1024.DFASDF234FDAS231.100.101/20200630163952/html_template/index.html"}','N','2020-06-11 00:00:00.000000','2020-06-11 00:00:00.000000'),
(3,3,'LGN','100','101','1024.ZBCDASDFGASDFASF.100.101',NULL,'{"a": "/resource/1024.ZBCDASDFGASDFASF.100.101/20200630174245/html_template/index.html"}','N','2020-06-11 00:00:00.000000','2020-06-11 00:00:00.000000'),
(4,4,'LGN','100','101','1024.WDYKUIYFJWPLJ6GS.100.101',NULL,'{"a": ""}','N','2020-06-11 00:00:00.000000','2020-06-11 00:00:00.000000'),
(5,5,'LGN','100','101','1024.H5TASKBCDASDFGAS.100.101',NULL,'{"a": "/resource/1024.H5TASKBCDASDFGAS.100.101/20200716161101/html_template/index.html"}','N','2020-06-11 00:00:00.000000','2020-06-11 00:00:00.000000');
`).Error
	if err != nil {
		fmt.Printf("保存 client_setting 数据失败: [%+v] \n", err)
		return 0
	}

	fmt.Println("初始化客户端数据完成")
	return 0
}
