package main

import (
	"encoding/json"
	"log"
	"net/http"

	"golang.org/x/net/websocket"
)

type MessageRepository struct {
	messages []string
}

func (r *MessageRepository) Put(m string) error {
	r.messages = append(r.messages, m)
	return nil
}

func (r *MessageRepository) Find() ([]string, error) {
	return r.messages, nil
}

var repository MessageRepository

func echoHandler(ws *websocket.Conn) {
	var m string
	var r Request
	if err := websocket.Message.Receive(ws, &m); err != nil {
		log.Fatalf("%+v\n", err)
	}

	if err := json.Unmarshal([]byte(m), &r); err != nil {
		log.Fatalf("%+v\n", err)
	}

	if err := repository.Put(r.Message); err != nil {
		log.Fatalf("%+v\n", err)
	}

	ms, err := repository.Find()
	if err != nil {
		log.Fatalf("%+v\n", err)
	}

	json, err := json.Marshal(Response{
		Messages: ms,
	})

	if err != nil {
		log.Fatalf("%+v\n", err)
	}

	if err = websocket.Message.Send(ws, string(json)); err != nil {
		log.Fatalf("%+v\n", err)
	}
}

type Request struct {
	Message string `json:"message"`
}

type Response struct {
	Messages []string `json:"messages"`
}

func main() {
	repository = MessageRepository{
		messages: []string{},
	}
	http.Handle("/", websocket.Handler(echoHandler))
	e := http.ListenAndServe(":8080", nil)
	if e != nil {
		log.Fatalf("%+v", e)
	}
}
