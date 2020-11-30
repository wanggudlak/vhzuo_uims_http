package uuid_test

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
	"uims/internal/service/uuid"
	"uims/pkg/tool"
)

func TestGenerateForCASS(t *testing.T) {
	var x, y uuid.ID
	for i := 0; i < 10000; i++ {
		y = uuid.GenerateForCASS()
		fmt.Println(y)
		if x == y {
			t.Errorf("x(%d) & y(%d) are the same", x, y)
		}
		x = y
	}
}

func TestGenerateForUIMD(t *testing.T) {
	var x, y uuid.ID
	for i := 0; i < 10000; i++ {
		y = uuid.GenerateForUIMS()
		fmt.Println(y)
		if x == y {
			t.Errorf("x(%d) & y(%d) are the same", x, y)
		}
		x = y
	}
}

func TestGenerateStringUUIDForUIMS(t *testing.T) {
	id := uuid.GenerateForUIMS().String()
	tool.Dump(id)
	assert.IsType(t, "", id)
}
