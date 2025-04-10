package store

import (
	"bytes"
	"fmt"
	"io"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGeneratePath(t *testing.T) {
	key := "testkeystring"
	path := GenerateCASPath(key)
	assert.Equal(t, path.dirPath, "590ef/df4d8/bda68/c12e0/daf97/9cbaf/54515/dbf96")
	assert.Equal(t, path.fileName, "590efdf4d8bda68c12e0daf979cbaf54515dbf96")
}

func TestStoreDelete(t *testing.T) {
	key := "newfiletotest"
	data := "this is some data"
	store := newStore()
	defer teardown(t, store)
	err := store.Write(key, bytes.NewReader([]byte(data)))
	if err != nil {
		t.Error(err)
	}
	if err := store.Delete(key); err != nil {
		t.Error(err)
	}
	if store.Has(key) {
		t.Errorf("file not deleted")
	}
}
func TestStoreReadWrite(t *testing.T) {
	store := newStore()
	defer teardown(t, store)

	for i := 0; i < 50; i++ {
		key := fmt.Sprintf("newfiletotest_%d", i)

		data := "this is some data"

		err := store.Write(key, bytes.NewReader([]byte(data)))
		if err != nil {
			t.Error(err)
		}

		r, err := store.Read(key)
		if err != nil {
			t.Error(err)
		}

		f, err := io.ReadAll(r)
		if err != nil {
			t.Error(err)
		}

		if string(f) != data {
			t.Errorf("read fail, expected %s,got %s", data, string(f))
		}
	}
}

func newStore() *Store {
	opts := StoreOptions{
		GeneratePath: GenerateCASPath,
	}
	return NewStore(opts)
}

func teardown(t *testing.T, s *Store) {
	if err := s.Clear(); err != nil {
		t.Error(err)
	}
}
