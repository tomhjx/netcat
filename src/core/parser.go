package core

import (
	"bytes"
	"errors"
	"io"
	"log"
	"strconv"

	mysqlAide "github.com/tomhjx/netcat/parser/mysql"
)

type Parser struct {
	aide *mysqlAide.Parser
}

func NewParser() *Parser {

	aide := mysqlAide.NewInstance()

	return &Parser{aide: aide}
}

func (per *Parser) Resolve(src *Source) *Resolved {

	//read packet
	var payload bytes.Buffer
	var seq uint8
	var err error
	if seq, err = per.resolvePacketTo(bytes.NewReader(src.payload), &payload); err != nil {
		return nil
	}

	src.payload = payload.Bytes()

	log.Println("start resolve")

	ret := Resolved{}

	if src.transport.Src().String() == strconv.Itoa(per.aide.Port) {
		ret.isClientFlow = false
	} else {
		ret.isClientFlow = true
	}
	code, content := per.aide.Resolve(int(seq), src.payload)
	if code != "" {
		ret.code = code
		ret.content = content
	}

	return &ret
}

func (per *Parser) resolvePacketTo(r io.Reader, w io.Writer) (uint8, error) {

	header := make([]byte, 4)
	if n, err := io.ReadFull(r, header); err != nil {
		if n == 0 && err == io.EOF {
			return 0, io.EOF
		}
		return 0, errors.New("ERR : Unknown stream")
	}

	length := int(uint32(header[0]) | uint32(header[1])<<8 | uint32(header[2])<<16)

	var seq uint8
	seq = header[3]

	if n, err := io.CopyN(w, r, int64(length)); err != nil {
		return 0, errors.New("ERR : Unknown stream")
	} else if n != int64(length) {
		return 0, errors.New("ERR : Unknown stream")
	} else {
		return seq, nil
	}

	return seq, nil
}
