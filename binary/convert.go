package binary

import (
	"bytes"
	"encoding/base64"
	"encoding/binary"
	"fmt"
	"log"
	"time"
)

// PutUint16 convert uint16 to big endian byte
func PutUint16(num uint16) []byte {
	b := make([]byte, 2)
	binary.BigEndian.PutUint16(b, num)
	return b
}

func GetUint16(buf []byte) uint16 {
	return binary.BigEndian.Uint16(buf)
}

// PutUint64 convert uint64 to big endian byte
func PutUint64(num uint64) []byte {
	b := make([]byte, 8)
	binary.BigEndian.PutUint64(b, num)
	return b
}

func GetUint64(buf []byte) uint64 {
	return binary.BigEndian.Uint64(buf)
}

func Base64Decoder(encode string) []byte {
	buf := make([]byte, 512)
	r := bytes.NewReader([]byte(encode))
	decoder := base64.NewDecoder(base64.StdEncoding, r)
	n, err := decoder.Read(buf)
	if err != nil {
		log.Println(err)
	}

	return buf[0:n]
}

func Timestamp() uint64 {
	now := time.Now()
	fmt.Println(now)

	tm := now.Unix()
	return uint64(tm)
}
