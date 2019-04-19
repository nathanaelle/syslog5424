package syslog5424 // import "github.com/nathanaelle/syslog5424/v2"

import (
	"github.com/nathanaelle/syslog5424/v2/sdata"
)

// GenericSD is a wrapper around sdata.GenericSD
func GenericSD(i sdata.SDIDLight) sdata.StructuredData {
	return sdata.GenericSD(i)
}

func emptyListSD() sdata.List {
	return sdata.EmptyList()
}
