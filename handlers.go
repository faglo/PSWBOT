package main

import (
	"fmt"
	tb "github.com/demget/telebot"
	"log"
	"regexp"
	"strconv"
	"time"
)

func onText() func(m *tb.Message) {
	return func(m *tb.Message) {
		if m.Private() {
			state, err := getState(m.Sender.ID)
			handleErr(err, m)
			switch state {
			case "enterEmail": {
				email, err := getEmail(m.Sender.ID)
				handleErr(err, m)
				re := regexp.MustCompile("(?:[a-z0-9!#$%&'*+/=?^_`{|}~-]+(?:\\.[a-z0-9!#$%&'*+/=?^_`{|}~-]+)*|\"(?:[\x01-\x08\x0b\x0c\x0e-\x1f\x21\x23-\x5b\x5d-\x7f]|\\[\x01-\x09\x0b\x0c\x0e-\x7f])*\")@(?:(?:[a-z0-9](?:[a-z0-9-]*[a-z0-9])?\\.)+[a-z0-9](?:[a-z0-9-]*[a-z0-9])?|\\[(?:(?:(2(5[0-5]|[0-4][0-9])|1[0-9][0-9]|[1-9]?[0-9]))\\.){3}(?:(2(5[0-5]|[0-4][0-9])|1[0-9][0-9]|[1-9]?[0-9])|[a-z0-9-]*[a-z0-9]:(?:[\x01-\x08\x0b\x0c\x0e-\x1f\x21-\x5a\x53-\x7f]|\\[\x01-\x09\x0b\x0c\x0e-\x7f])+)])")
				if !re.Match([]byte(m.Text)) {
					_, err = b.Send(m.Sender, b.Text("IncorrectMail"))
					handleErr(err, m)
				} else if email == m.Text {
					_, err = b.Send(m.Sender, b.Text("ExistMail"))
					handleErr(err, m)
				} else {
					err = setEmail(m.Sender.ID, m.Text)
					handleErr(err, m)
					_, err = b.Reply(m, "Вы успешно зарегистрировались!", b.Markup("main"))
					handleErr(err, m)
					err = setState(m.Sender.ID, "default")
					handleErr(err, m)
				}
			}
			case "enterFeedback":
				onFeedback(m)
			}
		} else {
			groupHandler(m)
		}

	}
}
func onPhoto() func(m *tb.Message) {
	return func(m *tb.Message) {
		if m.Private() {
			state, err := getState(m.Sender.ID)
			handleErr(err, m)
			switch state {
			case "enterFeedback":
				onFeedback(m)
			}
		} else {
			groupHandler(m)
		}
	}
}
func onVideo() func(m *tb.Message) {
	return func(m *tb.Message) {
		if m.Private() {
			state, err := getState(m.Sender.ID)
			handleErr(err, m)
			switch state {
			case "enterFeedback":
				onFeedback(m)
			}
		} else {
			groupHandler(m)
		}
	}
}
func onDocument() func(m *tb.Message) {
	return func(m *tb.Message) {
		if m.Private() {
			state, err := getState(m.Sender.ID)
			handleErr(err, m)
			switch state {
			case "enterFeedback":
				onFeedback(m)
			}
		} else {
			groupHandler(m)
		}
	}
}
func onChannelPost() func(m *tb.Message) {
	return func(m *tb.Message) {
		fmt.Println(m.Document)
	}
}

func onMainMenu() func(m *tb.Message) {
	return func(m *tb.Message) {
		err = setState(m.Sender.ID, "default")
		handleErr(err, m)
		_, err = b.Send(m.Sender, b.Text("YouInStart", m.Sender), mainKB(m), tb.ModeHTML)
		handleErr(err, m)
	}
}
func onFeedback(m *tb.Message)  {
	err = setState(m.Sender.ID, "default")
	handleErr(err, m)
	feedbackChatID, err := getConfig("FeedbackGroup")
	handleErr(err, m)
	feedbackChatIDint, _ := strconv.Atoi(feedbackChatID)
	FeedbackMessageID, err := getFeedbackMessage(m.Sender.ID)
	handleErr(err, m)

	if FeedbackMessageID != 0 {
		err = b.Delete(&tb.StoredMessage{
			ChatID:    int64(feedbackChatIDint),
			MessageID: strconv.Itoa(FeedbackMessageID),
		})
		if err != nil {
			log.Print("ERR: DELETE FAIL: "+err.Error())
		}
	}

	msg, err := b.Forward(&tb.Chat{ID: int64(feedbackChatIDint)}, m)
	handleErr(err, m)
	err = setFeedbackMessage(m.Sender.ID, msg.ID)
	handleErr(err, m)
	err = incrementFeedBackMessage(m.Sender.ID)
	handleErr(err, m)
	_, err = b.Send(m.Sender, b.Text("FeedbackSuccess"), mainKB(m), tb.ModeHTML)
	handleErr(err, m)
}

func groupHandler(m *tb.Message)  {
	if m.IsReply() && m.ReplyTo.Sender.IsBot {
		feedbackCIDRaw, err := getConfig("FeedbackGroup")
		handleErr(err, m)
		feedbackCID, _ := strconv.Atoi(feedbackCIDRaw)
		var user int
		if m.Chat.ID == int64(feedbackCID) {
			user, err = getUserByMessage(m.ReplyTo.ID)
			handleErr(err, m)
			if m.Photo != nil {
				m.Text = m.Caption
				_, err = b.Send(&tb.User{ID: user}, &tb.Photo{
					File:      m.Photo.File,
					Caption:   b.Text("SendAnswerFB", m),
					ParseMode: tb.ModeHTML,
				})
				handleErr(err, m)
			} else if m.Video != nil {
				m.Text = m.Caption
				_, err = b.Send(&tb.User{ID: user}, &tb.Video{
					File:      m.Video.File,
					Caption:   b.Text("SendAnswerFB", m),
				}, tb.ModeHTML)
				handleErr(err, m)
			} else if m.Document != nil {
				m.Text = m.Caption
				_, err = b.Send(&tb.User{ID: user}, &tb.Document{
					File:      m.Document.File,
					Caption:   b.Text("SendAnswerFB", m),
				}, tb.ModeHTML)
				handleErr(err, m)
			} else {
				_, err = b.Send(&tb.User{ID: user}, b.Text("SendAnswerFB", m), tb.ModeHTML)
				handleErr(err, m)
			}
			err = resetFeedbackMessage(user)
			handleErr(err, m)
			repliedMsg, err := b.Reply(m, b.Text("FBSent"))
			handleErr(err, m)
			go time.AfterFunc(time.Duration(deltime)*time.Second, func() {
				err = b.Delete(repliedMsg)
				handleErr(err, m)
			})

		}
	}
}

