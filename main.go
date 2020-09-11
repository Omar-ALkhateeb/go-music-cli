package main

import (
	"bufio"
	"fmt"
	h "go-music-cli/helperfunctions"
	"os"
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Enter a plylist id: ")
	id, _ := reader.ReadString('\n')
	//url := "https://www.youtube.com/playlist?list=PL4o29bINVT4EG_y-k5jGoOu3-Am8Nvi10"
	//fmt.Println(h.ExampleScrape(url))
	h.GetPlaylistItems(id)
}
