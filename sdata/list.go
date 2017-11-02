package sdata // import "github.com/nathanaelle/syslog5424/sdata"

import (
	"sync"
)

type (
	List []StructuredData

	registry_t struct {
		lock  *sync.Mutex
		sdids map[string]SDID
	}
)

var (
	registry = &registry_t{
		lock:  &sync.Mutex{},
		sdids: make(map[string]SDID),
	}

	emptyList = List([]StructuredData{})
)

func EmptyList() List {
	return emptyList
}

func (registry_t) String() string {
	return "internal registry"
}

func (registry *registry_t) Found(data []byte) (StructuredData, bool) {
	registry.lock.Lock()
	defer registry.lock.Unlock()

	for _, sdid := range registry.sdids {
		if sd, ok := sdid.Found(data); ok {
			return sd, ok
		}
	}

	return unknown_sdid.Found(data)
}

func Parse(data []byte) (StructuredData, bool) {
	return registry.Found(data)
}

func Register(s SDID) SDID {
	registry.lock.Lock()
	defer registry.lock.Unlock()

	key := s.String()
	if _, ok := registry.sdids[key]; !ok {
		registry.sdids[key] = s
	}
	return s
}

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

// Append data to a list of structured data
func (listSD List) Add(data ...StructuredData) List {

	return List(append([]StructuredData(listSD), data...))
}

func (listSD List) Marshal5424() ([]byte, error) {
	if len(listSD) == 0 {
		return []byte{'-'}, nil
	}

	ret_s := 0
	t_a := make([][]byte, len(listSD))
	for i, sd := range listSD {
		t, err := sd.Marshal5424()
		if err != nil {
			return nil, err
		}
		ret_s += len(t)
		t_a[i] = t
	}

	ret := make([]byte, 0, ret_s)
	for _, sd := range t_a {
		ret = append(ret, sd...)
	}

	return ret, nil
}

func (listSD List) String() (s string) {
	res, _ := listSD.Marshal5424()
	s = string(res)
	return
}

func (p List) MarshalText() ([]byte, error) {
	return p.Marshal5424()
}
