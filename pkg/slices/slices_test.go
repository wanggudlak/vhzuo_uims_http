package slices_test

import (
	"fmt"
	"testing"
	"uims/pkg/slices"
	"uims/pkg/tool"
)

func TestIsSlice(t *testing.T) {
	intSlice := []int{1, 2, 3, 4, 5, 6, 7, 8}
	//maps := map[string]string{"key1":"val1","key2":"val2","key3":"val3"}

	res, ok := slices.IsSlice(intSlice)
	fmt.Println(11, res, ok)
}

func TestRemoveSlice(t *testing.T) {

	//intSlice := []int{1,2,3,4,5,6,7,8}
	//strSlice := []string{"a","b","c","d"}
	boolSlice := []bool{true, true, false, true}

	b, ok := slices.CreateAnyTypeSlice(boolSlice)

	if ok {
		c := slices.RemoveSlice(b, false)

		fmt.Println(3333, c)
	}
}

func TestIsExistValue(t *testing.T) {
	//a := []int{1, 2, 3, 4, 5, 6, 7, 8, 9}
	//a := []string{"a","b","c","d"}
	var b []string
	isExist, _ := slices.IsExistValue("b", b)

	//b := map[string]string{"key1": "val1", "key2": "val2", "key3": "val4"}
	//isExist, _ := slices.IsExistValue("key2", b)

	fmt.Println(111, isExist)
	tool.Dump(isExist)
}
