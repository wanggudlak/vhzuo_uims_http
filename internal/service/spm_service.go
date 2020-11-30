package service

import "strings"

const (
	SPLIT_SEPARATOR = "."
)

type SPMcode struct {
	FullCode string
	Code1    string `comment:"client_spm1_code"`
	Code2    string `comment:"client_spm2_code"`
	Code3    string `comment:"频道ID"`
	Code4    string `comment:"页面ID"`
}

// ParseSPMstring 解析SPM字符串全编码到SPMcode对象
func ParseSPMstring(spmstr string) *SPMcode {
	if len(spmstr) == 0 {
		return nil
	}

	var pSpmCode = &SPMcode{
		FullCode: spmstr,
	}
	splitRes := strings.Split(spmstr, SPLIT_SEPARATOR)

	for i, code := range splitRes {
		switch i {
		case 0:
			pSpmCode.Code1 = code
		case 1:
			pSpmCode.Code2 = code
		case 2:
			pSpmCode.Code3 = code
		case 3:
			pSpmCode.Code4 = code
		}
	}
	return pSpmCode
}
