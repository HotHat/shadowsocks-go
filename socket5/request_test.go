package socket5

import (
	"crypto/aes"
	"crypto/cipher"
	"fmt"
	"shadowsocks-go/binary"
	"shadowsocks-go/utils"

	//"shadowsocks-go/ss2022"
	"testing"
)

const (
	KeySize  = 32
	NoneSize = 12
	TagSize  = 16
)

func TestRequest(t *testing.T) {
	//key := utils.EvpBytesToKey2("cdBIDV42DCwnfIN", KeySize)
	key := utils.EvpBytesToKey("cdBIDV42DCwnfIN", KeySize)
	fmt.Println("psk len:", len(key))
	fmt.Println("psk:", key)
	//key = append(key, []byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}...)
	//fmt.Println("psk len:", len(key2))
	//fmt.Println("psk2:", key2)

	//salt := utils.NewSalt(KeySize)
	salt := []byte{233, 28, 168, 191, 50, 47, 245, 2, 19, 164, 179, 1, 44, 150, 183, 121, 52, 96, 69, 29, 221, 255, 149, 207, 235, 146, 141, 183, 32, 136, 81, 41}

	fmt.Println("salt:", salt)

	subkey := make([]byte, KeySize)
	utils.KdfSHA1(key, salt, utils.ToByte("ss-subkey"), subkey)

	fmt.Println("subkey:", subkey)

	addr := NewSocket5DomainAddress([]byte("baidu.com"), 80)

	//conn, err := net.Dial("tcp", "85.208.108.60:8118")
	//if err != nil {
	//	fmt.Println("Dial fail")
	//}

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

	//data := "GET / HTTP/1.1\n\n"

	//aes.NewCipher(subkey)

	buf := make([]byte, 0)
	// add salt
	buf = append(buf, salt...)
	fmt.Printf("salt: % x\n", salt)

	// add address len, tag size 16
	aln := len(addr) + TagSize
	abl := binary.PutUint16(uint16(aln))
	fmt.Printf("nonce: % x\n", nonce)
	ablEn := aesgcm.Seal(nil, nonce, abl, nil)
	utils.IncrementNonce(nonce)

	buf = append(buf, ablEn...)
	fmt.Printf("address len: %d hex:% x\n", aln, abl)
	fmt.Printf("address len encode: % x\n", ablEn)

	// add address
	fmt.Printf("nonce: % x\n", nonce)
	addrEn := aesgcm.Seal(nil, nonce, addr, nil)
	utils.IncrementNonce(nonce)

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
	/*
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

	*/

}

func TestDecode(t *testing.T) {
	key := utils.EvpBytesToKey("cdBIDV42DCwnfIN", KeySize)

	buf := []byte{233, 28, 168, 191, 50, 47, 245, 2, 19, 164, 179, 1, 44, 150, 183, 121, 52, 96, 69, 29, 221, 255, 149, 207, 235, 146, 141, 183, 32, 136, 81, 41, 207, 215, 25, 124, 176, 59, 89, 88, 81, 40, 140, 86, 239, 191, 229, 161, 15, 127, 210, 12, 49, 135, 196, 141, 70, 184, 225, 84, 231, 60, 218, 254, 45, 151, 165, 16, 53, 224, 16, 52, 249, 100, 220, 169, 240, 222, 202}

	salt := buf[0:KeySize]

	subkey := make([]byte, KeySize)
	utils.KdfSHA1(key, salt, utils.ToByte("ss-subkey"), subkey)

	fmt.Println("subkey:", subkey)

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
	ln := buf[KeySize : KeySize+2+TagSize]

	dln, err := aesgcm.Open(nil, nonce, ln, nil)
	utils.IncrementNonce(nonce)
	if err != nil {
		panic("len open fail")
	}
	dn := binary.GetUint16(dln)
	fmt.Println("len:", dn)

	addrS := buf[KeySize+2+TagSize : KeySize+2+TagSize+dn]
	daddrs, err := aesgcm.Open(nil, nonce, addrS, nil)
	utils.IncrementNonce(nonce)
	if err != nil {
		panic("addr open fail")
	}

	addr, err := ParseSocket5Address(daddrs)
	if err != nil {
		panic("len open fail")
	}
	fmt.Println("Addr:", addr)

}
