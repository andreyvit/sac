package sac

import (
	"runtime"
	"sync"
	"testing"
)

// generate creates inputs data.
func generate() []item {
	var Inputs = []item{
		item{"key1igttgrthghsdfhdfhdfh", "val1"}, item{"keyhdsfhdfhdfbfhdhdf2", "val2"}, item{"keyhsdfhdsfhdfhdfsh3", "val3"}, item{"keysdfhdfhdfhdfhb4", "val4"},
		item{"key5cvcvbcxvncvxncxvncvn", "val5"}, item{"kexcvncvncxvnvcxny6", "val6"}, item{"kcxvncvncvnney7", "val7"}, item{"vvvnxcvncnvxcvnkey8", "val8"},
		item{"knxcvncvney9", "val9"}, item{"aabvbvcbbbbbbbbbbbbcxbcvbknaaey10", "val10"}, item{"cxvbcvbcvbcvb11", "11"}, item{"212151858555vbvcxb2xkey12", "val12"},
		item{"xcvbcvbbbbbbbbasdfgsagkey13", "val13"}, item{"asdgsadgasdgsdgkey14", "val14"}, item{"hthsdhdfhsdfhey15", "val15"}, item{"eeeeeeegadsfgasgkey159", "val16"},
		item{"cvbbbbbbbbasdfgsag1", "valeur1"}, item{"cvbbbbbbbbasdfgsag[3]int{1, 2, 3}", "4"}, item{"rfwerfwefweftest", "result"}, item{"nope", "val4"},
		item{"cvbbbbbbbbasdfgsagkey25", "val25"}, item{"kcvbbbbbbbbasdfgsagey26", "val26"}, item{"key27", "val27"}, item{"key28", "val28"},
		item{"keyjj35", "val35"}, item{"key36", "val36"}, item{"keycvbbbbbbbbasdfgsag37", "val37"}, item{"key38", "val38"},
		item{"cvbbbbbbbbasdfgsagkey45", "val45"}, item{"key46", "val46"}, item{"key47", "val47"}, item{"key48", "val48"},
		item{"keydddfdfdfewewgghujnm91", "val1"}, item{"cvbbbbbbbbasdfgsagkey92", "vcvbbbbbbbbasdfgsagal2"}, item{"keycvbbbbbbbbasdfgsag93", "val3"}, item{"9key4", "val4"},
		item{"key5cvbbbbbbbbasdfgsag8", "val5"}, item{"keyyydhsdhjssjgy86", "val6"}, item{"ke8y7", "val7"}, item{"k88ey8", "val8"},
		item{"key9cvbbbbbbbbasdfgsag", "val9"}, item{"keysfjgsfjgsfjg810", "val10"}, item{"cvbbbbbbbbasdfgsag1881", "11"}, item{"key1cvbbbbbbbbasdfgsag82", "val12"},
		item{"8keyjj13", "val13"}, item{"key814", "val14"}, item{"ke8y15", "val15"}, item{"key81cvbbbbbbbbasdfgsag59", "val16"},
		item{"81cvbbbbbbbbasdfgsag", "valeur1"}, item{"[3]i8nt{1, 2, 3}", "4"}, item{"8y45y45ysdfyw5y45wtest", "result"}, item{"n8ocvbbbbbbbbasdfgsagpe", "val4"},
		item{"8kevfhnfghfghfgjhy25", "val25"}, item{"cvbbbbbbbbasdfgsagkey826", "val26"}, item{"key287", "val27"}, item{"keycvbbbbbbbbasdfgsag828", "val28"},
		item{"8kbmbmey35", "val35"}, item{"key836", "val36"}, item{"key387", "val37"}, item{"ke8yjjjjjjjjdshsrttyw455638", "val38"},
		item{"8kecvbbbbbbbbasdfgsagy45", "val45"}, item{"kehthy846", "val46"}, item{"cvbbbbbbbbasdfgsagkey887", "val47"}, item{"ke8ttty48", "val48"},
	}
	return Inputs
}

// TESTS

func TestSet(t *testing.T) {
	var PL *sync.Pool = Pool()
	Inputs := generate()

	// For each input
	for _, in := range Inputs {
		// Fill the sac and retrieve from it for each input.
		Sac := New(PL)
		Sac.Put(in.Key, in.Value)
		v, err := Sac.Get(in.Key)
		// Test
		if err != nil {
			t.Log(err.Error())
			t.Errorf("Get() and Put() (a) : failed for Key %v Val %v", in.Key, in.Value)
		}
		if v != in.Value {
			t.Errorf("Get() and Put() (a) : Wanted %v but got %v", v, in.Value)
		}
	}

	// For all inputs
	Sac := New(PL)
	for l, in := range Inputs {
		// Fill the sac
		Sac.Put(in.Key, in.Value)

		// Test Length()
		slen := Sac.Length()
		if slen != l+1 {
			t.Errorf("Length() : Wanted %v but got %v", l+1, slen)
		}
	}
	for _, in := range Inputs {
		// Retrieve from the sac
		_, err := Sac.Get(in.Key)
		if err != nil {
			t.Logf("ERROR on %v : "+err.Error(), in.Key)
		}
	}
}

// TestSetsamekey test what happens when we insert a new value
// for an already entered key. The old value should be overwritten.
func TestSetsamekey(t *testing.T) {
	var PL *sync.Pool = Pool()
	var Inputs = []item{item{"key1", "val1"}, item{"key1", "val2"}}
	Sac := New(PL)

	for _, in := range Inputs {
		Sac.Put(in.Key, in.Value)

		// Test Length()
		slen := Sac.Length()
		if slen != 1 {
			t.Errorf("Length() : Wanted %v but got %v", 1, slen)
		}
	}
	v, err := Sac.Get(Inputs[0].Key)
	if err != nil {
		t.Logf("ERROR on %v : "+err.Error(), Inputs[0].Key)
	}
	if v != Inputs[1].Value {
		t.Errorf("Get() and Put() (b) : Wanted %v but got %v", Inputs[1].Value, v)
	}
}

// TestDelete tests the Delete() method of a Sac.
func TestDelete(t *testing.T) {
	var PL *sync.Pool = Pool()
	Inputs := generate()
	Sac := New(PL)

	for _, in := range Inputs {
		Sac.Put(in.Key, in.Value)
	}
	// Let's delete 12 items starting from "key6" included
	for _, in := range Inputs[5:17] {
		Sac.Delete(in.Key)
	}

	// Let's try to Get() the previously items from the sac. Expecting error "NOTFOUND"
	for _, in := range Inputs[5:17] {
		v, err := Sac.Get(in.Key)
		if err == nil {
			t.Errorf("Expected nothing but got %v", v)
		}
	}
	// Test Length()
	slen := Sac.Length()
	wlen := len(Inputs) - len(Inputs[5:17])
	if slen != wlen {
		t.Errorf("Length() : Wanted %v but got %v", wlen, slen)
	}
}

// TestClear tests the Clear() method of a Sac.
func TestClear(t *testing.T) {
	var PL *sync.Pool = Pool()
	Inputs := generate()
	Sac := New(PL)

	for _, in := range Inputs {
		Sac.Put(in.Key, in.Value)
	}
	// Let's Clear all
	Sac.Clear()

	// Let's try to Get() the previously items from the sac. Expecting error "NOTFOUND"
	for _, in := range Inputs {
		v, err := Sac.Get(in.Key)
		if err == nil {
			t.Errorf("Expecting nil but got %v", v)
		}
	}
	// Test Length()
	slen := Sac.Length()
	wlen := 0
	if slen != wlen {
		t.Errorf("Length() : Wanted %v but got %v", wlen, slen)
	}
}

/*

	BENCHMARKS

*/

// Thread safe hashmap  for comparison
// Hashmap is a wrapper for map[interface{}]interface{}
// it is RWMutex protected
type HashMap struct {
	mutex sync.RWMutex
	store map[interface{}]interface{}
}

func newHMap() *HashMap {
	h := new(HashMap)
	h.store = make(map[interface{}]interface{})
	return h
}

func (h *HashMap) Get(k interface{}) interface{} {
	h.mutex.RLock()
	defer h.mutex.RUnlock()
	return h.store[k]
}

func (h *HashMap) Put(k, v interface{}) {
	h.mutex.Lock()
	defer h.mutex.Unlock()
	h.store[k] = v
}

var Inputs []item
var Hmap *HashMap
var Sac *Instance

func init() {
	Inputs = generate()

	Sac = New(Pool())
	Hmap = newHMap()

	for l := 0; l < len(Inputs); l++ {
		Sac.Put(Inputs[l].Key, Inputs[l].Value)
		Hmap.Put(Inputs[l].Key, Inputs[l].Value)
	}
}

// Helper for parallel benchmarks
func MaxParallelism() int {
	maxProcs := runtime.GOMAXPROCS(0)
	numCPU := runtime.NumCPU()
	if maxProcs < numCPU {
		return maxProcs
	}
	return numCPU
}

// SET --

var result interface{}

func benchmarkSet_Sac(i int, b *testing.B) {
	b.ReportAllocs()
	b.SetParallelism(MaxParallelism())

	Inputs := generate()
	Sac := New(Pool())

	for l := 0; l < i-1; l++ {
		Sac.Put(Inputs[l].Key, Inputs[l].Value)
	}
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			Sac.Put(Inputs[i-1].Key, Inputs[i-1].Value)
		}
	})
}

func benchmarkSet_StdHMap(i int, b *testing.B) {
	b.ReportAllocs()
	b.SetParallelism(MaxParallelism())

	Hmap := newHMap()
	Inputs := generate()

	for l := 0; l < i-1; l++ {
		Hmap.Put(Inputs[l].Key, Inputs[l].Value)
	}
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			Hmap.Put(Inputs[i-1].Key, Inputs[i-1].Value)
		}
	})
}

// GET --

func benchmarkGet_Sac(i int, b *testing.B) {
	b.ReportAllocs()
	b.SetParallelism(MaxParallelism())

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			v, _ := Sac.Get(Inputs[i-1].Key)
			result = v
		}
	})
}

func benchmarkGet_StdHMap(i int, b *testing.B) {
	b.ReportAllocs()
	b.SetParallelism(MaxParallelism())

	Hmap := newHMap()
	Inputs := generate()

	for k := 0; k < b.N; k++ {
		for l := 0; l < i; l++ {
			Hmap.Put(Inputs[l].Key, Inputs[l].Value)
		}
	}
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			v := Hmap.Get(Inputs[i-1].Key)
			result = v
		}
	})
}

// Benchmark list below based on different values for the number of items

// Setter benchmark for Small Array Containers versus Standard Hashmap
func BenchmarkSetsac1(b *testing.B)  { benchmarkSet_Sac(1, b) }
func BenchmarkSethmap1(b *testing.B) { benchmarkSet_StdHMap(1, b) }

func BenchmarkSetsac5(b *testing.B)  { benchmarkSet_Sac(5, b) }
func BenchmarkSethmap5(b *testing.B) { benchmarkSet_StdHMap(5, b) }

func BenchmarkSetsac6(b *testing.B)  { benchmarkSet_Sac(6, b) }
func BenchmarkSethmap6(b *testing.B) { benchmarkSet_StdHMap(6, b) }

func BenchmarkSetsac7(b *testing.B)  { benchmarkSet_Sac(7, b) }
func BenchmarkSethmap7(b *testing.B) { benchmarkSet_StdHMap(7, b) }

func BenchmarkSetsac8(b *testing.B)  { benchmarkSet_Sac(8, b) }
func BenchmarkSethmap8(b *testing.B) { benchmarkSet_StdHMap(8, b) }

func BenchmarkSetsac10(b *testing.B)  { benchmarkSet_Sac(10, b) }
func BenchmarkSethmap10(b *testing.B) { benchmarkSet_StdHMap(10, b) }

func BenchmarkSetsac11(b *testing.B)  { benchmarkSet_Sac(11, b) }
func BenchmarkSethmap11(b *testing.B) { benchmarkSet_StdHMap(11, b) }

func BenchmarkSetsac16(b *testing.B)  { benchmarkSet_Sac(16, b) }
func BenchmarkSethmap16(b *testing.B) { benchmarkSet_StdHMap(16, b) }

func BenchmarkSetsac32(b *testing.B)  { benchmarkSet_Sac(32, b) }
func BenchmarkSethmap32(b *testing.B) { benchmarkSet_StdHMap(32, b) }

// Getter benchmark for Small Array Containers versus Standard Hashmap

func BenchmarkGetsac1(b *testing.B)  { benchmarkGet_Sac(1, b) }
func BenchmarkGethmap1(b *testing.B) { benchmarkGet_StdHMap(1, b) }

func BenchmarkGetsac5(b *testing.B)  { benchmarkGet_Sac(5, b) }
func BenchmarkGethmap5(b *testing.B) { benchmarkGet_StdHMap(5, b) }

func BenchmarkGetsac8(b *testing.B)  { benchmarkGet_Sac(8, b) }
func BenchmarkGethmap8(b *testing.B) { benchmarkGet_StdHMap(8, b) }

func BenchmarkGetsac10(b *testing.B)  { benchmarkGet_Sac(10, b) }
func BenchmarkGethmap10(b *testing.B) { benchmarkGet_StdHMap(10, b) }

func BenchmarkGetsac11(b *testing.B)  { benchmarkGet_Sac(11, b) }
func BenchmarkGethmap11(b *testing.B) { benchmarkGet_StdHMap(11, b) }

func BenchmarkGetsac16(b *testing.B)  { benchmarkGet_Sac(16, b) }
func BenchmarkGethmap16(b *testing.B) { benchmarkGet_StdHMap(16, b) }

func BenchmarkGetsac32(b *testing.B)  { benchmarkGet_Sac(32, b) }
func BenchmarkGethmap32(b *testing.B) { benchmarkGet_StdHMap(32, b) }

/* Benchmark results on Dual Core Intel Pentium CPU B950 2.1GHz
   6GB DDR3 RAM
   Ubuntu 14.04
   Size of Sac is 32. ()

// SET

BenchmarkSetsac1	 5000000	       264 ns/op	       0 B/op	       0 allocs/op
BenchmarkSethmap1	 5000000	       383 ns/op	       0 B/op	       0 allocs/op

BenchmarkSetsac5	 5000000	       336 ns/op	       0 B/op	       0 allocs/op
BenchmarkSethmap5	 5000000	       399 ns/op	       0 B/op	       0 allocs/op

BenchmarkSetsac6	 5000000	       347 ns/op	       0 B/op	       0 allocs/op
BenchmarkSethmap6	 5000000	       389 ns/op	       0 B/op	       0 allocs/op

BenchmarkSetsac7	 5000000	       360 ns/op	       0 B/op	       0 allocs/op
BenchmarkSethmap7	 5000000	       384 ns/op	       0 B/op	       0 allocs/op

BenchmarkSetsac8	 5000000	       388 ns/op	       0 B/op	       0 allocs/op
BenchmarkSethmap8	 5000000	       397 ns/op	       0 B/op	       0 allocs/op

BenchmarkSetsac10	 3000000	       424 ns/op	       0 B/op	       0 allocs/op
BenchmarkSethmap10	 3000000	       425 ns/op	       0 B/op	       0 allocs/op

BenchmarkSetsac11	 3000000	       437 ns/op	       0 B/op	       0 allocs/op
BenchmarkSethmap11	 5000000	       380 ns/op	       0 B/op	       0 allocs/op

BenchmarkSetsac16	 3000000	       528 ns/op	       0 B/op	       0 allocs/op
BenchmarkSethmap16	 5000000	       388 ns/op	       0 B/op	       0 allocs/op

BenchmarkSetsac32	 2000000	       870 ns/op	       0 B/op	       0 allocs/op
BenchmarkSethmap32	 5000000	       349 ns/op	       0 B/op	       0 allocs/op




// GET

BenchmarkGetsac1	10000000	       231 ns/op	       0 B/op	       0 allocs/op
BenchmarkGethmap1	 5000000	       334 ns/op	       0 B/op	       0 allocs/op

BenchmarkGetsac5	 5000000	       300 ns/op	       0 B/op	       0 allocs/op
BenchmarkGethmap5	 5000000	       346 ns/op	       0 B/op	       0 allocs/op

BenchmarkGetsac8	 5000000	       350 ns/op	       0 B/op	       0 allocs/op
BenchmarkGethmap8	 5000000	       336 ns/op	       0 B/op	       0 allocs/op

BenchmarkGetsac10	 5000000	       392 ns/op	       0 B/op	       0 allocs/op
BenchmarkGethmap10	 5000000	       368 ns/op	       0 B/op	       0 allocs/op

BenchmarkGetsac11	 5000000	       400 ns/op	       0 B/op	       0 allocs/op
BenchmarkGethmap11	 5000000	       323 ns/op	       0 B/op	       0 allocs/op

BenchmarkGetsac16	 3000000	       500 ns/op	       0 B/op	       0 allocs/op
BenchmarkGethmap16	 5000000	       336 ns/op	       0 B/op	       0 allocs/op

BenchmarkGetsac32	 2000000	       847 ns/op	       0 B/op	       0 allocs/op
BenchmarkGethmap32	 5000000	       297 ns/op	       0 B/op	       0 allocs/op

*/
