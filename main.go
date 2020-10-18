package main

import (
	"fmt"
	// "net/url"
	"flag"
)

var d_scrap dblp_scrap
var k_scrap kdd_scrap

var mode = flag.String("mode","conf","conf or journals")
var target = flag.String("target","kdd","target")
var volume = flag.String("volume","2020","year or volume")


func visit(){
	for i:= range d_scrap.link_channel{
		k_scrap.visit(i)
	}
}

func main(){
	flag.Parse()
	d_scrap.init(*mode, *target, *volume, 2)
	k_scrap.init(2)
	fmt.Println(d_scrap.base_link)
	go d_scrap.visit()
	go k_scrap.save("."+*mode+"_"+*target+"_"+*volume+".csv")
	visit()
	d_scrap.close()
	k_scrap.close()
}