# 背景:
配置数据提供的形式:
- 命令行(options) -> flag
- 参数(parameters)
- 环境变量(env vars) -> os
- 配置文件

***一个良好的应用配置层次应该是这样的***:
1.程序内内置配置项的初始默认值
2.配置文件中的配置项可以覆盖override程序内配置项的默认值
3.命令行选项和参数值具有最高优先级,可以override前两层的配置项值

下面就按作者的思路循序渐进探讨golang程序配置方案。

# 二、解析命令行选项和参数

从例子可以看出，简单情形下，你无需编写自己的命令行parser或使用第三方包，使用go内建的flag包即可以很好的完成工作。但是golang的 flag包与命令行Parser的事实标准：Posix getopt（C/C++/Perl/Shell脚本都可用）相比，还有较大差距，主要体现在：
1、无法支持区分long option和short option，比如：-h和–help。
2、不支持short options合并，比如：ls -l -h <=> ls -hl
3、命令行标志的位置不能任意放置，比如无法放在non-flag parameter的后面。


不过毕竟flag是golang内置标准库包，你无须付出任何cost，就能使用它的功能。另外支持bool型的flag也是其一大亮点。


# TOML,go配置文件的事实标准(这个可能不能得到认同)

不过toml也有其不足之处。想想如果你需要使用命令行选项的参数值来覆盖这些配置文件中的选项，你应该怎么做？事实上，我们常常会碰到类似下面这种三层配置结构的情况：
1、程序内内置配置项的初始默认值
2、配置文件中的配置项值可以覆盖(override)程序内配置项的默认值。
3、命令行选项和参数值具有最高优先级，可以override前两层的配置项值。
在go中，toml映射的结果体字段没有初始值。而且go内建flag包也没有将命令行参数值解析为一个go结构体，而是零散的变量。这些可以通过第三方工具来解决，但如果你不想用第三方工具，你也可以像下面这样自己解决，虽然难看一些。

```
func ConfigGet() *Config {
    var err error
    var cf *Config = NewConfig()
    // set default values defined in the program
    cf.ConfigFromFlag()
    //log.Printf("P: %d, B: '%s', F: '%s'\n", cf.MaxProcs, cf.Webapp.Path)
    // Load config file, from flag or env (if specified)
    _, err = cf.ConfigFromFile(*configFile, os.Getenv("APPCONFIG"))
    if err != nil {
        log.Fatal(err)
    }
    //log.Printf("P: %d, B: '%s', F: '%s'\n", cf.MaxProcs, cf.Webapp.Path)
    // Override values from command line flags
    cf.ConfigToFlag()
    flag.Usage = usage
    flag.Parse()
    cf.ConfigFromFlag()
    //log.Printf("P: %d, B: '%s', F: '%s'\n", cf.MaxProcs, cf.Webapp.Path)
    cf.ConfigApply()
    return cf
}
```


就像上面代码中那样，你需要：
1、用命令行标志默认值设置配置(cf)默认值。
2、接下来加载配置文件
3、用配置值(cf)覆盖命令行标志变量值
4、解析命令行参数
5、用命令行标志变量值覆盖配置(cf)值。
少一步你都无法实现三层配置能力。


