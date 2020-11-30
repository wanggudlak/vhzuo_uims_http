package migrate

import (
	"uims/migrate_file"
	"uims/pkg/db"
)

type Migration struct {
	Id        uint   `json:"id" gorm:"primary_key"`
	Migration string `json:"migration" gorm:"type:varchar(255)"`
	Batch     uint   `json:"batch"`
}

func (Migration) TableName() string {
	return "migrations"
}

type MigrateFile interface {
	Key() string
	Up() error
	Down() error
}

// 定义的迁移文件需要在这里注册
var MigrateFiles = []MigrateFile{
	migrate_file.CreateUsersTableMigrate{},
	migrate_file.CreateClientTableMigrate{},
	migrate_file.CreateResourcesTableMigrate{},
	migrate_file.CreateResourceGroupTableMigrate{},
	migrate_file.CreateRolesTableMigrate{},
	migrate_file.CreateUserRolesTableMigrate{},
	migrate_file.CreateRoleResMapTableMigrate{},
	migrate_file.CreateOrganizationTableMigrate{},
	migrate_file.CreateUserOrgTableMigrate{},
	migrate_file.CreateClientSettingsTableMigrate{},
	migrate_file.CreateUserInfoTableMigrate{},
	migrate_file.CreateSuperAdminUserMigrate{},
	migrate_file.AddWeChatTableMigrate{},
	migrate_file.MigrateTableDataMigrate{},
	migrate_file.InitWeChatTableData{},
	migrate_file.CreateUserInfoCursorTableMigrate{},
	migrate_file.UpdateUserAuthIdCard{},
	migrate_file.AddUserIdCardColumn{},
	migrate_file.AlterUserAuthPhoneColumnNotNull{},
}

// 获取需要迁移的 migrateFiles
// files 有, migrate_file 里没有的数据
func GetNeedMigrateFiles(migrateFiles []MigrateFile) []MigrateFile {
	var ans []MigrateFile
	var ms = GetAllMigrations()
	diff := map[string]string{}
	for _, migrateFile := range migrateFiles {
		diff[migrateFile.Key()] = ""
	}
	for _, migration := range ms {
		delete(diff, migration.Migration)
	}
	for _, migrateFile := range migrateFiles {
		if _, ok := diff[migrateFile.Key()]; ok {
			ans = append(ans, migrateFile)
		}
	}
	return ans
}

// 获取需要回滚的 migrateFiles
func GetNeedRollbackKeys(step int) []MigrateFile {
	var ans []MigrateFile
	var ms = GetAllMigrations()
	var keyMigrateFile = map[string]MigrateFile{}
	if step < 1 {
		return ans
	}
	for _, migrateFile := range MigrateFiles {
		keyMigrateFile[migrateFile.Key()] = migrateFile
	}
	cur := 0
	for _, migrate := range ms {
		if step < 1 {
			break
		}
		if m, ok := keyMigrateFile[migrate.Migration]; ok {
			ans = append(ans, m)
		}
		if int(migrate.Batch) != cur {
			step--
		}
	}
	return ans
}

// 获取所有迁移记录
func GetAllMigrations() []Migration {
	var migrations []Migration
	db.Def().Order("id desc").Find(&migrations)
	return migrations
}

// 获取下一个迁移版本号
func GetNextBatchNo() uint {
	m := Migration{}
	batch := uint(0)
	db.Def().Order("batch desc").Select("batch").First(&m)
	batch = m.Batch + 1
	return batch
}

func CreateMigrate(migration string, batch uint) (err error) {
	m := Migration{
		Migration: migration,
		Batch:     batch,
	}
	err = db.Def().Create(&m).Error
	return
}

func DeleteMigrate(migration string) (err error) {
	m := Migration{
		Migration: migration,
	}
	err = db.Def().Where(&m).Delete(Migration{}).Error
	return
}
