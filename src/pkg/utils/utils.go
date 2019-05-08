package utils

import (
	"crypto/md5"
	"crypto/rand"
	"encoding/json"
	"fmt"
	"log"
	"regexp"
	"strings"
	"unicode/utf8"
)

var isAscii = regexp.MustCompile(`^[a-zA-Z1-9 ]+$`).MatchString
var haveAscii = regexp.MustCompile(`[a-zA-Z1-9]+`).MatchString

func MD5Sum(content string) string {
	return fmt.Sprintf("%x", md5.Sum([]byte(content)))
}

func CopyAttrs(dst interface{}, src interface{}) {
	tmpByte, err := json.Marshal(src)
	if err != nil {
		return
	}
	json.Unmarshal(tmpByte, dst)
}

func PseudoUUID() (uuid string) {

	b := make([]byte, 16)
	_, err := rand.Read(b)
	if err != nil {
		log.Printf("PseudoUuid got error %s", err)
		return
	}

	uuid = fmt.Sprintf("%x-%x-%x-%x-%x", b[0:4], b[4:6], b[6:8], b[8:10], b[10:])

	return
}

func CountUnicodeString(str string) int {
	return utf8.RuneCountInString(str)
}

func CutUnicodeString(str string, stop int) string {
	if CountUnicodeString(str) <= stop {
		return str
	}
	temp := []rune(str)[:stop]
	return string(temp)
}

func IsAllAscii(raw string) bool {
	// tmp := strings.TrimSpace(raw)
	return isAscii(raw)
}

func IfHaveAscii(raw string) bool {
	tmp := strings.TrimSpace(raw)
	return haveAscii(tmp)
}
