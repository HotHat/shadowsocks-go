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
