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
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/gofrs/uuid"
)

// exit channel
var c = make(chan os.Signal, 1)

type Message struct {
	Name   string `json:"name"`
	Status int    `json:"status"`
}

// struct object of Mqtt
type Mqtt struct {
	client   mqtt.Client
	username string
	password string
	broker   string
	port     int
}

// init function of Mqtt struct
func (m *Mqtt) init(username, password, broker string, port int) {
	mqtt.DEBUG = log.New(os.Stdout, "", 0)
	mqtt.ERROR = log.New(os.Stdout, "", 0)

	// generate uuid for client id
	clientId := uuid.Must(uuid.NewV4()).String()
	opts := mqtt.NewClientOptions().AddBroker(fmt.Sprintf("tcp://%s:%d", broker, port)).SetClientID(clientId)
	opts.SetUsername(username)
	opts.SetPassword(password)
	opts.SetKeepAlive(0 * time.Second)
	opts.SetPingTimeout(0 * time.Second)

	// create a new instance of mqtt client
	m.client = mqtt.NewClient(opts)
	if token := m.client.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}

	m.username = username
	m.password = password
	m.broker = broker
	m.port = port
}

// NewMqtt function
func NewMqtt() *Mqtt {
	// as of streamr allow any username, we'll set username as server name
	username := GetEnv("MQTT_USERNAME", "server")
	// get password from environment variable
	password := GetEnv("MQTT_PASSWORD", "password")
	// get broker from environment variable
	broker := GetEnv("MQTT_BROKER", "localhost")
	// get port from environment variable
	port := GetEnvAsInt("MQTT_PORT", 1883)

	// create a new instance of mqtt
	mqtt := &Mqtt{}
	// init the mqtt
	mqtt.init(username, password, broker, port)

	// return the mqtt
	return mqtt
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
	q := NewMqtt()
	topic := GetEnv("MQTT_TOPIC", "mqtt")
	// fmt.Println(mqtt)

	// wait for exit signal
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		for sig := range c {
			cleanup(fmt.Sprintf("captured %v, exiting..", sig))
			os.Exit(0)
		}
	}()

	for {
		time.Sleep(5 * time.Second)
		fmt.Println("1")
		if q.client.IsConnected() {
			msg := Message{q.username, 1}
			payload, _ := json.Marshal(msg)
			token := q.client.Publish(topic, 0, false, payload)
			if token.Wait() {
				fmt.Println("published")
			} else {
				break
			}
		}
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

func NewHeartbeat(mqtt *Mqtt) {
	panic("unimplemented")
}
