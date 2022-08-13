package Utils

import (
	"bytes"
	"encoding/binary"
	"strings"
)

func Int2Bytes(number int) []byte {
	writer := bytes.NewBuffer([]byte{})
	_ = binary.Write(writer, binary.LittleEndian, int32(number))
	return writer.Bytes()
}

func CheckSubstring(str string, subs ...string) (bool, int) {
	matches := 0
	isCompleteMatch := true
	for _, sub := range subs {
		if strings.Contains(str, sub) {
			matches += 1
		} else {
			isCompleteMatch = false
		}
	}
	return isCompleteMatch, matches
}
