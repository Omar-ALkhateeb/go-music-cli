package helperfunctions

import (
	"encoding/csv"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	youtube "github.com/kkdai/youtube/v2"
)

//Data struct rep the json res
type Data struct {
	Kind  string `json:"kind"`
	Etag  string `json:"etag"`
	Items []item `json:"items"`
}

type item struct {
	Kind    string  `json:"kind"`
	Etag    string  `json:"etag"`
	ID      string  `json:"id"`
	Snippet snippet `json:"snippet"`
}
type snippet struct {
	Title      string     `json:"title"`
	ResourceID resourceID `json:"resourceId"`
}

type resourceID struct {
	VideoID string `json:"videoId"`
}

func readCSVFromURL(url string) ([][]string, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	// fmt.Print(resp.Body)
	reader := csv.NewReader(resp.Body)
	reader.Comma = ','
	reader.LazyQuotes = true
	reader.FieldsPerRecord = -1
	data, err := reader.ReadAll()
	if err != nil {
		return nil, err
	}

	return data, nil
}

//BlackListArtists return a list of top 200 artist of this year
func BlackListArtists() []string {
	url := "https://gist.githubusercontent.com/mbejda/9912f7a366c62c1f296c/raw/dd94a25492b3062f4ca0dc2bb2cdf23fec0896ea/10000-MTV-Music-Artists-page-1.csv"
	artists := []string{}
	data, err := readCSVFromURL(url)
	if err != nil {
		panic(err)
	}

	for idx, row := range data {
		if idx > 300 {
			break
		}
		// skip header
		if idx == 0 || idx == 1 {
			continue
		}
		// fmt.Println(row[0])
		artists = append(artists, row[0])
	}
	// fmt.Println(artists)
	return artists
}

//GetPlaylistItems gets and parses res from yt api
func GetPlaylistItems(id string) {

	client := youtube.Client{}

	playlist, err := client.GetPlaylist(id)
	if err != nil {
		panic(err)
	}

	/* ----- Enumerating playlist videos ----- */
	header := fmt.Sprintf("Playlist %s by %s", playlist.Title, playlist.Author)
	println(header)
	println(strings.Repeat("=", len(header)) + "\n")

	// for k, v := range playlist.Videos {
	// 	fmt.Printf("(%d) %s - '%s'\n", k+1, v.Author, v.Title)
	// }

	fmt.Println("black listed mainstream artists")
	var blackListed = BlackListArtists()
	// search for items to remove
	for _, item := range playlist.Videos {
		flag := false
		for _, str := range blackListed {
			// print(item.ID)
			if strings.Contains(item.Title, str) {
				flag = true
				break
			}
		}
		if !flag {
			// whiteListed = append(whiteListed, item.ID)
			fmt.Println(item.ID, item.Title)
			ExampleClient(item.ID, item.Title)
		}
	}
	//fmt.Println(result.Items[0].Snippet.Title)
	//return whiteListed
}

func recoverDownload() {
	if r := recover(); r != nil {
		fmt.Println("skipped a video")
	}
}

//ExampleClient : Example code for how to use this package for download video.
func ExampleClient(id string, title string) {
	defer recoverDownload()
	client := youtube.Client{}

	video, err := client.GetVideo(id)
	if err != nil {
		panic(err)
	}
	//fmt.Println(video.Formats)
	resp, _, err := client.GetStream(video, &video.Formats[len(video.Formats)-1]) // audio/mp3
	if err != nil {
		panic(err)
	}
	defer resp.Close()

	file, err := os.Create(title + ".mp3")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	_, err = io.Copy(file, resp)
	if err != nil {
		panic(err)
	}
	fmt.Println("done :)")

}
