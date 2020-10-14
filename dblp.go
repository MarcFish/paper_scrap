package dblp

import (
	"fmt"
	"time"
	"github.com/gocolly/colly"
	"encoding/csv"
	"os"
)

type dblp_scrap struct{
	c *colly.Collector
	base_link string
	mode string
	target string
	link_channel chan string
	content_channel chan [2]string
}

func (scrap *dblp_scrap) init(mode string, target string, year string){
	scrap.c = colly.NewCollector(
		colly.AllowedDomains("dblp.org"),
		colly.Async(true),
	)
	scrap.link_channel = make(chan string, 10240)
	scrap.content_channel = make(chan [2]string, 10240)
	scrap.mode = mode
	scrap.target = target
	scrap.base_link = "https://dblp.org/db/" + scrap.mode + "/" + scrap.target + "/" + scrap.target + year + ".html"

	scrap.c.OnHTML(".entry.inproceedings", func(e *colly.HTMLElement) {
		scrap.content_channel <- [2]string{e.ChildAttr("li.ee a", "href"),e.ChildText("span.title")}
		// fmt.Println(e.ChildAttr("li.ee a", "href"))  // doi and link to conf
		fmt.Println(e.ChildText("span.title"))  // title
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
	close(scrap.content_channel)
	close(scrap.link_channel)
}

func (scrap *dblp_scrap) save(fname string){
	file, err:=os.Create(fname)
	if err!=nil{
		fmt.Println("cannot creat file")
		return
	}
	defer file.Close()
	writer := csv.NewWriter(file)
	writer.Write([]string{"link", "title"})
	defer writer.Flush()
	for i:= range scrap.content_channel{
		writer.Write([]string{i[0], i[1]})
	}
}

var scrap dblp_scrap

// func main(){
// 	scrap.init("conf", "cvpr", "2019")
// 	fmt.Println(scrap.base_link)
// 	go scrap.visit()
// 	scrap.save(".dblp.csv")
// }