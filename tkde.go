package main

import (
	"fmt"
	"time"
	"github.com/gocolly/colly"
	"encoding/csv"
	"encoding/json"
	"os"
	"sync"
)


type tkde_scrap struct{
	c *colly.Collector
	paper_channel chan paper
	wg *sync.WaitGroup
}

func (scrap * tkde_scrap) init(paralle int){
	scrap.c = colly.NewCollector(
		colly.AllowedDomains("dl.acm.org"),
		colly.Async(true),
		colly.UserAgent("Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/48.0.2564.116 Safari/537.36"),
	)
	scrap.paper_channel = make(chan paper, 10240)
	scrap.wg = &sync.WaitGroup{}
	
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
		scrap.paper_channel <- p
	})

	scrap.c.Limit(&colly.LimitRule{
		Parallelism: paralle,
		RandomDelay: 5 * time.Second,
	})

	scrap.c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL.String())
	})
}

func (scrap *tkde_scrap) visit(link string){
	defer scrap.wg.Done()
	scrap.wg.Add(1)
	scrap.c.Visit(link)
	scrap.c.Wait()
}

func (scrap *tkde_scrap) close(){
	close(scrap.paper_channel)
	scrap.wg.Wait()
}

func (scrap *tkde_scrap) save(fname string){
	defer scrap.wg.Done()
	scrap.wg.Add(1)
	file, err:=os.Create(fname)
	if err!=nil{
		fmt.Println("cannot creat file")
		return
	}
	defer file.Close()
	writer := csv.NewWriter(file)
	writer.Write([]string{"doi", "title","abstract","authors","references"})
	defer writer.Flush()
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
		writer.Write([]string{i.doi,i.title,i.abstract,string(a),string(r)})
	}
}

// var scrap tkde_scrap

// func main(){
// 	scrap.init()
// 	scrap.visit("https://dl.acm.org/doi/10.1145/3394486.3403044")
// 	scrap.save(".acm.csv")	
// }