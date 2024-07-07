package util

import (
	"bytes"
	"crypto/md5"
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"github.com/google/uuid"
	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/transform"
	"hash"
	"io"
	"unicode"
)

const TimeFormat = "2006-01-02 15:04:05"

// MD5 md5
func MD5(str []byte) string {
	h := md5.New()
	h.Write(str)
	return hex.EncodeToString(h.Sum(nil))
}

// Sha1 sha1
func Sha1(str []byte) string {
	h := sha1.New()
	h.Write(str)
	return hex.EncodeToString(h.Sum(nil))
}

// FileHash 计算文件hash
func FileHash(reader io.Reader, tp string) string {
	var result []byte
	var h hash.Hash
	if tp == "md5" {
		h = md5.New()
	} else {
		h = sha1.New()
	}
	if _, err := io.Copy(h, reader); err != nil {
		return ""
	}
	return fmt.Sprintf("%x", h.Sum(result))
}

// GetUuid 生成uuid
func GetUuid() string {
	var u uuid.UUID
	var err error
	for i := 0; i < 3; i++ {
		u, err = uuid.NewUUID()
		if err == nil {
			return u.String()
		}
	}
	return ""
}

// GetUuidV4 生成uuid v4
func GetUuidV4() string {
	var u uuid.UUID
	var err error
	for i := 0; i < 3; i++ {
		u, err = uuid.NewRandom()
		if err == nil {
			return u.String()
		}
	}
	return ""
}
func isRomanChar(r rune) bool {
	switch unicode.ToUpper(r) {
	case 'I', 'V', 'X', ' ':
		return true
	}
	return false
}

func RemoveRomanNumber(s string) string {
	runes := []rune(s)
	for i := len(runes) - 1; i >= 0; i-- {
		if !isRomanChar(runes[i]) {
			return string(runes[:i+1])
		}
	}
	return ""
}

func IsGBK(data []byte) bool {
	reader := transform.NewReader(bytes.NewReader(data), simplifiedchinese.GB18030.NewDecoder())
	_, err := io.ReadAll(reader)
	return err == nil
}

func IsUTF8(data []byte) bool {
	return bytes.Equal([]byte(data), data)
}
