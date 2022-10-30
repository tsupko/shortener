package storage

import (
	"encoding/json"
	"os"
)

type Record struct {
	Hash string `json:"hash"`
	URL  string `json:"url"`
}

type producer struct {
	file    *os.File
	encoder *json.Encoder
}

func NewProducer(filename string) (*producer, error) {
	file, err := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0o777)
	if err != nil {
		return nil, err
	}

	return &producer{
		file:    file,
		encoder: json.NewEncoder(file),
	}, nil
}

func (p *producer) WriteRecord(record *Record) error {
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
	file, err := os.OpenFile(filename, os.O_RDONLY|os.O_CREATE, 0o777)
	if err != nil {
		return nil, err
	}

	return &consumer{
		file:    file,
		decoder: json.NewDecoder(file),
	}, nil
}

func (c *consumer) ReadRecord() (*Record, error) {
	event := Record{}
	if err := c.decoder.Decode(&event); err != nil {
		return nil, err
	}
	return &event, nil
}

func (c *consumer) Close() error {
	return c.file.Close()
}
