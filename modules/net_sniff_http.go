package modules

import (
	"fmt"

	"github.com/evilsocket/bettercap-ng/core"
	"regexp"

	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
)

var httpRe = regexp.MustCompile("(?s).*(GET|HEAD|POST|PUT|DELETE|CONNECT|OPTIONS|TRACE|PATCH) (.+) HTTP/\\d\\.\\d.+Host: ([^\\s]+)")
var uaRe = regexp.MustCompile("(?s).*User-Agent: ([^\\n]+).+")

func httpParser(ip *layers.IPv4, pkt gopacket.Packet, tcp *layers.TCP) bool {
	data := tcp.Payload
	dataSize := len(data)

	if dataSize < 20 {
		return false
	}

	m := httpRe.FindSubmatch(data)
	if len(m) != 4 {
		return false
	}

	ua := ""
	mu := uaRe.FindSubmatch(data)
	if len(mu) == 2 {
		ua = string(mu[1])
	}

	url := fmt.Sprintf("%s", core.Yellow(string(m[3])))
	if tcp.DstPort != 80 {
		url += fmt.Sprintf(":%s", vPort(tcp.DstPort))
	}
	url += fmt.Sprintf("%s", string(m[2]))

	SniffPrinter("[%s] %s %s %s %s %s\n",
		vTime(pkt.Metadata().Timestamp),
		core.W(core.BG_RED+core.FG_BLACK, "http"),
		vIP(ip.SrcIP),
		core.W(core.BG_LBLUE+core.FG_BLACK, vURL(string(m[1]))),
		vURL(url),
		core.Dim(ua))

	return true
}
