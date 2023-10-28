package server

import (
	"fmt"
	"sync"
)

type Log struct {
	mu      sync.Mutex
	records []Record
}

type Record struct {
	Value []byte `json:"value"`
	// This is being used to store the offset of the record from the start of the slice.
	Offset uint64 `json:"offset"`
}

var ErrOffsetNotFound = fmt.Errorf("offset not found")

func NewLog() *Log {
	return &Log{}
}

func (c *Log) Append(r Record) (uint64, error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	r.Offset = uint64(len(c.records))
	c.records = append(c.records, r)

	return r.Offset, nil
}

func (c *Log) Read(offset uint64) (Record, error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if offset >= uint64(len(c.records)) {
		return Record{}, ErrOffsetNotFound
	}

	return c.records[offset], nil
}
