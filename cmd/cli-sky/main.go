package main

import (
	"fmt"

	"github.com/nielsjaspers/cli-sky/bluesky"
)

func main(){
    _, responseBody := bluesky.CreateSession("")
    fmt.Println(responseBody)
}
