package websocket

import (
	"crypto/rand"
	"crypto/sha1"
	"encoding/base64"
	"encoding/binary"
	"errors"
	"fmt"
	"regexp"
	"shadowsocks-go/parser"
	"strings"
)

const (
	OpcodeContinue = 0x0
	OpcodeText     = 0x1
	OpcodeBinary   = 0x2
	OpcodeClose    = 0x8
	OpcodePing     = 0x9
	OpcodePong     = 0xA
)

func NewHttpUpgrade(path string, host, origin string) []byte {
	t := make([]byte, 16)
	_, _ = rand.Read(t)

	key := base64.StdEncoding.EncodeToString(t)

	b := "GET " + path + " HTTP/1.1\r\n" +
		"Host: " + host + "\r\n" +
		"origin: " + origin + "\r\n" +
		"Upgrade: websocket\r\n" +
		"Connection: Upgrade\r\n" +
		"Sec-WebSocket-Key: " + key + "\r\n" +
		"Sec-WebSocket-Version: 13\r\n" +
		"\r\n\r\n"

	return []byte(b)
}

func NewsSecWebSocketKey() string {
	t := make([]byte, 16)
	_, err := rand.Read(t)
	if err != nil {
		panic("NewsSecWebSocketKey fail")
	}

	return base64.StdEncoding.EncodeToString(t)
}

func HttpUpgradeKeyValidate(key, accept string) bool {
	key += "258EAFA5-E914-47DA-95CA-C5AB0DC85B11"
	s1 := sha1.New()
	s1.Write([]byte(key))
	k := base64.StdEncoding.EncodeToString(s1.Sum(nil))

	return k == accept
}

// NewFrame websocket frame struct
// 0                   1                   2                   3
// 0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1
// +-+-+-+-+-------+-+-------------+-------------------------------+
// |F|R|R|R| opcode|M| Payload len |    Extended payload length    |
// |I|S|S|S|  (4)  |A|     (7)     |             (16/64)           |
// |N|V|V|V|       |S|             |   (if payload len==126/127)   |
// | |1|2|3|       |K|             |                               |
// +-+-+-+-+-------+-+-------------+ - - - - - - - - - - - - - - - +
// |     Extended payload length continued, if payload len == 127  |
// + - - - - - - - - - - - - - - - +-------------------------------+
// |                               |Masking-key, if MASK set to 1  |
// +-------------------------------+-------------------------------+
// | Masking-key (continued)       |          Payload Data         |
// +-------------------------------- - - - - - - - - - - - - - - - +
// :                     Payload Data continued ...                :
// + - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - +
// |                     Payload Data continued ...                |
// +---------------------------------------------------------------+
func NewFrame(fin bool, opcode uint8, mask bool, data []byte) []byte {
	dataLen := len(data)
	payloadByteLen := 1

	// calculate payload byte number
	if dataLen <= 125 {
	} else if dataLen <= 0xFFFF {
		payloadByteLen += 2
	} else if dataLen <= 0x7FFFFFFFFFFFFFFF {
		payloadByteLen += 8
	} else {
		panic("frame data length too big")
	}

	if dataLen == 0 {
		panic("frame data length is 0")
	}

	// payload len byte number
	frameLen := 1 + payloadByteLen + dataLen
	if mask {
		frameLen += 4
	}

	// frame slice
	b := make([]byte, frameLen)

	// set fin
	if fin {
		b[0] = b[0] | 0b10000000
	}

	// set opcode
	b[0] = (0b11110000 & b[0]) | opcode

	// set payload length
	if payloadByteLen == 1 {
		b[1] = b[1] | uint8(dataLen)
	} else if payloadByteLen == 3 {
		b[1] = b[1] | 0x7F
		binary.BigEndian.PutUint16(b[2:4], uint16(dataLen))
	} else {
		b[1] = b[1] | 0x7F
		binary.BigEndian.PutUint64(b[2:10], uint64(dataLen))
	}

	// set mask
	dx := 1 + payloadByteLen
	mk := b[1+payloadByteLen : payloadByteLen+5]
	if mask {
		b[1] = b[1] | 0b10000000
		_, err := rand.Read(mk)
		if err != nil {
			panic("frame mask rand read fail")
		}

		// mask before payload data
		dx += 4
	}

	ds := b[dx:]
	copy(ds, data)

	// mask data
	if mask {
		for k, _ := range ds {
			m := mk[k%4]
			ds[k] ^= m
		}
	}

	return b
}

func ParseFrame(buf []byte) (fin bool, opcode uint8, data []byte, err error) {
	bufLen := len(buf)
	if bufLen < 3 {
		err = errors.New("frame length less than 3")
		return
	}

	fin = (buf[0] & 0b10000000) > 0
	opcode = buf[0] & 0b00001111

	mask := (buf[1] & 0b10000000) > 0
	maskLen := 0
	maskStart := 2
	if mask {
		maskLen = 4
	}

	payloadLenMark := int(buf[1] & 0b01111111)

	if payloadLenMark == 0 {
		err = errors.New("payload length is 0")
		return
	}

	if payloadLenMark <= 125 {
		s := 2 + maskLen
		if bufLen < s+payloadLenMark {
			err = errors.New(fmt.Sprintf("payload length less then %d", s+payloadLenMark))
			return
		}
		data = buf[s : s+payloadLenMark]
	} else if payloadLenMark == 126 {
		s := 4 + maskLen
		dl := int(binary.BigEndian.Uint16(buf[2:4]))
		if bufLen < s+dl {
			err = errors.New(fmt.Sprintf("payload length less then %d", s+dl))
			return
		}

		data = buf[s : s+dl]
		maskStart += 2

	} else { // 127
		s := 10 + maskLen

		if (buf[3] & 0b10000000) > 0 {
			err = errors.New("64-bit unsigned integer the most significant bit must be 0")
			return
		}

		dl := int(binary.BigEndian.Uint64(buf[2:10]))
		if bufLen < s+dl {
			err = errors.New(fmt.Sprintf("payload length less then %d", s+dl))
			return
		}

		data = buf[s : s+dl]
		maskStart += 8
	}
	var mk []byte
	if mask {
		mk = buf[maskStart : maskStart+4]

		for k, _ := range data {
			m := mk[k%4]
			data[k] ^= m
		}
	}

	return
}

func ParseFramePayloadLength(buf []byte) (fin bool, opcode uint8, mask []byte, payload uint64, err error) {
	bufLen := len(buf)
	if bufLen < 2 {
		err = parser.ParseContinue.WithReason("frame length less than 2")
		return
	}

	fin = (buf[0] & 0b10000000) > 0
	opcode = buf[0] & 0b00001111

	isMask := (buf[1] & 0b10000000) > 0
	maskLen := 0
	maskStart := 2
	if isMask {
		maskLen = 4
	}

	payloadLenMark := int(buf[1] & 0b01111111)

	if payloadLenMark == 0 {
		err = parser.ParseFatal.WithReason("payload length is 0")
		return
	}

	if payloadLenMark <= 125 {
		payload = uint64(payloadLenMark)

	} else if payloadLenMark == 126 {
		s := 2 + 2 + maskLen
		if bufLen < s {
			err = parser.ParseContinue.WithReason(fmt.Sprintf("frame length less than %d", s))
			return
		}

		payload = uint64(binary.BigEndian.Uint16(buf[2:4]))
		maskStart += 2
	} else { // 127
		s := 2 + 8 + maskLen
		if bufLen < s {
			err = parser.ParseContinue.WithReason(fmt.Sprintf("frame length less than %d", s))
			return
		}

		if (buf[3] & 0b10000000) > 0 {
			err = parser.ParseFatal.WithReason("64-bit unsigned integer the most significant bit must be 0")
			return
		}

		payload = binary.BigEndian.Uint64(buf[2:10])
		maskStart += 8
	}

	// mask
	if isMask {
		if bufLen < maskStart+4 {
			err = parser.ParseContinue.WithReason(fmt.Sprintf("frame length less than %d", maskStart+4))
			return
		}
		mask = buf[maskStart : maskStart+4]
	}

	return
}
func ParseHttpHeaders(buf []byte) (headerMap map[string]string, err error) {
	str := string(buf)
	idx := strings.Index(str, "\r\n\r\n")

	fmt.Println(idx)
	headers := strings.Split(str, "\r\n")
	// two \r\n and http request line
	if len(headers) < 3 {
		return nil, parser.ParseContinue.WithReason("http request line required")
	}
	re := regexp.MustCompile(" +")

	requestLine := headers[0]
	ra := re.Split(requestLine, -1)
	//fmt.Println("request line len:", len(ra), "content:", ra)
	if len(ra) != 3 {
		return nil, parser.ParseFatal.WithReason("http request line required")
	}

	re1 := regexp.MustCompile(" *: *")

	for i := 1; i < len(headers)-2; i++ {
		//fmt.Println("index:", i, " header:", headers[i])
		sp := re1.Split(headers[i], -1)
		k := strings.Trim(sp[0], " ")
		headerMap[k] = strings.Trim(sp[1], " ")
	}

	//fmt.Println(headerMap)
	return
}
