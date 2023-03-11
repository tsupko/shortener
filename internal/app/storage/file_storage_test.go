package storage

import (
	"log"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/tsupko/shortener/internal/app/util"
)

func TestReadFromFileWhenCreated(t *testing.T) {
	hash := util.GenerateUniqueID()

	fileStorage := NewFileStorage("file.log")
	fileStorage.writeToFile(hash, "URL")

	url, _ := fileStorage.Get(hash)
	assert.Equal(t, User{}, url)

	anotherStorage := NewFileStorage("file.log")
	url, _ = anotherStorage.Get(hash)
	assert.Equal(t, User{UserID: "URL"}, url)
}

func TestDoubleSave(t *testing.T) {
	hash := util.GenerateUniqueID()

	fileStorage := NewFileStorage("file.log")
	_, err := fileStorage.Save(hash, "URL", "user1")
	if err != nil {
		log.Printf("error saving user data: %v", err)
	}
	_, err = fileStorage.Save(hash, "url2", "user2")
	if err != nil {
		log.Printf("error saving user data: %v", err)
	}

	url, _ := fileStorage.Get(hash)
	assert.Equal(t, User{UserID: "url2"}, url)

	anotherStorage := NewFileStorage("file.log")
	url, _ = anotherStorage.Get(hash)
	assert.Equal(t, User{UserID: "url2"}, url)
}

func Test(t *testing.T) {
	assert.NotEmpty(t, NewFileStorage("/tmp/shortener/shortener.log"))
}

func TestDirNotExist(t *testing.T) {
	dir := util.GenerateUniqueID()
	assert.NotEmpty(t, NewFileStorage("/tmp/shortener/"+dir+"/log.file"))
}
