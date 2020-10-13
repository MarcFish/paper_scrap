package main

import (
	"fmt"
	"time"
	"github.com/gocolly/colly"
	// "encoding/csv"
	// "os"
)

type dblp_scrap struct{
	c *colly.Collector
	base_link string
	mode string
	target string
	link_channel chan string
}

func (scrap *dblp_scrap) init(mode string, target string, year string){
	scrap.c = colly.NewCollector(
		colly.AllowedDomains("dblp.org"),
		colly.Async(true),
	)
	scrap.link_channel = make(chan string, 1024)
	scrap.mode = mode
	scrap.target = target
	scrap.base_link = "https://dblp.org/db/" + scrap.mode + "/" + scrap.target + "/" + scrap.target + year + ".html"

	scrap.c.OnHTML(".entry.inproceedings", func(e *colly.HTMLElement) {
		// fmt.Println(e.ChildAttr("li.ee a", "href"))  // doi and link to conf
		// fmt.Println(e.ChildText("span.title"))  // title
		scrap.link_channel <- e.ChildAttr("li.ee a", "href")
	})

	scrap.c.Limit(&colly.LimitRule{
		Parallelism: 2,
		RandomDelay: 5 * time.Second,
	})

	scrap.c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL.String())
	})
}

func (scrap *dblp_scrap) visit(){
	scrap.c.Visit(scrap.base_link)
	scrap.c.Wait()
}

var scrap dblp_scrap

func main(){
	scrap.init("conf", "kdd", "2020")
	// fmt.Println(scrap.base_link)
	go scrap.visit()
	for i:= range scrap.link_channel{
		fmt.Println(i)
	}
	time.Sleep(1e9)
}