package settings

import (
	"fmt"
	"sync"
)

type ActiveLeader interface {
	Get() bool
	Set(val bool)
}

type _actve_leader struct {
	_ok    bool
	_mutex sync.Mutex
}

func (a *_actve_leader) Get() bool {
	a._mutex.Lock()
	defer a._mutex.Unlock()
	return a._ok
}

func (a *_actve_leader) Set(val bool) {
	a._mutex.Lock()
	defer a._mutex.Unlock()
	a._ok = val
}

func logo_print() {
	L.Warn("###############################")
	L.Warn(fmt.Sprintf("##  %s [%s]        ##", promoter, version))
	L.Warn("###############################")
}
