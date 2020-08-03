package main

import (
	tb "github.com/demget/telebot"
	"github.com/panjf2000/ants/v2"
	"go.uber.org/atomic"
	"sync"
)

func Distribution(mailing TempMailing)  {
	var wg sync.WaitGroup
	pool, err := ants.NewPool(4)
	handleErr(err, mailing.mailing)
	defer pool.Release()
	var success atomic.Int64
	var fail atomic.Int64
	var overall int
	users, err := getUsers("default")
	handleErr(err, mailing.mailing)


	for i, user := range users {
		_ = pool.Submit(func() {
			defer wg.Done()
			var what interface{}

			if mailing.mailing.Photo != nil{
				what = &tb.Photo{
					File:      tb.File{FileID: mailing.mailing.Document.FileID},
					Caption:   mailing.mailing.Text,
					ParseMode: tb.ModeHTML,
				}
			} else if mailing.mailing.Video != nil {
				what = &tb.Video{
					File:    tb.File{FileID: mailing.mailing.Video.FileID},
					Caption: mailing.mailing.Text,
				}
			} else {
				what = mailing.mailing.Text
			}
			_, err = b.Send(&tb.User{ID: user}, what, tb.ModeHTML)
			if err != nil {
				fail.Add(1)
			} else {
				success.Add(1)
			}
			overall = i
		})
	}
	type result struct {
		success int
		fail int
		all int
	}
	_, err = b.Edit(mailing.last_msg, b.Text("Sent", result{
		success: int(success.Load()),
		fail:    int(fail.Load()),
		all:     overall,
	}), tb.ModeHTML)
}
