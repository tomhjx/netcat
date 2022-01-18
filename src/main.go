package main

import (
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/tomhjx/netcat/core"
)

func init() {
	// 设置当前时区
	var cstZone = time.FixedZone("CST", 8*3600) // 东八
	time.Local = cstZone

	log.SetLevel(log.TraceLevel)
	log.SetReportCaller(true)
	log.SetFormatter(&log.TextFormatter{
		TimestampFormat: "2006-01-02 15:04:05",
		FullTimestamp:   true})

}

func main() {
	core.NewProcessor().Run()
}
