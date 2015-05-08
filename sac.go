// sac package defines `small array containers` a data structure of the associative array family.
// It is useful when one needs a very small Key/Value store to transit between handlers of each http request.
// It is not supposed to hold a large number of objects. Typically below ~16.
// It is useful to control the amount of garbage generated.
// TODO (Tay) : think about optional map fallback.
package sac

import (
	"errors"
	"sync"
)

const SIZE = 32 // Optimal size in synthetic tests.

type Instance struct {
	items  [SIZE]item
	pool   *sync.Pool
	mutex  sync.RWMutex
	length int
	next   *Instance
	pad    [(SIZE<<1 + 4) & 7]uint64
}

type item struct {
	Key   interface{}
	Value interface{}
}

// New returns a new sac
func New(p *sync.Pool) *Instance {
	res := new(Instance)
	res.pool = p
	return res
}

// Pool initializes a new empty sac Pool ready for use.
// It is possible to call Get()
func Pool() *sync.Pool {
	v := new(sync.Pool)
	v.New = func() interface{} {
		return New(v)
	}
	return v
}

// Get retrieves an item from the sac.
// The key must be suitable for comparison or else it will panic.
func (i *Instance) Get(key interface{}) (interface{}, error) {
	i.mutex.RLock()
	defer i.mutex.RUnlock()

	for k := 0; k < i.length; k++ {
		if i.items[k].Key == key {
			return i.items[k].Value, nil
		}
	}
	if i.length < SIZE {
		return nil, errors.New("NOTFOUND")
	}
	if i.next == nil {
		return nil, errors.New("NOTFOUND")
	}
	return i.next.Get(key)
}

// Put puts an item into the sac.
// The key must be suitable for comparison or else it will panic.
func (i *Instance) Put(key, value interface{}) {
	i.mutex.Lock()
	defer i.mutex.Unlock()

	for k := 0; k < i.length; k++ {
		if i.items[k].Key == key {
			i.items[k].Value = value
			return
		}
	}
	if i.length == SIZE {
		if i.next == nil {
			i.next = i.pool.Get().(*Instance)
		}
		i.next.Put(key, value)
		return
	}
	i.items[i.length] = item{key, value}
	i.length++
	return
}

// Delete removes an item from the sac.
// The key must be suitable for comparison or else it will panic.
func (i *Instance) Delete(key interface{}) {
	i.mutex.Lock()
	defer i.mutex.Unlock()
	lgt := i.length

	if i.next == nil {
		for k := 0; k < lgt; k++ {
			if i.items[k].Key == key {
				for l := k + 1; l < lgt; l++ {
					i.items[l-1].Value = i.items[l].Value
					i.items[l-1].Key = i.items[l].Key
				}
				i.length--
				return
			}
		}
	}

	if i.next != nil {
		for k := 0; k < lgt; k++ {
			if i.items[k].Key == key {
				for l := k + 1; l < lgt; l++ {
					i.items[l-1].Value = i.items[l].Value
					i.items[l-1].Key = i.items[l].Key
				}
				i.items[lgt-1] = i.next.items[0]
				i.next.Delete(i.items[lgt-1].Key)
			}
		}
		i.next.Delete(key)
		if i.next.length == 0 {
			i.pool.Put(i.next)
			i.next = nil
			return
		}
	}
	return
}

// Length will return the total number of element in a sac.
// A sac can nest several sac Instances.
// Length will add the items from those nested sac instances too.
func (i *Instance) Length() int {
	i.mutex.Lock()
	defer i.mutex.Unlock()
	c := i
	var l int
	for c != nil {
		l += c.length
		c = c.next
	}
	return l
}

// Clear will empty a sac and all its nested sacs.
func (i *Instance) Clear() {
	i.mutex.Lock()
	defer i.mutex.Unlock()

	for i.next != nil {
		if i.next.length == 0 {
			i.pool.Put(i.next)
			i.next = nil
			break
		}
		s := i.next
		for s.next != nil {
			if s.next.length == 0 {
				s.pool.Put(s.next)
				s.next = nil
				break
			}
			s = s.next
			break
		}
		for k := SIZE - 1; k >= 0; k-- {
			s.Delete(s.items[k].Key)
		}

		s.length = 0
		s.pool.Put(s)
	}

	for _, v := range i.items {
		v.Key = nil
		v.Value = nil
	}
	i.length = 0
	return
}
