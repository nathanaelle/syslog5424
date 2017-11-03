package syslog5424 // import "github.com/nathanaelle/syslog5424"

func valid_host(d string) string {
	return validator_generic([]byte(d), 255)
}

func valid_app(d string) string {
	return validator_generic([]byte(d), 48)
}

func valid_procid(d string) string {
	return validator_generic([]byte(d), 128)
}

func valid_msgid(d string) string {
	return validator_generic([]byte(d), 32)
}

func validator_generic(data []byte, maxlen int) string {
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
