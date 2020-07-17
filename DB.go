package main

import (
	"database/sql"
	"errors"
	"fmt"
	tb "github.com/demget/telebot"
	"github.com/lib/pq"
	_ "github.com/lib/pq"
	"log"
	"strings"
	"time"
)

var db = dbInit()
var err error
var connString = "postgres://urxjvjvplarwhw:5826200bbb15c18a04f4bcbebbd0f4ea510aace28faa52d840670e90acd19cc8@ec2-46-137-100-204.eu-west-1.compute.amazonaws.com:5432/d844hr65721plp"
var void interface{}

// Init connection to db
func dbInit() *sql.DB {
	dbr, err := sql.Open("postgres", connString)
	if err != nil {
		log.Fatal(err)
	}
	return dbr
}

// Create users table
func initTable() error {
	_, err = db.Exec(`
	CREATE TABLE users (
	    UID SERIAL,
	    userID INTEGER,
	    username VARCHAR(64),
	    mail VARCHAR(64),
	    permissionLevel INTEGER,
	    premium BOOLEAN,
	    groupID INTEGER,
	    hwQuestionMessage INTEGER,
	    feedbackMessage INTEGER,
	    hwMessage INTEGER,
	    messagesSent INTEGER,
	    hwCount INTEGER,
	    courses INTEGER[],
	    state VARCHAR(64),
	    isBlocked BOOLEAN,
	    isBanned BOOLEAN
	)`)

	if err != nil {
		return err
	}
	return nil
}

func initTableConfig() error {
	_, err = db.Exec(`
	CREATE TABLE configs (
		KEY VARCHAR(64),
		VALUE VARCHAR(64)
	)`)

	if err != nil {
		return err
	}
	return nil
}

// State machine
func getState(userID int) (string, error) {
	var state string

	row := db.QueryRow(`SELECT state FROM users WHERE userID = $1`, userID)
	err = row.Scan(&state)
	if err != nil {
		return "", err
	}
	return state, nil
}

func setState(userID int, state string) error {
	_, err = db.Exec(`UPDATE users SET state = $1 WHERE userID = $2`, state, userID)
	if err != nil {
		return err
	}
	return nil
}

// ============================USERS==========================

func checkUser(userID int) (bool, error) {
	var user string

	row := db.QueryRow(`SELECT userid FROM users WHERE userid = $1`, userID)
	err = row.Scan(&user)
	if err != nil {
		if strings.Contains(err.Error(), "sql: no rows in result set") {
			return false, nil
		} else {
			return false, err
		}
	}
	return true, nil
}

func addUser(userID int, username string) error {
	_, err = db.Exec(`INSERT INTO users (userID, username, mail, permissionLevel, premium, groupID, hwQuestionMessages, feedbackMessage, hwMessage,
                   messagesSent, hwCount, courses, state, isBlocked, isBanned) values ($1, $2, 'none', 0, false, 0, 0, 0, 0, 1, 0, NULL, 'default', false, false)`, userID, username)
	if err != nil {
		return err
	}
	return nil
}

func getPermLevel(userID int) (int, error) {
	var permLevel int

	row := db.QueryRow(`SELECT permissionlevel FROM users WHERE userID = $1`, userID)
	err = row.Scan(&permLevel)
	if err != nil {
		return -1, err
	}
	return permLevel, nil
}

func getFeedbackMessage(userID int) (int, error) {
	var fbMessage int

	row := db.QueryRow(`SELECT feedbackmessage FROM users WHERE userID = $1`, userID)
	err = row.Scan(&fbMessage)
	if err != nil {
		return -1, err
	}
	return fbMessage, nil
}

func setFeedbackMessage(userID int, messageID int) error {
	_, err = db.Exec(`UPDATE users SET feedbackmessage = $1 WHERE userID = $2`, messageID, userID)
	if err != nil {
		return err
	}
	return nil
}

func setHWQuestion(userID int, messageID int) error {
	_, err = db.Exec(`UPDATE users SET hwquestionmessages = array_append(hwquestionmessages, $1) WHERE userid = $2`, messageID, userID)
	if err != nil {
		return err
	}
	return nil
}

func resetFeedbackMessage(userID int) error {
	_, err = db.Exec(`UPDATE users SET feedbackmessage = 0 WHERE userID = $1`, userID)
	if err != nil {
		return err
	}
	return nil
}

func isBlocked(userID int) (bool, error) {
	var blocked bool

	row := db.QueryRow(`SELECT isblocked FROM users WHERE userid = $1`, userID)
	err = row.Scan(&blocked)
	if err != nil {
		return false, err
	}
	return blocked, nil
}

func getUserByMessage(messageID int) (int, error) {
	var userID int

	row := db.QueryRow(`SELECT userid FROM users WHERE feedbackmessage = $1`, messageID)
	err = row.Scan(&userID)
	if err != nil {
		return -1, err
	}
	return userID, nil
}

func getEmail(userID int) (string, error) {
	var mail string

	row := db.QueryRow(`SELECT mail FROM users WHERE userid = $1`, userID)
	err = row.Scan(&mail)
	if err != nil {
		return "", err
	}
	if mail == "none" {
		return "", nil
	}
	return mail, nil
}

func setEmail(userID int, mail string) error {
	_, err = db.Exec(`UPDATE users SET mail = $1 WHERE userID = $2`, mail, userID)
	if err != nil {
		return err
	}
	return nil
}

func getUsers(userType string) ([]int, error) {
	var userID int
	var userIDS []int
	var rows *sql.Rows

	if userType == "default" {
		rows, err = db.Query(`SELECT userid FROM users WHERE isbanned = false AND isblocked = false`)
		if err != nil {
			return []int{}, err
		}
	} else if userType == "premium" {
		rows, err = db.Query(`SELECT userid FROM users WHERE isbanned = false AND isblocked = false AND premium = true`)
		if err != nil {
			return []int{}, err
		}
	}

	for rows.Next() {
		err = rows.Scan(&userID)
		if err != nil {
			return []int{}, err
		}
		userIDS = append(userIDS, userID)
	}
	return userIDS, nil
}

func addBlocked(userID int) error {
	_, err = db.Exec(`UPDATE users SET isblocked = $1 WHERE userID = $2`, true, userID)
	if err != nil {
		return err
	}
	return nil
}

func removeBlocked(userID int) error {
	_, err = db.Exec(`UPDATE users SET isblocked = $1 WHERE userID = $2`, false, userID)
	if err != nil {
		return err
	}
	return nil
}

func getUnReplied() (int, error) {
	var unreplied int
	err = db.QueryRow(`SELECT COUNT(*) FROM users WHERE feedbackmessage > 0`).Scan(&unreplied)
	if err != nil {
		return 0, err
	}
	return unreplied, nil
}

func isPremium(userID int) (bool, error) {
	var premium bool

	row := db.QueryRow(`SELECT premium FROM users WHERE userid = $1`, userID)
	err = row.Scan(&premium)
	if err != nil {
		return false, err
	}

	return premium, nil
}

func isNowPremium(userID int) (bool, error) {
	var arrlen int

	row := db.QueryRow(`SELECT array_length(courses, 1) FROM users WHERE userid = $1`, userID)
	err = row.Scan(&arrlen)
	if err != nil {
		return false, err
	}
	if arrlen > 0 {
		return true, nil
	} else if arrlen <= 0 {
		return false, nil
	}
	return false, nil
}

func getHWChat(userID int) (int, error) {
	var chatId int

	row := db.QueryRow(`SELECT groupid FROM users WHERE userid = $1`, userID)
	err = row.Scan(&chatId)
	if err != nil {
		return 0, err
	}
	return chatId, nil
}

func getUserByHWQ(messageID int, chatID int64) (int, error) {
	var userID int

	row := db.QueryRow(`SELECT userid FROM users WHERE groupid = $1 AND $2 = ANY (hwquestionmessages)`, chatID, messageID)
	err = row.Scan(&userID)
	if err != nil {
		return 0, err
	}
	return userID, nil
}

func removeHWQ(messageID int, userID int) error {
	_, err = db.Exec(`UPDATE users SET hwquestionmessages = array_remove(hwquestionmessages, $1) WHERE userID = $2`, messageID, userID)
	if err != nil {
		return err
	}
	return nil
}

// =========================STATISTICS========================

func incrementMessage(userID int) error {
	_, err = db.Exec(`UPDATE users SET messagessent = messagessent + 1 WHERE userid = $1`, userID)
	if err != nil {
		return err
	}
	return nil
}

func incrementFeedBackMessage(userID int) error {
	_, err = db.Exec(`UPDATE users SET messagesfeedbacked = messagesfeedbacked + 1 WHERE userid = $1`, userID)
	if err != nil {
		return err
	}
	return nil
}

func getBlocked() (int, error) {
	var userID int
	var usersCount int

	rows, err := db.Query(`SELECT userid FROM users WHERE isblocked = true`)
	if err != nil {
		return 0, err
	}
	for rows.Next() {
		err = rows.Scan(&userID)
		if err != nil {
			return 0, err
		}
		usersCount += 1
	}
	return usersCount, nil
}

// =========================CONFIGS===========================

func getConfig(key string) (string, error) {
	var value string
	row := db.QueryRow(`SELECT value FROM configs WHERE key = $1`, key)
	err = row.Scan(&value)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", err
		} else {
			log.Panicf("PANIC: UNABLE TO GET CONFIG (%v)", err.Error())
		}
	}
	return value, nil
}

func setConfig(key, value string) error {
	_, err = db.Exec(`INSERT INTO configs(key, value) VALUES ($1, $2)`, key, value)
	if err != nil {
		return err
	}
	return nil
}

// =========================MAILINGS==========================

func getMailings() ([]Mailing, error) {
	var mailing Mailing
	var mailings []Mailing
	var id int

	rows, err := db.Query(`SELECT * FROM mailings`)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		err = rows.Scan(&mailing.UsersType, &mailing.Time, &mailing.Text, &mailing.MailingStatus, &mailing.Photo, &id)
		if err != nil {
			return nil, err
		}
		mailings = append(mailings, mailing)
	}
	return mailings, nil
}

func getAdmins() ([]int, error) {
	var userID int
	var userIDS []int

	rows, err := db.Query(`SELECT userid FROM users WHERE permissionlevel = 3`)
	if err != nil {
		return []int{}, err
	}
	for rows.Next() {
		err = rows.Scan(&userID)
		if err != nil {
			return []int{}, err
		}
		userIDS = append(userIDS, userID)
	}
	return userIDS, nil
}

func addMailing(m *tb.Message, timestamp time.Time, userType string, status string) error {
	if userType == "всем" {
		userType = "default"
	} else if userType == "покупателям" {
		userType = "premium"
	}

	defer db.Close()
	if m.Photo == nil {
		_, err = db.Exec(`INSERT INTO mailings(type, time, text, status, photo) VALUES('default', $1, $2, $3, '')`, timestamp, stringBuilder(m), status)
		if err != nil {
			return err
		}
	} else {
		_, err = db.Exec(`INSERT INTO mailings(type, time, text, status, photo) VALUES('default', $1, $2, $3, $4)`, timestamp, stringBuilder(m), status, m.Photo.FileID)
		if err != nil {
			return err
		}
	}
	return nil
}

func getMailing(text string) (Mailing, error) {
	var mailing Mailing
	var id int
	err = db.QueryRow(`SELECT * FROM mailings WHERE text = $1`, text).Scan(&mailing.UsersType, &mailing.Time, &mailing.Text, &mailing.MailingStatus, &mailing.Photo, &id)
	if err != nil {
		if strings.Contains(err.Error(), "sql: no rows in result set") {
			return Mailing{}, errors.New("Рассылки не существует")
		} else {
			return Mailing{}, err
		}
	}
	return mailing, nil
}

func removeMailing(text string) error {
	_, err = db.Exec(`DELETE FROM mailings WHERE text = $1`, text)
	if err != nil {
		return err
	}
	return nil
}

func changeStatus(text string, status string) error {
	_, err = db.Exec(`UPDATE mailings SET status = $1 WHERE text = $2`, status, text)
	if err != nil {
		return err
	}
	return nil
}


// ========================SERVICES===========================

func getServices(userID int) ([]Service, error) {
	var serviceIDs []string
	var services []Service

	err = db.QueryRow("SELECT courses from  users WHERE userid = $1", userID).Scan((*pq.StringArray)(&serviceIDs))
	if err != nil {
		return []Service{}, err
	}

	for _, serviceID := range serviceIDs {
		var service Service
		err = db.QueryRow("SELECT * FROM courses WHERE serviceid = $1", serviceID).Scan(&service.ServiceID, &service.Type, &service.Description, &service.FileURI, &service.Price, &service.ArticleURL, &service.Name, &void, &service.Bought)
		if err != nil {
			if strings.Contains(err.Error(), "no rows") {
				continue
			} else {
				return []Service{}, err
			}
		}
		services = append(services, service)
	}
	return services, nil
}

func getService(serviceID string) (Service, error) {
	var service Service
	err = db.QueryRow("SELECT * FROM courses WHERE serviceid = $1", serviceID).Scan(&service.ServiceID, &service.Type, &service.Description, &service.FileURI, &service.Price, &service.ArticleURL, &service.Name, &void, &service.Bought)
	if err != nil {
		return Service{}, err
	}

	return service, nil

}

func getCourses() ([]Service, error) {
	var service Service
	var services []Service

	rows, err := db.Query("SELECT * FROM courses WHERE type = 'course' ORDER BY id")
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		err = rows.Scan(&service.ServiceID, &service.Type, &service.Description, &service.FileURI, &service.Price, &service.ArticleURL, &service.Name, &void, &service.Bought)
		if err != nil {
			return nil, err
		}
		services = append(services, service)
	}
	return services, nil
}

func getServicesByCourse(courseID string) ([]Service, error) {
	var service Service
	var services []Service

	rows, err := db.Query("SELECT * FROM courses WHERE serviceid LIKE $1 ORDER BY id", fmt.Sprintf("%v.%%", courseID))
	if err != nil {
		return []Service{}, err
	}
	for rows.Next() {
		err = rows.Scan(&service.ServiceID, &service.Type, &service.Description, &service.FileURI, &service.Price, &service.ArticleURL, &service.Name, &void, &service.Bought)
		if err != nil {
			return []Service{}, err
		}
		services = append(services, service)
	}

	return services, nil
}

func getBoughtByCourse(courseID string, userID int) ([]Service, error) {
	var serviceIDs []string
	var services []Service

	err = db.QueryRow("SELECT courses from  users WHERE userid = $1", userID).Scan((*pq.StringArray)(&serviceIDs))
	if err != nil {
		return []Service{}, err
	}

MainLoop:
	for _, serviceID := range serviceIDs {

		if !strings.Contains(serviceID, ".") {
			if serviceID == courseID {
				courseServices, err := getServicesByCourse(serviceID)
				if err != nil {
					if strings.Contains(err.Error(), "no rows") {
						continue
					} else {
						return []Service{}, err
					}
				}
				services = append(services, courseServices...)
			} else {
				continue MainLoop
			}
		} else if strings.Split(serviceID, ".")[0] != courseID {
			continue MainLoop
		}
		var service Service

		err = db.QueryRow("SELECT * FROM courses WHERE serviceid = $1", serviceID).Scan(&service.ServiceID, &service.Type, &service.Description, &service.FileURI, &service.Price, &service.ArticleURL, &service.Name, &void, &service.Bought)
		if err != nil {
			if strings.Contains(err.Error(), "no rows") {
				continue
			} else {
				return []Service{}, err
			}
		}
		if service.Type == "lesson" {
			services = append(services, service)
		}
	}
	return services, nil
}

func isBought(serviceID string, userID int) (bool, error) {
	var serviceIDs []string

	err = db.QueryRow("SELECT courses from  users WHERE userid = $1", userID).Scan((*pq.StringArray)(&serviceIDs))
	if err != nil {
		return false, err
	}

	for _, SID := range serviceIDs {
		if strings.Contains(serviceID, ".") {
			if strings.Split(serviceID, ".")[0] == SID {
				return true, nil
			}
		}
		if SID == serviceID {
			return true, nil
		}
	}

	return false, nil
}

func addService(service Service) error {
	_, err = db.Exec(`INSERT INTO courses(serviceid, type, description, fileurl, price, "articleURL", name) VALUES ($1, $2, $3, $4, $5, $6, $7)`, service.ServiceID, "pending", service.Description, service.FileURI, service.Price, service.ArticleURL, service.Name)
	if err != nil {
		return err
	}
	return nil
}

// =========================STRUCTS===========================

type Mailing struct {
	UsersType     string
	Time          time.Time
	Text          string
	MailingStatus string
	Photo         string
}

type Service struct {
	ServiceID   string
	Type        string
	Description string
	FileURI     string
	Price       int
	ArticleURL  string
	Name        string
	Bought      int
}
