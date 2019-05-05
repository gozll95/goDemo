package main

import (
	"sync"
)

var loadIconsOnce sync.Once
var icons map[string]image.Image

func Icon(name string) image.Image {
	loadIconsOnce.Do(loadIcons)
	return icons[name]
}

/*
每一次对Do(loadIcons)的调用都会锁定mutex，并会检查boolean变量。
在第一次调用时，变量的值是false，Do会调用loadIcons并会将boolean设置为true。
随后的调用什么都不会做，但是mutex同步会保证loadIcons对内存(这里其实就是指icons
变量啦)产生的效果能够对所有goroutine可见。用这种方式来使用sync.Once的话，
我们能够避免在变量被构建完成之前和其它goroutine共享该变量。
*/

/*
sync.Once 就当成
if boolean==false{
	lock
	...
	boolean=true
	unlock
}

if boolean==true{
	什么都不做
}
*/
