package bot

import (
	"bytes"
	"encoding/json"
	"fmt"
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

func SendStartMessage(chatId int, name string) {
	message := fmt.Sprintf("Xin chào %s \n\nGõ <code>/help</code> để xem danh sách các lệnh mà bot hỗ trợ nhé.\n\nBạn cũng có thể truy cập nhanh các chức năng bằng cách nhấn nút Menu bên dưới.", name)
	SendMessage(chatId, message)
}

func SendHelpMessage(chatId int) {
	messages := ""
	botCommands := GetBotCommands()

	for _, command := range botCommands.Result {
		messages += fmt.Sprintf("<code>/%s</code> - %s\n\n", command.Command, command.Description)
	}

	SendMessage(chatId, messages)
}

func SendGroupId(chatId int) {
	SendMessage(chatId, fmt.Sprintf("Group ID: <code>%v</code>", chatId))
}

func SendMessage(chatId int, message string) {
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
		log.Fatalln(err)
	}

	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)

	if err != nil {
		log.Fatalln(err)
	}

	var status structs.Status

	err = json.Unmarshal(body, &status)
	if err != nil {
		log.Fatalln(err)
	}

	if status.Ok == false {
		log.Fatalln(string(body))
	}

	log.Println("SendMessage OK")
}

func SendAPhoto(chatId int, path string) {
	uri := API_URL + "/sendPhoto"
	method := "GET"

	payload := &bytes.Buffer{}
	writer := multipart.NewWriter(payload)
	writer.WriteField("chat_id", strconv.Itoa(chatId))

	file, errFile2 := os.Open(path)
	defer file.Close()

	part2, errFile2 := writer.CreateFormFile("photo", filepath.Base(path))
	_, errFile2 = io.Copy(part2, file)
	if errFile2 != nil {
		log.Fatalln(errFile2)
	}

	err := writer.Close()
	if err != nil {
		log.Fatalln(err)
	}

	client := &http.Client{}
	req, err := http.NewRequest(method, uri, payload)

	if err != nil {
		log.Fatalln(err)
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())
	res, err := client.Do(req)
	if err != nil {
		log.Fatalln(err)
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Fatalln(err)
	}

	var status structs.Status

	err = json.Unmarshal(body, &status)
	if err != nil {
		log.Fatalln(err)
	}

	if status.Ok == false {
		log.Fatalln(string(body))
	}

	log.Println("SendAPhoto OK")
}

type BodyReplyMarkup struct {
	ReplyMarkup struct {
		InlineKeyboard [][]ButtonCallback `json:"inline_keyboard"`
	} `json:"reply_markup"`
}

type ButtonCallback struct {
	Text         string `json:"text"`
	CallbackData string `json:"callback_data"`
}

func SendMessageWithReplyMarkup(chatId int, message string, replyMark []ButtonCallback) {
	uri := API_URL + "/sendMessage"

	var markup BodyReplyMarkup
	markup.ReplyMarkup.InlineKeyboard = append(markup.ReplyMarkup.InlineKeyboard, replyMark)
	marshalled, err := json.Marshal(markup)

	req, err := http.NewRequest("GET", uri, bytes.NewReader(marshalled))
	req.Header.Add("Content-Type", "application/json")
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
		log.Fatalln(err)
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Fatalln(err)
	}

	var status structs.Status

	err = json.Unmarshal(body, &status)
	if err != nil {
		log.Fatalln(err)
	}

	if status.Ok == false {
		log.Fatalln(string(body))
	}

	log.Println("SendMessageWithReplyMarkup OK")
}

func SetTypingAction(chatId int) {
	uri := API_URL + "/sendChatAction"
	req, err := http.NewRequest("GET", uri, nil)

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
		log.Fatalln(err)
	}

	if body != nil {
		log.Println("SetTypingAction OK")
	}
}

func GetBotCommands() structs.BotCommands {
	uri := API_URL + "/getMyCommands"
	req, err := http.NewRequest("GET", uri, nil)

	if err != nil {
		log.Fatalln(err)
	}

	client := &http.Client{}

	res, err := client.Do(req)

	if err != nil {
		log.Fatalln(err)
	}

	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)

	if err != nil {
		log.Fatalln(err)
	}

	var botCommands structs.BotCommands

	err = json.Unmarshal(body, &botCommands)
	if err != nil {
		log.Fatalln(err)
	}

	if botCommands.Ok == false {
		log.Fatalln(string(body))
	}

	log.Println("GetBotCommands OK")
	return botCommands
}
