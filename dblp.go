package main

import (
	"fmt"
	"time"
	"github.com/gocolly/colly"
)

type dblp_scrap struct{
	c *colly.Collector
	base_link string
	mode string
	target string
	link_channel chan string
	content_channel chan [2]string
}

func (scrap *dblp_scrap) init(mode string, target string, year string, paralle int){
	scrap.c = colly.NewCollector(
		colly.AllowedDomains("dblp.org"),
		colly.Async(true),
		colly.UserAgent("Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/48.0.2564.116 Safari/537.36"),
	)
	scrap.link_channel = make(chan string, 10240)
	scrap.content_channel = make(chan [2]string, 10240)
	scrap.mode = mode
	scrap.target = target
	scrap.base_link = "https://dblp.org/db/" + scrap.mode + "/" + scrap.target + "/" + scrap.target + year + ".html"
	var c string
	if mode == "journals"{
		c = ".entry.article"
	} else {
		c = ".entry.inproceedings"
	}
	scrap.c.OnHTML(c, func(e *colly.HTMLElement) {
		scrap.content_channel <- [2]string{e.ChildAttr("li.ee a", "href"),e.ChildText("span.title")}
		// fmt.Println(e.ChildAttr("li.ee a", "href"))  // doi and link to conf
		// fmt.Println(e.ChildText("span.title"))  // title
		scrap.link_channel <- e.ChildAttr("li.ee a", "href")
	})

	scrap.c.Limit(&colly.LimitRule{
		Parallelism: paralle,
		RandomDelay: 5 * time.Second,
	})

	scrap.c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL.String())
	})
}

func (scrap *dblp_scrap) visit(){
	scrap.c.Visit(scrap.base_link)
}

func (scrap *dblp_scrap) close(){
	close(scrap.content_channel)
	close(scrap.link_channel)
	scrap.c.Wait()
}

// func (scrap *dblp_scrap) save(fname string){
// 	file, err:=os.Create(fname)
// 	if err!=nil{
// 		fmt.Println("cannot creat file")
// 		return
// 	}
// 	defer file.Close()
// 	writer := csv.NewWriter(file)
// 	defer writer.Flush()
// 	writer.Write([]string{"link", "title"})
// 	for i:= range scrap.content_channel{
// 		writer.Write([]string{i[0], i[1]})
// 	}
// }

// var scrap dblp_scrap

// func main(){
// 	scrap.init("conf", "cvpr", "2019")
// 	fmt.Println(scrap.base_link)
// 	go scrap.visit()
// 	scrap.save(".dblp.csv")
// }