/*
HTTP请求包(浏览器信息)
GET /domains/example/ HTTP/1.1      //请求行: 请求方法 请求URI HTTP协议/协议版本
Host：www.iana.org               //服务端的主机名
User-Agent：Mozilla/5.0 (Windows NT 6.1) AppleWebKit/537.4 (KHTML, like Gecko) Chrome/22.0.1229.94 Safari/537.4          //浏览器信息
Accept：text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8  //客户端能接收的mine
Accept-Encoding：gzip,deflate,sdch       //是否支持流压缩
Accept-Charset：UTF-8,*;q=0.5        //客户端字符编码集
//空行,用于分割请求头和消息体
//消息体,请求资源参数,例如POST传递的参数
*/