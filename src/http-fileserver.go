用go写一个文件服务器很简单：

    http.handle(“/”,  http.FileServer(http.Dir(“doc”))

   http.ListenAndServe(":8888”, nil)

   打来localhost:8888，就能看到doc目录下的所有文件。

   但如果，你想用localhost:8888/doc来显示进入文件目录，则需要

   http.Handle(“/doc", http.StripPrefix(“/doc", http.FileServer(http.Dir(“doc"))))

   http.StripPrefix用于过滤request，参数里的handler的request过滤掉特定的前序，只有这样，才能正确显示文件目录。