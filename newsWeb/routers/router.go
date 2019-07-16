package routers

import (
	"newsWeb/controllers"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/context"
)

func init() {
	beego.InsertFilter("/article/*",beego.BeforeExec,filters)

    beego.Router("/", &controllers.MainController{})
    beego.Router("/register",&controllers.RegisterController{},"get:ShowRegister;post:DoRegister")
	beego.Router("/login", &controllers.LoginController{},"post:DoLogin;get:ShowLogin")
	beego.Router("/index",&controllers.ArticleController{},"get:ShowIndex;post:HandleIndex")
	beego.Router("/addArticle",&controllers.ArticleController{},"get:ShowAddArticle;post:HandleAddArticle")
	beego.Router("/content",&controllers.ArticleController{},"get:ShowContent")
	beego.Router("/update",&controllers.ArticleController{},"get:ShowUpDate;post:DoUpDate")
	beego.Router("/delete",&controllers.ArticleController{},"get:Delete")
	beego.Router("/addType",&controllers.ArticleController{},"get:ShowAddType;post:DoAddType")
	beego.Router("/exit",&controllers.LoginController{},"get:Exit")
	beego.Router("/deleteArticleType",&controllers.ArticleController{},"get:DeleteArticleType")

	}
func filters(ctx *context.Context){
	userName:=ctx.Input.Session("userName")
	if userName==nil{
		ctx.Redirect(302,"/login")
		return
	}

}