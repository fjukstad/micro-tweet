package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"html/template"
	"net/http"
	"strconv"
	"time"

	"github.com/dghubble/go-twitter/twitter"
	"github.com/dghubble/oauth1"
)

var client *twitter.Client
var Messages []Message

type Message struct {
	From    string
	Message string
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
	b, err := json.Marshal(messages)
	if err != nil {
		fmt.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Write(b)
}

func GenerateMessages() {
	i := 0
	for {
		from := "From " + strconv.Itoa(i)
		message := "Message " + strconv.Itoa(i)
		messages = append(messages, Message{from, message})
		time.Sleep(10 * time.Second)
	}
}

func IndexHandler(w http.ResponseWriter, r *http.Request) {
	indexTemplate := template.Must(template.ParseFiles("public/index.html"))
	err := indexTemplate.Execute(w, nil)
	if err != nil {
		fmt.Println(err)
	}
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
	mux.HandleFunc("/", IndexHandler)
	mux.Handle("/public/", http.FileServer(http.Dir(".")))

	port := "8080"

	fmt.Println("Server started on port", port)
	err := http.ListenAndServe(":"+port, mux)

	if err != nil {
		fmt.Println(err)
		return
	}

}
