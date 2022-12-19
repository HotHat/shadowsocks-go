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
	key := utils.Kdf("cdBIDV42DCwnfIN", KeySize)
	fmt.Println("psk len:", len(key))
	//key = append(key, []byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}...)
	fmt.Println("psk len:", len(key))
	fmt.Println("psk:", key)

	//salt := utils.NewSalt(KeySize)
	salt := []byte{233, 28, 168, 191, 50, 47, 245, 2, 19, 164, 179, 1, 44, 150, 183, 121, 52, 96, 69, 29, 221, 255, 149, 207, 235, 146, 141, 183, 32, 136, 81, 41}

	fmt.Println("salt:", salt)

	subkey, err := NewSHA1(key, salt, utils.ToByte("ss-subkey"), KeySize)
	if err != nil {

	}
	fmt.Println("subkey:", subkey)

	addr := NewSocket5DomainAddress([]byte("baidu.com"), 80)

	conn, err := net.Dial("tcp", "85.208.108.60:8118")
	if err != nil {
		fmt.Println("Dial fail")
	}

	fmt.Println("Shadowsocks server connected")

	block, err := aes.NewCipher(subkey)
	if err != nil {
		panic(err.Error())
	}
	fmt.Printf("subkey: % x\n", subkey)

	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		panic(err.Error())
	}

	nonce := make([]byte, NoneSize)

	addrEn := aesgcm.Seal(nil, nonce, addr, nil)
	utils.IncrementNonce(nonce)

	//data := "GET / HTTP/1.1\n\n"

	//aes.NewCipher(subkey)

	buf := make([]byte, 0)
	// add salt
	buf = append(buf, salt...)
	fmt.Printf("salt: % x\n", salt)

	// add address len
	aln := len(addrEn)
	abl := binary.PutUint16(uint16(aln))
	ablEn := aesgcm.Seal(nil, nonce, abl, nil)
	utils.IncrementNonce(nonce)
	buf = append(buf, ablEn...)
	fmt.Printf("address len: %d hex:% x\n", aln, abl)
	fmt.Printf("address len encode: % x\n", ablEn)

	// add address
	buf = append(buf, addrEn...)
	fmt.Printf("address : % x\n", addr)
	fmt.Printf("address encode: % x\n", addrEn)

	/*
		// add chunk len
		ln := len(data)
		bl := binary.PutUint16(uint16(ln))
		blEn := aesgcm.Seal(nil, nonce, bl, nil)
		utils.IncrementNonce(nonce)
		buf = append(buf, blEn...)
		fmt.Printf("chunk: %d hex:% x\n", ln, bl)
		fmt.Printf("chunk encode: % x\n", blEn)

		// add chunk data
		dataEn := aesgcm.Seal(nil, nonce, []byte(data), nil)
		utils.IncrementNonce(nonce)
		buf = append(buf, dataEn...)
		fmt.Printf("data : % x\n", data)
		fmt.Printf("data encode: % x\n", dataEn)

	*/

	fmt.Println("send buffer:", buf)
	fmt.Printf("send buffer: % x\n", buf)
	n, err := conn.Write(buf)
	if err != nil {
		fmt.Println("read fail:", err)
	}
	fmt.Println("write byte number:", n)

	rb := make([]byte, 4096)
	n, err = conn.Read(rb)
	if err != nil {
		fmt.Println("read fail:", err)
	}

	fmt.Println("read data:", rb[0:n])

}
