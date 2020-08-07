package main

import (
	tb "github.com/demget/telebot"
)

type Result struct {
	Success int
	Fail    int
	All     int
}

func Distribution(mailing TempMailing) {
	var success int
	var fail int
	var overall int
	users, err := getUsers("default")
	handleErr(err, mailing.mailing)

	for _, user := range users {
		var what interface{}

		if mailing.mailing.Photo != nil {
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
			fail += 1
		} else {
			success += 1
		}
		overall += 1
	}

	fields := Result{
		Success: success,
		Fail:    fail,
		All:     overall,
	}

	_, err = b.Edit(mailing.last_msg, b.Text("Sent", fields), tb.ModeHTML)
}
