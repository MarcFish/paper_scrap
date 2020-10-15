package main

import (
	"fmt"
	// "time"
	"github.com/gocolly/colly"
	// "encoding/csv"
	// "os"
)

type author struct{
	name string
	id string
}

type paper struct{
	title string
	doi string
	abstract string
	authors []author
}

type kdd_scrap struct{
	c *colly.Collector
	paper_channel chan paper
}

func (scrap * kdd_scrap) init(){
	scrap.c = colly.NewCollector(
		colly.AllowedDomains("dl.acm.org"),
		colly.Async(true),
	)
	scrap.paper_channel = make(chan paper, 10240)
	
}

func main(){
	

}