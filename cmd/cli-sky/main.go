package main

import (
	"fmt"

	"github.com/nielsjaspers/cli-sky/bluesky"
    "github.com/nielsjaspers/cli-sky/internal/datahandler"
)

func main(){
    //// COMMENT THIS IF YOU DONT WANT TO HIT THE API
    //// YOU'LL BE RATE LIMITED :(
    _, authResponse, responseBody := bluesky.CreateSession("")
    fmt.Println(responseBody)

    err := datahandler.WriteAuthResponseToFile(authResponse)
    if err != nil {
        panic(1)
    }

    // bluesky.Post("this is a test from the terminal! i am making a cli tool for bluesky in go!\nhttps://github.com/nielsjaspers/cli-sky\n\n#computerscience #buildinpublic #programming #golang #bluesky #opensource", authResponse)
    // bluesky.Post("testing from the #terminal", authResponse)
}
