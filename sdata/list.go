package sdata // import "github.com/nathanaelle/syslog5424/v2/sdata"

import (
	"sync"
)

type (
	// List represent a list of Structured Data
	List []StructuredData

	registryT struct {
		lock  *sync.Mutex
		sdids map[string]SDID
	}
)

var (
	registry = &registryT{
		lock:  &sync.Mutex{},
		sdids: make(map[string]SDID),
	}

	emptyList = List([]StructuredData{})
)

// EmptyList return an empty []StructuredData
func EmptyList() List {
	return emptyList
}

func (registry *registryT) String() string {
	return "internal registry"
}

func (registry *registryT) Found(data []byte) (StructuredData, bool) {
	registry.lock.Lock()
	defer registry.lock.Unlock()

	for _, sdid := range registry.sdids {
		if sd, ok := sdid.Found(data); ok {
			return sd, ok
		}
	}

	return unknownSDID.Found(data)
}

// Parse decode a []byte to a Structured Data
func Parse(data []byte) (StructuredData, bool) {
	return registry.Found(data)
}

// Register add an new SDID to a global registry
func Register(s SDID) SDID {
	registry.lock.Lock()
	defer registry.lock.Unlock()

	key := s.String()
	if _, ok := registry.sdids[key]; !ok {
		registry.sdids[key] = s
	}
	return s
}

// RegisterGroup is a wrapper around Register
func RegisterGroup(sdids ...SDID) {
	registry.lock.Lock()
	defer registry.lock.Unlock()

	for _, s := range sdids {
		key := s.String()
		if _, ok := registry.sdids[key]; !ok {
			registry.sdids[key] = s
		}
	}
}

// Add append data to a list of structured data
func (listSD List) Add(data ...StructuredData) List {

	return List(append([]StructuredData(listSD), data...))
}

// Marshal5424 encode a Structured Data to syslog 5424 format
func (listSD List) Marshal5424() ([]byte, error) {
	if len(listSD) == 0 {
		return []byte{'-'}, nil
	}

	retSize := 0
	tmpArray := make([][]byte, len(listSD))
	for i, sd := range listSD {
		t, err := sd.Marshal5424()
		if err != nil {
			return nil, err
		}
		retSize += len(t)
		tmpArray[i] = t
	}

	ret := make([]byte, 0, retSize)
	for _, sd := range tmpArray {
		ret = append(ret, sd...)
	}

	return ret, nil
}

func (listSD List) String() (s string) {
	res, _ := listSD.Marshal5424()
	s = string(res)
	return
}

// MarshalText implements encoding.TextMarshaller
func (listSD List) MarshalText() ([]byte, error) {
	return listSD.Marshal5424()
}
