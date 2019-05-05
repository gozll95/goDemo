#反射
1.反射: 可以在运行时动态获取变量的相关信息

两个函数:
    a. reflect.TypeOf,获取变量的类型,返回reflect.Type类型
    b. reflect.ValueOf,获取变量的值,返回relect.Value类型
    c. reflect.Value.Kind,获取变量的类别,返回一个常量
    d. reflect.Value.Interface(),转换成interface{}类型

变量 <--> interface{} <---> Reflect.Value


通过反射来改变变量的值:

reflect.Value.SetXX 相关方法,比如:
reflect.Value.SetFoat()
reflect.Value.SetInt()
reflect.Value.SetString()


