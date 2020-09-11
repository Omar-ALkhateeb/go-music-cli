package helperfunctions

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
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
	reader := csv.NewReader(resp.Body)
	reader.Comma = ','
	data, err := reader.ReadAll()
	if err != nil {
		return nil, err
	}

	return data, nil
}

//BlackListArtists return a list of top 200 artist of this year
func BlackListArtists() []string {
	url := "https://spotifycharts.com/regional/global/daily/latest/download"
	artists := []string{}
	data, err := readCSVFromURL(url)
	if err != nil {
		panic(err)
	}

	for idx, row := range data {
		// skip header
		if idx == 0 || idx == 1 {
			continue
		}
		//fmt.Println(row[2])
		artists = append(artists, row[2])
	}
	fmt.Println(artists)
	return artists
}

//GetPlaylistItems gets and parses res from yt api
func GetPlaylistItems(id string) {
	url := "https://www.googleapis.com/youtube/v3/playlistItems?part=snippet&maxResults=50&playlistId=" + id + "&key=AIzaSyCIE10uSul8S-MVftkpsyPurgdc8O-4MNY"
	res, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()
	var result Data
	//whiteListed := []string{}
	var blackListed = BlackListArtists()
	byteValue, _ := ioutil.ReadAll(res.Body)
	json.Unmarshal(byteValue, &result)
	// search for items to remove
	for _, item := range result.Items {
		flag := false
		for _, str := range blackListed {
			if strings.Contains(item.Snippet.Title, str) {
				flag = true
				break
			}
		}
		if !flag {
			//whiteListed = append(whiteListed, item.ID)
			fmt.Println(item.Snippet.ResourceID.VideoID)
			ExampleClient(item.Snippet.ResourceID.VideoID)
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
func ExampleClient(id string) {
	defer recoverDownload()
	client := youtube.Client{}

	video, err := client.GetVideo(id)
	if err != nil {
		panic(err)
	}
	//fmt.Println(video.Formats)
	resp, err := client.GetStream(video, &video.Formats[len(video.Formats)-1])
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	file, err := os.Create(id + ".mp3")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	_, err = io.Copy(file, resp.Body)
	if err != nil {
		panic(err)
	}
}
