package main

import (
	"fmt"

	"github.com/nielsjaspers/cli-sky/bluesky"
)

func main(){
    //// COMMENT THIS IF YOU DONT WANT TO HIT THE API
    //// YOU'LL BE RATE LIMITED :(
    _, authResponse, responseBody := bluesky.CreateSession("")
    fmt.Println(responseBody)

    bluesky.Post("hashtag test #helloworld", authResponse)
}
