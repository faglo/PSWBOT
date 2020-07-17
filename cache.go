package main

import (
	"fmt"
	"sync"
)

type Entry struct {
	sync.RWMutex
	data map[string]interface{}
}

type Cache struct {
	CanWrite *Entry
	Message *Entry
}

func (e *Entry) Set(key string, val interface{}) {
	e.Lock()
	e.data[key] = val
	e.Unlock()
}

func (e *Entry) Get(key string) (interface{}, bool)  {
	e.RLock()
	defer e.RUnlock()
	val, ok := e.data[key]
	if !ok {
		fmt.Println("Cache corrupted")
	}
	return val, ok
}

func NewUserCache() *Cache {
	entry := func() *Entry {
		data := make(map[string]interface{})
		return &Entry{data: data}
	}
	return &Cache{
		CanWrite: entry(),
		Message: entry(),
	}
}
