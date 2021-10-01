package discord

import "hash/crc32"

func Hashcode(s string) int {
	v := int(crc32.ChecksumIEEE([]byte(s)))
	if v >= 0 {
		return v
	}
	if -v >= 0 {
		return -v
	}
	// v == MinInt
	return 0
}

func contains(arr [3]string, str string) bool {
	for _, a := range arr {
		if a == str {
			return true
		}
	}
	return false
}

func fixPermission(v interface{}) uint64 {
	switch val := v.(type) {
	case int:
		return uint64(val)
	case uint64:
		return val
	default:
		return 0
	}
}
