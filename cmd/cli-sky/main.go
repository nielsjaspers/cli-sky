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

    bluesky.Post("this is a test from the terminal! i am making a cli tool for bluesky in go!\nhttps://github.com/nielsjaspers/cli-sky\n\n#computerscience #dev #programming #golang #bluesky #opensource", authResponse)
}
