
ctx, cancel := context.WithCancel(context.TODO())  
timer := time.AfterFunc(5*time.Second, func() {  
    cancel()
})
req, err := http.NewRequest("GET", "http://httpbin.org/range/2048?duration=8&chunk_size=256", nil)  
if err != nil {  
    log.Fatal(err)
}
req = req.WithContext(ctx)