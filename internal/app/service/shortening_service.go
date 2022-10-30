package service

type ShorteningService interface {
	Get(hash string) (string, error)
	GetAll() (map[string]string, error)
	Save(url string) (string, error)
	SaveBatch(hashes []string, urls []string) ([]string, error)
}
