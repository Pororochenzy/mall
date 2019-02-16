package routers

import (
	"dailyFresh/controllers"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/context"
)

func init() {
	beego.InsertFilter("/goods/*",beego.BeforeExec,FilterFunc)
	// beego.Router("/", &controllers.MainController{})
	//首页
	beego.Router("/",&controllers.GoodsController{},"get:ShowIndex")
	beego.Router("/register",&controllers.UserController{},"get:ShowRegister;post:HandleRegister")
	//激活用户
	beego.Router("/active",&controllers.UserController{},"get:HandleActive")
	//登录业务
	beego.Router("/login",&controllers.UserController{},"get:ShowLogin;post:HandleLogin")
	//退出登录
	beego.Router("/logout",&controllers.UserController{},"get:Logout")
	//用户中心信息
	beego.Router("/goods/userCenterInfo",&controllers.UserController{},"get:ShowUserCenterInfo")
	//用户中心订单
	beego.Router("/goods/userCenterOrder",&controllers.UserController{},"get:ShowUserCenterOrder")
	//用户中心地址
	beego.Router("/goods/userCenterSite",&controllers.UserController{},"get:ShowUserCenterSite;post:HandleSite")
}

func FilterFunc(ctx *context.Context)  {
	if ctx.Input.Session("userName")==nil{
		ctx.Redirect(302,"/login")
		return
	}
}