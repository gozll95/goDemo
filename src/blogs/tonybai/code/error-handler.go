func handlerError(c net.Conn,err error){
    if err==nil{
        return
    }
    // repeated error handling
}

func DoStuff(c net.Conn)error{
    defer func(){handleError(c,err)}()
    f,err:=downloadFile(c,path)
    if err!=nil{
        return err
    }
    f,err:=doOtherThing(c)
    if err!=nil{
        return err
    }
}