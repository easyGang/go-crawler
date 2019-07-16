package main

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"regexp"
)

// Get document body as a string
func getDocumentBody(url string) string {
	resp, err := http.Get(url)
	if err != nil {
		fmt.Printf("%v - error", err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("%v - error", err)
	}
	return string(body)
}

// Get md5 Hash from a string
func GetMD5Hash(text string) string {
	hasher := md5.New()
	hasher.Write([]byte(text))
	return hex.EncodeToString(hasher.Sum(nil))
}

// Get domain from url
func getDomain(sUrl string) string {
	u, err := url.Parse(sUrl)
	if err != nil {
		fmt.Println(err)
	}
	return u.Host
}

// Func for close a body of a doc
func bodyClose(body io.ReadCloser) {
	_ = body.Close()
}

// Returns a map of links
func crawl(sUrl string, linkMap map[string]bool, visitedLinksMap map[string]bool) (map[string]bool, map[string]bool) {
	if visitedLinksMap[sUrl] == false {
		visitedLinksMap[sUrl] = true
		resp, err := http.Get(sUrl)
		if err != nil {
			fmt.Printf("%v - error", err)
		}
		defer bodyClose(resp.Body)
		if resp.Header.Get("Content-Type") == "text/html; charset=UTF-8" {
			doc, err := goquery.NewDocumentFromReader(resp.Body)
			if err != nil {
				log.Fatal(err)
			}
			domain := getDomain(sUrl)
			doc.Find("a").Each(func(i int, s *goquery.Selection) {
				link, _ := s.Attr("href")
				isMatched, err := regexp.Match(`^(http:|https:)\/\/`+domain+`\/`, []byte(link))
				if err != nil {
					fmt.Println(err)
				}
				if isMatched {
					if linkMap[sUrl] == false {
						linkMap[sUrl] = true
						fmt.Println("- ", sUrl)
					}
					crawl(link, linkMap, visitedLinksMap)
				}
			})
		}
	}
	return linkMap, visitedLinksMap
}

func main() {
	var linkMap = make(map[string]bool)
	var visitedLinksMap = make(map[string]bool)
	sUrl := "https://www.your-site.ru/"
	linkSlice, visitedLinksMap := crawl(sUrl, linkMap, visitedLinksMap)
	b, err := json.MarshalIndent(linkSlice, "", "  ")
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("\n ", string(b), "\n ")
	fmt.Println("\n", len(linkSlice), " - коль-во страниц \n", len(visitedLinksMap), " - коль-во посещенных ссылок")
}
