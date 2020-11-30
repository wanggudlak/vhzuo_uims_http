package service_test

import (
	"fmt"
	"testing"
	"uims/internal/service"
	"uims/pkg/db"
)

func TestGetRolesMapByGroupId(t *testing.T) {

	tx := db.Def()
	roles, err := service.GetRoleResMapService().GetRolesMapByGroupId(8, 1, tx)

	if err != nil {
		t.Fatal(err)
	}

	fmt.Println("roles=", roles)
}

func TestDelResGroupNeedUpdateRoleResMapByGroupID(t *testing.T) {
	tx := db.Def()
	err := service.GetRoleResMapService().DelResGroupNeedUpdateRoleResMapByGroupID(uint(1), uint(12), tx)
	if err != nil {
		t.Fatal(err)
	}

}
