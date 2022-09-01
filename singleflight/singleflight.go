package singleflight

import "sync"

type callResult struct {
	wg  sync.WaitGroup
	val interface{}
	err error
}

type CallOnce struct {
	sync.Mutex

	onces map[string]*callResult
}

func NewCallOnce() *CallOnce {
	return &CallOnce{
		onces: make(map[string]*callResult),
	}
}

func (v *CallOnce) Call(key string, fn func() (interface{}, error)) (interface{}, error) {
	v.Lock()

	once, ok := v.onces[key]
	if ok {
		v.Unlock()
		once.wg.Wait()
		return once.val, once.err
	}

	cr := &callResult{}
	v.onces[key] = cr
	cr.wg.Add(1)
	v.Unlock()

	cr.val, cr.err = fn()

	v.Lock()
	delete(v.onces, key)
	v.Unlock()

	cr.wg.Done()

	return cr.val, cr.err
}
