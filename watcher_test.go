package globnotify

import (
	"log"
	"os"
	"path/filepath"
	"reflect"
	"sort"
	"testing"
	"time"

	assert "github.com/pilu/miniassert"
)

func mkdirAll(root, path string) string {
	fullPath := filepath.Join(root, path)
	err := os.MkdirAll(fullPath, os.ModePerm)
	if err != nil {
		log.Fatal(err)
	}

	return fullPath
}

func createFile(path, name string) string {
	fullPath := filepath.Join(path, name)
	file, err := os.Create(fullPath)
	if err != nil {
		log.Fatal(err)
	}
	file.Close()

	return fullPath
}

func TestWatcher(t *testing.T) {
	tmpRoot := filepath.Join(os.TempDir(), "golang.globnotify.tests")
	defer os.RemoveAll(tmpRoot)

	w, err := New(filepath.Join(tmpRoot, "**/*.css"))
	assert.Nil(t, err)
	defer func() {
		w.Close()
	}()

	path1 := mkdirAll(tmpRoot, "foo")
	path2 := mkdirAll(tmpRoot, "foo/bar")
	path3 := mkdirAll(tmpRoot, "foo/bar/baz")

	events, err := w.Watch()
	assert.Nil(t, err)

	createFile(tmpRoot, "foo.txt")
	ok1 := createFile(tmpRoot, "foo.css")

	createFile(path1, "foo.txt")
	ok2 := createFile(path1, "foo.css")

	createFile(path2, "foo.txt")
	ok3 := createFile(path2, "foo.css")

	createFile(path3, "foo.txt")
	ok4 := createFile(path3, "foo.css")

	received := make([]string, 0)

	waiting := true
	for waiting {
		select {
		case e := <-events:
			received = append(received, e.Name)
		case <-time.After(time.Second):
			waiting = false
		}
	}

	expected := []string{ok1, ok2, ok3, ok4}

	doSort := func(s []string) {
		sort.Slice(s, func(i, j int) bool { return s[i] < s[j] })
	}

	doSort(received)
	doSort(expected)

	if !reflect.DeepEqual(expected, received) {
		t.Fail()
		dump := func(s []string) {
			for _, path := range s {
				t.Error(path)
			}
		}
		t.Error("Expected:")
		dump(expected)
		t.Error("Got:")
		dump(received)
	}
}
