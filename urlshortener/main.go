// +build !solution

package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"math/rand"
	"net/http"
	"strconv"
	"strings"
	"sync"
)

var port = flag.Int("port", 8080, "port of server")

var (
	keysToUrls map[string]string
	urlsToKeys map[string]string
	mtx        sync.Mutex
)

func main() {
	keysToUrls = make(map[string]string)
	urlsToKeys = make(map[string]string)

	flag.Parse()
	http.HandleFunc("/shorten", handlerShorten)
	http.HandleFunc("/go/", handlerGo)
	serverAddress := ":" + strconv.Itoa(*port)
	err := http.ListenAndServe(serverAddress, nil)
	fmt.Println(err)
}

type shortenUrl struct {
	Url string `json:"url"`
	Key string `json:"key"`
}

var (
	randomMax = 10 + 26*2
	keyLen    = 10
)

func genRandomKey() string {
	sb := strings.Builder{}

	for i := 0; i < keyLen; i++ {
		n := rand.Intn(randomMax)
		if n < 10 {
			sb.WriteByte('0' + byte(n))
		} else if n < 10+26 {
			sb.WriteByte('a' + byte(n-10))
		} else {
			sb.WriteByte('A' + byte(n-10-26))
		}
	}

	return sb.String()
}

func handlerShorten(w http.ResponseWriter, r *http.Request) {
	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(400)
		_, _ = io.WriteString(w, "invalid request")
	}

	var su shortenUrl
	err = json.Unmarshal(data, &su)
	if err != nil {
		w.WriteHeader(400)
		_, _ = io.WriteString(w, "invalid request")
	}

	mtx.Lock()

	key, urlsIsExist := urlsToKeys[su.Url]

	if !urlsIsExist {
		key = genRandomKey()
		for _, c := keysToUrls[key]; c; _, c = keysToUrls[key] {
			key = genRandomKey()
		}

		keysToUrls[key] = su.Url
		urlsToKeys[su.Url] = key
	}

	su.Key = key

	mtx.Unlock()

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	res, err := json.Marshal(su)
	_, _ = w.Write(res)
}

func handlerGo(w http.ResponseWriter, r *http.Request) {
	splittedPath := strings.Split(r.URL.Path, "/")
	queryKey := splittedPath[2]

	fmt.Println(r.URL.Path)
	fmt.Println(queryKey)

	mtx.Lock()

	url, contain := keysToUrls[queryKey]

	mtx.Unlock()

	if !contain {
		w.WriteHeader(404)
		_, _ = io.WriteString(w, "key not found")
	}

	responseBody := fmt.Sprintf("<a href=%q>Found</a>.", url)

	w.Header().Set("Location", url)
	w.WriteHeader(302)
	_, _ = io.WriteString(w, responseBody)
}
