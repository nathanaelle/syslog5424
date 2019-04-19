package syslog5424 // import "github.com/nathanaelle/syslog5424/v2"

import (
	"bytes"
	"fmt"
	"testing"
)

func TestTransportZeroEnded(t *testing.T) {
	msg := []byte("<0>1 1970-01-01T01:00:00Z bla bli blu blo - message")

	buf1 := []byte("<0>1 1970-01-01T01:00:00Z bla bli blu blo - message\x00<0>1 1970-01-01T01:00:00Z bla bli blu blo - message\x00")
	buf2 := []byte("<0>1 1970-01-01T01:00:00Z bla bli blu blo - message\x00")

	transportPrefixTester(t, TransportZeroEnded, false, nil, []byte(``), nil, nil)
	transportSuffixTester(t, TransportZeroEnded, false, nil, []byte(``), nil, nil)

	transportPrefixTester(t, TransportZeroEnded, false, buf1, buf1, nil, nil)
	transportSuffixTester(t, TransportZeroEnded, false, buf1, msg, buf2, nil)
	transportSuffixTester(t, TransportZeroEnded, false, buf2, msg, nil, nil)
	transportSuffixTester(t, TransportZeroEnded, false, msg, msg, nil, nil)
	transportSuffixTester(t, TransportZeroEnded, true, msg, msg, nil, nil)
}

func TestTransportLFEnded(t *testing.T) {
	msg := []byte("<0>1 1970-01-01T01:00:00Z bla bli blu blo - message")

	buf1 := []byte("<0>1 1970-01-01T01:00:00Z bla bli blu blo - message\n<0>1 1970-01-01T01:00:00Z bla bli blu blo - message\n")
	buf2 := []byte("<0>1 1970-01-01T01:00:00Z bla bli blu blo - message\n")

	transportSuffixTester(t, TransportLFEnded, false, nil, []byte(``), nil, nil)
	transportPrefixTester(t, TransportLFEnded, false, nil, []byte(``), nil, nil)

	transportPrefixTester(t, TransportLFEnded, false, buf1, buf1, nil, nil)
	transportSuffixTester(t, TransportLFEnded, false, buf1, msg, buf2, nil)
	transportSuffixTester(t, TransportLFEnded, false, buf2, msg, nil, nil)
	transportSuffixTester(t, TransportLFEnded, false, msg, msg, nil, nil)
	transportSuffixTester(t, TransportLFEnded, true, msg, msg, nil, ErrTransportIncomplete)
}

func TestTransportRFC5425(t *testing.T) {
	msg := []byte("<0>1 1970-01-01T01:00:00Z bla bli blu blo - message")

	buf1 := []byte("51 <0>1 1970-01-01T01:00:00Z bla bli blu blo - message51 <0>1 1970-01-01T01:00:00Z bla bli blu blo - message")
	buf2 := []byte("51 <0>1 1970-01-01T01:00:00Z bla bli blu blo - message")

	transportSuffixTester(t, TransportRFC5425, false, nil, []byte(``), nil, nil)
	transportSuffixTester(t, TransportRFC5425, false, msg, msg, nil, nil)
	transportSuffixTester(t, TransportRFC5425, true, msg, msg, nil, nil)

	transportPrefixTester(t, TransportRFC5425, false, nil, []byte(``), nil, nil)

	transportPrefixTester(t, TransportRFC5425, false, buf1, msg, buf2, nil)
	transportPrefixTester(t, TransportRFC5425, false, buf1, msg, buf2, nil)
	transportPrefixTester(t, TransportRFC5425, false, buf2, msg, nil, nil)

	transportPrefixTester(t, TransportRFC5425, false, buf2[:len(buf2)-10], nil, nil, nil)
	transportPrefixTester(t, TransportRFC5425, true, buf1[:len(buf1)-10], msg, buf2[:len(buf2)-10], nil)
	transportPrefixTester(t, TransportRFC5425, true, buf2[:len(buf2)-10], msg[:len(msg)-10], nil, ErrTransportIncomplete)
}

func transportPrefixTester(t *testing.T, transport Transport, atEOF bool, buffer, data, rest []byte, err error) {
	d, r, e := transport.PrefixStrip(buffer, atEOF)
	if !bytes.Equal(data, d) || !bytes.Equal(rest, r) || err != e {
		t.Error(fmt.Errorf("For Prefix {%q}\n  Expected {%q} {%q} %v\n  Got: {%q} {%q} %v", buffer, data, rest, err, d, r, e))
		t.Fail()
	}
}

func transportSuffixTester(t *testing.T, transport Transport, atEOF bool, buffer, data, rest []byte, err error) {
	d, r, e := transport.SuffixStrip(buffer, atEOF)
	if !bytes.Equal(data, d) || !bytes.Equal(rest, r) || err != e {
		t.Error(fmt.Errorf("For Suffix {%q}\n  Expected {%q} {%q} %v\n  Got: {%q} {%q} %v", buffer, data, rest, err, d, r, e))
		t.Fail()
	}
}
