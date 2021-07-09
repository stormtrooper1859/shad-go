// +build !solution

package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
)

var port = flag.Int("port", 8080, "port of server")

var dao DAO

func main() {
	dao = NewDAO()

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

	su.Key = dao.GetShortener(su.Url)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	res, err := json.Marshal(su)
	_, _ = w.Write(res)
}

func handlerGo(w http.ResponseWriter, r *http.Request) {
	splittedPath := strings.Split(r.URL.Path, "/")
	queryKey := splittedPath[2]

	url, err := dao.GetFullURL(queryKey)
	if err != nil {
		w.WriteHeader(404)
		_, _ = io.WriteString(w, "key not found")
	}

	responseBody := fmt.Sprintf("<a href=%q>Found</a>.", url)

	w.Header().Set("Location", url)
	w.WriteHeader(302)
	_, _ = io.WriteString(w, responseBody)
}
