package migrate_file

import "uims/pkg/db"

type AlterUserAuthPhoneColumnNotNull struct {
}

func (AlterUserAuthPhoneColumnNotNull) Key() string {
	return "20200804_184815_alter_user_auth_phone_column_not_null.go"
}

func (AlterUserAuthPhoneColumnNotNull) Up() (err error) {
	return db.Def().Exec("alter table uims_user_auth modify column phone char(11) default null comment '手机号'").Error
}

func (AlterUserAuthPhoneColumnNotNull) Down() (err error) {
	return
}
