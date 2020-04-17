package flushable

import (
	"fmt"
	"io/ioutil"
	"math/rand"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/Fantom-foundation/go-lachesis/common/bigendian"
	"github.com/Fantom-foundation/go-lachesis/kvdb/leveldb"
	"github.com/Fantom-foundation/go-lachesis/kvdb/table"
)

func TestFlushableParallelXXX(t *testing.T) {
	testPairsNum := uint64(1000)

	dir, err := ioutil.TempDir("", "test-flushable")
	if err != nil {
		panic(fmt.Sprintf("can't create temporary directory %s: %v", dir, err))
	}
	disk := leveldb.NewProducer(dir)

	// open raw databases
	ldb := disk.OpenDb("1")
	defer ldb.Drop()
	defer ldb.Close()

	flushableDb := Wrap(ldb)

	tableMutable1 := table.New(flushableDb, []byte("1"))
	tableImmutable := table.New(flushableDb, []byte("2"))
	tableMutable2 := table.New(flushableDb, []byte("3"))

	// fill data
	for i := uint64(0); i < testPairsNum; i++ {
		_ = tableImmutable.Put(bigendian.Int64ToBytes(i), bigendian.Int64ToBytes(i))
		if i == testPairsNum/2 { // a half of data is flushed, other half isn't
			_ = flushableDb.Flush()
		}
	}

	r := rand.New(rand.NewSource(0))
	assertar := assert.New(t)

	// iterate over tableImmutable and check its content
	it := tableImmutable.NewIterator()
	defer it.Release()
	i := uint64(0)
	for ; it.Next(); i++ {
		assertar.Equal(bigendian.Int64ToBytes(i), it.Key(), i)
		assertar.Equal(bigendian.Int64ToBytes(i), it.Value(), i)
		assertar.NoError(it.Error(), i)

		_ = tableMutable1.Put(bigendian.Int64ToBytes(r.Uint64()%testPairsNum), bigendian.Int64ToBytes(r.Uint64()))
		_ = tableMutable2.Put(bigendian.Int64ToBytes(r.Uint64() % testPairsNum)[:7], bigendian.Int64ToBytes(r.Uint64()))
		// if i == 0x15f-1 {
		_ = flushableDb.Flush()
		//}
	}
	assertar.Equal(testPairsNum, i)
}

func TestFlushableParallel(t *testing.T) {
	testDuration := 2 * time.Second
	testPairsNum := uint64(1000)

	dir, err := ioutil.TempDir("", "test-flushable")
	if err != nil {
		panic(fmt.Sprintf("can't create temporary directory %s: %v", dir, err))
	}
	disk := leveldb.NewProducer(dir)

	// open raw databases
	ldb := disk.OpenDb("1")
	defer ldb.Drop()
	defer ldb.Close()

	flushableDb := Wrap(ldb)

	tableMutable1 := table.New(flushableDb, []byte("1"))
	tableImmutable := table.New(flushableDb, []byte("2"))
	tableMutable2 := table.New(flushableDb, []byte("3"))

	// fill data
	for i := uint64(0); i < testPairsNum; i++ {
		_ = tableImmutable.Put(bigendian.Int64ToBytes(i), bigendian.Int64ToBytes(i))
		if i == testPairsNum/2 { // a half of data is flushed, other half isn't
			_ = flushableDb.Flush()
		}
	}

	done := make(chan struct{})
	wg := sync.WaitGroup{}
	wg.Add(2)

	go func() {
		defer wg.Done()
		assertar := assert.New(t)
		cicle := 0
		for {
			select {
			case <-done:
				return
			default:
			}
			// iterate over tableImmutable and check its content
			it := tableImmutable.NewIterator()
			// defer it.Release()
			i := uint64(0)
			for ; it.Next(); i++ {
				assertar.Equal(bigendian.Int64ToBytes(i), it.Key(), i)
				assertar.Equal(bigendian.Int64ToBytes(i), it.Value(), i)
				assertar.NoError(it.Error(), i)
			}
			it.Release()
			assertar.Equalf(testPairsNum, i, "cicle %d", cicle)
			cicle++
		}
	}()

	go func() {
		defer wg.Done()
		r := rand.New(rand.NewSource(0))
		for {
			select {
			case <-done:
				return
			default:
			}
			// try to spoil data in tableImmutable by updating other tables
			_ = tableMutable1.Put(bigendian.Int64ToBytes(r.Uint64()%testPairsNum), bigendian.Int64ToBytes(r.Uint64()))
			_ = tableMutable2.Put(bigendian.Int64ToBytes(r.Uint64() % testPairsNum)[:7], bigendian.Int64ToBytes(r.Uint64()))
			if r.Int63n(10) == 0 {
				_ = flushableDb.Flush() // flush with 1% chance
			}
		}
	}()

	time.Sleep(testDuration)
	close(done)
	wg.Wait()
}
