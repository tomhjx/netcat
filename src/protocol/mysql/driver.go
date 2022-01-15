package mysql

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"log"
	"strconv"
	"strings"

	"github.com/tomhjx/netcat/protocol"
)

type Driver struct {
	protocol.BaseDriver
	stmtMap  map[uint32]*Stmt
	stmtTMap map[uint8]*Stmt
}

func NewDriver() protocol.Driver {

	return &Driver{
		BaseDriver: protocol.BaseDriver{
			SpecPort:          3306,
			SpecTransportType: "tcp",
		},
		stmtMap:  make(map[uint32]*Stmt),
		stmtTMap: make(map[uint8]*Stmt),
	}
}

func (me *Driver) ResolveClient(payload0 []byte) (string, string) {
	return me.resolve(payload0)
}

func (me *Driver) ResolveServer(payload0 []byte) (string, string) {
	return me.resolve(payload0)
}

func (me *Driver) resolve(payload0 []byte) (string, string) {

	//read packet
	var (
		payloadWB bytes.Buffer
		seq       uint8
		err       error
	)
	if seq, err = me.resolvePacketTo(bytes.NewReader(payload0), &payloadWB); err != nil {
		return "", ""
	}

	payload := payloadWB.Bytes()

	if me.continuePrepareStmt(seq, payload) {
		return "", ""
	}

	switch payload[0] {

	case 0xff:
		errorCode := int(binary.LittleEndian.Uint16(payload[1:3]))
		errorMsg, _ := ReadStringFromByte(payload[4:])

		return "MYSQL_RESP_ERR", fmt.Sprintf("Err code:%s,Err msg:%s", strconv.Itoa(errorCode), strings.TrimSpace(errorMsg))

	case 0x00:
		var pos = 1
		l, _, _ := LengthEncodedInt(payload[pos:])
		affectedRows := int(l)
		return "MYSQL_RESP_EFFECT", fmt.Sprintf("Effect Row:%s", strconv.Itoa(affectedRows))

	case COM_INIT_DB:

		return "MYSQL_REQ_USE_DB", fmt.Sprintf("USE %s;", payload[1:])
	case COM_DROP_DB:

		return "MYSQL_REQ_DROP_DB", fmt.Sprintf("Drop DB %s;", payload[1:])
	case COM_CREATE_DB, COM_QUERY:

		return "MYSQL_REQ_QUERY", string(payload[1:])

	case COM_STMT_PREPARE:

		stmt := &Stmt{
			Query: string(payload[1:]),
		}
		me.stmtTMap[seq+1] = stmt

		// return "MYSQL_REQ_PREPARE_QUERY", stmt.Query

	case COM_STMT_SEND_LONG_DATA:

		stmtID := binary.LittleEndian.Uint32(payload[1:5])
		paramId := binary.LittleEndian.Uint16(payload[5:7])
		stmt, _ := me.stmtMap[stmtID]

		if stmt.Args[paramId] == nil {
			stmt.Args[paramId] = payload[7:]
		} else {
			if b, ok := stmt.Args[paramId].([]byte); ok {
				b = append(b, payload[7:]...)
				stmt.Args[paramId] = b
			}
		}
		return "", ""
	case COM_STMT_RESET:

		stmtID := binary.LittleEndian.Uint32(payload[1:5])
		stmt, _ := me.stmtMap[stmtID]
		stmt.Args = make([]interface{}, stmt.ParamCount)
		return "", ""
	case COM_STMT_EXECUTE:
		return me.resolveExecuteStmt(payload)
	default:
		return "", ""
	}
	return "", ""
	// fmt.Println(GetNowStr(true) + msg + "\n")
}

func (me *Driver) continuePrepareStmt(seq uint8, payload []byte) bool {
	stmt, ok := me.stmtTMap[seq]
	if !ok {
		return false
	}
	delete(me.stmtTMap, seq)
	stmtID := binary.LittleEndian.Uint32(payload[1:5])
	stmt.ID = stmtID
	//record stm sql
	me.stmtMap[stmtID] = stmt
	stmt.FieldCount = binary.LittleEndian.Uint16(payload[5:7])
	stmt.ParamCount = binary.LittleEndian.Uint16(payload[7:9])
	stmt.Args = make([]interface{}, stmt.ParamCount)
	return true
}

func (me *Driver) resolveExecuteStmt(payload []byte) (string, string) {

	var pos = 1
	stmtID := binary.LittleEndian.Uint32(payload[pos : pos+4])
	pos += 4
	var stmt *Stmt
	var ok bool
	if stmt, ok = me.stmtMap[stmtID]; !ok {
		log.Println("ERR : Not found stm id", stmtID)
		return "", ""
	}

	//params
	pos += 5
	if stmt.ParamCount > 0 {

		//（Null-Bitmap，len = (paramsCount + 7) / 8 byte）
		step := int((stmt.ParamCount + 7) / 8)
		nullBitmap := payload[pos : pos+step]
		pos += step

		//Parameter separator
		flag := payload[pos]

		pos++

		var pTypes []byte
		var pValues []byte

		//if flag == 1
		//n （len = paramsCount * 2 byte）
		if flag == 1 {
			pTypes = payload[pos : pos+int(stmt.ParamCount)*2]
			pos += int(stmt.ParamCount) * 2
			pValues = payload[pos:]
		}

		//bind params
		err := stmt.BindArgs(nullBitmap, pTypes, pValues)
		if err != nil {
			log.Println("ERR : Could not bind params", err)
		}
	}
	return "MYSQL_REQ_PREPARE_EXEC", string(stmt.WriteToText())
}

func (me *Driver) resolvePacketTo(r io.Reader, w io.Writer) (uint8, error) {

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
