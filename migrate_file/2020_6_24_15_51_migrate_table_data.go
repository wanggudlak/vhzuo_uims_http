package migrate_file

import (
	"uims/command"
	"uims/command/commands/migrate_data"
)

type MigrateTableDataMigrate struct {
}

func (MigrateTableDataMigrate) Key() string {
	return "2020_6_24_15_51_migrate_table_data"
}

// migrate
func (MigrateTableDataMigrate) Up() (err error) {
	commands := [...]*command.Command{
		migrate_data.CMDInitClientData,           //初始化客户端数据
		migrate_data.CMDMigrateCassRegisterUser,  //迁移结算系统注册用户
		migrate_data.CMDMigrateCassRoles,         //迁移结算系统角色数据
		migrate_data.CMDMigrateCassGroup,         //迁移结算系统用户组数据
		migrate_data.CMDMigrateCassResource,      //迁移结算系统资源点数据
		migrate_data.CMDMigrateCassUserRole,      //迁移结算系统用户关联角色数据
		migrate_data.CMDMigrateCassUserGroup,     //迁移结算系统资源用户关联组数据
		migrate_data.CMDMigrateCassResourceGroup, //迁移结算系统资源组数据
		migrate_data.CMDInitRoleResMap,           //初始化角色资源组数据（在角色资源数据全部同步完成之后执行）
	}

	for _, cmd := range commands {
		command.CMD.Call(cmd, command.Args{})
	}
	return
}

// rollback
func (MigrateTableDataMigrate) Down() (err error) {
	return
}
