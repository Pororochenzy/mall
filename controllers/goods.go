package controllers

import "github.com/astaxie/beego"

type GoodsController struct {
	beego.Controller
}

func(this*GoodsController)ShowIndex(){
	userName := this.GetSession("userName")
	if userName == nil{
		this.Data["userName"] = ""
	}else{
		this.Data["userName"] = userName.(string)  //????
	}
	beego.Info(userName)
	//指定视图
	this.TplName = "index.html"
}
