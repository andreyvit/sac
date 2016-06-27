// Package sac defines `small array containers`, an implementation
// of an unrolled linked-list that is safe for concurrent use.
//
// It is useful when one needs a small, recyclable Key/Value store
// that does not constrain the types that can be used as keys.
// In comparison, a standard map[interface{}]interface{} is
// susceptible to constraints on the types it accepts due to
// its reliance on hashing.
//
// It is not supposed to hold a large number of objects.
// Typically up to ~16 will compete with a RWMutex protected regular map.
// It is still possible to group multiple values under a single key to save
// space, however.
package sac

import (
	"errors"
	"sync"
)

// TODO (tay) : create generator to parameterize Key/Value types.

// NUM is the optimal sac size in synthetic tests.
const NUM = 32

// for friendliness on both 32bit and 64bit
const padOffset = 2 - int(2*(^uintptr(0)>>63<<1))

var (
	// ErrNOTFOUND is returned when no value could be found for a given key.
	ErrNOTFOUND = errors.New("NOTFOUND")
)

// Instance defines the sac datatype.
// It is an opaque datastructure to the user.
// Internally, it is simply a linked list of same-NUM arrays.
type Instance struct {
	items  [NUM]item
	pool   *sync.Pool
	mutex  *sync.RWMutex
	length int
	next   *Instance
	pad    [padOffset + (NUM<<1+4)&7]uint64
}

type item struct {
	Key   interface{}
	Value interface{}
}

// New returns a new sac.
func New(p *sync.Pool) *Instance {
	res := new(Instance)
	res.mutex = new(sync.RWMutex)
	res.pool = p
	return res
}

// Pool initializes a new empty sac Pool ready for use.
func Pool() *sync.Pool {
	v := new(sync.Pool)
	v.New = func() interface{} {
		return New(v)
	}
	return v
}

// Get retrieves an item from the sac.
// Safe for concurrent use by multiple goroutines.
func (i *Instance) Get(key interface{}) (interface{}, error) {
	i.mutex.RLock()

	for k := 0; k < i.length; k++ {
		if i.items[k].Key == key {
			i.mutex.RUnlock()
			return i.items[k].Value, nil
		}
	}
	if i.length < NUM {
		i.mutex.RUnlock()
		return nil, ErrNOTFOUND
	}
	if i.next == nil {
		i.mutex.RUnlock()
		return nil, ErrNOTFOUND
	}
	i.mutex.RUnlock()
	return i.next.Get(key)
}

// Put puts an item into the sac.
// Safe for concurrent use by multiple goroutines.
func (i *Instance) Put(key, value interface{}) {
	i.mutex.Lock()

	for k := 0; k < i.length; k++ {
		if i.items[k].Key == key {
			i.items[k].Value = value
			i.mutex.Unlock()
			return
		}
	}
	if i.length == NUM {
		if i.next == nil {
			i.next = i.pool.Get().(*Instance)
		}
		i.next.Put(key, value)
		i.mutex.Unlock()
		return
	}
	i.items[i.length] = item{key, value}
	i.length++
	i.mutex.Unlock()
	return
}

// Delete removes an item from the sac.
// Safe for concurrent use by multiple goroutines.
func (i *Instance) Delete(key interface{}) {
	i.mutex.Lock()

	lgt := i.length

	if i.next == nil {
		for k := 0; k < lgt; k++ {
			if i.items[k].Key == key {
				for l := k + 1; l < lgt; l++ {
					i.items[l-1].Value = i.items[l].Value
					i.items[l-1].Key = i.items[l].Key
				}
				i.length--
				i.mutex.Unlock()
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
			i.mutex.Unlock()
			return
		}
	}
	i.mutex.Unlock()
	return
}

// Length will return the total number of element in a sac.
// Safe for concurrent use by multiple goroutines.
func (i *Instance) Length() int {
	i.mutex.Lock()
	c := i
	var l int
	for c != nil {
		l += c.length
		c = c.next
	}
	i.mutex.Unlock()
	return l
}

// Clear will get rid of every value in a sac.
// Safe for concurrent use by multiple goroutines.
func (i *Instance) Clear() {
	i.mutex.Lock()

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
		for k := NUM - 1; k >= 0; k-- {
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
	i.mutex.Unlock()
	return
}

// Clone will make and return a copy of a sac.
func (i *Instance) Clone() *Instance {
	i.mutex.RLock()
	vj := *i
	vj.mutex = new(sync.RWMutex)
	if vj.next != nil {
		vj.next = vj.next.Clone()
	}
	i.mutex.Unlock()
	return &vj
}
