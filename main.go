package main

import (
	"flag"
	"fmt"

	"github.com/dghubble/go-twitter/twitter"
	"github.com/dghubble/oauth1"
)

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

	client := twitter.NewClient(httpClient)

	// Send a Tweet
	tweet, resp, err := client.Statuses.Update("Hello World!", nil)

	if err != nil {
		fmt.Println("Could not tweet :(")
		return
	}

	fmt.Println(tweet, resp)
}
