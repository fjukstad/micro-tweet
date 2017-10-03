package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"html/template"
	"net/http"

	"github.com/dghubble/go-twitter/twitter"
	"github.com/dghubble/oauth1"
)

var client *twitter.Client

func TweetHandler(w http.ResponseWriter, r *http.Request) {
	// Home Timeline
	tweets, _, err := client.Timelines.HomeTimeline(&twitter.HomeTimelineParams{
		Count: 20,
	})

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	b, err := json.Marshal(tweets)
	if err != nil {
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

	// Send a Tweet
	//tweet, resp, err := client.Statuses.Update("Hello World!", nil)

	//if err != nil {
	//	fmt.Println("Could not tweet :(")
	//	return
	//}

	port := "8080"

	fmt.Println("Server started on port", port)
	err := http.ListenAndServe(":"+port, mux)

	if err != nil {
		fmt.Println(err)
		return
	}

}
