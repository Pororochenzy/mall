package models

import (
	"github.com/astaxie/beego/orm"
	_ "github.com/go-sql-driver/mysql"
	"time"
)

//用户表
type User struct {
	Id int
	UserName string `orm:"unique;size(50)"`
	Pwd string `orm:"size(100)"`
	Email string `orm:"size(50)"`
	Active bool `orm:"default(false)"`
	Power int `orm:"default(1)"`
	Addresses []*Address `orm:"reverse(many)"`
}

//收件人表
type Address struct {
	Id int
	Receiver string `orm:"size(50)"`
	Addr string `orm:"size(200)"`
	ZipCode string
	Phone string
	Default bool `orm:"default(false)"`
	User *User `orm:"rel(fk)"`
}
type Goods struct { //商品SPU表
	Id 		int
	Name 	string`orm:"size(20)"`  //商品名称
	Detail 	string`orm:"size(200)"` //详细描述
	GoodsSKU []*GoodsSKU `orm:"reverse(many)"`
}

type GoodsType struct{//商品类型表
	Id int
	Name string			//种类名称
	Logo string			//logo
	Image string   		//图片
	GoodsSKU []*GoodsSKU `orm:"reverse(many)"`
	IndexTypeGoodsBanner  []*IndexTypeGoodsBanner  `orm:"reverse(many)"`
}

type GoodsSKU struct { //商品SKU表
	Id int
	Goods     *Goods 	 `orm:"rel(fk)"` //商品SPU
	GoodsType *GoodsType `orm:"rel(fk)"`  //商品所属种类
	Name       string					 //商品名称
	Desc       string					 //商品简介
	Price      int						 //商品价格
	Unite      string					 //商品单位
	Image      string				 	 //商品图片
	Stock      int	`orm:"default(1)"`	 //商品库存
	Sales      int	`orm:"default(0)"`	 //商品销量
	Status     int	 `orm:"default(1)"`	 //商品状态
	Time       time.Time `orm:"auto_now_add"`  //添加时间
	GoodsImage []*GoodsImage `orm:"reverse(many)"`
	IndexGoodsBanner   []*IndexGoodsBanner `orm:"reverse(many)"`
	IndexTypeGoodsBanner []*IndexTypeGoodsBanner  `orm:"reverse(many)"`
}

type GoodsImage struct { //商品图片表
	Id 			int
	Image 		string					//商品图片
	GoodsSKU 	*GoodsSKU   `orm:"rel(fk)"` //商品SKU
}
type IndexGoodsBanner struct { //首页轮播商品展示表
	Id 		  int
	GoodsSKU  *GoodsSKU	`orm:"rel(fk)"`	//商品sku
	Image     string					//商品图片
	Index     int  `orm:"default(0)"`   //展示顺序
}

type IndexTypeGoodsBanner struct {//首页分类商品展示表
	Id 				int
	GoodsType 		*GoodsType 	`orm:"rel(fk)"`			//商品类型
	GoodsSKU  		*GoodsSKU  	`orm:"rel(fk)"`			//商品sku
	DisplayType 	int   		`orm:"default(1)"`		//展示类型 0代表文字，1代表图片
	Index 			int   		`orm:"default(0)"`		//展示顺序
}

type IndexPromotionBanner struct {//首页促销商品展示表
	Id 		int
	Name 	string	`orm:"size(20)"`				//活动名称
	Url 	string	`orm:"size(50)"`				//活动链接
	Image 	string						//活动图片
	Index 	int  `orm:"default(0)"` //展示顺序
}



func init(){
	//注册库
	orm.RegisterDataBase("default","mysql","root:123456@tcp(127.0.0.1:3306)/ttsx?charset=utf8")
	//注册表
	orm.RegisterModel(new(User),new(Address),new(Goods),new(GoodsSKU),new(GoodsImage),new(GoodsType),new(IndexGoodsBanner),new(IndexPromotionBanner),new(IndexTypeGoodsBanner))

	//跑起来
	orm.RunSyncdb("default",false,true)
}