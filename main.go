package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strings"

	"github.com/dghubble/go-twitter/twitter"
	"github.com/dghubble/oauth1"
	"github.com/tarm/serial"
)

var client *twitter.Client
var Messages []Message

type Message struct {
	From    string
	Message string
	Invalid bool
	Raw     string
}

func TweetHandler(w http.ResponseWriter, r *http.Request) {
	tweets, resp, err := client.Timelines.UserTimeline(&twitter.UserTimelineParams{
		ScreenName: "kodingforvoksne",
	})

	if err != nil {
		fmt.Println(err)
		fmt.Println(resp)

		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	b, err := json.Marshal(tweets)
	if err != nil {
		fmt.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Write(b)
}

func MessageHandler(w http.ResponseWriter, r *http.Request) {
	b, err := json.Marshal(Messages)
	if err != nil {
		fmt.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Write(b)
}

func IndexHandler(w http.ResponseWriter, r *http.Request) {
	indexTemplate := template.Must(template.ParseFiles("public/index.html"))
	err := indexTemplate.Execute(w, nil)
	if err != nil {
		fmt.Println(err)
	}
}
func microMessages() error {
	port := "/dev/tty.usbmodem1412"
	c := &serial.Config{Name: port, Baud: 115200}
	s, err := serial.OpenPort(c)
	if err != nil {
		log.Fatal(err)
	}
	defer s.Close()

	fmt.Println("Waiting for messages on port", port)
	for {

		var msg string
		for {
			buf := make([]byte, 500)
			n, err := s.Read(buf)
			if err != nil {
				return err
			}

			buf = buf[:n]

			msg += string(buf)
			if strings.Contains(msg, "\n") {
				break
			}
		}

		b := []byte(msg)

		var m Message
		err = json.Unmarshal(b, &m)
		if err != nil {
			fmt.Println("Invalid message:", string(b))
			m = Message{Invalid: true, Raw: string(b)}
		}

		Messages = append([]Message{m}, Messages...)
		if len(Messages) > 100 {
			Messages = Messages[:99]
		}

		_, _, err = client.Statuses.Update(m.From+": "+m.Message, nil)
		if err != nil {
			fmt.Println(err)
		}
	}
	return nil
}
func main() {

	var consumerKey = flag.String("consumer-key", "", "Twitter Consumer Key (API Key).")
	var consumerSecret = flag.String("consumer-secret", "", "Twitter Consumer Secret (API Secret).")

	var accessToken = flag.String("access-token", "", "Twitter Access Token.")
	var accessTokenSecret = flag.String("access-token-secret", "", "Twitter Access Token Secret.")

	flag.Parse()

	if *consumerKey == "" || *consumerSecret == "" ||
		*accessToken == "" || *accessTokenSecret == "" {
		fmt.Println("Please enter Twitter Consumer Key and Consumer Secret. Visit apps.twitter.com to create one.")
		return
	}
	config := oauth1.NewConfig(*consumerKey, *consumerSecret)
	token := oauth1.NewToken(*accessToken, *accessTokenSecret)
	httpClient := config.Client(oauth1.NoContext, token)

	client = twitter.NewClient(httpClient)

	mux := http.NewServeMux()
	mux.HandleFunc("/tweets", TweetHandler)
	mux.HandleFunc("/messages", MessageHandler)
	mux.HandleFunc("/", IndexHandler)
	mux.Handle("/public/", http.FileServer(http.Dir(".")))

	go microMessages()

	port := "8080"

	fmt.Println("Server started on port", port)
	err := http.ListenAndServe(":"+port, mux)

	if err != nil {
		fmt.Println(err)
		return
	}

}
