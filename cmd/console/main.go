package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"sync"
)

type Urls struct {
	AStream string `json:"aStream"`
	Audio   string `json:"audio"`
	Video   string `json:"video"`
}
type Chapter struct {
	Id           int    `json:"id"`
	Order        string `json:"order"`
	FileName     string `json:"filename"`
	TitleArabic  string `json:"title-arabic"`
	TitleEnglish string `json:"title-english"`
	Updated      int    `json:"updated"`
	Urls         Urls   `json:"urls"`
}

type Response struct {
	Chapters []Chapter `json:"chapters"`
}

func main() {
	url := "http://psmlive.psm.mv/quran/api/"
	response, err := http.Get(url)

	if err != nil {
		log.Fatal("could not download the data")
	}

	jsonDecoder := json.NewDecoder(response.Body)
	r := Response{}
	jsonDecoder.Decode(&r)

	var wg sync.WaitGroup

	for _, chapter := range r.Chapters {
		wg.Add(1)
		go func(chapter *Chapter) {
			defer wg.Done()
			save(chapter)
		}(&chapter)
	}

	wg.Wait()
}

func save(chapter *Chapter) {
	file, err := os.Create(fmt.Sprintf("tmp/%s_%s", chapter.Order, chapter.FileName))

	if err != nil {
		log.Fatal("could not create file for writing")
		return
	}

	audioResponse, err := http.Get(chapter.Urls.Audio)

	if err != nil {
		log.Fatal("could not download audio file", chapter.FileName)
		return
	}

	audioContent, err := io.ReadAll(audioResponse.Body)

	if err != nil {
		log.Fatal("could not get the contents of the audio file", chapter.FileName)
		return
	}

	_, err = file.Write(audioContent)

	if err != nil {
		log.Fatal("could not write file", chapter.FileName)
		return
	}

	fmt.Println("file downloaded and saved", chapter.FileName)
}
