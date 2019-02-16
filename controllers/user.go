package controllers

import (
	"github.com/astaxie/beego"
	"regexp"
	"github.com/astaxie/beego/orm"
	"dailyFresh/models"
	"github.com/astaxie/beego/utils"
	"strconv"
)

type UserController struct {
	beego.Controller
}

//展示注册页面
func(this*UserController)ShowRegister(){
	//指定注册页面
	this.TplName = "register.html"
}

//处理注册业务
func(this*UserController)HandleRegister(){
	//获取数据
	userName := this.GetString("user_name")
	pwd := this.GetString("pwd")
	cpwd := this.GetString("cpwd")
	email := this.GetString("email")
	//校验数据
	if userName == "" || pwd == "" || cpwd == "" || email == ""{
		this.Data["errmsg"] = "输入数据不能为空"
		this.TplName = "register.html"
		return
	}
	if pwd != cpwd{
		this.Data["errmsg"] = "两次密码输入不一致"
		this.TplName = "register.html"
		return
	}

	reg,_ := regexp.Compile(`^[a-zA-Z0-9_-]+@[a-zA-Z0-9_-]+(\.[a-zA-Z0-9_-]+)+$`)
	result := reg.FindString(email)
	if result == ""{
		this.Data["errmsg"] = "邮箱格式不正确"
		this.TplName = "register.html"
		return
	}

	//处理数据
	o := orm.NewOrm()
	//获取插入对象
	var user models.User
	//给插入对象赋值
	user.UserName = userName
	user.Pwd = pwd
	user.Email = email
	//插入
	o.Insert(&user)

	//发送邮件   注意password是授权码
	emailConfig := `{"username":"13160549408@163.com","password":"chen3316204","host":"smtp.163.com","port":25}`
	emailSend := utils.NewEMail(emailConfig)//得到的是一个邮箱的对象
	emailSend.From = "13160549408@163.com"
	emailSend.To = []string{email}      //切片代表着 可以群发

	emailSend.Subject = "天天生鲜用户激活"   //发过去的邮件标题
	emailSend.HTML = `<a href="http://192.168.42.139:8080/active?id=`+strconv.Itoa(user.Id)+`">点击激活用户</a>`

	emailSend.Send()

	//返回数据
	//this.Redirect("/login",302)
	this.Ctx.WriteString("注册成功，请去邮箱激活当前用户")
}

//激活当前用户
func(this*UserController)HandleActive(){
	//获取数据
	id,err :=this.GetInt("id")
	//校验数据
	if err != nil {
		this.Data["errmsg"] = "激活用户失败"
		this.TplName = "register.html"
		return
	}
	//处理数据
	//更新操作
	o := orm.NewOrm()
	//获取一个更新对象
	var user models.User
	//给更新对象赋值
	user.Id = id
	//查询
	err = o.Read(&user)
	//更新
	if err != nil{
		this.Data["errmsg"] = "激活用户失败"
		this.TplName = "register.html"
		return
	}
	user.Active = true
	o.Update(&user)


	//返回数据
	this.Redirect("/login",302)
}

//展示登录页面
func(this*UserController)ShowLogin(){
	userName:=this.Ctx.GetCookie("userName")
	if userName==""{

		this.Data["userName"]=""
		this.Data["checked"]=""
	}else{
		this.Data["userName"]=userName
		this.Data["checked"]="checked"
	}
	this.TplName = "login.html"
}

//处理登录业务
func(this*UserController)HandleLogin(){


	//获取数据
	userName := this.GetString("username")
	pwd :=this.GetString("pwd")
	//校验数据
	if userName == "" || pwd == ""{
		this.Data["errmsg"] = "登录失败"
		this.TplName = "login.html"
		return
	}



	//处理数据
	//查询校验
	o := orm.NewOrm()
	//获取查询对象
	var user models.User
	//给查询条件赋值
	user.UserName = userName
	//查询
	err := o.Read(&user,"UserName")
	if err != nil{
		this.Data["errmsg"] = "用户名不存在"
		this.TplName = "login.html"
		return
	}

	if user.Pwd != pwd{
		this.Data["errmsg"] = "密码错误"
		this.TplName = "login.html"
		return
	}

	if user.Active == false{
		this.Data["errmsg"] = "当前用户未激活,请先去邮箱激活"
		this.TplName = "login.html"
		return
	}
	//是否有勾选记住用户名

	remember:= this.GetString("remember")
	if remember=="on" {
		this.Ctx.SetCookie("userName",userName,60*60*24)
	}else{
		this.Ctx.SetCookie("userName",userName,-1)
	}
	//返回数据
	//存入session
	this.SetSession("userName",userName)

	this.Redirect("/",302)
}
//退出用户
func(this*UserController)Logout() {
	this.DelSession("userName")
	this.Redirect("/login",302)
}

//用户中心信息
func(this*UserController)ShowUserCenterInfo(){
	this.Layout="layout.html"
	this.TplName="user_center_info.html"
}
//用户中心订单
func(this*UserController)ShowUserCenterOrder(){
	this.Layout="layout.html"
	this.TplName="user_center_order.html"
}


//用户中心地址
func(this*UserController)ShowUserCenterSite(){
	userName:=this.GetSession("userName")
	//获取默认地址
	o:=orm.NewOrm()
	//获取查询对象
	var address models.Address
	//查询
	o.QueryTable("address").RelatedSel("user").Filter("user__UserName",userName.(string)).Filter("Default",true).One(&address)
	//返回数据
	this.Data["address"] = address

	this.Layout="layout.html"
	this.TplName="user_center_site.html"
}

//处理地址信息
func(this*UserController)HandleSite(){
	//获取数据
	receiver:=this.GetString("receiver")
	addr:=this.GetString("addr")
	zipCode:=this.GetString("zipCode")
	phone:=this.GetString("phone")

	//校验数据
	if receiver==""||addr==""||zipCode==""||phone==""{
		this.Data["errmsg"]="地址信息输入不完整"
		this.TplName="user_center_site.html"
		return
	}
	//处理数据
	o:=orm.NewOrm()
	var address models.Address
	address.Receiver=receiver
	address.Addr=addr
	address.ZipCode=zipCode
	address.Phone=phone

	userName:=this.GetSession("userName")
	var user models.User
	user.UserName=userName.(string)
	err :=o.Read(&user,"UserName")
	if err!=nil {
		beego.Info("没有此用户")
		return
	}
	address.Default=true
	address.User=&user
	//查询是否有默认地址，如果有，更新为非默认地址
	var oldAddress models.Address
	qs:=o.QueryTable("address").RelatedSel("user").Filter("user__UserName",userName.(string))
	err=qs.Filter("Default",true).One(&oldAddress)
	if err ==nil {
		oldAddress.Default=false
		o.Update(&oldAddress)
	}

	_,err=o.Insert(&address)
	if err ==nil{
		beego.Info("insert success")
	}
	//返回数据
	this.Redirect("/goods/userCenterSite",302)
}
