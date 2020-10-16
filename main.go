package main

import (
	"fmt"
	"runtime"
	"flag"
)

var d_scrap dblp_scrap
var k_scrap kdd_scrap

var mode = flag.String("mode","conf","conf or journals")
var target = flag.String("target","kdd","target")
var volume = flag.String("volume","2020","year or volume")

func main(){
	flag.Parse()
	runtime.GOMAXPROCS(10)
	d_scrap.init(mode, target, volume)
	k_scrap.init()
	fmt.Println(d_scrap.base_link)
	go d_scrap.visit()
	go d_scrap.save(".dblp_conf_kdd_2020.csv")
	go k_scrap.save(".conf_kdd_2020.csv")
	for i:= range d_scrap.link_channel{
		k_scrap.visit(i)
	}
	d_scrap.close()
	k_scrap.close()
}