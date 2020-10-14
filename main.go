package main

import (
	"fmt"
)

var scrap dblp_scrap

func main(){
	scrap.init("conf", "cvpr", "2019")
	fmt.Println(scrap.base_link)
	go scrap.visit()
	scrap.save(".dblp.csv")
}