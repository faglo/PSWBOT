package main

import (
	"fmt"
	tb "github.com/demget/telebot"
	"sync"
)

type Entry struct {
	sync.RWMutex
	data map[string]interface{}
}

type Cache struct {
	CanWrite    *Entry
	Message     *Entry
	Preview     *Entry
	TempMailing *Entry
	UserCache 	*Entry
	AdminCache  *Entry
}

type UserCache struct {
	ReservedService string
}

type AdminCache struct {
	CheckingHW HomeworkResult
	PreviewMsg *tb.Message
	Comment string
	Reject bool
}

func (e *Entry) Set(key string, val interface{}) {
	e.Lock()
	e.data[key] = val
	e.Unlock()
}

func (e *Entry) Get(key string) (interface{}, bool) {
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
		CanWrite:    entry(),
		Message:     entry(),
		Preview:     entry(),
		TempMailing: entry(),
		UserCache:	 entry(),
		AdminCache:  entry(),
	}
}
