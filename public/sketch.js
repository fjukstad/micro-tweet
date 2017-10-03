var url  = "/tweets" 
var tweets = []; 

function draw(){
    background(255);
    for(var i in tweets){
        if (tweets.hasOwnProperty(i)) {
            tweet = tweets[i]; 
            text(tweet.text, 10, 10+10*i);
        }
    }
}

function setup(){
    createCanvas(600,600);
    setTimeout(getTweets, 10);
}

function getTweets(){
    tweets = loadJSON(url);
    console.log(tweets);
    //setTimeout(getTweets, 10000);
}
