package mlib

/*
**** 管理对象为音乐，定义音乐结构体类型
****
**** 音乐结构体：id,name,artist,source,type
**** 音乐id、音乐名、艺术家名、音乐位置、类型
 */

type MusicEntry struct {
	Id     string
	Name   string
	Artist string
	Source string
	Type   string
}
