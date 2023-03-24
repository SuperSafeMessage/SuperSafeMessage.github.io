package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"sync"
)

var messages sync.Map = sync.Map{}

func addMessage(receiver string, message string) {
	empty := &sync.Map{}

	value, _ := messages.LoadOrStore(receiver, empty)
	messageSet := value.(*sync.Map)
	for i := 0; ; i++ {
		_, loaded := messageSet.LoadOrStore(i, message)
		if !loaded {
			break
		}
	}
}

func getMessage(receiver string) []string {
	value, ok := messages.Load(receiver)
	if !ok {
		return nil
	}

	messageSet := value.(*sync.Map)
	messageList := make([]string, 0)
	for i := 0; ; i++ {
		messageValue, ok := messageSet.Load(i)
		if ok {
			messageList = append(messageList, messageValue.(string))
		} else {
			break
		}
	}
	return messageList
}

func send(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	receiver := req.URL.Query()["receiver"][0]
	var message string
	switch req.Method {
	case "GET":
		message = req.URL.Query()["message"][0]
	case "POST":
		body, err := ioutil.ReadAll(req.Body)
		if err != nil {
			return
		}
		defer req.Body.Close()
		message = string(body)
	}

	if len(message) > 256*1024 {
		return
	}

	addMessage(receiver, message)
	fmt.Fprintf(w, "ok")
}

func receive(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	receiver := req.URL.Query()["receiver"][0]
	messages := getMessage(receiver)
	b, _ := json.Marshal(messages)
	w.Write(b)
}

func main() {
	http.HandleFunc("/receive", receive)
	http.HandleFunc("/send", send)
	http.ListenAndServe(":8090", nil)
}
