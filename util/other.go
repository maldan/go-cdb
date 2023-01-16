package util

import "strconv"

func StrToUInt(s string) uint {
	n, _ := strconv.Atoi(s)
	return uint(n)
}
