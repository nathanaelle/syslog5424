package syslog5424 // import "github.com/nathanaelle/syslog5424/v2"

func validHost(d string) string {
	return validatorGeneric([]byte(d), 255)
}

func validApp(d string) string {
	return validatorGeneric([]byte(d), 48)
}

func validProcid(d string) string {
	return validatorGeneric([]byte(d), 128)
}

func validMsgid(d string) string {
	return validatorGeneric([]byte(d), 32)
}

func validatorGeneric(data []byte, maxlen int) string {
	l := len(data)

	if l == 0 {
		return "-"
	}

	if l > maxlen {
		l = maxlen
	}

	i := l
	for i > 0 {
		i--
		if data[i] < 33 || data[i] > 126 {
			data[i] = '%'
		}
	}

	return string(data[0:l])
}
