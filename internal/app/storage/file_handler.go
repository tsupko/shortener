package storage

import (
	"encoding/json"
	"os"
)

type record struct {
	Hash string `json:"hash"`
	URL  string `json:"url"`
}

type producer struct {
	file    *os.File
	encoder *json.Encoder
}

func NewProducer(filename string) (*producer, error) {
	file, err := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0777)
	if err != nil {
		return nil, err
	}

	return &producer{
		file:    file,
		encoder: json.NewEncoder(file),
	}, nil
}

func (p *producer) WriteRecord(record *record) error {
	return p.encoder.Encode(&record)
}

func (p *producer) Close() error {
	return p.file.Close()
}

type consumer struct {
	file    *os.File
	decoder *json.Decoder
}

func NewConsumer(filename string) (*consumer, error) {
	file, err := os.OpenFile(filename, os.O_RDONLY|os.O_CREATE, 0777)
	if err != nil {
		return nil, err
	}

	return &consumer{
		file:    file,
		decoder: json.NewDecoder(file),
	}, nil
}

func (c *consumer) ReadRecord() (*record, error) {
	event := record{}
	if err := c.decoder.Decode(&event); err != nil {
		return nil, err
	}
	return &event, nil
}

func (c *consumer) Close() error {
	return c.file.Close()
}
