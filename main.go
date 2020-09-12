package main

import (
	"bufio"
	"fmt"
	h "go-music-cli/helperfunctions"
	"os"
	"strings"
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Enter a playlist url: ")
	url, _ := reader.ReadString('\n')
	url = strings.Replace(url, "\r\n", "", -1)
	getIDFromURL(url)
	//url := "https://www.youtube.com/playlist?list=PL4o29bINVT4EG_y-k5jGoOu3-Am8Nvi10"
	//url2 := "https://www.youtube.com/watch?v=ZLPiYZrwAzU&list=RD1vhFnTjia_I&index=2"
}

func getIDFromURL(url string) {
	//url = strings.TrimLeft(url, "https://www.youtube.com/playlist?list=")
	url = strings.Split(url, "&list=")[1]
	id := strings.Split(url, "&")[0]
	fmt.Println("id is " + id)
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("error parsing link!")
		}
	}()

	h.GetPlaylistItems(id)
}
