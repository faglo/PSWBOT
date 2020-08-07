package main

import (
	tb "github.com/demget/telebot"
	"log"
	"strconv"
	"strings"
	"time"
)

type Preview struct {
	PreviewMsg     *tb.Message
	PreviewService Service
}
type TempMailing struct {
	mailing  *tb.Message
	stime    time.Time
	isNow    bool
	last_msg *tb.Message
}

func onText() func(m *tb.Message) {
	return func(m *tb.Message) {
		if m.Private() {
			state, err := getState(m.Sender.ID)
			handleErr(err, m)
			switch state {
			case "enterFeedback":
				onFeedback(m)
			case "enterQuestion":
				onQuestion(m)
			case "enterMailing":
				onMailing(m)
			case "sendZip":
			case "writeCommentReject":
				preview, err := b.Send(m.Sender, b.Text("HWCommentPreview", m.Text), b.InlineMarkup("comment_kb"))
				handleErr(err, m)
				adminCacheRaw, _ := cache.AdminCache.Get(strconv.Itoa(m.Sender.ID))
				adminCache := adminCacheRaw.(AdminCache)
				adminCache.PreviewMsg = preview
				adminCache.Comment = m.Text
				adminCache.CheckingHW.AdminComment = m.Text
				adminCache.Reject = true
				cache.AdminCache.Set(strconv.Itoa(m.Sender.ID), adminCache)
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
			case "sendZip":
				onSentZip(m)
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
func onEditHW(m *tb.Message) {
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
		handleErr(err, m)
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

	_, err = b.Send(m.Chat, b.Text("ServicePreview", course))
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

	err = delService(course.ServiceID)
	handleErr(err, m)
	err = addService(course)
	handleErr(err, m)
	_, err = b.Send(m.Sender, b.Text("editHWSuccess"))
	handleErr(err, m)
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
func onSentZip(m *tb.Message) {
	hwGroupRaw, err := getConfig("HWGroup")
	handleErr(err, m)
	hwGroup, _ := strconv.Atoi(hwGroupRaw)
	serviceID, _ := cache.UserCache.Get(strconv.Itoa(m.Sender.ID))
	rawService := strings.Split(serviceID.(string), ".")
	course, _ := strconv.Atoi(rawService[0])
	lesson, _ := strconv.Atoi(rawService[1])
	service, err := getService(serviceID.(string))
	handleErr(err, m)
	result := HomeworkResult{
		UserID:      m.Sender.ID,
		Course:      course,
		Lesson:      lesson,
		Grade:       0,
		IsGraded:    false,
		MessageID:   0,
		ResultID:    0,
		UserComment: m.Caption,
		Username: m.Sender.Username,
		CourseName: service.Name,
	}
	if m.Caption != "" {
		m.Caption = "Комментраий: "+m.Caption
	}
	m.Document.Caption = "<i>Обработка</i>"
	gradeMessage, err := b.Send(&tb.Chat{ID:int64(hwGroup)}, m.Document, tb.ModeHTML)
	handleErr(err, m)
	result.MessageID = gradeMessage.ID
	_, err = b.Send(m.Sender, b.Text("HWSent", result))
	handleErr(err, m)
	m.Document.Caption = b.Text("ForGraduate", result)
	_, err = b.Edit(gradeMessage, m.Document, b.InlineMarkup("pre_rate", result.MessageID), tb.ModeHTML)
	handleErr(err, m)
	err = SetResult(result)
	handleErr(err, m)
	err = removeService(serviceID.(string), m.Sender.ID)
	handleErr(err, m)
	err = setState(m.Sender.ID, "default")
	handleErr(err, m)
}
