package storage

type Storage interface {
	Put(id string, url string) string
	Get(id string) (string, bool)
}
