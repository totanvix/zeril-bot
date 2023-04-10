package shorturl

import (
	"fmt"
	"net/http"
	"strings"
	"unicode"
	"zeril-bot/utils/redis"
)

func Handler(w http.ResponseWriter, r *http.Request) {
	fmt.Println(r.URL.Path)
	path := r.URL.Path

	f := func(c rune) bool {
		return unicode.IsPunct(c) == unicode.IsPunct('/')
	}

	els := strings.FieldsFunc(path, f)
	if els[0] == "url" {
		id := els[1]
		url := redis.Get(id).Val()
		if url != "" {
			http.Redirect(w, r, url, http.StatusSeeOther)
			return
		}
	}

	fmt.Fprintf(w, "<h1>Không có dữ liệu</h1>")
}
