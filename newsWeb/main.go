package main

import (
	_ "newsWeb/models"
	_ "newsWeb/routers"
	"github.com/astaxie/beego"
)

func main() {
	beego.AddFuncMap("prePage",ShowPrePage)
	beego.AddFuncMap("afterPage",ShowAfterPage)
	beego.Run()
}

func ShowPrePage(pageIndex int) int {
	if pageIndex <= 1 {
		return 1
	}
	return pageIndex - 1
}

func ShowAfterPage(pageIndex int,pageAll int)int{
	if  pageIndex>=pageAll{
		return pageAll
	}
	return pageIndex+1
}