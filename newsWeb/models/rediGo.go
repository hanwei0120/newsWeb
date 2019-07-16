package models

import (
	"github.com/gomodule/redigo/redis"
	"fmt"
)

func init() {
	//链接
	conn,err:=redis.Dial("tcp",":6379")
	if err != nil {
		fmt.Println("redis 链接失败",err)
		return
	}

	//操作数据
	//1send方法
	//conn.Send("set","aa","bb")
	//conn.Flush()
	//conn.Receive()

	//2 DO方法
	//rep,err:=conn.Do("set","aa","bb")
	////回复助手函数--》类型转换
	//result,err:=redis.String(rep,err)
	//fmt.Println("result:",result)

	//多数据
	rep,err:=conn.Do("mget","fm1","fm2","fm3")
	result,err:=redis.Values(rep,err)
	var fm1,fm2 string
	var fm3 int
	redis.Scan(result,&fm1,&fm2,&fm3)
	fmt.Println(fm1,fm2,fm3)

	//自定义数据获取 结构体 map
	//序列化3种 1.json 2.binary 3.gob





}
