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

const BASE_URL = "https://api.telegram.org"

type Telegram struct {
	client   *http.Client
	endpoint string
}

func New(http *http.Client, baseUrl string) *Telegram {
	endpoint := baseUrl + "/bot" + os.Getenv("TELE_BOT_TOKEN")
	return &Telegram{client: http, endpoint: endpoint}
}

func (t Telegram) SendMessage(data structs.DataTele) error {
	message := data.ReplyMessage
	if data.ChatType == "group" {
		message += "\n@" + data.Username
	}

	url := t.endpoint + "/sendMessage"
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return err
	}

	q := req.URL.Query()
	q.Add("chat_id", strconv.Itoa(data.ChatId))
	q.Add("text", message)
	q.Add("parse_mode", "html")

	req.URL.RawQuery = q.Encode()

	res, err := t.client.Do(req)

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

func (t Telegram) SendPhoto(data structs.DataTele, path string) error {
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

	url := t.endpoint + "/sendPhoto"
	req, _ := http.NewRequest("GET", url, payload)
	req.Header.Set("Content-Type", writer.FormDataContentType())

	res, err := t.client.Do(req)
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

func (t Telegram) SendMessageWithReplyMarkup(data structs.DataTele, replyMark []structs.ButtonCallback) error {
	var markup structs.BodyReplyMarkup
	markup.ReplyMarkup.InlineKeyboard = append(markup.ReplyMarkup.InlineKeyboard, replyMark)
	marshalled, err := json.Marshal(markup)
	if err != nil {
		return err
	}

	url := t.endpoint + "/sendMessage"
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

	res, err := t.client.Do(req)
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

func (t Telegram) SetTypingAction(data structs.DataTele) error {
	chatId := data.ChatId

	url := t.endpoint + "/sendChatAction"
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return err
	}

	q := req.URL.Query()
	q.Add("chat_id", strconv.Itoa(chatId))
	q.Add("action", "typing")

	req.URL.RawQuery = q.Encode()

	res, err := t.client.Do(req)
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

func (t Telegram) GetBotCommands() (*structs.BotCommands, error) {
	url := t.endpoint + "/getMyCommands"
	req, err := http.NewRequest("GET", url, nil)

	if err != nil {
		return nil, err
	}

	res, err := t.client.Do(req)
	if err != nil {
		return nil, err
	}

	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	var botCommands structs.BotCommands

	err = json.Unmarshal(body, &botCommands)
	if err != nil {
		return nil, err
	}

	if botCommands.Ok {
		log.Println("GetBotCommands OK")

		return &botCommands, nil
	}

	return nil, errors.New(string(body))
}
