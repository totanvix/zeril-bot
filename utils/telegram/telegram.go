package telegram

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"io/ioutil"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"zeril-bot/utils/structs"
)

var API_URL string = "https://api.telegram.org/bot" + os.Getenv("TELE_BOT_TOKEN")

func getApiURL() string {
	return "https://api.telegram.org/bot" + os.Getenv("TELE_BOT_TOKEN")
}

type Data struct {
	ChatId   int
	ChatType string
	Username string
	Message  string
}

func SendMessage(data Data) error {
	message := data.Message
	if data.ChatType == "group" {
		message = message + "\n@" + data.Username
	}

	uri := API_URL + "/sendMessage"
	req, err := http.NewRequest("GET", uri, nil)

	if err != nil {
		return err
	}

	q := req.URL.Query()
	q.Add("chat_id", strconv.Itoa(data.ChatId))
	q.Add("text", message)
	q.Add("parse_mode", "html")

	req.URL.RawQuery = q.Encode()

	client := &http.Client{}

	res, err := client.Do(req)

	if err != nil {
		return err
	}

	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)

	if err != nil {
		return err
	}

	var status structs.Status

	err = json.Unmarshal(body, &status)
	if err != nil {
		return err
	}

	if status.Ok {
		log.Println("SendMessage OK")
		return nil
	}

	return errors.New(string(body))
}

func SendPhoto(chatId int, path string) {
	uri := API_URL + "/sendPhoto"

	file, _ := os.Open(path)
	defer file.Close()

	payload := &bytes.Buffer{}
	writer := multipart.NewWriter(payload)
	writer.WriteField("chat_id", strconv.Itoa(chatId))

	// if chatType == "group" {
	// 	writer.WriteField("caption", "@"+chatFrom.Username)
	// }

	part, _ := writer.CreateFormFile("photo", filepath.Base(path))
	io.Copy(part, file)

	writer.Close()

	req, _ := http.NewRequest("GET", uri, payload)
	req.Header.Set("Content-Type", writer.FormDataContentType())

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		log.Panic(err)
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Panic(err)
	}

	var status structs.Status

	err = json.Unmarshal(body, &status)
	if err != nil {
		log.Panic(err)
	}

	if status.Ok == false {
		log.Panic(string(body))
	}

	log.Println("SendPhoto OK")
}

func SendMessageWithReplyMarkup(chatId int, message string, replyMark []structs.ButtonCallback) {
	uri := API_URL + "/sendMessage"

	var markup structs.BodyReplyMarkup
	markup.ReplyMarkup.InlineKeyboard = append(markup.ReplyMarkup.InlineKeyboard, replyMark)
	marshalled, err := json.Marshal(markup)

	req, err := http.NewRequest("GET", uri, bytes.NewReader(marshalled))
	req.Header.Add("Content-Type", "application/json")
	if err != nil {
		log.Println(err)
		return
	}
	// if chatType == "group" {
	// 	message = message + "\n@" + chatFrom.Username
	// }

	q := req.URL.Query()
	q.Add("chat_id", strconv.Itoa(chatId))
	q.Add("text", message)
	q.Add("parse_mode", "html")

	req.URL.RawQuery = q.Encode()

	client := &http.Client{}

	res, err := client.Do(req)
	if err != nil {
		log.Panic(err)
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Panic(err)
	}

	var status structs.Status

	err = json.Unmarshal(body, &status)
	if err != nil {
		log.Panic(err)
	}

	if status.Ok == false {
		log.Panic(string(body))
	}

	log.Println("SendMessageWithReplyMarkup OK")
}
