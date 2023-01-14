package storage

type Storage interface {
	Save(hash, url, userID string) (string, error)
	SaveBatch(hashes, urls, userIds []string) ([]string, error)
	Get(hash string) (User, error)
	GetAll(userID string) (map[string]string, error)
}
