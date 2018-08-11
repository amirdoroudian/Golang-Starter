package main

import (
	"fmt"
	"net/http"
	"html/template"
	"encoding/xml"
	"io/ioutil"
)

type NewsMap struct {
	Link    string
	Creator string
}

type NewsAggPage struct {
	Title string
	News map[string]NewsMap
}

type News struct {
	Titles   []string `xml:"channel>item>title"`
	Links    []string `xml:"channel>item>link"`
	Creators []string `xml:"channel>item>creator"`
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "<h1>چقدر خفنیم ما</h1>")
}

func newsAggHandler(w http.ResponseWriter, r *http.Request) {	
	var n News

	resp, _ := http.Get("https://www.theguardian.com/world/rss")
	bytes, _ := ioutil.ReadAll(resp.Body)
	resp.Body.Close()

	xml.Unmarshal(bytes, &n)

	news_map := make(map[string]NewsMap)

	for idx, _ := range n.Titles {
		news_map[n.Titles[idx]] = NewsMap{n.Links[idx], n.Creators[idx]}
	}

	p := NewsAggPage{Title: "اخبار روز", News: news_map}
	t, _ := template.ParseFiles("newsaggtemplate.html")
	t.Execute(w, p)
}

func main() {
	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/agg/", newsAggHandler)
	http.ListenAndServe(":8000", nil)
}