package main

import (
	"fmt"
	// "net/url"
	"flag"
	"sync"
)

var d_scrap dblp_scrap
var k_scrap kdd_scrap

var mode = flag.String("mode","conf","conf or journals")
var target = flag.String("target","kdd","target")
var volume = flag.String("volume","2020","year or volume")

func main(){
	flag.Parse()
	wg := sync.WaitGroup{}
	d_scrap.init(*mode, *target, *volume, 5)
	k_scrap.init(5)
	fmt.Println(d_scrap.base_link)
	go d_scrap.visit()
	go func(){
		wg.Add(1)
		defer wg.Done()
		k_scrap.save("."+*mode+"_"+*target+"_"+*volume+".csv")
	}()
	for i:= range d_scrap.link_channel{
		k_scrap.visit(i)
	}
	d_scrap.close()
	k_scrap.close()
	wg.Wait()
}