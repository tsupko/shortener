package storage

import "github.com/tsupko/shortener/internal/app/exceptions"

type TestStorage struct{}

var _ Storage = &TestStorage{}

func (t TestStorage) Save(_, url, _ string) (string, error) {
	if url == "https://ya.ru" || url == "https://github.com/tsupko/shortener/runs/7826862296?check_suite_focus=true" {
		return "12345", nil
	}
	return "67890", nil
}

func (t *TestStorage) SaveBatch(hashes, urls, userIds []string) ([]string, error) {
	values := make([]string, 0, len(hashes))
	for i := range hashes {
		hash, err := t.Save(hashes[i], urls[i], userIds[i])
		if err != nil {
			panic("unexpected behavior")
		}
		values = append(values, hash)
	}
	return values, nil
}

func (t TestStorage) Get(hash string) (User, error) {
	if hash == "12345" {
		return User{"https://ya.ru", ""}, nil
	}
	return User{"", ""}, exceptions.ErrURLNotFound
}

func (t TestStorage) GetAll(string) (map[string]string, error) {
	data := make(map[string]string)
	data["12345"] = "https://ya.ru"
	return data, nil
}
