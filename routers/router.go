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
	//支付宝
	beego.Router("/alipay",&controllers.GoodsController{},"get:HandleAlipay")
	//支付成功
	//beego.Router("/payok",&controllers.GoodsController{},"get:PayOk")
	//发短信
	beego.Router("/sendMsg",&controllers.GoodsController{},"get:SendMsg")
	//商品详情页
	beego.Router("/goodsDetail",&controllers.GoodsController{},"get:ShowDetail")
	//列表页展示
	beego.Router("/goodsList",&controllers.GoodsController{},"get:ShowGoodsList")
	//搜索功能
	beego.Router("/searchGoods",&controllers.GoodsController{},"post:HandleSearch")

}

func FilterFunc(ctx *context.Context)  {
	if ctx.Input.Session("userName")==nil{
		ctx.Redirect(302,"/login")
		return
	}
}