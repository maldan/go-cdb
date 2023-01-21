package util

import (
	"strconv"
)

func StrToUInt(s string) uint {
	n, _ := strconv.Atoi(s)
	return uint(n)
}

func IntToStr(i uint32) string {
	return strconv.Itoa(int(i))
}
