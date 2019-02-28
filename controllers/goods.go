package controllers

import (
	"github.com/astaxie/beego"
	"fmt"
	"github.com/smartwalle/alipay"
	"github.com/KenmyZhang/aliyun-communicate"
	"encoding/json"
	"dailyFresh/models"
	"github.com/astaxie/beego/orm"
	"strconv"
	"github.com/gomodule/redigo/redis"

	"math"
)

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
	//查询商品类型   -->因为GoodsType 里还有个IndexTypeGoodsBanner成员
	////////////
	var goodsTypes []models.GoodsType
	 o:=  orm.NewOrm()
	 o.QueryTable("goods_type").All(&goodsTypes)
	 // goods[1] = map[type] goodsType 用map 里值存map
	this.Data["types"]=goodsTypes

	goods := make([]map[string]interface{},len(goodsTypes))
	for index,goodsType := range goodsTypes{
		temp := make(map[string]interface{})
		temp["type"] = goodsType
		goods[index] = temp
	}
	//接着把商品存储到我们的容器中，注意这里是分为图片商品和文字商品，需要分开来存储
	var goodsImage []models.IndexTypeGoodsBanner
	var goodsText []models.IndexTypeGoodsBanner
	for _,temp := range goods {
		o.QueryTable("IndexTypeGoodsBanner").RelatedSel("GoodsSKU", "GoodsType").Filter("GoodsType", temp["type"]).Filter("DisplayType", 1).OrderBy("Index").All(&goodsImage)
		o.QueryTable("IndexTypeGoodsBanner").RelatedSel("GoodsSKU", "GoodsType").Filter("GoodsType", temp["type"]).Filter("DisplayType", 0).OrderBy("Index").All(&goodsText)

		temp["goodsText"] = goodsText
		temp["goodsImage"] = goodsImage
	}
	this.Data["goods"] = goods


	 ////////////

	 //查询轮播商品图片
	 var goodsBanner []models.IndexGoodsBanner
	 o.QueryTable("IndexGoodsBanner").OrderBy("Index").All(&goodsBanner)
	 this.Data["goodsBanners"] = goodsBanner

	 //查询促销商品数据
	var promotionBanner []models.IndexPromotionBanner
	o.QueryTable("IndexPromotionBanner").OrderBy("Index").All(&promotionBanner)
	this.Data["proBanner"] = promotionBanner

	//指定视图
	this.TplName = "index.html"
}
func(this*GoodsController)HandleAlipay(){
	var privateKey =`MIIEogIBAAKCAQEAwRZh7xJOvERbDeQtZoV+W4fKuc7V1CbN2U2Od15DSrRECSUY
LAEvkNAeUBtjl3ZLXT3TwUC3IzIEISU8EOk+V6X5t+vXVyPy2WoZQWbkq8hRo16s
4qzuAYGT3dCuuhvP8PeU9Q3IpDIQIa+XjB2ur1EqOEtuZhvq3VZLiB8gEp9i2LL0
6qWpJf1mfrNcjPT+9xp39AORf2UM0rKGY/YwhZsH0Aisx15ZnMsYaW2d/3ehL++H
LxG2IqnzPid9jRjJmEiSXLt8FRPd6T122Yi9wUNH+WBhKM8ykWKK0GS5E0N2DX7q
NXBVAm+EC8xn4ZfNROBUnC5602hU4mcsghkrdwIDAQABAoIBADOmhB4KnKs58c1+
ezKQhSOA6JbZoFN26du2JmUB/ygtnoF/vb6PtqSbN3CgUvzCNRjFIC7y28p7Z6Vi
K3YunnGdwXYUjW8O+7hy7DyVhJf7JvN3sPGV5rjaa66LUyIPrIk+AUeoH0Lk7XHF
YdmmWwMkyBF1BBwmXaZFnkDUaqTwLONPwiB2pWxSsj57g6APfQyIXGWfYaiECNk6
abv8ekYLRamPBis1QwfnhEpBYsu2TcLoFbBMUcCmvyXtrWqVZDBu6+8O7Y5bVXhd
PcHZN6sIDG8aF/QL0ndP7aVTSF6jPiu7eHctlXU8Ypm8RWQQfGRFljrIsVqh3R1t
igjrOykCgYEA5a6GojqYVdq88UUQH+KkKtBKrN5ik7ju4hP8KlA1SdG/3au10qYT
2t0Epvck7ahovGjaHk8vV6eu1nb2urizW9jxkaZuKlvlS76IC0HYReAZpOufPpZx
wggdBPeR0K2Rue5gV/27Z1RyZ1tcS39MV3HVQPYpw3zVkck1jcEn2AMCgYEA1zZm
mJKwd+sA7YjCNS4pcZGFs0DMNhH2SSiXvQTIW/+bYmEeiNkWxNM3zRTrQo42DRt1
1lSw8CwICeqIFIcqUf2Af5DrztmLCa4FOOe94XBYpp1Vyi31yDk2ux0qYehSNO4a
dwSU+6g1VbnBnw8EcM35+Dvbn0SnEET6UByG5n0CgYBdfTCoEBm5uJN30ZjCocoY
8zeyLcMKRhhWRbQ6tPM73PiwDhiwaZFjYNtn6ulJI2eeaT9/XtPyZfqwqTO8xTmc
hS2vD4OAEm++6QsPKfoSVymZC5+CJlKfnBXT08GyozPR7smgh1MkuCbpEzL6OBKm
9VrMWmadf86ezLvTu2+waQKBgFXUW1xz/C8HKUqSZSnCbELfz5uqtwbFaRzKNKHs
u199dGHq67uWIK+EsVd4BU942fOmRPuisSJH3TjfMUBGm8kxOcMmb/UB66KDpY+J
VMAJB0IDs4g7hi9BKiO7wQBlIAv9/c4DpMsszYCP4blmytWWQCAQ90jwn1Qsvkue
5OodAoGAd3DyV/+3N0BcLrM/xy0RnYlRgxWPvQ9TJsG6cIBdXPXadrud3dscBfrD
EukYccE3/UnT1iEv0I7bXLGjO5mg7lDsjAzoOOxSN2Vuga/a/eSe7Hmy0CZ93DJK
UlBaoeOXawjXlJbCtHYHznh3indx7QWCXWOrMJada4Jdoz1Ed1Q=`

	var appId = "2016092800615273"

	var aliPublicKey =`MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEAwRZh7xJOvERbDeQtZoV+
W4fKuc7V1CbN2U2Od15DSrRECSUYLAEvkNAeUBtjl3ZLXT3TwUC3IzIEISU8EOk+
V6X5t+vXVyPy2WoZQWbkq8hRo16s4qzuAYGT3dCuuhvP8PeU9Q3IpDIQIa+XjB2u
r1EqOEtuZhvq3VZLiB8gEp9i2LL06qWpJf1mfrNcjPT+9xp39AORf2UM0rKGY/Yw
hZsH0Aisx15ZnMsYaW2d/3ehL++HLxG2IqnzPid9jRjJmEiSXLt8FRPd6T122Yi9
wUNH+WBhKM8ykWKK0GS5E0N2DX7qNXBVAm+EC8xn4ZfNROBUnC5602hU4mcsghkr
dwIDAQAB`
	var client = alipay.New(appId, aliPublicKey, privateKey, false)

	//alipay.trade.page.pay
	var p = alipay.AliPayTradePagePay{}
	p.NotifyURL = "http://192.168.42.142:8080/user/payOk"
	p.ReturnURL = "http://192.168.42.142:8080/user/payOk"
	p.Subject = "天天生鲜"
	p.OutTradeNo = "1234567811"
	p.TotalAmount = "1000.00"
	p.ProductCode = "FAST_INSTANT_TRADE_PAY"

	var url, err = client.TradePagePay(p)
	if err != nil {
		fmt.Println(err)
	}

	var payURL = url.String()

	this.Redirect(payURL,302)
}
//发短信
func(this*GoodsController)SendMsg(){
	var (
		gatewayUrl      = "http://dysmsapi.aliyuncs.com/"
		accessKeyId     = "LTAIh83X7bYYTIXw"
		accessKeySecret = "fYSLqA3BI8jNviNhURKT9T9TmHeOuP"
		phoneNumbers    = "19951750429"
		signName        = "天天生鲜"
		templateCode    = "SMS_149101793"
		//templateParam   = "{\"code\":\"1234\"}"
		templateParam   = "{\"code\":\"IMissYouZhuTou\"}"
	)
	smsClient := aliyunsmsclient.New(gatewayUrl)
	result, err := smsClient.Execute(accessKeyId, accessKeySecret, phoneNumbers, signName, templateCode, templateParam)
	fmt.Println("Got raw response from server:", string(result.RawResponse))
	if err != nil {
		panic("Failed to send Message: " + err.Error())
	}

	resultJson, err := json.Marshal(result)
	if err != nil {
		panic(err)
	}
	if result.IsSuccessful() {
		fmt.Println("A SMS is sent successfully:", resultJson)
	} else {
		fmt.Println("Failed to send a SMS:", resultJson)
	}
}
//商品详情页
func(this*GoodsController)ShowDetail(){
	//获取数据
	id,err:= this.GetInt("id")
	//校验数据
	if err != nil{
		beego.Error("获取数据不存在")
		this.Redirect("/",302)
		return
	}
	//数据处理
	//获取商品类型
	var goodsTypes []models.GoodsType
	o:=orm.NewOrm()
	o.QueryTable("GoodsType").All(&goodsTypes)
	this.Data["types"]=goodsTypes
	//获取商品详情
	var goods models.GoodsSKU
	o.QueryTable("GoodsSKU").RelatedSel("GoodsType","Goods").Filter("Id",id).One(&goods)
	this.Data["goods"] = goods

	//添加历史浏览记录，需要先查询有没有登陆，只有登陆之后可以添加历史浏览记录
	userName := this.GetSession("userName")


	if userName != nil{
		//查询用户信息
		var user models.User
		user.UserName = userName.(string)
		o.Read(&user,"UserName")
		conn,_:=redis.Dial("tcp",":6379")
		//先清空以前的记录
		reply,err:=conn.Do("lrem","history"+strconv.Itoa(user.Id),0,id)//？
		reply,_ = redis.Bool(reply,err)
		if reply == false{
			beego.Info("插入浏览数据错误")
		}
		//插入历史纪录
		conn.Do("lpush","history"+strconv.Itoa(user.Id),id)
	}
this.TplName="detail.html"

}

//列表页商品内容
func(this*GoodsController)ShowGoodsList(){


	//获取类型id
	typeId,err := this.GetInt("typeId")
	o:=orm.NewOrm()
	//校验数据
	if err != nil{
		beego.Info("获取类型ID错误")
		this.Redirect("/",302)
		return
	}
	//类型数据
	var types []models.GoodsType
	o.QueryTable("GoodsType").All(&types)
	this.Data["types"] = types



	//获取新品数据 这里获取的是同类型，时间靠前的两个商品数据
	var goodsNew []models.GoodsSKU
	o.QueryTable("GoodsSKU").RelatedSel("GoodsType").Filter("GoodsType__Id",typeId).
		OrderBy("Time").Limit(2,0).All(&goodsNew)

	this.Data["goodsNew"] = goodsNew

	//分页处理
	//指定每页显示多少个数据
	pageSize :=2
	pageIndex ,err:=this.GetInt("pageIndex")
	if err !=nil{
		pageIndex = 1
	}
	start := pageSize * (pageIndex - 1)

	count,_ := o.QueryTable("GoodsSKU").RelatedSel("GoodsType").Filter("GoodsType__Id",typeId).Count()
	//获取总页码
	pageCount := math.Ceil(float64(count)/ float64(pageSize))

	//获取当前类型的商品
	var goodsSKus []models.GoodsSKU

	//根据不同的选项获取不同的数据
	sort := this.GetString("sort")
	//如果sort等于空，就按照默认排序
	if sort == ""{	o.QueryTable("GoodsSKU").RelatedSel("GoodsType").Filter("GoodsType__Id",typeId).Limit(pageSize,start).All(&goodsSKus)
	}else if sort == "price"{	o.QueryTable("GoodsSKU").RelatedSel("GoodsType").Filter("GoodsType__Id",typeId).OrderBy("Price").Limit(pageSize,start).All(&goodsSKus)
	}else{
		o.QueryTable("GoodsSKU").RelatedSel("GoodsType").Filter("GoodsType__Id",typeId).OrderBy("Sales").Limit(pageSize,start).All(&goodsSKus)
	}
	this.Data["goods"]=goodsSKus
	//调用分页助手函数
	page:=5
	pageData:=PageTool(int(pageCount),page,pageIndex)//pageData是map[string]
	//把数据传递给视图
	this.Data["pagePre"] = pageData["pagePre"]
	this.Data["pageNext"] = pageData["pageNext"]
	this.Data["pages"] = pageData["pageIndex"].([]int) //存放着应该显示的页码
	this.Data["pageIndex"] = pageIndex //告诉视图当前的是第几页

	this.Data["typeId"] = typeId
	this.TplName="list.html"
}


//算出分页后显示的是哪个些页码数
func GetPage(pageCount int,pageIndex int ){
	/*
	if pageCount < 5{
		pageIndexBuffer := make([]int,pageCount)
		for index,_ := range  pageIndexBuffer{
			pageIndexBuffer[index] = index + 1
		}
	}else if pageIndex < 3  {      //显示前三页
		pageIndexBuffer := make([]int,5)
		for index,_ := range pageIndexBuffer{
			pageIndexBuffer[index] = index + 1
		}
	}else if pageIndex >= pageCount -3{   //最后三页
		pageIndexBuffer := make([]int,page)
		for index,_ := range  pageIndexBuffer{
			pageIndexBuffer[index] = pageCount - 5 + index
		}
	}else {
		pageIndexBuffer := make([]int, 5)
		for index,_ := range pageIndexBuffer{
			pageIndexBuffer[index] = pageIndex - 3 + index
		}
		*//* 假如index是4
			1 2 3 4 5 6 7
			a[0] 4-3+0
			a[1]  4-3+1
					3
				   4
			a[4]    5
		*//*
	}*/
}

//**获取显示的页码
// 封装一个获取相应页码的函数,需要把总页码，当前页码，和要显示多少个页码当参数传递
// **
//指定显示多少个页码
//page := 5
//函数定义如下

//分页助手函数
func PageTool(pageCount int,page int,pageIndex int)map[string]interface{}{
	//获取应该显示的页码
	var pageIndexBuffer []int
	if pageCount < page{   //第一种：总页码数不到五页
		pageIndexBuffer = make([]int,pageCount)
		for index,_ := range  pageIndexBuffer{
			pageIndexBuffer[index] = index + 1
		}
	}else if pageIndex < ( page + 1)/2  {    //判断是否小于中间页
		pageIndexBuffer = make([]int,page)
		for index,_ := range pageIndexBuffer{
			pageIndexBuffer[index] = index + 1
		}
	}else if pageIndex > pageCount -( page + 1)/2{
		pageIndexBuffer = make([]int,page)
		for index,_ := range  pageIndexBuffer{
			pageIndexBuffer[index] = pageCount - page + index
		}
	}else {
		pageIndexBuffer = make([]int, page)
		for index,_ := range pageIndexBuffer{
			pageIndexBuffer[index] = pageIndex - (page - 1)/2 + index
		}
	}

	//上一页页码
	pagePre := pageIndex - 1
	if pageIndex == 1{
		pagePre = 1
	}
	pageNext := pageIndex + 1
	if pageIndex == pageCount{
		pageNext = pageIndex
	}
	//把数据返回
	pageData := make(map[string]interface{})
	pageData["pagePre"] = pagePre
	pageData["pageNext"] = pageNext
	pageData["pageIndex"] = pageIndexBuffer

	return pageData
}

func(this*GoodsController)HandleSearch(){
	//获取数据
	goodsName := this.GetString("goodsName")
	o := orm.NewOrm()
	var goods []models.GoodsSKU

	//校验数据
	if goodsName == ""{
		beego.Info("查找的数据为空")
		o.QueryTable("GoodsSKU").All(&goods)
		this.Data["goods"] = goods
		//ShowLayout(&this.Controller)
		this.TplName = "search.html"
	}
	//根据拿到的数据进数据库查询

	o.QueryTable("GoodsSKU").Filter("Name__icontains",goodsName).All(&goods)
	//返回数据
	this.Data["goods"] = goods
	this.TplName = "search.html"
}
