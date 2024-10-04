package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
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

	printChapterList(r.Chapters)

	var prompt string

	for strings.ToLower(prompt) != "n" {
		var chapterNumber string

		fmt.Print("Enter the Surah number to download: ")
		_, err = fmt.Scanln(&chapterNumber)

		if err != nil {
			log.Println("Could not get the Surah number", err)
			return
		}

		for _, chapter := range r.Chapters {
			if chapter.Order == chapterNumber {
				save(&chapter)
			}
		}

		fmt.Print("Want to continue downloading? (y/n): ")
		_, err = fmt.Scanln(&prompt)

		if err != nil {
			log.Println("could not understand your answer", err)
			return
		}
	}
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

	log.Println("writing the audio file: ", chapter.FileName)

	_, err = file.Write(audioContent)

	if err != nil {
		log.Fatal("could not write file", chapter.FileName)
		return
	}

	log.Println("file downloaded and saved", chapter.FileName)
}

func printChapterList(chapters []Chapter) {
	fmt.Println("List of Surahs:")
	fmt.Println("---------------")
	for _, chapter := range chapters {
		fmt.Printf("\t%s. %s\n", chapter.Order, chapter.TitleEnglish)
	}
}
