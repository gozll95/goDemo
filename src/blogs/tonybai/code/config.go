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
/*
就像上面代码中那样，你需要：
1、用命令行标志默认值设置配置(cf)默认值。
2、接下来加载配置文件
3、用配置值(cf)覆盖命令行标志变量值
4、解析命令行参数
5、用命令行标志变量值覆盖配置(cf)值。
少一步你都无法实现三层配置能力。
*/