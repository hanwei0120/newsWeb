package models

import(
	"github.com/astaxie/beego/orm"
	_"github.com/go-sql-driver/mysql"
	"fmt"
	"time"
)

//当没有设置主键的时候 会以名为Id 类型为int,int8,int64的在字段当主键
type User struct {
	Id int	`orm:"pk;auto"`
	UserName string `orm:"size(20);unique"`
	Pwd string	`orm:"size(20)"`

	Articles []*Article `orm:"reverse(many)"`
}
//orm建表默认非空
type Article struct {
	Id int `orm:"pk;auto"`
	Title string `orm:"size(50);unique"`
	Content string `orm:"size(500)"`
	Time time.Time `orm:"type(datetime);auto_now_add"`
	ReadCount int `orm:"default(0)"`
	Image string `orm:"null"`

	ArticleType *ArticleType `orm:"rel(fk)"`
	Users []*User `orm:"rel(m2m)"`
}

type ArticleType struct {
	Id int
	TypeName string `orm:"size(50)"`

	Articles []*Article `orm:"reverse(many)"`

}


func init()  {
	orm.RegisterDataBase("default","mysql","root:123456@tcp(127.0.0.1:3306)/newsWeb")
	orm.RegisterModel(new(User),new(Article),new(ArticleType))
	err:=orm.RunSyncdb("default",false,true)
	if err != nil {
		fmt.Println("RunSyncdb Err:",err)
		return
	}
}
