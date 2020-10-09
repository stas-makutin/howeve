package httpsrv

import (
	"strconv"
	"sync/atomic"

	"github.com/stas-makutin/howeve/log"
)

type wsOrdinal uint32
type wsConnOrdinal wsOrdinal
type wsMsgOrdinal wsOrdinal

var wsCurrentConnOrdinal wsConnOrdinal
var wsCurrentMsgOrdinal wsMsgOrdinal

func nextWsConnOrdinal() wsConnOrdinal {
	return wsConnOrdinal(atomic.AddUint32((*uint32)(&wsCurrentConnOrdinal), 1))
}

func nextWsMsgOrdinal() wsMsgOrdinal {
	return wsMsgOrdinal(atomic.AddUint32((*uint32)(&wsCurrentMsgOrdinal), 1))
}

func (o wsConnOrdinal) logOpen() {
	log.Report(log.SrcWS, "S", strconv.FormatUint(uint64(o), 36))
}

func (o wsConnOrdinal) logClose() {
	log.Report(log.SrcWS, "F", strconv.FormatUint(uint64(o), 36))
}

func (o wsConnOrdinal) logMsg(mo wsMsgOrdinal, id string, incoming bool, msgType string, size int64) {
	dir := "O"
	if incoming {
		dir = "I"
	}
	log.Report(log.SrcWS, dir, strconv.FormatUint(uint64(o), 36), strconv.FormatUint(uint64(mo), 36), id, strconv.FormatInt(size, 10), msgType)
}

func (o wsConnOrdinal) logErr(mo wsMsgOrdinal, id string, incoming bool, size int64, report string) {
	dir := "O"
	if incoming {
		dir = "I"
	}
	log.Report(log.SrcWS, dir, strconv.FormatUint(uint64(o), 36), strconv.FormatUint(uint64(mo), 36), id, strconv.FormatInt(size, 10), report)
}
