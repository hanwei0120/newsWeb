package controllers

import (
	"github.com/astaxie/beego"
	"fmt"
	"path"
	"time"
	"github.com/astaxie/beego/orm"
	"newsWeb/models"
	"math"
	"strconv"
	"github.com/gomodule/redigo/redis"
	"encoding/gob"
	"bytes"
)

type ArticleController struct {
	beego.Controller
}

func (a *ArticleController) ShowIndex() {
	//session判断登录
	userName := a.GetSession("userName")
	if userName == nil {
		a.Redirect("/login", 302)
		return
	}

	//获取数据并现实
	//高级查询
	o := orm.NewOrm()

	typeName := a.GetString("select")

	qs := o.QueryTable("Article") //查询集合
	var articles []models.Article
	var count int64
	if typeName == "" {
		//总记录数
		count, _ = qs.RelatedSel("ArticleType").Count()
	} else {
		count, _ = qs.RelatedSel("ArticleType").Filter("ArticleType__TypeName", typeName).Count()
	}

	//每页条数
	pagecount := 2
	//总页数
	pageAll := math.Ceil(float64(count) / float64(pagecount))

	pageIndex, err := a.GetInt("pageIndex")
	if err != nil {
		pageIndex = 1
	}
	if typeName == "" {
		qs.Limit(pagecount, (pageIndex-1)*pagecount).RelatedSel("ArticleType").All(&articles)
	} else {
		qs.Limit(pagecount, (pageIndex-1)*pagecount).RelatedSel("ArticleType").Filter("ArticleType__TypeName", typeName).All(&articles)
	}
	//qs.All(&articles)
	a.Data["articles"] = articles
	a.Data["count"] = count
	a.Data["pageAll"] = int(pageAll)
	a.Data["pageIndex"] = pageIndex


	//存到redis中
	conn,err:=redis.Dial("tcp",":6379")
	if err != nil {
		fmt.Println("redis存储错误",err)
		return
	}

	result,err:=redis.Bytes(conn.Do("get","articleTypes"))
	//获取文章类型
	var articleTypes []models.ArticleType

	if len(result)==0{
		o.QueryTable("ArticleType").All(&articleTypes)
		a.Data["articleTypes"] = articleTypes

		//编码器
		var buffer bytes.Buffer
		enc:=gob.NewEncoder(&buffer)
		enc.Encode(articleTypes)
		conn.Do("set","articleTypes",buffer.Bytes())
		fmt.Println("从mysql中取数据")
	}else {
		//解码器
		dec:=gob.NewDecoder(bytes.NewReader(result))
		dec.Decode(&articleTypes)
		fmt.Println("从redis种获取数据",articleTypes)
	}


	a.Data["typeName"] = typeName
	errmsg := a.GetString("errmsg")
	a.Data["errmsg"] = errmsg
	a.Layout = "layout.html"

	a.TplName = "index.html"
}

func (a *ArticleController) ShowAddArticle() {
	o := orm.NewOrm()
	var articleTypes []models.ArticleType
	o.QueryTable("ArticleType").All(&articleTypes)
	a.Data["articleTypes"] = articleTypes
	a.Layout = "layout.html"
	a.TplName = "add.html"
}

func (a *ArticleController) HandleAddArticle() {
	articleName := a.GetString("articleName")
	content := a.GetString("content")
	_, header, err := a.GetFile("uploadname")

	typeName := a.GetString("select")

	if articleName == "" || content == "" || err != nil {
		fmt.Println("获取数据错误！", err)
		a.Layout = "layout.html"
		a.TplName = "add.html"
		return
	}
	//上传文件要校验
	//1.文件大小
	if header.Size > 10000000 {
		a.Data["errmsg"] = "图片太大，清重新选择！"
		a.Layout = "layout.html"
		a.TplName = "add.html"
		return
	}
	//2.文件格式
	ext := path.Ext(header.Filename)
	if ext != ".jpg" && ext != ".png" && ext != ".jpeg" {
		a.Data["errmsg"] = "图片格式错误，清重新选择！"
		a.Layout = "layout.html"
		a.TplName = "add.html"
		return
	}
	//3.防止重名
	fileName := time.Now().Format("20060102150405")
	ImagePath := "static/img/" + fileName + ext
	//err = a.SaveToFile("uploadname", "./static/img/"+fileName+ext)
	err = a.SaveToFile("uploadname", ImagePath)
	if err != nil {
		fmt.Println("存储文件失败，清重来！")
		a.Layout = "layout.html"
		a.TplName = "add.html"
		return
	}
	//插入数据库
	o := orm.NewOrm()
	var article models.Article
	article.Title = articleName
	article.Content = content
	article.Image = ImagePath

	var articleType models.ArticleType
	articleType.TypeName = typeName
	o.Read(&articleType, "TypeName")

	article.ArticleType = &articleType

	id, err := o.Insert(&article)
	if err != nil {
		fmt.Println("插入失败！！", err)
		a.Layout = "layout.html"
		a.TplName = "add.html"
		return
	}
	fmt.Println("插入成功！！id为：", id)
	a.Layout = "layout.html"
	//成功回到首页
	a.Redirect("index.html", 302)

}

//展示详情
func (a *ArticleController) ShowContent() {

	id, err := a.GetInt("id")
	if err != nil {
		fmt.Println("获取id失败！！", err)
		a.Layout = "layout.html"
		a.Redirect("/index?errmsg=文章详情获取失败", 302)
		return
	}
	o := orm.NewOrm()
	var article models.Article
	article.Id = id
	err = o.Read(&article)
	if err != nil {
		fmt.Println("读取错误！！！", err)
		a.Layout = "layout.html"
		a.Redirect("index.html", 302)
		return
	}
	//关联多对多关系（简单查询不去重）
	//o.LoadRelated(&article,"Users")
	var users []models.User
	o.QueryTable("User").Filter("Articles__Article__id", id).Distinct().All(&users)

	a.Data["users"]= users
	a.Data["article"] = article

	//获取最近阅读--多对多插入
	//在登录情况下点击查看详情
	userName := a.GetSession("userName")

	m2m := o.QueryM2M(&article, "Users")

	var user models.User

	user.UserName = userName.(string)
	o.Read(&user, "userName")

	m2m.Add(user)
	a.Layout = "layout.html"
	a.TplName = "content.html"

}

//展示编辑页面
func (a *ArticleController) ShowUpDate() {

	id, err := a.GetInt("id")
	//fmt.Println(id,"`````````````````````````````````````````")
	if err != nil {
		fmt.Println("编辑获取id失败！！", err)
		//a.TplName = "index.html"//空白
		a.Layout = "layout.html"
		a.Redirect("/index?errmsg=编辑页面获取id失败", 302)
		return
	}
	o := orm.NewOrm()
	var article models.Article
	article.Id = id
	err = o.Read(&article)
	if err != nil {
		fmt.Println("编辑时查询错误！！", err)
		a.Layout = "layout.html"
		a.Redirect("/index?errmsg=编辑时查询错误", 302)
		return
	}
	a.Data["article"] = article
	a.Layout = "layout.html"
	a.TplName = "update.html"
}

//编辑提交
func (a *ArticleController) DoUpDate() {
	id, err := a.GetInt("id")
	if err != nil {
		fmt.Println("更新获取id失败！！", err)
		a.Layout = "layout.html"
		a.TplName = "update.html"
		return
	}
	articleName := a.GetString("articleName")
	content := a.GetString("content")

	_, header, err := a.GetFile("uploadname")
	if articleName == "" || content == "" || err != nil {
		fmt.Println("文件获取失败！！", err)
		a.Layout = "layout.html"
		a.TplName = "update.html"
		return
	}
	if header.Size > 10000000 {
		errmsg := "上传图片过大！！！"
		fmt.Println(errmsg)
		a.Data["errmsg"] = errmsg
		a.Layout = "layout.html"
		a.TplName = "update.html"
		return
	}
	ext := path.Ext(header.Filename)
	if ext != ".jpg" && ext != ".png" && ext != ".jpeg" {
		errmsg := "图片格式错误！！"
		fmt.Println(errmsg)
		a.Data["errmsg"] = errmsg
		//a.TplName = "update.html"
		a.Layout = "layout.html"
		a.Redirect("update.html", 302)
		return
	}
	//文件名去重
	filename := time.Now().Format("20060102150405")
	imagePath := "static/img/" + filename + ext
	err = a.SaveToFile("uploadname", imagePath)
	if err != nil {
		errmsg := "文件上传失败！！！"
		fmt.Println(errmsg, err)
		a.Data["errmsg"] = errmsg
		a.Layout = "layout.html"
		a.Redirect("/update?id="+strconv.Itoa(id), 302)
		//a.TplName = "update.html"
		return
	}

	//更新
	o := orm.NewOrm()
	var article models.Article
	article.Id = id
	err = o.Read(&article, "Id")
	if err != nil {
		fmt.Println("找不到原始数据，更新失败！！", err)
		a.Layout = "layout.html"
		a.Redirect("/update?id="+strconv.Itoa(id), 302)
		//a.TplName = "update.html"
		return
	}

	article.Title = articleName
	article.Content = content
	article.Image = imagePath
	fmt.Println(article.Title, article.Content, article.Image)

	_, err = o.Update(&article)
	if err != nil {
		errmsg := "更新失败！！！"
		fmt.Println(errmsg, err)
		a.Layout = "layout.html"
		a.Redirect("/update?id="+strconv.Itoa(id), 302)
		//a.TplName = "update.html"
		return
	}

	fmt.Println("编辑图片成功！！！", id)
	a.Layout = "layout.html"
	a.Redirect("/index", 302)
}

func (a *ArticleController) Delete() {
	id, err := a.GetInt("id")
	if err != nil {
		fmt.Println("获取id失败", err)
		a.Layout = "layout.html"
		a.Redirect("/index?errmsg=删除失败", 302)
		return
	}

	o := orm.NewOrm()
	var article models.Article
	article.Id = id
	_, err = o.Delete(&article)
	if err != nil {
		fmt.Println("删除失败", err)
		a.Layout = "layout.html"
		a.Redirect("/index?errmsg=删除失败", 302)
		return
	}

	fmt.Println("删除成功")
	a.Layout = "layout.html"
	a.Redirect("/index", 302)

}

//展示类型
func (a *ArticleController) ShowAddType() {
	o := orm.NewOrm()
	var articleTypes []models.ArticleType
	o.QueryTable("ArticleType").All(&articleTypes)

	a.Data["articleTypes"] = articleTypes
	a.Layout = "layout.html"
	a.TplName = "addType.html"
}

//添加类型
func (a *ArticleController) DoAddType() {
	typeName := a.GetString("typeName")

	if typeName == "" {
		a.Data["errmsg"] = "不能为空"
		a.Layout = "layout.html"
		a.TplName = "addType.html"
		return
	}

	o := orm.NewOrm()
	var articleType models.ArticleType
	articleType.TypeName = typeName

	_, err := o.Insert(&articleType)
	if err != nil {
		fmt.Println("插入失败")
		a.Data["errmsg"] = "插入失败"
		a.Layout = "layout.html"
		a.TplName = "addType.html"
		return
	}

	fmt.Println("插入成功")
	a.Layout = "layout.html"
	a.Redirect("/addType", 302)

}

//处理首页数据(选择标签)
func (a *ArticleController) HandleIndex() {

}
//删除类型
func (a *ArticleController) DeleteArticleType(){
	id,err := a.GetInt("id")
	if err != nil {
		fmt.Println("获取到的id为",id)
		a.Layout = "layout.html"
		a.Redirect("/addType",302)
		return
	}
	o:=orm.NewOrm()
	var articleType models.ArticleType

	articleType.Id = id
	o.Delete(&articleType,"Id")
	a.Layout = "layout.html"

	a.Redirect("/addType",302)
}