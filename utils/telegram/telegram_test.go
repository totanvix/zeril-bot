package telegram

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"reflect"
	"testing"
	"time"
	"zeril-bot/utils/structs"
)

func TestNew(t *testing.T) {
	token := "123123:test"
	t.Setenv("TELE_BOT_TOKEN", token)

	endpoint := BASE_URL + "/bot" + os.Getenv("TELE_BOT_TOKEN")

	type args struct {
		http    *http.Client
		baseUrl string
	}

	tests := []struct {
		name string
		args args
		want *Telegram
	}{
		{
			name: "Test Nil",
			args: args{
				http:    nil,
				baseUrl: BASE_URL,
			},
			want: &Telegram{
				client:   nil,
				endpoint: endpoint,
			},
		},
		{
			name: "Test happy case",
			args: args{
				http:    &http.Client{},
				baseUrl: BASE_URL,
			},
			want: &Telegram{
				client:   &http.Client{},
				endpoint: endpoint,
			},
		},
		{
			name: "Test with client config",
			args: args{
				http: &http.Client{
					Timeout: 30 * time.Second,
				},
				baseUrl: BASE_URL,
			},
			want: &Telegram{
				client: &http.Client{
					Timeout: 30 * time.Second,
				},
				endpoint: endpoint,
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got := New(test.args.http, test.args.baseUrl)
			if !reflect.DeepEqual(got, test.want) {
				t.Errorf("New() = %v, want %v", got, test.want)
			}
		})
	}
}

func TestSendMessage(t *testing.T) {
	token := "123123:test"
	t.Setenv("TELE_BOT_TOKEN", token)

	t.Run("Test Happy Case", func(t *testing.T) {
		server := createHappyServer(t)
		defer server.Close()

		telegram := New(server.Client(), server.URL)
		err := telegram.SendMessage(structs.DataTele{
			ChatId:       123,
			Username:     "totanvix",
			ChatType:     "private",
			ReplyMessage: "test",
		})

		if err != nil {
			t.Errorf("SendMessage() got %v, want = %v", err.Error(), nil)
		}
	})

	t.Run("Test Telegram Bad Request", func(t *testing.T) {
		server := createBadServer(t)
		defer server.Close()

		telegram := New(server.Client(), server.URL)
		err := telegram.SendMessage(structs.DataTele{
			ChatId:       123,
			Username:     "totanvix",
			ChatType:     "private",
			ReplyMessage: "test",
		})

		if err == nil {
			t.Errorf("SendMessage() not return error, want error")
		}
	})
}

func createHappyServer(t *testing.T) *httptest.Server {
	t.Helper()
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(200)
		res := make(map[string]any)
		res["ok"] = true
		res["result"] = map[string]any{
			"message_id": 123,
			"chat": map[string]int{
				"id": 1,
			},
			"text": "test",
		}
		response, _ := json.Marshal(res)
		w.Write(response)
	}))

	return server
}

func createBadServer(t *testing.T) *httptest.Server {
	t.Helper()
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(400)
		res := make(map[string]any)
		res["ok"] = false
		res["error_code"] = 400
		res["description"] = ""
		response, _ := json.Marshal(res)
		w.Write(response)
	}))

	return server
}
