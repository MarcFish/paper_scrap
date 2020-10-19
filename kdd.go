package main

import (
	"fmt"
	"time"
	"github.com/gocolly/colly"
	"github.com/gocolly/colly/extensions"
	"encoding/csv"
	"encoding/json"
	"os"
	"net/url"
	// "sync"
)

type author struct{
	Name string `json:"name"`
	Id string `json:"id"`
}

type paper struct{
	title string
	doi string
	abstract string
	authors []author
	references []string
}

type kdd_scrap struct{
	c *colly.Collector
	paper_channel chan paper
	// wg *sync.WaitGroup
}

func (scrap * kdd_scrap) init(paralle int){
	scrap.c = colly.NewCollector(
		colly.AllowedDomains("dl.acm.org"),
		colly.Async(false),
		// colly.UserAgent("Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/48.0.2564.116 Safari/537.36"),
	)
	extensions.RandomUserAgent(scrap.c)
    extensions.Referer(scrap.c)
	scrap.paper_channel = make(chan paper, 10240)
	// scrap.wg = &sync.WaitGroup{}
	scrap.c.OnHTML("article", func(e *colly.HTMLElement) {
		var p paper
		p.title = e.ChildText(".citation__title")
		p.doi = e.ChildAttr("a.issue-item__doi", "href")
		p.abstract = e.ChildText("div.abstractSection.abstractInFull p")
		e.ForEach(".loa__item", func(_ int, el *colly.HTMLElement){
			var a author
			a.Name = el.ChildAttr(".author-name", "title")
			a.Id = el.ChildAttr(".btn.blue.stretched","href")
			p.authors = append(p.authors, a)
		})
		// reference
		e.ForEach(".references__item", func(_ int, el *colly.HTMLElement){
			ref := el.ChildText("span.references__note")
			p.references = append(p.references, ref)
		})
		fmt.Println("parse:"+p.title)
		scrap.paper_channel <- p
	})

	scrap.c.Limit(&colly.LimitRule{
		// Parallelism: paralle,
		RandomDelay: 5 * time.Second,
	})

	scrap.c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visitng:"+ r.URL.String())
	})

	scrap.c.OnResponse(func(r *colly.Response){
		fmt.Println("Get:"+r.Request.URL.String())
	})
}

func (scrap *kdd_scrap) visit(link string){
	u, err := url.Parse(link)
	if err!= nil{
		fmt.Println("url parse error")
	}
	if u.Host == "doi.org"{
		link = "https://dl.acm.org/doi" + u.Path
	}
	scrap.c.Visit(link)
}

func (scrap *kdd_scrap) close(){
	close(scrap.paper_channel)
	// scrap.wg.Wait()
	scrap.c.Wait()
}

func (scrap *kdd_scrap) save(fname string){
	// scrap.wg.Add(1)
	// defer scrap.wg.Done()
	defer fmt.Println("over")
	file, err:=os.Create(fname)
	if err!=nil{
		fmt.Println("cannot creat file")
		return
	} else {
		fmt.Println("creat:"+fname)
	}
	defer file.Close()
	writer := csv.NewWriter(file)
	defer writer.Flush()
	writer.Write([]string{"doi", "title","abstract","authors","references"})
	for i:= range scrap.paper_channel{
		a, err := json.Marshal(i.authors)
		if err != nil {
			fmt.Println("cannot marchal authors")
			return
		}
		r, err := json.Marshal(i.references)
		if err != nil{
			fmt.Println("cannot marchal references")
			return
		}
		fmt.Println("write:"+i.title)
		writer.Write([]string{i.doi,i.title,i.abstract,string(a),string(r)})
	}
}

// var scrap kdd_scrap

// func main(){
// 	scrap.init()
// 	scrap.visit("https://dl.acm.org/doi/10.1145/3394486.3403044")
// 	scrap.save(".acm.csv")	
// }