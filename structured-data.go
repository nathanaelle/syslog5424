package syslog5424 // import "github.com/nathanaelle/syslog5424"

import (
	"reflect"
	"strings"
	"unicode"
)

//	http://www.iana.org/assignments/syslog-parameters/syslog-parameters.xhtml#syslog-parameters-4

type (
	listStructuredData []interface{}

	SDUnmarshaler interface {
		Unmarshal5424([]byte) error
	}

	SDMarshaler interface {
		Marshal5424() ([]byte, error)
	}

	SDPEN interface {
		GetPEN() uint64
		SetPEN(uint64)
	}
)

var (
	marshalerType   = reflect.TypeOf(new(SDMarshaler)).Elem()
	unmarshalerType = reflect.TypeOf(new(SDUnmarshaler)).Elem()
	sdpenType       = reflect.TypeOf(new(SDPEN)).Elem()
	emptyListSD     = listStructuredData([]interface{}{})
)


// Append data to a list of structured data
func (listSD listStructuredData) Add(data ...interface{}) listStructuredData {
	return listStructuredData(append([]interface{}(listSD), data...))
}


func (listSD listStructuredData) marshal5424() []byte {
	if len(listSD) == 0 {
		return []byte{'-'}
	}

	ret_s	:= 0
	t_a	:= make([][]byte, len(listSD))
	for i, sd := range listSD {
		t:= MarshalSD(sd)
		ret_s	+= len(t)
		t_a[i]	= t
	}

	ret	:= make([]byte, 0, ret_s)
	for _,sd := range t_a {
		ret = append( ret, sd... )
	}

	return ret
}

func (listSD listStructuredData) String() string {
	return string(listSD.marshal5424())
}



// tagOptions is the string following a comma in a struct field's "sd5424"
// tag, or the empty string. It does not include the leading comma.
type tagOptions []string

// parseTag splits a struct field's sd5424 tag into its name and array of options.
func parseTag(tag string) (string, tagOptions) {
	if idx := strings.Index(tag, ","); idx != -1 {
		return tag[:idx], tagOptions(strings.Split(tag[idx+1:], ","))
	}
	return tag, tagOptions([]string{})
}

func (o tagOptions) Contains(optionName string) bool {
	if len(o) == 0 {
		return false
	}
	flags := []string(o)
	for _, f_o := range flags {
		if f_o == optionName {
			return true
		}
	}
	return false
}

func isValidTag(s string) bool {
	if s == "" {
		return false
	}
	for _, c := range s {
		switch {
		case strings.ContainsRune("!#$%&()*+-./:<=>?@[]^_{|}~ ", c):
			// Backslash and quote chars are reserved, but
			// otherwise any punctuation chars are allowed
			// in a tag name.
		default:
			if !unicode.IsLetter(c) && !unicode.IsDigit(c) {
				return false
			}
		}
	}
	return true
}
