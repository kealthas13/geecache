package http

import (
	"hash/crc32"
)

func String_hash(s string) int {//根据hash(key)决定1-3
	v := int(crc32.ChecksumIEEE([]byte(s)))//hash
	if -v >= 0 {//如果为负数
		v = -v
	}
	v = (v % 3)+1
	return v
}
