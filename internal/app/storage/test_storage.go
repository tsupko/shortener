package storage

type TestStorage struct {
}

func (t TestStorage) GetAll() (interface{}, interface{}) {
	//TODO implement me
	panic("implement me")
}

func NewTestStorage() *TestStorage {
	return &TestStorage{}
}

func (t TestStorage) Put(string, string) string {
	return "12345"
}

func (t TestStorage) Get(id string) (string, bool) {
	if id == "12345" {
		return "https://ya.ru", true
	}
	return "", false
}
