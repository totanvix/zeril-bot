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

func getApiURL(t string) string {
	url := "https://api.telegram.org/bot" + os.Getenv("TELE_BOT_TOKEN")

	switch t {
	case "sendMessage":
		return url + "/sendMessage"
	case "sendPhoto":
		return url + "/sendPhoto"
	case "sendChatAction":
		return url + "/sendChatAction"
	case "getMyCommands":
		return url + "/getMyCommands"
	default:
		return ""
	}
}

func SendMessage(data structs.DataTele) error {
	message := data.ReplyMessage
	if data.ChatType == "group" {
		message += "\n@" + data.Username
	}

	url := getApiURL("sendMessage")
	req, err := http.NewRequest("GET", url, nil)

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

func SendPhoto(data structs.DataTele, path string) error {
	file, _ := os.Open(path)
	defer file.Close()

	payload := &bytes.Buffer{}
	writer := multipart.NewWriter(payload)
	writer.WriteField("chat_id", strconv.Itoa(data.ChatId))

	if data.ChatType == "group" {
		writer.WriteField("caption", "@"+data.Username)
	}

	part, _ := writer.CreateFormFile("photo", filepath.Base(path))
	io.Copy(part, file)

	writer.Close()

	url := getApiURL("sendPhoto")
	req, _ := http.NewRequest("GET", url, payload)
	req.Header.Set("Content-Type", writer.FormDataContentType())

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
		log.Println("SendPhoto OK")
		return nil
	}

	return errors.New(string(body))
}

func SendMessageWithReplyMarkup(data structs.DataTele, replyMark []structs.ButtonCallback) error {
	var markup structs.BodyReplyMarkup
	markup.ReplyMarkup.InlineKeyboard = append(markup.ReplyMarkup.InlineKeyboard, replyMark)
	marshalled, err := json.Marshal(markup)
	if err != nil {
		return err
	}

	url := getApiURL("sendMessage")
	req, err := http.NewRequest("GET", url, bytes.NewReader(marshalled))
	req.Header.Add("Content-Type", "application/json")
	if err != nil {
		return err
	}

	message := data.ReplyMessage
	if data.ChatType == "group" {
		message += "\n@" + data.Username
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
		log.Println("SendMessageWithReplyMarkup OK")
		return nil
	}

	return errors.New(string(body))
}

func SetTypingAction(data structs.DataTele) error {
	chatId := data.ChatId

	url := getApiURL("sendChatAction")
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return err
	}

	q := req.URL.Query()
	q.Add("chat_id", strconv.Itoa(chatId))
	q.Add("action", "typing")

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

	if body != nil {
		log.Println("SetTypingAction OK")
	}

	return nil
}

func GetBotCommands() structs.BotCommands {
	url := getApiURL("getMyCommands")
	req, err := http.NewRequest("GET", url, nil)

	if err != nil {
		log.Panic(err)
	}

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

	var botCommands structs.BotCommands

	err = json.Unmarshal(body, &botCommands)
	if err != nil {
		log.Panic(err)
	}

	if botCommands.Ok == false {
		log.Fatalln(string(body))
	}

	log.Println("GetBotCommands OK")
	return botCommands
}