package ss2022

import (
	"fmt"
	"shadowsocks-go/binary"
	"shadowsocks-go/slices"
	"shadowsocks-go/socket5"
	"testing"
)

func TestRequestStream(t *testing.T) {
	protocol := "2022-blake3-aes-128-gcm"
	psk := "JIloOlaO1V506UnRV521mg=="
	salt := NewSalt(CryptoGCMBlockLength(protocol))

	requestStream := NewStream(binary.Base64Decoder(psk))
	requestStream.AddSalt(salt)

	variableHeader := NewRequestVariableHeader(socket5.Socket5AddressDomain, []byte("baidu.com"), 80, 10, nil)

	b := NewRequestFixedHeader(0, binary.Timestamp(), uint16(len(variableHeader)+16))
	fmt.Println(b)

	eb := requestStream.Encryption(b)

	vb := requestStream.Encryption(variableHeader)

	nb := slices.ConcatMultipleSlices([][]byte{salt, eb, vb})

	fmt.Println("salt:", salt)
	fmt.Println(nb)
}

func TestParseRequestStream(t *testing.T) {
	protocol := "2022-blake3-aes-128-gcm"
	psk := "JIloOlaO1V506UnRV521mg=="
	salt := NewSalt(CryptoGCMBlockLength(protocol))

	requestStream := NewStream(binary.Base64Decoder(psk))
	requestStream.AddSalt(salt)

	variableHeader := NewRequestVariableHeader(socket5.Socket5AddressDomain, []byte("baidu.com"), 80, 10, []byte("data to baidu.com"))

	b := NewRequestFixedHeader(0, binary.Timestamp(), uint16(len(variableHeader)+16))
	fmt.Println(b)

	eb := requestStream.Encryption(b)

	vb := requestStream.Encryption(variableHeader)

	nb := slices.ConcatMultipleSlices([][]byte{salt, eb, vb})

	fmt.Println("salt:", salt)
	fmt.Println("fixed head:", b)
	fmt.Println("variable head:", variableHeader)
	fmt.Println(nb)

	requestStream2 := NewStream(binary.Base64Decoder(psk))
	saltLength := CryptoGCMBlockLength(protocol)
	salt2 := nb[0:saltLength]
	requestStream2.AddSalt(salt2)
	fixedHead, err := requestStream2.Decryption(nb[saltLength : saltLength+27])
	if err != nil {
		fmt.Println("fixed head decryption fail", err)
	}
	_, ln, err := ParseRequestFixedHeader(fixedHead)
	if err != nil {
		fmt.Println("parse fixed head fail", err)
	}

	variableHead, err := requestStream2.Decryption(nb[saltLength+27 : saltLength+27+int(ln)])
	if err != nil {
		fmt.Println("variable  head decryption fail", err)
	}
	fmt.Println("fixed head:", fixedHead)
	fmt.Println("variable head:", variableHead)

	addr, payload, err := ParseRequestVariableHeader(variableHead)
	if err != nil {
		fmt.Println("parse variable head fail", err)
	}

	fmt.Println("addr:", addr)
	fmt.Println("payload:", string(payload))
}
