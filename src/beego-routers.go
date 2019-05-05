        // 注册 beego 路由
        beego.Router("/", &controllers.HomeController{})
        beego.Router("/category", &controllers.CategoryController{})
        beego.Router("/topic", &controllers.TopicController{})
        beego.AutoRouter(&controllers.TopicController{})
        beego.Router("/reply", &controllers.ReplyController{})
        beego.Router("/reply/add", &controllers.ReplyController{}, "post:Add")
        beego.Router("/reply/delete", &controllers.ReplyController{}, "get:Delete")
        beego.Router("/login", &controllers.LoginController{})

        // 附件处理
        os.Mkdir("attachment", os.ModePerm)
        beego.Router("/attachment/:all", &controllers.AttachController{})