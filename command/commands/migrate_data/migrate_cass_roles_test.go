package migrate_data_test

import (
	"testing"
	"time"
	"uims/boot"
	"uims/command"
	"uims/command/commands/migrate_data"
	"uims/pkg/db"
)

func TestMain(m *testing.M) {
	boot.SetInTest()
	boot.Boot()
	m.Run()
}

type CassRoles struct {
	ID          int       `gorm:"column:id" json:"id"`
	Name        string    `gorm:"column:name" json:"name"`
	GuardName   string    `gorm:"column:guard_name" json:"guard_name"`
	Platform    int       `gorm:"column:platform" json:"platform"`
	Title       string    `gorm:"column:title" json:"title"`
	Description string    `gorm:"column:description" json:"description"`
	CreatedAt   time.Time `gorm:"column:updated_at" json:"created_at"`
	UpdatedAt   time.Time `gorm:"column:updated_at" json:"updated_at"`
	GroupID     int       `gorm:"column:groupID" json:"groupID"`
}

func TestDB(t *testing.T) {
	cassDBConn := db.Conn("cass")

	cassRoles := []CassRoles{}
	rows := cassDBConn.
		Table("vz_roles").
		Model(CassRoles{})
	rows.Scan(&cassRoles)

	t.Logf("查到了 %d 条数据", len(cassRoles))

	rows2 := cassDBConn.
		Table("vz_roles").
		Model(CassRoles{})
	rows2.Scan(&cassRoles)
	t.Logf("查到了 %d 条数据", len(cassRoles))
}

func TestRoleCommandCall(t *testing.T) {
	command.CMD.Call(migrate_data.CMDMigrateCassRoles, command.Args{})
	command.CMD.Call(migrate_data.CMDMigrateCassRoles, command.Args{})
}
