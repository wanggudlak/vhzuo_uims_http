package service_test

import (
	"fmt"
	"testing"
	"uims/internal/service"
	"uims/pkg/db"
)

func TestGetResourceMapByIDs(t *testing.T) {

	tx := db.Def()
	ids := []int{30}

	resources, err := service.GetResourceService().GetResourceMapByIDs(ids, tx)

	if err != nil {
		t.Fatal(err)
	}

	fmt.Println("resources=", resources)

}
