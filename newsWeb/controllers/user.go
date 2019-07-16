package controllers

import (
	"github.com/astaxie/beego"
	"fmt"
	"github.com/astaxie/beego/orm"
	"newsWeb/models"
	"encoding/base64"
)

type LoginController struct {
	beego.Controller
}

type RegisterController struct {
	beego.Controller
}


func (r *RegisterController) ShowRegister() {
	r.TplName = "register.html"
}
func (r *RegisterController) DoRegister() {
	userName := r.GetString("userName")
	pwd := r.GetString("pwd")
	if userName == "" || pwd == "" {
		r.Data["errmsg"] = "用户名或者密码不能为空"
		r.TplName = "login.html"
		return
	}
	//fmt.Println(userName,pwd)
	o := orm.NewOrm()
	var user models.User
	user.UserName = userName
	user.Pwd = pwd

	id, err := o.Insert(&user)
	if err != nil {
		fmt.Println("注册失败！！！", err)
		r.TplName = "register.html"
		return
	}
	fmt.Println("注册成功！！！id为：", id)

	//r.TplName = "login.html"
	r.Redirect("/login", 302)
}

func (l *LoginController) ShowLogin() {
	userName := l.Ctx.GetCookie("userName")
	if userName!=""{
		decUserName,_:=base64.StdEncoding.DecodeString(userName)
		l.Data["userName"] = string(decUserName)
		l.Data["checked"] = "checked"
	}else {
		l.Data["userName"] = ""
		l.Data["checked"] = ""
	}


	l.TplName = "login.html"
}
func (l *LoginController) DoLogin() {
	userName := l.GetString("userName")
	pwd := l.GetString("pwd")
	if userName == "" || pwd == "" {
		l.Data["errmsg"] = "用户名或者密码不能为空"
		l.TplName = "login.html"
		return
	}

	o := orm.NewOrm()
	var user models.User
	user.UserName = userName

	err := o.Read(&user, "UserName")
	if err != nil {
		l.Data["errmsg"] = "账号不存在登录失败！！"
		fmt.Println("账号不存在登录失败！！", err)
		l.TplName = "login.html"
		return
	}
	//fmt.Println(user)
	if user.Pwd != pwd {
		l.Data["errmsg"] = "密码错误，登陆失败！！！"
		fmt.Println("密码错误，登陆失败！！！")
		l.TplName = "login.html"
		return
	} else {
		remember := l.GetString("remember")
		if remember == "on" {
			enc:=base64.StdEncoding.EncodeToString([]byte(userName))
			l.Ctx.SetCookie("userName", enc, 3600)
		}else {
			l.Ctx.SetCookie("userName", userName, -1)
		}

		l.SetSession("userName",userName)

		l.Redirect("/index", 302)
	}

}

func (l *LoginController) Exit(){
	l.DelSession("userName")
	l.Redirect("/login",302)
}
