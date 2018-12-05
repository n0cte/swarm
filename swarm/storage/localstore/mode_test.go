// Copyright 2018 The go-ethereum Authors
// This file is part of the go-ethereum library.
//
// The go-ethereum library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The go-ethereum library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the go-ethereum library. If not, see <http://www.gnu.org/licenses/>.

package localstore

import (
	"bytes"
	"context"
	"fmt"
	"sort"
	"testing"
	"time"

	"github.com/syndtr/goleveldb/leveldb"

	"github.com/ethereum/go-ethereum/swarm/shed"
	"github.com/ethereum/go-ethereum/swarm/storage"
)

// TestModeSyncing validates internal data operations and state
// for ModeSyncing on DB with default configuration.
func TestModeSyncing(t *testing.T) {
	db, cleanupFunc := newTestDB(t, nil)
	defer cleanupFunc()

	testModeSyncingValues(t, db)
}

// TestModeSyncing_useRetrievalCompositeIndex validates internal
// data operations and state for ModeSyncing on DB with
// retrieval composite index enabled.
func TestModeSyncing_useRetrievalCompositeIndex(t *testing.T) {
	db, cleanupFunc := newTestDB(t, &Options{UseRetrievalCompositeIndex: true})
	defer cleanupFunc()

	testModeSyncingValues(t, db)
}

// testModeSyncingValues validates ModeSyncing index values on the provided DB.
func testModeSyncingValues(t *testing.T, db *DB) {
	a := db.Accessor(ModeSyncing)

	chunk := generateRandomChunk()

	wantTimestamp := time.Now().UTC().UnixNano()
	defer func(n func() int64) { now = n }(now)
	now = func() (t int64) {
		return wantTimestamp
	}

	wantSize, err := db.sizeCounter.Get()
	if err != nil {
		t.Fatal(err)
	}

	err = a.Put(context.Background(), chunk)
	if err != nil {
		t.Fatal(err)
	}

	wantSize++

	t.Run("retrieve indexes", testRetrieveIndexesValues(db, chunk, wantTimestamp, wantTimestamp))

	t.Run("pull index", testPullIndexValues(db, chunk, wantTimestamp, nil))

	t.Run("size counter", testSizeCounter(db, wantSize))
}

// TestModeUpload validates internal data operations and state
// for ModeUpload on DB with default configuration.
func TestModeUpload(t *testing.T) {
	db, cleanupFunc := newTestDB(t, nil)
	defer cleanupFunc()

	testModeUploadValues(t, db)
}

// TestModeUpload_useRetrievalCompositeIndex validates internal
// data operations and state for ModeUpload on DB with
// retrieval composite index enabled.
func TestModeUpload_useRetrievalCompositeIndex(t *testing.T) {
	db, cleanupFunc := newTestDB(t, &Options{UseRetrievalCompositeIndex: true})
	defer cleanupFunc()

	testModeUploadValues(t, db)
}

// testModeUploadValues validates ModeUpload index values on the provided DB.
func testModeUploadValues(t *testing.T, db *DB) {
	a := db.Accessor(ModeUpload)

	chunk := generateRandomChunk()

	wantTimestamp := time.Now().UTC().UnixNano()
	defer func(n func() int64) { now = n }(now)
	now = func() (t int64) {
		return wantTimestamp
	}

	wantSize, err := db.sizeCounter.Get()
	if err != nil {
		t.Fatal(err)
	}

	err = a.Put(context.Background(), chunk)
	if err != nil {
		t.Fatal(err)
	}

	wantSize++

	t.Run("retrieve indexes", testRetrieveIndexesValues(db, chunk, wantTimestamp, wantTimestamp))

	t.Run("pull index", testPullIndexValues(db, chunk, wantTimestamp, nil))

	t.Run("push index", testPushIndexValues(db, chunk, wantTimestamp, nil))

	t.Run("size counter", testSizeCounter(db, wantSize))
}

// TestModeRequest validates internal data operations and state
// for ModeRequest on DB with default configuration.
func TestModeRequest(t *testing.T) {
	db, cleanupFunc := newTestDB(t, nil)
	defer cleanupFunc()

	testModeRequestValues(t, db)
}

// TestModeRequest_useRetrievalCompositeIndex validates internal
// data operations and state for ModeRequest on DB with
// retrieval composite index enabled.
func TestModeRequest_useRetrievalCompositeIndex(t *testing.T) {
	db, cleanupFunc := newTestDB(t, &Options{UseRetrievalCompositeIndex: true})
	defer cleanupFunc()

	testModeRequestValues(t, db)
}

// testModeRequestValues validates ModeRequest index values on the provided DB.
func testModeRequestValues(t *testing.T, db *DB) {
	a := db.Accessor(ModeRequest)

	chunk := generateRandomChunk()

	wantTimestamp := time.Now().UTC().UnixNano()
	defer func(n func() int64) { now = n }(now)
	now = func() (t int64) {
		return wantTimestamp
	}

	err := a.Put(context.Background(), chunk)
	if err != nil {
		t.Fatal(err)
	}

	t.Run("retrieve indexes", testRetrieveIndexesValuesWithAccess(db, chunk, wantTimestamp, wantTimestamp))

	t.Run("gc index", testGCIndexValues(db, chunk, wantTimestamp, wantTimestamp))
}

// TestModeSynced validates internal data operations and state
// for ModeSynced on DB with default configuration.
func TestModeSynced(t *testing.T) {
	db, cleanupFunc := newTestDB(t, nil)
	defer cleanupFunc()

	testModeSyncedValues(t, db)
}

// TestModeSynced_useRetrievalCompositeIndex validates internal
// data operations and state for ModeSynced on DB with
// retrieval composite index enabled.
func TestModeSynced_useRetrievalCompositeIndex(t *testing.T) {
	db, cleanupFunc := newTestDB(t, &Options{UseRetrievalCompositeIndex: true})
	defer cleanupFunc()

	testModeSyncedValues(t, db)
}

// testModeSyncedValues validates ModeSynced index values on the provided DB.
func testModeSyncedValues(t *testing.T, db *DB) {
	a := db.Accessor(ModeSyncing)

	chunk := generateRandomChunk()

	wantTimestamp := time.Now().UTC().UnixNano()
	defer func(n func() int64) { now = n }(now)
	now = func() (t int64) {
		return wantTimestamp
	}

	err := a.Put(context.Background(), chunk)
	if err != nil {
		t.Fatal(err)
	}

	a = db.Accessor(ModeSynced)

	err = a.Put(context.Background(), chunk)
	if err != nil {
		t.Fatal(err)
	}

	t.Run("retrieve indexes", testRetrieveIndexesValues(db, chunk, wantTimestamp, wantTimestamp))

	t.Run("push index", testPushIndexValues(db, chunk, wantTimestamp, leveldb.ErrNotFound))

	t.Run("gc index", testGCIndexValues(db, chunk, wantTimestamp, wantTimestamp))
}

// TestModeAccess validates internal data operations and state
// for ModeAccess on DB with default configuration.
func TestModeAccess(t *testing.T) {
	db, cleanupFunc := newTestDB(t, nil)
	defer cleanupFunc()

	testModeAccessValues(t, db)
}

// TestModeAccess_useRetrievalCompositeIndex validates internal
// data operations and state for ModeAccess on DB with
// retrieval composite index enabled.
func TestModeAccess_useRetrievalCompositeIndex(t *testing.T) {
	db, cleanupFunc := newTestDB(t, &Options{UseRetrievalCompositeIndex: true})
	defer cleanupFunc()

	testModeAccessValues(t, db)
}

// testModeAccessValues validates ModeAccess index values on the provided DB.
func testModeAccessValues(t *testing.T, db *DB) {
	a := db.Accessor(ModeUpload)

	chunk := generateRandomChunk()

	uploadTimestamp := time.Now().UTC().UnixNano()
	defer func(n func() int64) { now = n }(now)
	now = func() (t int64) {
		return uploadTimestamp
	}

	err := a.Put(context.Background(), chunk)
	if err != nil {
		t.Fatal(err)
	}

	a = db.Accessor(modeAccess)

	t.Run("first get", func(t *testing.T) {
		got, err := a.Get(context.Background(), chunk.Address())
		if err != nil {
			t.Fatal(err)
		}

		if !bytes.Equal(chunk.Address(), got.Address()) {
			t.Errorf("got chunk address %x, want %s", chunk.Address(), got.Address())
		}

		if !bytes.Equal(chunk.Data(), got.Data()) {
			t.Errorf("got chunk data %x, want %s", chunk.Data(), got.Data())
		}

		t.Run("retrieve indexes", testRetrieveIndexesValuesWithAccess(db, chunk, uploadTimestamp, uploadTimestamp))

		t.Run("gc index", testGCIndexValues(db, chunk, uploadTimestamp, uploadTimestamp))

		t.Run("gc index count", testIndexItemsCount(db.gcIndex, 1))
	})

	t.Run("second get", func(t *testing.T) {
		accessTimestamp := time.Now().UTC().UnixNano()
		now = func() (t int64) {
			return accessTimestamp
		}

		got, err := a.Get(context.Background(), chunk.Address())
		if err != nil {
			t.Fatal(err)
		}

		if !bytes.Equal(chunk.Address(), got.Address()) {
			t.Errorf("got chunk address %x, want %s", chunk.Address(), got.Address())
		}

		if !bytes.Equal(chunk.Data(), got.Data()) {
			t.Errorf("got chunk data %x, want %s", chunk.Data(), got.Data())
		}

		t.Run("retrieve indexes", testRetrieveIndexesValuesWithAccess(db, chunk, uploadTimestamp, accessTimestamp))

		t.Run("gc index", testGCIndexValues(db, chunk, uploadTimestamp, accessTimestamp))

		t.Run("gc index count", testIndexItemsCount(db.gcIndex, 1))
	})
}

// TestModeRemoval validates internal data operations and state
// for ModeRemoval on DB with default configuration.
func TestModeRemoval(t *testing.T) {
	db, cleanupFunc := newTestDB(t, nil)
	defer cleanupFunc()

	testModeRemovalValues(t, db)
}

// TestModeRemoval_useRetrievalCompositeIndex validates internal
// data operations and state for ModeRemoval on DB with
// retrieval composite index enabled.
func TestModeRemoval_useRetrievalCompositeIndex(t *testing.T) {
	db, cleanupFunc := newTestDB(t, &Options{UseRetrievalCompositeIndex: true})
	defer cleanupFunc()

	testModeRemovalValues(t, db)
}

// testModeRemovalValues validates ModeRemoval index values on the provided DB.
func testModeRemovalValues(t *testing.T, db *DB) {
	a := db.Accessor(ModeUpload)

	chunk := generateRandomChunk()

	err := a.Put(context.Background(), chunk)
	if err != nil {
		t.Fatal(err)
	}

	a = db.Accessor(modeRemoval)

	wantSize, err := db.sizeCounter.Get()
	if err != nil {
		t.Fatal(err)
	}

	wantSize--

	err = a.Put(context.Background(), chunk)
	if err != nil {
		t.Fatal(err)
	}

	t.Run("retrieve indexes", func(t *testing.T) {
		wantErr := leveldb.ErrNotFound
		if db.useRetrievalCompositeIndex {
			_, err := db.retrievalCompositeIndex.Get(addressToItem(chunk.Address()))
			if err != wantErr {
				t.Errorf("got error %v, want %v", err, wantErr)
			}
			t.Run("retrieve index count", testIndexItemsCount(db.retrievalCompositeIndex, 0))
		} else {
			_, err := db.retrievalDataIndex.Get(addressToItem(chunk.Address()))
			if err != wantErr {
				t.Errorf("got error %v, want %v", err, wantErr)
			}
			t.Run("retrieve data index count", testIndexItemsCount(db.retrievalDataIndex, 0))

			// access index should not be set
			_, err = db.retrievalAccessIndex.Get(addressToItem(chunk.Address()))
			if err != wantErr {
				t.Errorf("got error %v, want %v", err, wantErr)
			}
			t.Run("retrieve access index count", testIndexItemsCount(db.retrievalAccessIndex, 0))
		}
	})

	t.Run("pull index", testPullIndexValues(db, chunk, 0, leveldb.ErrNotFound))

	t.Run("pull index count", testIndexItemsCount(db.pullIndex, 0))

	t.Run("gc index count", testIndexItemsCount(db.gcIndex, 0))

	t.Run("size counter", testSizeCounter(db, wantSize))
}

// TestDB_pullIndex validates the ordering of keys in pull index.
// Pull index key contains PO prefix which is calculated from
// DB base key and chunk address. This is not an IndexItem field
// which are checked in Mode tests.
// This test uploads chunks, sorts them in expected order and
// validates that pull index iterator will iterate it the same
// order.
func TestDB_pullIndex(t *testing.T) {
	db, cleanupFunc := newTestDB(t, nil)
	defer cleanupFunc()

	a := db.Accessor(ModeUpload)

	chunkCount := 50

	// a wrapper around Chunk to keep
	// store timestamp for sorting
	type testChunk struct {
		storage.Chunk
		storeTimestamp int64
	}

	chunks := make([]testChunk, chunkCount)

	// upload random chunks
	for i := 0; i < chunkCount; i++ {
		chunk := generateRandomChunk()

		err := a.Put(context.Background(), chunk)
		if err != nil {
			t.Fatal(err)
		}

		chunks[i] = testChunk{
			Chunk: chunk,
			// this timestamp is not the same as in
			// the index, but given that uploads
			// are sequential and that only ordering
			// of events matter, this information is
			// sufficient
			storeTimestamp: now(),
		}
	}

	// check if all chunks are stored
	testIndexItemsCount(db.pullIndex, chunkCount)

	// sort uploaded chunk is an expected pull index keys order
	// "PO|StoredTimestamp|Hash"
	sort.Slice(chunks, func(i, j int) (less bool) {
		poi := storage.Proximity(db.baseKey, chunks[i].Address())
		poj := storage.Proximity(db.baseKey, chunks[j].Address())
		if poi < poj {
			return true
		}
		if poi > poj {
			return false
		}
		if chunks[i].storeTimestamp < chunks[j].storeTimestamp {
			return true
		}
		if chunks[i].storeTimestamp > chunks[j].storeTimestamp {
			return false
		}
		return bytes.Compare(chunks[i].Address(), chunks[j].Address()) == -1
	})

	// iterate over all items
	var cursor int
	err := db.pullIndex.IterateAll(func(item shed.IndexItem) (stop bool, err error) {
		want := chunks[cursor].Address()
		got := item.Address
		if !bytes.Equal(got, want) {
			return true, fmt.Errorf("got address %x at position %v, want %x", got, cursor, want)
		}
		cursor++
		return false, nil
	})
	if err != nil {
		t.Fatal(err)
	}
}

// testRetrieveIndexesValues returns a test function that validates if the right
// chunk values are in the retrieval indexes.
func testRetrieveIndexesValues(db *DB, chunk storage.Chunk, storeTimestamp, accessTimestamp int64) func(t *testing.T) {
	return func(t *testing.T) {
		if db.useRetrievalCompositeIndex {
			item, err := db.retrievalCompositeIndex.Get(addressToItem(chunk.Address()))
			if err != nil {
				t.Fatal(err)
			}
			validateItem(t, item, chunk.Address(), chunk.Data(), storeTimestamp, accessTimestamp)
		} else {
			item, err := db.retrievalDataIndex.Get(addressToItem(chunk.Address()))
			if err != nil {
				t.Fatal(err)
			}
			validateItem(t, item, chunk.Address(), chunk.Data(), storeTimestamp, 0)

			// access index should not be set
			wantErr := leveldb.ErrNotFound
			item, err = db.retrievalAccessIndex.Get(addressToItem(chunk.Address()))
			if err != wantErr {
				t.Errorf("got error %v, want %v", err, wantErr)
			}
		}
	}
}

// testRetrieveIndexesValuesWithAccess returns a test function that validates if the right
// chunk values are in the retrieval indexes when access time must be stored.
func testRetrieveIndexesValuesWithAccess(db *DB, chunk storage.Chunk, storeTimestamp, accessTimestamp int64) func(t *testing.T) {
	return func(t *testing.T) {
		if db.useRetrievalCompositeIndex {
			item, err := db.retrievalCompositeIndex.Get(addressToItem(chunk.Address()))
			if err != nil {
				t.Fatal(err)
			}
			validateItem(t, item, chunk.Address(), chunk.Data(), storeTimestamp, accessTimestamp)
		} else {
			item, err := db.retrievalDataIndex.Get(addressToItem(chunk.Address()))
			if err != nil {
				t.Fatal(err)
			}
			validateItem(t, item, chunk.Address(), chunk.Data(), storeTimestamp, 0)

			// access index should not be set
			item, err = db.retrievalAccessIndex.Get(addressToItem(chunk.Address()))
			if err != nil {
				t.Fatal(err)
			}
			validateItem(t, item, chunk.Address(), nil, 0, accessTimestamp)
		}
	}
}

// testPullIndexValues returns a test function that validates if the right
// chunk values are in the pull index.
func testPullIndexValues(db *DB, chunk storage.Chunk, storeTimestamp int64, wantError error) func(t *testing.T) {
	return func(t *testing.T) {
		item, err := db.pullIndex.Get(shed.IndexItem{
			Address:        chunk.Address(),
			StoreTimestamp: storeTimestamp,
		})
		if err != wantError {
			t.Errorf("got error %v, want %v", err, wantError)
		}
		if err == nil {
			validateItem(t, item, chunk.Address(), nil, storeTimestamp, 0)
		}
	}
}

// testPushIndexValues returns a test function that validates if the right
// chunk values are in the push index.
func testPushIndexValues(db *DB, chunk storage.Chunk, storeTimestamp int64, wantError error) func(t *testing.T) {
	return func(t *testing.T) {
		item, err := db.pushIndex.Get(shed.IndexItem{
			Address:        chunk.Address(),
			StoreTimestamp: storeTimestamp,
		})
		if err != wantError {
			t.Errorf("got error %v, want %v", err, wantError)
		}
		if err == nil {
			validateItem(t, item, chunk.Address(), nil, storeTimestamp, 0)
		}
	}
}

// testGCIndexValues returns a test function that validates if the right
// chunk values are in the push index.
func testGCIndexValues(db *DB, chunk storage.Chunk, storeTimestamp, accessTimestamp int64) func(t *testing.T) {
	return func(t *testing.T) {
		item, err := db.gcIndex.Get(shed.IndexItem{
			Address:         chunk.Address(),
			StoreTimestamp:  storeTimestamp,
			AccessTimestamp: accessTimestamp,
		})
		if err != nil {
			t.Fatal(err)
		}
		validateItem(t, item, chunk.Address(), nil, storeTimestamp, accessTimestamp)
	}
}

// testIndexItemsCount returns a test function that validates if
// an index contains expected number of key/value pairs.
func testIndexItemsCount(i shed.Index, want int) func(t *testing.T) {
	return func(t *testing.T) {
		var c int
		i.IterateAll(func(item shed.IndexItem) (stop bool, err error) {
			c++
			return
		})
		if c != want {
			t.Errorf("got %v items in index, want %v", c, want)
		}
	}
}

// testSizeCounter returns a test function that validates the expected
// value from sizeCounter field.
func testSizeCounter(db *DB, wantSize uint64) func(t *testing.T) {
	return func(t *testing.T) {
		got, err := db.sizeCounter.Get()
		if err != nil {
			t.Fatal(err)
		}
		if got != wantSize {
			t.Errorf("got size counter value %v, want %v", got, wantSize)
		}
	}
}

// validateItem is a helper function that checks IndexItem values.
func validateItem(t *testing.T, item shed.IndexItem, address, data []byte, storeTimestamp, accessTimestamp int64) {
	t.Helper()

	if !bytes.Equal(item.Address, address) {
		t.Errorf("got item address %x, want %x", item.Address, address)
	}
	if !bytes.Equal(item.Data, data) {
		t.Errorf("got item data %x, want %x", item.Data, data)
	}
	if item.StoreTimestamp != storeTimestamp {
		t.Errorf("got item store timestamp %v, want %v", item.StoreTimestamp, storeTimestamp)
	}
	if item.AccessTimestamp != accessTimestamp {
		t.Errorf("got item access timestamp %v, want %v", item.AccessTimestamp, accessTimestamp)
	}
}
