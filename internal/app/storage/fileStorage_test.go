package storage

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/tsupko/shortener/internal/app/util"
)

func TestReadFromFileWhenCreated(t *testing.T) {
	hash := util.GenerateUniqueID()

	fileStorage := NewFileStorage("file.log")
	fileStorage.writeToFile(hash, "url")

	url, _ := fileStorage.Get(hash)
	assert.Equal(t, "", url)

	anotherStorage := NewFileStorage("file.log")
	url, _ = anotherStorage.Get(hash)
	assert.Equal(t, "url", url)
}

func TestDoubleSave(t *testing.T) {
	hash := util.GenerateUniqueID()

	fileStorage := NewFileStorage("file.log")
	fileStorage.Put(hash, "url")
	fileStorage.Put(hash, "url2")

	url, _ := fileStorage.Get(hash)
	assert.Equal(t, "url2", url)

	anotherStorage := NewFileStorage("file.log")
	url, _ = anotherStorage.Get(hash)
	assert.Equal(t, "url2", url)
}

func Test(t *testing.T) {
	assert.NotEmpty(t, NewFileStorage("/tmp/shortener/shortener.log"))
}

func TestDirNotExist(t *testing.T) {
	dir := util.GenerateUniqueID()
	assert.NotEmpty(t, NewFileStorage("/tmp/shortener/"+dir+"/log.file"))
}
