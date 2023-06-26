package bot

import (
	"bytes"
	"encoding/json"
	"io"
	"io/ioutil"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	// "zeril-bot/utils/quote"
	"zeril-bot/utils/quote"
	"zeril-bot/utils/structs"
	"zeril-bot/utils/telegram"
)

var API_URL string = "https://api.telegram.org/bot" + os.Getenv("TELE_BOT_TOKEN")

var chatType string
var chatFrom structs.From

type Bot struct {
	HookData  structs.HookData
	typingCh  chan struct{}
	commandCh chan error
}

func NewBot(hookData structs.HookData) *Bot {
	tCh := make(chan struct{})
	commandCh := make(chan error)

	return &Bot{hookData, tCh, commandCh}
}

func (b Bot) ResolveHook() error {
	go b.setTypingAction()

	var err error

	switch {
	case b.isCommand():
		go b.resolveCommand()
	case b.isCallbackCommand():
		go b.resolveCallbackCommand()
	}

	// to do refactor
	<-b.typingCh
	err = <-b.commandCh

	return err
}

func (b Bot) resolveCommand() error {
	var err error

	defer func() {
		b.commandCh <- err
	}()

	data := b.HookData
	name := data.Message.From.FirstName
	chatId := data.Message.Chat.ID
	text := data.Message.Text
	arr := strings.Fields(text)

	log.Printf("Yêu cầu từ bạn %s: %s", name, text)

	tData := telegram.Data{
		ChatId:   chatId,
		ChatType: "",
		Username: "",
	}

	command := arr[0]

	switch command {
	case "/quote", "/quote@zerill_bot":
		err = quote.SendAQuote(tData)
	default:
		// channel.SendMessage(chatId, "Tôi không hiểu câu lệnh của bạn !!!")
	}

	return nil
}

func (b Bot) resolveCallbackCommand() error {
	return nil
}

func (b Bot) isCommand() bool {
	return b.HookData.CallbackQuery.Data == ""
}

func (b Bot) isCallbackCommand() bool {
	return b.HookData.CallbackQuery.Data != ""
}

func (b Bot) getApiURL() string {
	return "https://api.telegram.org/bot" + os.Getenv("TELE_BOT_TOKEN")
}

func SendMessage(chatId int, message string) {
	if chatType == "group" {
		message = message + "\n@" + chatFrom.Username
	}

	uri := API_URL + "/sendMessage"
	req, err := http.NewRequest("GET", uri, nil)

	if err != nil {
		log.Println(err)
		return
	}

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

	log.Println("SendMessage OK")
}

func SendPhoto(chatId int, path string) {
	uri := API_URL + "/sendPhoto"

	file, _ := os.Open(path)
	defer file.Close()

	payload := &bytes.Buffer{}
	writer := multipart.NewWriter(payload)
	writer.WriteField("chat_id", strconv.Itoa(chatId))

	if chatType == "group" {
		writer.WriteField("caption", "@"+chatFrom.Username)
	}

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
	if chatType == "group" {
		message = message + "\n@" + chatFrom.Username
	}

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

func (b Bot) setTypingAction() {
	defer close(b.typingCh)

	url := b.getApiURL()
	chatId := b.HookData.Message.Chat.ID

	req, err := http.NewRequest("GET", url+"/sendChatAction", nil)
	if err != nil {
		log.Println(err)
		return
	}

	q := req.URL.Query()
	q.Add("chat_id", strconv.Itoa(chatId))
	q.Add("action", "typing")

	req.URL.RawQuery = q.Encode()

	client := &http.Client{}

	res, err := client.Do(req)

	if err != nil {
		log.Println(err)
		return
	}

	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)

	if err != nil {
		log.Panic(err)
	}

	if body != nil {
		log.Println("SetTypingAction OK")
	}

	// channel.GetWg().Done()
}

func GetBotCommands() structs.BotCommands {
	uri := API_URL + "/getMyCommands"
	req, err := http.NewRequest("GET", uri, nil)

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

func SetChatFrom(chat structs.From) {
	chatFrom = chat
}
func SetChatType(t string) {
	chatType = t
}
