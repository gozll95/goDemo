尽量使用

slice  []*struct 

这样

for _,v:=range slice {
	v.xx=yy
}

就直接修改了
不然修改的是副本