package service_test

import (
	"fmt"
	"testing"
	"uims/internal/service"
	"uims/pkg/db"
)

func TestGetResGroupByClientId(t *testing.T) {

	resGroup, err := service.GetResGroupService().GetResGroupByClientId(1)

	if err != nil {
		t.Fatal(err)
	}

	fmt.Println("resGroup=", resGroup)
	fmt.Println("resGroup=", resGroup[0].ResOfCurr.ResourceIDs)
}

func TestDelResNeedUpdateResGroupByResOfCurr(t *testing.T) {

	tx := db.Def()

	resGroup, err := service.GetResGroupService().DelResNeedUpdateResGroupByResOfCurr(1, 28, tx)

	if err != nil {
		t.Fatal(err)
	}

	fmt.Println("resGroup=", resGroup)
}

func TestGetRoleWithResourceMap(t *testing.T) {
	tx := db.Def()

	// delete
	delResourceIDs, addResourceIDs, err := service.GetResGroupService().FilterResourceIDs("delete", 8, 12, nil, []int{31, 32}, tx)

	// update
	//delResourceIDs, addResourceIDs, err := service.GetResGroupService().
	//	FilterResourceIDs("update",8, 12, []int{31, 32}, []int{31, 33}, tx)

	fmt.Println("ResourceIDs=", delResourceIDs, addResourceIDs, err)
}
