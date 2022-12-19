package socket5

import (
	"crypto/aes"
	"crypto/cipher"
	"fmt"
	"net"
	"shadowsocks-go/binary"
	"shadowsocks-go/utils"

	//"shadowsocks-go/ss2022"
	"testing"
)

const KeySize = 32
const NoneSize = 12

func TestRequest(t *testing.T) {
	key := []byte("cdBIDV42DCwnfIN")

	salt := utils.NewSalt(KeySize)

	fmt.Println("salt:", salt)

	subkey, err := NewSHA1(key, salt, utils.ToByte("ss-subkey"), KeySize)
	if err != nil {

	}
	fmt.Println(subkey)

	addr := NewSocket5DomainAddress([]byte("baidu.com"), 80)

	conn, err := net.Dial("tcp", "85.208.108.60:8118")
	if err != nil {
		fmt.Println("Dial fail")
	}

	block, err := aes.NewCipher(subkey)
	if err != nil {
		panic(err.Error())
	}

	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		panic(err.Error())
	}
	nonce := make([]byte, NoneSize)
	addrEn := aesgcm.Seal(nil, nonce, addr, nil)
	utils.IncrementNonce(nonce)

	data := "GET / HTTP/1.1\n\n"

	//aes.NewCipher(subkey)

	buf := make([]byte, 0)
	// add salt
	buf = append(buf, salt...)

	// add address
	buf = append(buf, addrEn...)

	// add chunk len
	ln := len(data)
	bl := binary.PutUint16(uint16(ln))
	blEn := aesgcm.Seal(nil, nonce, bl, nil)
	utils.IncrementNonce(nonce)
	buf = append(buf, blEn...)

	// add chunk data
	dataEn := aesgcm.Seal(nil, nonce, []byte(data), nil)
	utils.IncrementNonce(nonce)
	buf = append(buf, dataEn...)

	conn.Write(buf)
}
