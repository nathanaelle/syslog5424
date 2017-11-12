package syslog5424 // import "github.com/nathanaelle/syslog5424"

import (
	"bytes"
	"fmt"
	"testing"
)

func TestTransport_T_ZEROENDED(t *testing.T) {
	msg := []byte("<0>1 1970-01-01T01:00:00Z bla bli blu blo - message")

	buf1 := []byte("<0>1 1970-01-01T01:00:00Z bla bli blu blo - message\x00<0>1 1970-01-01T01:00:00Z bla bli blu blo - message\x00")
	buf2 := []byte("<0>1 1970-01-01T01:00:00Z bla bli blu blo - message\x00")

	transportPrefixTester(t, T_ZEROENDED, false, nil, []byte(``), nil, nil)
	transportSuffixTester(t, T_ZEROENDED, false, nil, []byte(``), nil, nil)

	transportPrefixTester(t, T_ZEROENDED, false, buf1, buf1, nil, nil)
	transportSuffixTester(t, T_ZEROENDED, false, buf1, msg, buf2, nil)
	transportSuffixTester(t, T_ZEROENDED, false, buf2, msg, nil, nil)
	transportSuffixTester(t, T_ZEROENDED, false, msg, msg, nil, nil)
	transportSuffixTester(t, T_ZEROENDED, true, msg, msg, nil, nil)
}

func TestTransport_T_LFENDED(t *testing.T) {
	msg := []byte("<0>1 1970-01-01T01:00:00Z bla bli blu blo - message")

	buf1 := []byte("<0>1 1970-01-01T01:00:00Z bla bli blu blo - message\n<0>1 1970-01-01T01:00:00Z bla bli blu blo - message\n")
	buf2 := []byte("<0>1 1970-01-01T01:00:00Z bla bli blu blo - message\n")

	transportSuffixTester(t, T_LFENDED, false, nil, []byte(``), nil, nil)
	transportPrefixTester(t, T_LFENDED, false, nil, []byte(``), nil, nil)

	transportPrefixTester(t, T_LFENDED, false, buf1, buf1, nil, nil)
	transportSuffixTester(t, T_LFENDED, false, buf1, msg, buf2, nil)
	transportSuffixTester(t, T_LFENDED, false, buf2, msg, nil, nil)
	transportSuffixTester(t, T_LFENDED, false, msg, msg, nil, nil)
	transportSuffixTester(t, T_LFENDED, true, msg, msg, nil, ERR_TRANSPORT_INCOMPLETE)
}

func TestTransport_T_RFC5425(t *testing.T) {
	msg := []byte("<0>1 1970-01-01T01:00:00Z bla bli blu blo - message")

	buf1 := []byte("51 <0>1 1970-01-01T01:00:00Z bla bli blu blo - message51 <0>1 1970-01-01T01:00:00Z bla bli blu blo - message")
	buf2 := []byte("51 <0>1 1970-01-01T01:00:00Z bla bli blu blo - message")

	transportSuffixTester(t, T_RFC5425, false, nil, []byte(``), nil, nil)
	transportSuffixTester(t, T_RFC5425, false, msg, msg, nil, nil)
	transportSuffixTester(t, T_RFC5425, true, msg, msg, nil, nil)

	transportPrefixTester(t, T_RFC5425, false, nil, []byte(``), nil, nil)

	transportPrefixTester(t, T_RFC5425, false, buf1, msg, buf2, nil)
	transportPrefixTester(t, T_RFC5425, false, buf1, msg, buf2, nil)
	transportPrefixTester(t, T_RFC5425, false, buf2, msg, nil, nil)

	transportPrefixTester(t, T_RFC5425, false, buf2[:len(buf2)-10], nil, nil, nil)
	transportPrefixTester(t, T_RFC5425, true, buf1[:len(buf1)-10], msg, buf2[:len(buf2)-10], nil)
	transportPrefixTester(t, T_RFC5425, true, buf2[:len(buf2)-10], msg[:len(msg)-10], nil, ERR_TRANSPORT_INCOMPLETE)
}

func transportPrefixTester(t *testing.T, transport Transport, atEOF bool, buffer, data, rest []byte, err error) {
	d, r, e := transport.PrefixStrip(buffer, atEOF)
	if !bytes.Equal(data, d) || !bytes.Equal(rest, r) || err != e {
		t.Error(fmt.Errorf("For Prefix {%q}\n  Expected {%q} {%q} %v\n  Got: {%q} {%q} %v\n", buffer, data, rest, err, d, r, e))
		t.Fail()
	}
}

func transportSuffixTester(t *testing.T, transport Transport, atEOF bool, buffer, data, rest []byte, err error) {
	d, r, e := transport.SuffixStrip(buffer, atEOF)
	if !bytes.Equal(data, d) || !bytes.Equal(rest, r) || err != e {
		t.Error(fmt.Errorf("For Suffix {%q}\n  Expected {%q} {%q} %v\n  Got: {%q} {%q} %v\n", buffer, data, rest, err, d, r, e))
		t.Fail()
	}
}
