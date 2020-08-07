package main

import (
	"fmt"
	tb "github.com/demget/telebot"
	"log"
	"strconv"
	"strings"
)

var b *tb.Bot
var deltime = 5
var btnCount = 9
var cache = NewUserCache()

func main() {
	// Bot init
	tmplEngine := &tb.TemplateText{Dir: "data"}
	pref, err := tb.NewSettings("bot.json", tmplEngine)
	if err != nil {
		log.Fatal("Unable to start bot FE-1")
	}
	pref.Token = "1039019679:AAH2oVWii5zNMhrlVJCM_D0ySUYgM4_zTRw"
	b, err = tb.NewBot(pref)
	if err != nil {
		log.Fatal("Unable to start bot FE-2")
	}

	// Handlers
	b.Handle(tb.OnText, onText())
	b.Handle(tb.OnPhoto, onPhoto())
	b.Handle(tb.OnVideo, onVideo())
	b.Handle(tb.OnDocument, onDocument())
	b.Handle("/start", func(m *tb.Message) {
		if m.Private() {
			exist, err := checkUser(m.Sender.ID)
			handleErr(err, m)
			if !exist {
				err = addUser(m.Sender.ID, m.Sender.Username)
				handleErr(err, m)
				_, err = b.Send(m.Sender, b.Text("Start", m.Sender), b.Markup("main"), tb.ModeHTML)
				handleErr(err, m)
			} else {
				err = setState(m.Sender.ID, "default")
				handleErr(err, m)
				_, err = b.Send(m.Sender, b.Text("Start", m.Sender), mainKB(m), tb.ModeHTML)
				handleErr(err, m)
			}
		} else {
			_, err = b.Send(m.Chat, "ID чата: "+strconv.Itoa(int(m.Chat.ID)))
		}
	})
	b.Handle("/addchat", func(m *tb.Message) {
		permLevel, err := getPermLevel(m.Sender.ID)
		handleErr(err, m)
		if permLevel == 3 {
			groupID := strings.Split(m.Text, " ")[1]
			err := setConfig("FeedbackGroup", groupID)
			handleErr(err, m)
			groupIDint, err := strconv.Atoi(groupID)
			handleErr(err, m)
			_, err = b.Send(&tb.Chat{ID: int64(groupIDint)}, "Чат используется для обработки фидбэков")
			handleErr(err, m)
			_, err = b.Send(m.Sender, "Чат для фидбэков установлен")
			handleErr(err, m)
		}
	})
	b.Handle("/hwqchat", func(m *tb.Message) {
		permLevel, err := getPermLevel(m.Sender.ID)
		handleErr(err, m)
		if permLevel == 3 {
			groupID := strings.Split(m.Text, " ")[1]
			err := setConfig("HWQGroup", groupID)
			handleErr(err, m)
			groupIDint, err := strconv.Atoi(groupID)
			handleErr(err, m)
			_, err = b.Send(&tb.Chat{ID: int64(groupIDint)}, "Чат используется для обработки вопросов")
			handleErr(err, m)
			_, err = b.Send(m.Sender, "Чат для вопросов установлен")
			handleErr(err, m)
		}
	})
	b.Handle("/hwchat", func(m *tb.Message) {
		permLevel, err := getPermLevel(m.Sender.ID)
		handleErr(err, m)
		if permLevel == 3 {
			groupID := strings.Split(m.Text, " ")[1]
			err := setConfig("HWGroup", groupID)
			handleErr(err, m)
			groupIDint, err := strconv.Atoi(groupID)
			handleErr(err, m)
			_, err = b.Send(&tb.Chat{ID: int64(groupIDint)}, "Здесь дз")
			handleErr(err, m)
			_, err = b.Send(m.Sender, "Чат для дз установлен")
			handleErr(err, m)
		}
	})
	b.Handle("/adminchat", func(m *tb.Message) {
		permLevel, err := getPermLevel(m.Sender.ID)
		handleErr(err, m)
		if permLevel == 3 {
			groupID := strings.Split(m.Text, " ")[1]
			err := setConfig("AdminGroup", groupID)
			handleErr(err, m)
			groupIDint, err := strconv.Atoi(groupID)
			handleErr(err, m)
			_, err = b.Send(&tb.Chat{ID: int64(groupIDint)}, "хорошо")
			handleErr(err, m)
			_, err = b.Send(m.Sender, "я понял")
			handleErr(err, m)
		}
	})
	b.Handle("/msg", func(m *tb.Message) {
		if !m.Private() {
			unReplied, err := getUnReplied()
			handleErr(err, m)
			_, err = b.Send(m.Chat, fmt.Sprintf("Неотвеченных сообщений: %d", unReplied))
			handleErr(err, m)
		}
	})
	b.Handle("/add", func(m *tb.Message) { onAddHW(m) })
	b.Handle("/del", func(m *tb.Message) { onDelHW(m) })
	b.Handle("/edit", func(m *tb.Message) { onEditHW(m) })
	b.Handle("/test", func(m *tb.Message) {
		_, err = b.Send(m.Sender, "iwiopefi", b.InlineMarkup("grading"))
		if err != nil {
			handleErr(err, m)
		}
	})

	b.Handle(b.Button("lk"), func(m *tb.Message) {
		_, err = b.Send(m.Sender, b.Text("YouInLK"), b.InlineMarkup("lk_kb"))
		handleErr(err, m)
	})
	b.Handle(b.Button("FAQ"), func(m *tb.Message) {
		_, err = b.Send(m.Sender, b.Text("FAQ"), tb.ModeHTML)
		handleErr(err, m)
	})
	b.Handle(b.Button("feedback"), func(m *tb.Message) {
		_, err = b.Send(m.Sender, b.Text("EnterFeedback"), b.Markup("cancel"), tb.ModeHTML)
		handleErr(err, m)
		err = setState(m.Sender.ID, "enterFeedback")
		handleErr(err, m)
	})
	b.Handle(b.Button("cancel"), onMainMenu())
	b.Handle(b.Button("admin_panel"), func(m *tb.Message) {
		_, err := b.Send(m.Sender, "Вы в админ панели", b.Markup("admin_panelkb"))
		handleErr(err, m)
		err = setState(m.Sender.ID, "default")
		handleErr(err, m)
	})
	b.Handle(b.Button("bot_stats"), func(m *tb.Message) {
		_, err = b.Send(m.Sender, getBotStats(m), tb.ModeHTML)
		handleErr(err, m)
	})
	b.Handle(b.Button("course_stats"), func(m *tb.Message) {
		_, err = b.Send(m.Sender, getCourseStats(m), tb.ModeHTML)
		handleErr(err, m)
	})
	b.Handle(b.Button("new_mailing"), func(m *tb.Message) {
		_, err = b.Send(m.Sender, b.Text("CreateMailing"), b.Markup("to_admin"))
		handleErr(err, m)
		err = setState(m.Sender.ID, "enterMailing")
		handleErr(err, m)
	})
	b.Handle(b.Button("to_main"), func(m *tb.Message) {
		_, err = b.Send(m.Sender, b.Text("Start", m.Sender), mainKB(m), tb.ModeHTML)
		handleErr(err, m)
	})

	b.Handle(b.InlineButton("buy_homework"), buyHWboard())
	b.Handle(b.InlineButton("to_courses"), buyHWboard())
	b.Handle(b.InlineButton("to_lk"), func(c *tb.Callback) {
		err = setState(c.Sender.ID, "default")
		_, err = b.Edit(c.Message, b.Text("YouInLK"), b.InlineMarkup("lk_kb"))
		handleErr(err, c.Message)
	})
	b.Handle(b.InlineButton("to_feedback"), func(c *tb.Callback) {
		_, err = b.Send(c.Sender, b.Text("EnterFeedback"), b.Markup("cancel"), tb.ModeHTML)
		handleErr(err, c.Message)
		err = setState(c.Sender.ID, "enterFeedback")
		handleErr(err, c.Message)
	})
	b.Handle(b.InlineButton("ask_hw"), func(c *tb.Callback) {
		premium, err := isPremium(c.Sender.ID)
		handleErr(err, c.Message)
		if !premium {
			_, err = b.Edit(c.Message, b.Text("NoPremium"), b.InlineMarkup("premium_error"))
			handleErr(err, c.Message)
			return
		}
		_, err = b.Edit(c.Message, b.Text("EnterQuestion"), b.InlineMarkup("to_lk_kb"), tb.ModeHTML)
		handleErr(err, c.Message)
		err = setState(c.Sender.ID, "enterQuestion")
		handleErr(err, c.Message)
		cache.Message.Set(strconv.Itoa(c.Sender.ID), c.Message)
	})
	b.Handle(b.InlineButton("del_service"), func(c *tb.Callback) {
		err = b.Delete(c.Message)
		handleErr(err, c.Message)
		err = b.Respond(c, &tb.CallbackResponse{Text: "Отменено", CallbackID: c.ID})
	})
	b.Handle(b.InlineButton("add_service"), func(c *tb.Callback) {
		prev, _ := cache.Preview.Get(strconv.Itoa(c.Sender.ID))
		err = addService(prev.(Preview).PreviewService)
		handleErr(err, c.Message)
		_, err = b.Edit(c.Message, b.Text("AddHWSuccess"))
		handleErr(err, c.Message)
	})
	b.Handle(b.InlineButton("erase_service"), func(c *tb.Callback) {
		msg, _ := cache.Message.Get(strconv.Itoa(c.Sender.ID))
		err = b.Delete(msg.(*tb.Message))
		handleErr(err, c.Message)
		err = delService(c.Data)
		handleErr(err, c.Message)
		err = b.Respond(c, &tb.CallbackResponse{Text: "Удалено"})
		handleErr(err, c.Message)
	})
	b.Handle(b.InlineButton("send_now"), func(c *tb.Callback) {
		temp, _ := cache.TempMailing.Get(strconv.Itoa(c.Sender.ID))
		Distribution(temp.(TempMailing))
		_, err := b.Send(c.Sender, "Вы в админ панели", b.Markup("admin_panelkb"))
		handleErr(err, c.Message)
		err = setState(c.Sender.ID, "default")
		handleErr(err, c.Message)
	})
	b.Handle(b.InlineButton("send_hw"), func(c *tb.Callback) {
		premium, err := isNowPremium(c.Sender.ID)
		handleErr(err, c.Message)
		if premium == false {
			_, err = b.Edit(c.Message, b.Text("NoBought"), b.InlineMarkup("quick_buy"))
			handleErr(err, c.Message)
			return
		} else {
			_, err = b.Edit(c.Message, b.Text("SelectHWtoSend"), genSendHWKB(c, 0))
			handleErr(err, c.Message)
		}
	})
	b.Handle(b.InlineButton("dismiss"), func(c *tb.Callback) {
		hwID, _ := strconv.Atoi(c.Data)
		result, err := GetResult(hwID)
		handleErr(err, c.Message)
		cache.AdminCache.Set(strconv.Itoa(c.Sender.ID), AdminCache{CheckingHW: result})
		c.Message.Document.Caption = c.Message.Caption + "\n"
		_, err = b.Edit(c.Message, c.Message.Document, b.InlineMarkup("send_comment_msg"))
		handleErr(err, c.Message)
		_, err = b.Send(c.Sender, "Введите комментарий")
		err = setState(c.Sender.ID, "writeCommentReject")
		handleErr(err, c.Message)
	})
	b.Handle(b.InlineButton("apply_comment"), func(c *tb.Callback) {
		resultRaw, _ := cache.AdminCache.Get(strconv.Itoa(c.Sender.ID))
		result := resultRaw.(AdminCache)
		if result.Reject {
			_, err = b.Send(&tb.User{ID: result.CheckingHW.UserID}, b.Text("RejectedHW", result.CheckingHW))
			handleErr(err, c.Message)
			err = RemoveResult(result.CheckingHW.MessageID)
			handleErr(err, c.Message)
			err = addLesson(strconv.Itoa(result.CheckingHW.Course) + "." + strconv.Itoa(result.CheckingHW.Lesson), result.CheckingHW.UserID)
			handleErr(err, c.Message)
			_, err = b.Edit(result.PreviewMsg, b.Text("CommentSent", result.CheckingHW))
			handleErr(err, c.Message)
			err = setState(c.Sender.ID, "default")
			handleErr(err, c.Message)
		}
	})
 	b.Handle(tb.OnUserJoined, func(m *tb.Message) {
		permLevel, err := getPermLevel(m.Sender.ID)
		if err != nil {
			err = b.Ban(m.Chat, &tb.ChatMember{User: m.Sender})
			handleErr(err, m)
		}
		if permLevel < 3 {
			err = b.Ban(m.Chat, &tb.ChatMember{User: m.Sender})
			handleErr(err, m)
		}
	})

	b.Start()

}

func handleErr(err error, m *tb.Message) {
	if err != nil {
		_, _ = b.Send(m.Chat, b.Text("Error", err.Error()), tb.ModeHTML)
		panic(err.Error())
	}
}

func stringBuilder(m *tb.Message) string {
	var msgString = m.Text
	if m.Photo != nil {
		msgString = m.Caption
	}
	runes := []rune(msgString)
	for _, entity := range m.Entities {
		entityString := string(runes[entity.Offset:(entity.Offset + entity.Length)])
		if entity.Type == tb.EntityBold {
			msgString = strings.Replace(msgString, fmt.Sprintf("%v", entityString), fmt.Sprintf("<b>%v</b>", entityString), 1)
		} else if entity.Type == tb.EntityItalic {
			msgString = strings.Replace(msgString, fmt.Sprintf("%v", entityString), fmt.Sprintf("<i>%v</i>", entityString), 1)
		} else if entity.Type == tb.EntityUnderline {
			msgString = strings.Replace(msgString, fmt.Sprintf("%v", entityString), fmt.Sprintf("<u>%v</u>", entityString), 1)
		} else if entity.Type == tb.EntityStrikethrough {
			msgString = strings.Replace(msgString, fmt.Sprintf("%v", entityString), fmt.Sprintf("<s>%v</s>", entityString), 1)
		} else if entity.Type == tb.EntityCode {
			msgString = strings.Replace(msgString, fmt.Sprintf("%v", entityString), fmt.Sprintf("<code>%v</code>", entityString), 1)
		}
	}
	return msgString
}

func mainKB(m *tb.Message) *tb.ReplyMarkup {
	level, err := getPermLevel(m.Sender.ID)
	handleErr(err, m)
	if level == 3 {
		return b.Markup("admin")
	} else {
		return b.Markup("main")
	}
}

func buyHWboard() func(c *tb.Callback) {
	return func(c *tb.Callback) {
		services, err := getServices(c.Sender.ID)
		handleErr(err, c.Message)
		var sb string
		for _, service := range services {
			sb += service.Name + "\n"
		}
		_, err = b.Edit(c.Message, b.Text("HWStatus", sb), tb.ModeHTML)
		handleErr(err, c.Message)
		_, err = b.Send(c.Sender, b.Text("SelectCourse"), genKBoard(0, c.Message))
		handleErr(err, c.Message)
	}
}

func genKBoard(offset int, message *tb.Message) *tb.ReplyMarkup {
	courses, err := getCourses()
	handleErr(err, message)
	var keyboard [][]tb.InlineButton
	var crange []Service
	var rallow = true
	var lallow = true

	if len(courses) < btnCount*offset+btnCount {
		crange = courses[btnCount*offset:]
		rallow = false
	} else {
		crange = courses[btnCount*offset : btnCount*offset+btnCount]
	}

	if btnCount*offset-btnCount <= 0 {
		lallow = false
	}

	for _, course := range crange {
		cb := tb.InlineButton{
			Unique: course.ServiceID,
			Text:   course.Name,
			Data:   course.ServiceID,
		}
		keyboard = append(keyboard, []tb.InlineButton{cb})
		b.Handle(&cb, func(c *tb.Callback) {
			ce, _ := getService(c.Data)
			_, err := b.Edit(c.Message, b.Text("SelectLesson", ce), genCourseKBoard(0, c.Message, c.Data))
			handleErr(err, c.Message)
		})
	}

	if len(courses) <= btnCount {

	} else {
		keyboard = append(keyboard, genCoursesNav(rallow, lallow, offset))
	}

	keyboard = append(keyboard, []tb.InlineButton{*b.InlineButton("to_lk")})

	return &tb.ReplyMarkup{InlineKeyboard: keyboard}
}
func genCoursesNav(rallow bool, lallow bool, lastOffset int) []tb.InlineButton {
	var buttons []tb.InlineButton
	rightArr := tb.InlineButton{
		Unique: "rbBtn",
		Text:   ">>",
	}
	leftArr := tb.InlineButton{
		Unique: "lbBtn",
		Text:   "<<",
	}
	if !rallow {
		buttons = append(buttons, leftArr)
	} else if !lallow {
		buttons = append(buttons, rightArr)
	} else if !lallow && !rallow {
		return buttons
	} else {
		buttons = append(buttons, leftArr)
		buttons = append(buttons, rightArr)
	}

	b.Handle(&rightArr, func(c *tb.Callback) {
		_, err := b.Edit(c.Message, b.Text("SelectCourse"), genKBoard(lastOffset+1, c.Message))
		handleErr(err, c.Message)
	})
	b.Handle(&leftArr, func(c *tb.Callback) {
		_, err := b.Edit(c.Message, b.Text("SelectCourse"), genKBoard(lastOffset-1, c.Message))
		handleErr(err, c.Message)

	})

	return buttons
}

func genCourseKBoard(offset int, message *tb.Message, courseID string) *tb.ReplyMarkup {
	lessons, err := getServicesByCourse(courseID)
	handleErr(err, message)
	var keyboard [][]tb.InlineButton
	var crange []Service
	var rallow = true
	var lallow = true

	if len(lessons) < btnCount*offset+btnCount {
		crange = lessons[btnCount*offset:]
		rallow = false
	} else {
		crange = lessons[btnCount*offset : btnCount*offset+btnCount]
	}

	if btnCount*offset-btnCount <= 0 {
		lallow = false
	}

	for _, course := range crange {
		btn := tb.InlineButton{
			Unique: course.ServiceID,
			Text:   course.Name,
			Data:   course.ServiceID,
		}

		keyboard = append(keyboard, []tb.InlineButton{btn})

		b.Handle(&btn, func(c *tb.Callback) {
			fmt.Println("Handled")
			ce, err := getService(c.Data)
			_, err = b.Send(c.Sender, b.Text("buyLesson", ce), tb.ReplyMarkup{InlineKeyboard: [][]tb.InlineButton{{*b.InlineButton("buy_hw", c.Data)}}})
			handleErr(err, c.Message)
		})
	}

	if len(lessons) <= btnCount {

	} else {
		keyboard = append(keyboard, genLessonsNav(rallow, lallow, offset, courseID))
	}

	keyboard = append(keyboard, []tb.InlineButton{*b.InlineButton("to_courses")})

	return &tb.ReplyMarkup{InlineKeyboard: keyboard}
}
func genLessonsNav(rallow bool, lallow bool, lastOffset int, courseId string) []tb.InlineButton {
	var buttons []tb.InlineButton
	rightArr := tb.InlineButton{
		Unique: "rbBtn",
		Text:   ">>",
	}
	leftArr := tb.InlineButton{
		Unique: "lbBtn",
		Text:   "<<",
	}
	if !rallow {
		buttons = append(buttons, leftArr)
	} else if !lallow {
		buttons = append(buttons, rightArr)
	} else if !lallow && !rallow {
		return buttons
	} else {
		buttons = append(buttons, leftArr)
		buttons = append(buttons, rightArr)
	}

	b.Handle(&rightArr, func(c *tb.Callback) {
		_, err := b.Edit(c.Message, b.Text("SelectCourse"), genCourseKBoard(lastOffset+1, c.Message, courseId))
		handleErr(err, c.Message)
	})
	b.Handle(&leftArr, func(c *tb.Callback) {
		_, err := b.Edit(c.Message, b.Text("SelectCourse"), genCourseKBoard(lastOffset-1, c.Message, courseId))
		handleErr(err, c.Message)

	})

	return buttons
}

func genSendHWNav(rallow bool, lallow bool, lastOffset int) []tb.InlineButton {
	var buttons []tb.InlineButton
	rightArr := tb.InlineButton{
		Unique: "rbBtn",
		Text:   ">>",
	}
	leftArr := tb.InlineButton{
		Unique: "lbBtn",
		Text:   "<<",
	}
	if !rallow {
		buttons = append(buttons, leftArr)
	} else if !lallow {
		buttons = append(buttons, rightArr)
	} else if !lallow && !rallow {
		return buttons
	} else {
		buttons = append(buttons, leftArr)
		buttons = append(buttons, rightArr)
	}

	b.Handle(&rightArr, func(c *tb.Callback) {
		_, err = b.Edit(c.Message, b.Text("SelectHWtoSend"), genSendHWKB(c, lastOffset+1))
		handleErr(err, c.Message)
	})
	b.Handle(&leftArr, func(c *tb.Callback) {
		_, err := b.Edit(c.Message, b.Text("SelectCourse"), genSendHWKB(c, lastOffset-1))
		handleErr(err, c.Message)

	})
	return buttons
}
func genSendHWKB(c *tb.Callback, offset int) *tb.ReplyMarkup {
	services, err := getServices(c.Sender.ID)
	handleErr(err, c.Message)

	var keyboard [][]tb.InlineButton
	var crange []Service
	var rallow = true
	var lallow = true

	if len(services) < btnCount*offset+btnCount {
		crange = services[btnCount*offset:]
		rallow = false
	} else {
		crange = services[btnCount*offset : btnCount*offset+btnCount]
	}

	if btnCount*offset-btnCount <= 0 {
		lallow = false
	}

	for i, lesson := range crange {
		cb := tb.InlineButton{
			Unique:          "slesbtn"+strconv.Itoa(i),
			Text:            lesson.Name,
			Data:            lesson.ServiceID,
		}
		keyboard = append(keyboard, []tb.InlineButton{cb})
		b.Handle(&cb, func(c *tb.Callback) {
			cache.UserCache.Set(strconv.Itoa(c.Sender.ID), c.Data)
			err = setState(c.Sender.ID, "sendZip")
			handleErr(err, c.Message)
			_, err = b.Edit(c.Message, b.Text("AttachZip"), b.InlineMarkup("to_lk_kb"))
			handleErr(err, c.Message)
		})
	}

	if len(services) <= btnCount {

	} else {
		keyboard = append(keyboard, genSendHWNav(rallow, lallow, offset))
	}

	keyboard = append(keyboard, []tb.InlineButton{*b.InlineButton("to_lk")})

	return &tb.ReplyMarkup{InlineKeyboard: keyboard}

}

func getBotStats(m *tb.Message) string {
	var usersCount int
	var active int
	var canSendHW int
	var premium int
	var stopped int
	var sentFeedback int

	err = db.QueryRow(`SELECT COUNT(*) FROM users`).Scan(&usersCount)
	handleErr(err, m)

	err = db.QueryRow(`SELECT COUNT(*) FROM users WHERE cardinality(courses) > 0`).Scan(&canSendHW)
	handleErr(err, m)

	err = db.QueryRow(`SELECT COUNT(*) FROM users WHERE premium = true`).Scan(&premium)
	handleErr(err, m)

	err = db.QueryRow(`SELECT COUNT(*) FROM users WHERE isblocked = true`).Scan(&stopped)
	handleErr(err, m)

	err = db.QueryRow(`SELECT COUNT(*) FROM users WHERE feedbackmessage > 0`).Scan(&active)
	handleErr(err, m)

	err = db.QueryRow(`SELECT SUM(messagesfeedbacked) FROM users`).Scan(&sentFeedback)
	handleErr(err, m)

	stringBuilder := fmt.Sprintf(`<b>Статистика бота</b>:
<b>Количество пользователей:</b> %d

<b>Пользователи:</b>
Которые оплачивали услуги: %d
С возможностью отправки дз: %d

<b>Сообщения:</b>
Отправленные в обратную связь: %d
Активных сообщений: %d

Забанили бота: %d
`,
		usersCount, premium, canSendHW, sentFeedback, active, stopped)
	return stringBuilder
}
func getCourseStats(m *tb.Message) string {
	var sb string
	sb += "Статистика по курсам\n\n"
	courses, err := getCourses()
	handleErr(err, m)
	for _, course := range courses {
		sb += fmt.Sprintf("<b>%v. %v</b>\nВсего: %d\n", course.ServiceID, course.Name, course.Bought)
		lessons, err := getServicesByCourse(course.ServiceID)
		handleErr(err, m)
		for _, lesson := range lessons {
			sb += fmt.Sprintf("%v - %d\n", lesson.ServiceID, lesson.Bought)
		}
		sb += "\n"
	}
	return sb
}
