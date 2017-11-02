package syslog5424 // import "github.com/nathanaelle/syslog5424"

import (
	"./sdata"
)

func GenericSD(i sdata.SDIDLight) sdata.StructuredData {
	return sdata.GenericSD(i)
}

func emptyListSD() sdata.List {
	return sdata.EmptyList()
}
