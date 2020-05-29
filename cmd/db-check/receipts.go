package main

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"math"
	"reflect"
	"unsafe"

	"github.com/Fantom-foundation/go-lachesis/app"
	"github.com/Fantom-foundation/go-lachesis/inter/idx"
	"github.com/Fantom-foundation/go-lachesis/kvdb"
	"github.com/Fantom-foundation/go-lachesis/kvdb/flushable"
	"github.com/Fantom-foundation/go-lachesis/poset"
)

func checkReceipts(db kvdb.DbProducer) {

	dbs := flushable.NewSyncedPool(db)

	cfg := app.DefaultStoreConfig()
	cfg.ReceiptsCacheSize = 150000
	s := app.NewStore(dbs, cfg)
	defer s.Close()

	p := poset.NewStore(dbs, poset.DefaultStoreConfig())
	defer p.Close()

	//g := gossip.NewStore(dbs, gossip.DefaultStoreConfig())
	//defer g.Close()

	lastBlock := p.GetCheckpoint().LastBlockN

	/*
		var allocated uint64
		runtime.GC()
		var m0, m1 runtime.MemStats
		runtime.ReadMemStats(&m0)
	*/
	for n := idx.Block(0); n <= lastBlock; n++ {
		//for n := lastBlock; n > idx.Block(0); n-- {
		rr := s.GetReceipts(n)
		if rr == nil {
			continue
		}

		/*
			runtime.ReadMemStats(&m1)
			if m1.Alloc > m0.Alloc {
				allocated = (m1.Alloc - m0.Alloc) / 1000000
			} else {
				allocated = 0
			}
		*/
		size0 := memSizeOf(rr)
		size1 := DeepSize(rr)
		if len(rr) < 20 && size0 < 8000 && size1 < 8000 {
			continue
		}
		fmt.Printf("%d [%d] - \t%d\t%d\n", n, len(rr), size0, size1)
	}

}

func memSizeOf(v interface{}) int {
	b := new(bytes.Buffer)
	if err := gob.NewEncoder(b).Encode(v); err != nil {
		panic(err)
	}
	return b.Len()
}

func DeepSize(v interface{}) int64 {
	return int64(valueSize(reflect.ValueOf(v), make(map[uintptr]bool)))
}

func valueSize(v reflect.Value, seen map[uintptr]bool) uintptr {
	base := v.Type().Size()
	switch v.Kind() {
	case reflect.Ptr:
		p := v.Pointer()
		if !seen[p] && !v.IsNil() {
			seen[p] = true
			base += valueSize(v.Elem(), seen)
		}

	case reflect.Slice:
		n := v.Len()
		for i := 0; i < n; i++ {
			base += valueSize(v.Index(i), seen)
		}

		// Account for the parts of the array not covered by this slice.  Since
		// we can't get the values directly, assume they're zeroes. That may be
		// incorrect, in which case we may underestimate.
		if cap := v.Cap(); cap > n {
			base += v.Type().Size() * uintptr(cap-n)
		}

	case reflect.Map:
		// A map m has len(m) / 6.5 buckets, rounded up to a power of two, and
		// a minimum of one bucket. Each bucket is 16 bytes + 8*(keysize + valsize).
		//
		// We can't tell which keys are in which bucket by reflection, however,
		// so here we count the 16-byte header for each bucket, and then just add
		// in the computed key and value sizes.
		// TODO: use runtime.MapBucketsCount()
		nb := uintptr(math.Pow(2, math.Ceil(math.Log(float64(v.Len())/6.5)/math.Log(2))))
		if nb == 0 {
			nb = 1
		}
		base = 16 * nb
		for _, key := range v.MapKeys() {
			base += valueSize(key, seen)
			base += valueSize(v.MapIndex(key), seen)
		}

		// We have nb buckets of 8 slots each, and v.Len() slots are filled.
		// The remaining slots we will assume contain zero key/value pairs.
		zk := v.Type().Key().Size()  // a zero key
		zv := v.Type().Elem().Size() // a zero value
		base += (8*nb - uintptr(v.Len())) * (zk + zv)

	case reflect.Struct:
		// Chase pointer and slice fields and add the size of their members.
		for i := 0; i < v.NumField(); i++ {
			f := v.Field(i)
			base += valueSize(f, seen)
		}

	case reflect.String:
		base += uintptr(v.Len())

	default:
		base += unsafe.Sizeof(v)

	}
	return base
}
