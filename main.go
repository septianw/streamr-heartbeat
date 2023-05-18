/**
 *	@file main.go
 *	This is the main file for the application.
 *
 *  This heartbeat program is a program to check the status of a server.
 */

// This heartbeat program will connected to mqtt to transmit the status of the server.
// the status of the server will identified by 1 or 0. 1 means the server is running and 0 means the server is not running.

package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"
)

// exit channel
var c = make(chan os.Signal, 1)

type Metadata struct {
	Timestamp  int64  `json:"timestamp"`
	ServerName string `json:"server_name"`
}

type Message struct {
	Status   int      `json:"status"`
	Metadata Metadata `json:"metadata"`
}

// struct object of Mqtt
type Publisher struct {
	url           string
	Broker        string
	Port          int
	StreamId      string
	ClientName    string
	Authorization string
	mode          string
}

// init function of Mqtt struct
func (h *Publisher) init(broker string, port int, streamId, ClientName, Auth string) {
	h.url = fmt.Sprintf("http://%s:%d/streams/%s", broker, port, streamId)
	h.Broker = broker
	h.Port = port
	h.StreamId = streamId
	h.ClientName = ClientName
	h.Authorization = Auth
	h.mode = GetEnv("HTTP_MODE", "production")
}

// function to publish data to http
func (h *Publisher) Publish() {
	meta := Metadata{time.Now().Unix(), h.ClientName}
	msg := Message{1, meta}
	data, err := json.Marshal(msg)
	if err != nil {
		log.Printf("error marshall msg: %v\n", err)
		return
	}
	// send post request to http
	req, err := http.NewRequest("POST", h.url, bytes.NewBuffer(data))
	if err != nil {
		log.Printf("error: %v\n", err)
		return
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("bearer %s", h.Authorization))
	client := &http.Client{}
	resp, err := client.Do(req)
	// resp, err := http.Post(h.url, "application/json", bytes.NewBuffer(data))

	if h.mode == "debug" {
		fmt.Printf("init: publishing to %v\n", h.url)
		fmt.Printf("resp.Body: %d\n", resp.StatusCode)
	}
	if err != nil {
		log.Printf("error: %v\n", err)
		resp.Body.Close()
		return
	}
}

// NewPublisher function
func NewPublisher() *Publisher {
	// get broker from environment variable
	broker := GetEnv("HTTP_BROKER", "localhost")
	// get port from environment variable
	port := GetEnvAsInt("HTTP_PORT", 7171)
	// get stream id from environment variable
	// streamId := strings.ReplaceAll(url.QueryEscape(GetEnv("HTTP_STREAM_ID", "mqtt")), "%2F", "/")
	streamId := url.QueryEscape(GetEnv("HTTP_STREAM_ID", "mqtt"))
	// get client name from environment variable
	clientName := GetEnv("HTTP_CLIENT_NAME", "mqtt")
	// get authorization from environment variable
	auth := GetEnv("HTTP_AUTH", "dummy")

	// create a new instance of HTTP
	p := &Publisher{}

	// init the http
	p.init(broker, port, streamId, clientName, auth)

	// return the mqtt
	return p
}

// function to get environment variable as integer
func GetEnvAsInt(s string, i int) int {
	// get enironment variable
	// if the environment variable is not set, return the default value
	es := os.Getenv(s)
	if es == "" {
		return i
	}

	// return the environment variable as integer
	result, err := strconv.Atoi(es)
	if err != nil {
		return i
	}

	return result
}

// function to get environment variable
func GetEnv(s1, s2 string) string {
	// get enironment variable
	// if the environment variable is not set, return the default value
	es := os.Getenv(s1)
	if es == "" {
		return s2
	}
	return es
}

// cleanup function
func cleanup(msg string) {
	log.Printf("log: %v\n", msg)
	fmt.Println("cleanup")
}

// main function
func main() {
	// create a new instance of mqtt
	q := NewPublisher()

	// wait for exit signal
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		for sig := range c {
			cleanup(fmt.Sprintf("captured %v, exiting..", sig))
			os.Exit(0)
		}
	}()

	for {
		time.Sleep(1 * time.Second)
		q.Publish()
	}

	// loop in a thread
	// for {
	// 	// check the status of the server
	// 	// if the server is running, send 1 to mqtt+.0

	// 	// if the server is not running, send 0 to mqtt
	// 	// send the status to mqtt
	// 	// sleep for 5 seconds
	// 	time.Sleep(5 * time.Second)
	// }

	// create a new instance of heartbeat
	// heartbeat := NewHeartbeat(mqtt)

	// create a new instance of server
	// server := NewServer(heartbeat)

	// create a new instance of signal
	// signal := NewSignal()

	// create a new instance of logger
	// logger := NewLogger()

	// create a new instance of config
	// config := NewConfig()

	// create a new instance of app
	// app := NewApp(server, signal, logger, config)

	// start the application
	// app.Start()

	// wait for the signal
	// signal.WaitForSignal()

	// stop the application
	// app.Stop()
}

// func NewApp(server, signal, logger, config invalid type) {
// 	panic("unimplemented")
// }
