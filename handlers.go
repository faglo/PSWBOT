package main

import (
	tb "github.com/demget/telebot"
	"log"
	"regexp"
	"strconv"
	"strings"
	"time"
)

type Preview struct {
	PreviewMsg     *tb.Message
	PreviewService Service
}
type TempMailing struct {
	 mailing *tb.Message
	 stime time.Time
	 isNow bool
	 last_msg *tb.Message
}

func onText() func(m *tb.Message) {
	return func(m *tb.Message) {
		if m.Private() {
			state, err := getState(m.Sender.ID)
			handleErr(err, m)
			switch state {
			case "enterEmail":
				{
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
			case "enterQuestion":
				onQuestion(m)
			case "enterMailing":
				onMailing(m)
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
			case "enterMailing":
				onMailing(m)

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
			case "enterMailing":
				onMailing(m)
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
			cid, err := getConfig("HWGroup")
			handleErr(err, m)
			chatid, err := strconv.Atoi(cid)
			if m.Chat.ID == int64(chatid) {
				onAddHW(m)
			} else {
				groupHandler(m)
			}

		}
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
func onFeedback(m *tb.Message) {
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
			log.Print("ERR: DELETE FAIL: " + err.Error())
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
func onQuestion(m *tb.Message) {
	err = setState(m.Sender.ID, "default")
	handleErr(err, m)
	hwChatID, err := getConfig("HWQGroup")
	handleErr(err, m)
	hwChatIDint, _ := strconv.Atoi(hwChatID)
	msg, err := b.Forward(&tb.Chat{ID: int64(hwChatIDint)}, m)
	handleErr(err, m)
	err = setHWQuestion(m.Sender.ID, msg.ID)
	handleErr(err, m)
	lmsg, _ := cache.Message.Get(strconv.Itoa(m.Sender.ID))
	err = b.Delete(lmsg.(*tb.Message))
	handleErr(err, m)
	_, err = b.Send(m.Sender, b.Text("HWQSuccess", msg.ID), mainKB(m), tb.ModeHTML)
	handleErr(err, m)

}

func groupHandler(m *tb.Message) {
	if m.IsReply() && m.ReplyTo.Sender.IsBot {
		feedbackCIDRaw, err := getConfig("FeedbackGroup")
		handleErr(err, m)
		HWQCIDRaw, err := getConfig("HWQGroup")
		feedbackCID, _ := strconv.Atoi(feedbackCIDRaw)
		hwCID, _ := strconv.Atoi(HWQCIDRaw)
		var user int
		if m.Chat.ID == int64(feedbackCID) {
			user, err = getUserByMessage(m.ReplyTo.ID)
			handleErr(err, m)
			mediaHandler(m, user, "feedback")
			err = resetFeedbackMessage(user)
			handleErr(err, m)
			repliedMsg, err := b.Reply(m, b.Text("FBSent"))
			handleErr(err, m)
			go time.AfterFunc(time.Duration(deltime)*time.Second, func() {
				err = b.Delete(repliedMsg)
				handleErr(err, m)
			})

		} else if m.Chat.ID == int64(hwCID) {
			user, err = getUserByHWQ(m.ReplyTo.ID)
			handleErr(err, m)
			mediaHandler(m, user, "question")
			err := removeHWQ(m.ReplyTo.ID, user)
			repliedMsg, err := b.Reply(m, b.Text("FBSent"))
			handleErr(err, m)
			go time.AfterFunc(time.Duration(deltime)*time.Second, func() {
				err = b.Delete(repliedMsg)
				handleErr(err, m)
			})
		}
	}
}
func mediaHandler(m *tb.Message, user int, typeOf string) {
	var template string
	if typeOf == "feedback" {
		template = b.Text("SendAnswerFB", m)
	} else if typeOf == "question" {
		template = b.Text("SendAnswerHWQ", m)
	}
	if m.Photo != nil {
		m.Text = m.Caption
		_, err = b.Send(&tb.User{ID: user}, &tb.Photo{
			File:      m.Photo.File,
			Caption:   template,
			ParseMode: tb.ModeHTML,
		})
		handleErr(err, m)
	} else if m.Video != nil {
		m.Text = m.Caption
		_, err = b.Send(&tb.User{ID: user}, &tb.Video{
			File:    m.Video.File,
			Caption: template,
		}, tb.ModeHTML)
		handleErr(err, m)
	} else if m.Document != nil {
		m.Text = m.Caption
		_, err = b.Send(&tb.User{ID: user}, &tb.Document{
			File:    m.Document.File,
			Caption: template,
		}, tb.ModeHTML)
		handleErr(err, m)
	} else {
		_, err = b.Send(&tb.User{ID: user}, template, tb.ModeHTML)
		handleErr(err, m)
	}
}

func onAddHW(m *tb.Message) {
	chatRaw, err := getConfig("HWGroup")
	handleErr(err, m)
	chat, _ := strconv.Atoi(chatRaw)

	if m.Chat.ID != int64(chat) {
		return
	}

	if m.Document != nil {
		m.Text = m.Caption
	}

	fieldsRaw := strings.Split(m.Text, "#")
	var fields []string
	var typeOf string
	var article string
	var price int
	var fileID string
	var template string

	for i, fieldRaw := range fieldsRaw {
		if i == 2 {
			fields = append(fields, fieldRaw)
			continue
		}
		field := strings.Trim(fieldRaw, "\n")
		fields = append(fields, field)
	}

	sid := strings.Split(fields[0], " ")[1]
	_, err = getService(sid)
	if err != nil {
		if strings.Contains(err.Error(), "no rows") {
		} else {
			handleErr(err, m)
		}
	} else {
		_, err = b.Send(m.Chat, b.Text("ServiceExistErr"))
		return
	}

	if strings.Contains(sid, ".") {
		typeOf = "lesson"
	} else {
		typeOf = "course"
	}

	if len(fields) == 5 {
		article = fields[4]
	}

	price, err = strconv.Atoi(fields[3])
	handleErr(err, m)

	if m.Document != nil {
		fileID = m.Document.FileID
	}

	course := Service{
		ServiceID:   sid,
		Type:        typeOf,
		Description: fields[2],
		FileURI:     fileID,
		Price:       price,
		ArticleURL:  article,
		Name:        fields[1],
		Bought:      0,
	}

	prev, err := b.Send(m.Chat, b.Text("ServicePreview", course))
	handleErr(err, m)
	if course.Type == "lesson" {
		template = "buyLesson"
	} else if course.Type == "course" {
		template = "CourseInfo"
	}
	if course.FileURI != "" && course.ArticleURL != "" {
		_, err = b.Send(m.Chat, &tb.Document{
			File:    tb.File{FileID: course.FileURI},
			Caption: b.Text(template, course),
		}, b.InlineMarkup("article_btn", course))
	} else if course.ArticleURL != "" {
		_, err = b.Send(m.Chat, b.Text(template, course), b.InlineMarkup("article_btn", course))
	} else if course.FileURI != "" {
		_, err = b.Send(m.Chat, &tb.Document{
			File:    tb.File{FileID: course.FileURI},
			Caption: b.Text(template, course),
		})
	} else {
		_, err = b.Send(m.Chat, b.Text(template, course))
	}

	cache.Preview.Set(strconv.Itoa(m.Sender.ID), Preview{
		PreviewMsg:     prev,
		PreviewService: course,
	})

	_, err = b.Send(m.Sender, b.Text("sendConfirm"), b.InlineMarkup("confirm_service_kb"))
}
func onDelHW(m *tb.Message) {
	chatRaw, err := getConfig("HWGroup")
	handleErr(err, m)
	chat, _ := strconv.Atoi(chatRaw)

	if m.Chat.ID != int64(chat) {
		return
	}

	fieldsRaw := strings.Split(m.Text, " ")
	service, err := getService(fieldsRaw[1])
	handleErr(err, m)
	confirmMsg, err := b.Send(m.Sender, b.Text("HWConfirmDelete", service), b.InlineMarkup("delete_service", service.ServiceID))
	handleErr(err, m)
	cache.Message.Set(strconv.Itoa(m.Sender.ID), confirmMsg)
}

func onMailing(m *tb.Message) {
	var preview *tb.Message
	if m.Photo != nil {
		m.Photo.Caption = m.Caption
		preview, err = b.Send(m.Sender, m.Photo, b.InlineMarkup("pre_send"))
		handleErr(err, m)
	} else if m.Video != nil {
		m.Video.Caption = m.Caption
		preview, err = b.Send(m.Sender, m.Video, b.InlineMarkup("pre_send"))
		handleErr(err, m)
	} else {
		preview, err = b.Send(m.Sender, m.Text, b.InlineMarkup("pre_send"))
		handleErr(err, m)
	}
	temp := TempMailing{
		mailing:  m,
		stime:    time.Time{},
		isNow:    true,
		last_msg: preview,
	}
	cache.TempMailing.Set(strconv.Itoa(m.Sender.ID), temp)

}
