package main

import (
		"log/syslog"
		syslogv2 "gopkg.in/mcuadros/go-syslog.v2"
		"flag"
)

var syslogserver string
var listenaddress string

func init() {
	flag.StringVar(&syslogserver, "syslogserver", "", "Set remote syslogserver")
	flag.StringVar(&listenaddress, "listenaddress", "127.0.0.1:514", "Set address to listen on")
	flag.Parse()
}

func main() {
	channel := make(syslogv2.LogPartsChannel)
	handler := syslogv2.NewChannelHandler(channel)

	server := syslogv2.NewServer()
	server.SetFormat(syslogv2.RFC3164)
	server.SetHandler(handler)

	server.ListenUDP(listenaddress)
	server.ListenTCP(listenaddress)

	server.Boot()

	go func(channel syslogv2.LogPartsChannel) {
		for logParts := range channel {
			tag := logParts["tag"].(string)
			priority, _ := logParts["priority"].(syslog.Priority)
			content := logParts["content"].(string)

			logwriter, e := syslog.Dial("tcp", syslogserver, priority, tag)

			if e == nil {
				logwriter.Notice(content)
			}
		}
	}(channel)

	server.Wait()
}
