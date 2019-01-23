package main

import (
	"bufio"
	"flag"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"path"
	"strconv"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"golang.org/x/text/encoding/japanese"
	"golang.org/x/text/transform"
)

var word string
var number = 20
var directory = `C:\temp`

// 実行例:go run crawler.go -w "ピカチュウ" -n 30 -o "C:\temp"
func main() {
	flag.StringVar(&word, "w", "", "検索キーワード")
	flag.IntVar(&number, "n", 30, "ダウンロード数")
	flag.StringVar(&directory, "o", `C:\temp`, "画像の保存先")
	flag.Parse()

	query := strings.Join(strings.Split(word, " "), "+")
	baseURL, _ := url.Parse(strings.Join([]string{"https://www.google.co.jp/search?q=", query, "&source=lnms&tbm=isch"}, ""))
	searchurl := baseURL.String()

	counter := 1
	page := number / 20
	for index := 0; index <= page; index++ {
		resp, err := http.Get(searchurl)
		if err != nil {
			return
		}

		utfBody := transform.NewReader(bufio.NewReader(resp.Body), japanese.ShiftJIS.NewDecoder())
		doc, _ := goquery.NewDocumentFromReader(utfBody)

		doc.Find("img").EachWithBreak(func(_ int, s *goquery.Selection) bool {
			imgurl, exists := s.Attr("src")
			if exists {
				response, e := http.Get(imgurl)
				if e != nil {
					log.Fatal(e)
				}

				file, err := os.Create(path.Join(directory, strconv.Itoa(counter)+`.jpg`))
				if err != nil {
					log.Fatal(err)
				}

				_, err = io.Copy(file, response.Body)
				if err != nil {
					log.Fatal(err)
				}
			}
			if counter >= number {
				return false
			}
			counter++
			time.Sleep(100 * time.Millisecond)
			return true
		})

		next := doc.Find("a.fl").Last()
		nexturl, exists := next.Attr("href")
		if exists {
			println(nexturl)
		}
		searchurl = "https://www.google.co.jp/" + nexturl
	}
}
