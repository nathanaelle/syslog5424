package sdata // import "github.com/nathanaelle/syslog5424/v2/sdata"

import (
	"fmt"
)

// NextHeader try to decode a Structured Data header from a []byte
func NextHeader(data []byte) (header string, ret []byte, err error) {
	if data == nil || len(data) < 4 {
		return "", nil, fmt.Errorf("invalid length for [%q]", string(data))
	}

	length := 0
	ret = data
	for len(ret) > 0 {
		if ret[0] < ' ' || data[length] == '=' || data[length] == byte(127) || data[length] == ']' || data[length] == '"' {
			return "", nil, fmt.Errorf("invalid char %v", data[length])

		}
		if data[length] == ' ' {
			if length == 0 {
				return "", nil, fmt.Errorf("unexpected space")
			}
			header = string(data[0:length])
			return
		}
		ret = ret[1:]
		length++
	}
	return "", nil, fmt.Errorf("unexpected EOF")
}

// NextNonSpace try to decode a Structured Data part from a []byte
func NextNonSpace(data []byte) (ret []byte, err error) {
	if data == nil || len(data) < 1 {
		return nil, fmt.Errorf("invalid length for [%q]", string(data))
	}

	ret = data
	if ret[0] != ' ' {
		return nil, fmt.Errorf("unexpected non space")
	}

	for len(ret) > 0 {
		switch ret[0] {
		case ' ':
			ret = ret[1:]
		default:
			return
		}
	}
	return nil, fmt.Errorf("unexpected EOF")
}

// NextSDName try to decode a Structured Data Name from a []byte
func NextSDName(data []byte) (name string, ret []byte, err error) {
	if data == nil || len(data) < 4 {
		return "", nil, fmt.Errorf("invalid length for [%q]", string(data))
	}

	length := 0
	ret = data
	for len(ret) > 0 {
		if data[length] <= ' ' || data[length] == byte(127) || data[length] == ']' || data[length] == '"' {
			return "", nil, fmt.Errorf("invalid char %v", data[length])

		}
		if data[length] == '=' {
			ret = ret[1:]
			name = string(data[0:length])
			return
		}
		ret = ret[1:]
		length++
	}
	return "", nil, fmt.Errorf("unexpected EOF")
}

// NextSDValue try to decode a Structured Data Value from a []byte
func NextSDValue(data []byte) (value string, ret []byte, err error) {
	if data == nil || len(data) < 2 {
		return "", nil, fmt.Errorf("invalid length for [%q]", string(data))
	}
	if data[0] != '"' {
		return "", nil, fmt.Errorf("no double quote at begin of [%q]", string(data))
	}

	length := 1
	ret = data[1:]
	for len(ret) > 0 {
		if data[length] == '\\' {
			if len(ret) < 2 {
				return "", nil, fmt.Errorf("invalid char %v", data[length])
			}
			ret = ret[2:]
			length += 2
			continue
		}

		if data[length] == ']' {
			return "", nil, fmt.Errorf("invalid char %v", data[length])
		}

		if data[length] == '"' {
			value = string(data[1:length])
			ret = ret[1:]
			return
		}
		ret = ret[1:]
		length++
	}
	return "", nil, fmt.Errorf("unexpected EOF")
}
