beego

基础模块:
cache
config
context
httplibs
logs
orm
session
toolbox


beego的执行逻辑:

main文件监听端口接收请求---->路由功能----->参数过滤----->controller <----> 辅助工具包
                    ----->视图输出<----输出过滤<------           <----> Model <-------> 数据库
                                                              <----> Session管理
                                                              <----> 日志处理
                                                              <----> 缓存处理


bee new
bee api
bee run



基于 beego 的 Controller 设计，只需要匿名组合 beego.Controller 就可以了
beego.Controller 拥有很多方法，其中包括 Init、Prepare、Post、Get、Delete、Head等 方法。我们可以通过重写的方式来实现这些方法，而我们上面的代码就是重写了 Get 方法。

# goconvey - 课时 1：优雅的单元测试


我们经常需要获取用户传递的数据,包括Get、POST等方式的请求,beego里面会自动解析这些数据,你可以通过如下方式获取数据:

GetString(key string)string
GetStrings(key string)[]string
GetInt(key string)(int64,error)
GetBool(key string)(bool,error)
GetFloat(key string)(float64,error)


beego 
# cookie session


this.Ctx.SetCookie("name", name, maxage, "/")
this.Ctx.SetCookie("pwd", Md5([]byte(pwd)), maxage, "/")
this.Ctx.GetCookie



beego 内置了 session 模块，目前 session 模块支持的后端引擎包括 memory、cookie、file、mysql、redis、couchbase、memcache、postgres，用户也可以根据相应的 interface 实现自己的引擎。

SetSession(name string,value interface{})
GetSession(name string)interface{}
DelSession(name string)
SessionRegenateID()
DestroySession()


/*
package controllers

import (
	"github.com/astaxie/beego"
)

type TestLoginController struct {
	beego.Controller
}

type UserInfoV2 struct{
	Username string
	Password string
}

func (c *TestLoginController) Login(){
	name := c.Ctx.GetCookie("name")
	password := c.Ctx.GetCookie("password")

	//do verify work
	if name != ""{
		c.Ctx.WriteString("Username:" + name + " Password:" + password)
	}else{
		c.Ctx.WriteString(`<html><form action="http://127.0.0.1:8080/test_login" method="post"> 
							<input type="text" name="Username"/>
							<input type="password" name="Password"/>
							<input type="submit" value="提交"/>
					   </form></html>`)
	}
}


func (c *TestLoginController) Post(){
	u := UserInfoV2{}
	if err:=c.ParseForm(&u) ; err != nil{
		//process error
	}

	c.Ctx.SetCookie("name", u.Username, 100, "/")
	c.Ctx.SetCookie("password", u.Password, 100, "/")
	c.SetSession("name", u.Username)
	c.SetSession("password", u.Password)
	c.Ctx.WriteString("Username:" + u.Username + " Password:" + u.Password)
}

*/

/*
package controllers

import (
	"github.com/astaxie/beego"
)

type TestInputController struct {
	beego.Controller
}

type User struct{
	Username string
	Password string
}

func (c *TestInputController) Get(){
	//id := c.GetString("id")
	//c.Ctx.WriteString("<html>" + id + "<br/>")

	//name := c.Input().Get("name")
	//c.Ctx.WriteString(name + "</html>")
	name := c.GetSession("name")
	password := c.GetSession("password")

	if nameString, ok := name.(string); ok && nameString != ""{
		c.Ctx.WriteString("Name:" + name.(string) + " password:" + password.(string))
	}else{
		c.Ctx.WriteString(`<html><form action="http://127.0.0.1:8080/test_input" method="post"> 
							<input type="text" name="Username"/>
							<input type="password" name="Password"/>
							<input type="submit" value="提交"/>
					   </form></html>`)
	}
}


func (c *TestInputController) Post(){
	u := User{}
	if err:=c.ParseForm(&u) ; err != nil{
		//process error
	}

	c.Ctx.WriteString("Username:" + u.Username + " Password:" + u.Password)
}

* /


# beego框架之config/httplib/context

go get github.com/astaxie/beego/config

iniconfig,err:=NewConfig("ini","testini.conf")
if err!=nil{
    t.Fatal(err)
}
iniconf.String("appname")

解析器对象支持的函数有如下:
Set(key,val string)error
String(key string)string
Int(key string)(int,error)
Int64(key string)(int64,error)
Bool(key string)(bool,error)
Float(key string)(float64,error)
DIY(key string)(interface{},error)


解析器对象支持的函数有如下：

            ini 配置文件支持 section 操作，key通过 section::key 的方式获取

            例如下面这样的配置文件

            [demo]
            key1 = "asta"
            key2 = "xie"

            那么可以通过 iniconf.String("demo::key2") 获取值

# httplib
httplib 库主要用来模拟客户端发送 HTTP 请求，类似于 Curl 工具，支持 JQuery 类似的链式操作。使用起来相当的方便；通过如下方式进行安装：


#context:
context对象
是对input和output的封装,里面封装了几个方法:

Redirect
Abort
WriteString
GetCookie
SetCookie


# 爬虫
去重处理:
- 布隆过滤器 ?????
- 哈希存储

标签匹配:
- 正则表达式
- beautiful soup/lxml这种标签提取库

动态内容:
- phantomjs
- selenium

