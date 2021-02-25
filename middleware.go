package requests

import "sync"

type reqMiddleware struct {
	index      int
	middleware map[int]func(req *SRequest)
	lock       sync.Mutex
}

func (rm *reqMiddleware) Add(f func(req *SRequest)) (id int) {
	if f == nil {
		return
	}
	rm.lock.Lock()
	defer rm.lock.Unlock()
	if rm.middleware == nil {
		rm.Clear()
	}
	rm.index++
	rm.middleware[rm.index] = f
	return rm.index
}

func (rm *reqMiddleware) Clear() {
	rm.middleware = make(map[int]func(req *SRequest))
}

func (rm *reqMiddleware) Remove(id int) {
	if rm.middleware == nil {
		return
	}
	rm.lock.Lock()
	defer rm.lock.Unlock()
	delete(rm.middleware, id)
}

type respMiddleware struct {
	index      int
	middleware map[int]func(resp *SResponse)
	lock       sync.Mutex
}

func (rm *respMiddleware) Add(f func(resp *SResponse)) (id int) {
	if f == nil {
		return
	}
	rm.lock.Lock()
	defer rm.lock.Unlock()
	if rm.middleware == nil {
		rm.Clear()
	}
	rm.index++
	rm.middleware[rm.index] = f
	return rm.index
}

func (rm *respMiddleware) Clear() {
	rm.middleware = make(map[int]func(resp *SResponse))
}

func (rm *respMiddleware) Remove(id int) {
	if rm.middleware == nil {
		return
	}
	rm.lock.Lock()
	defer rm.lock.Unlock()
	delete(rm.middleware, id)
}
