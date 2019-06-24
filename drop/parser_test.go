package drop

import (
	"fmt"
	"github.com/wix-playground/kube-iptables-tailer/util"
	"testing"
	"time"
)

const (
	testHostname  = "hostname"
	testLogPrefix = "log-prefix"
	testSrcIP     = "11.111.11.111"
	testDstIP     = "22.222.22.222"
	testSrcPort   = "32157"
	testDstPort   = "445"
)

// Test if PacketDrop.IsExpired() works
func TestPacketDropIsExpired(t *testing.T) {
	expiredTime := util.GetExpiredTimeInString(util.DefaultPacketDropExpirationMinutes, PacketDropLogTimeLayout)
	expiredPacketDrop := PacketDrop{LogTime: expiredTime}
	if !expiredPacketDrop.IsExpired() {
		t.Fatal("Expected IsExpired() return true, got false")
	}

	curTime := time.Now().Format(PacketDropLogTimeLayout)
	curPacketDrop := PacketDrop{LogTime: curTime}
	if curPacketDrop.IsExpired() {
		t.Fatal("Expected IsExpired() return false, got true")
	}

}

// Test if packet parser works for packet drop
func TestParsingDropLog(t *testing.T) {
	channel := make(chan PacketDrop, 100)
	// need to use curTime because parse() will not insert expired packetDrop
	curTime := time.Now().Format(PacketDropLogTimeLayout)
	testLog := fmt.Sprintf("%s %s %s SRC=%s DST=%s SPT=%s DPT=%s", curTime, testHostname, testLogPrefix, testSrcIP, testDstIP, testSrcPort, testDstPort)
	expected := PacketDrop{
		LogTime:  curTime,
		HostName: testHostname,
		SrcIP:    testSrcIP,
		DstIP:    testDstIP,
		SrcPort:  testSrcPort,
		DstPort:  testDstPort,
	}
	parse(testLogPrefix, testLog, channel)

	result := <-channel
	if result != expected {
		t.Fatalf("Expected %+v, but got result %+v", expected, result)
	}
}

func TestParsingDropKernelLogTimeFormat(t *testing.T) {
	channel := make(chan PacketDrop, 100)
	eventTime := time.Now().In(time.UTC)
	eventTimeStr := eventTime.Format(time.Stamp)
	// Our timestamp does not have nanseconds
	eventTimeToExpect := eventTime.Truncate(time.Second).Format(PacketDropLogTimeLayout)
	//2019, 6, 23, 18, 17, 16, 0, time.UTC)
	testLog := fmt.Sprintf("%s %s kernel: [383376.014436] calico-packet: IN=eni3a0f5e4322a OUT=eth1 MAC=86:52:e0:ab:82:52:22:05:e3:af:1b:3c:08:00 SRC=%s DST=%s LEN=60 TOS=0x00 PREC=0x00 TTL=63 ID=5082 DF PROTO=TCP SPT=%s DPT=%s WINDOW=26883 RES=0x00 SYN URGP=0",
		eventTimeStr, testHostname, testSrcIP, testDstIP, testSrcPort, testDstPort)

	expected := PacketDrop{
		LogTime:  eventTimeToExpect,
		HostName: testHostname,
		SrcIP:    testSrcIP,
		DstIP:    testDstIP,
		SrcPort:  testSrcPort,
		DstPort:  testDstPort,
	}
	parse("calico-packet:", testLog, channel)

	result := <-channel
	if result != expected {
		t.Fatalf("Expected %+v, but got result %+v", expected, result)
	}
}

// Test if packet parser works for outdated packet drop (should not add it to channel)
func TestParsingExpiredPacketDropLog(t *testing.T) {
	channel := make(chan PacketDrop, 100)
	expiredTime := util.GetExpiredTimeInString(util.DefaultPacketDropExpirationMinutes, PacketDropLogTimeLayout)
	expiredLog := fmt.Sprintf("%s %s %s SRC=%s DST=%s",
		expiredTime, testHostname, testLogPrefix, testSrcIP, testDstIP)
	parse(testLogPrefix, expiredLog, channel)

	select {
	case result := <-channel:
		t.Fatalf("expected channel empty, but got result %v", result)
	default:
		return
	}
}

// Test if packet parser works for bad packet drop (should return error)
func TestParsingBadPacketDropLog(t *testing.T) {
	channel := make(chan PacketDrop)
	// testing bad log without source IP
	curTime := time.Now().Format(PacketDropLogTimeLayout)
	testLog1 := fmt.Sprintf("%s %s %s %s", curTime, testHostname, testLogPrefix, testDstIP)
	err := parse(testLogPrefix, testLog1, channel)
	if err == nil {
		t.Fatalf("Expected error, but got error nil!")
	}
	// testing bad log without destination IP
	testLog2 := fmt.Sprintf("%s %s %s %s", curTime, testHostname, testLogPrefix, testSrcIP)
	err = parse(testLogPrefix, testLog2, channel)
	if err == nil {
		t.Fatalf("Expected error, but got error nil!")
	}
}

// Test if packet parser works for none packet drop log (should just ignore and not return error)
func TestParsingNonePacketDropLog(t *testing.T) {
	channel := make(chan PacketDrop)
	curTime := time.Now().Format(PacketDropLogTimeLayout)
	testLog := fmt.Sprintf("%s %s None Packet Drop Log", curTime, testHostname)
	err := parse(testLogPrefix, testLog, channel)

	if err != nil {
		t.Fatalf("Expected error nil, but got error %s", err)
	}
}

// Test if getPacketDropLogFields() function works
func TestGetPacketDropLogFields(t *testing.T) {
	emptyLog := " "
	_, err := getPacketDropLogFields(emptyLog)
	if err == nil {
		t.Fatalf("Expected error from empty log, but got nil")
	}

	curTime := time.Now().Format(PacketDropLogTimeLayout)
	// missing IPs
	packetDropLogMissingField := fmt.Sprintf("%s %s ", curTime, testHostname)
	_, err = getPacketDropLogFields(packetDropLogMissingField)
	if err == nil {
		t.Fatalf("Expected error from log %s, but got nil", packetDropLogMissingField)
	}
}
