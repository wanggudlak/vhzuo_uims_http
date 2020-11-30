package service

import "strings"

func SplitURI(fullURI, sep string) []string {
	return strings.Split(fullURI, sep)
}
