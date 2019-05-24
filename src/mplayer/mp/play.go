package mp

import "fmt"

//音乐播放模块 接口
type Player interface {
	Play(source string)
}

//音乐位置及类型
func Play(source, mtype string) {
	var p Player

	switch mtype {
	case "MP3":
		p = &MP3Player{}
	case "WAV":
		p = &WAVPlayer{}
	default:
		fmt.Println("unsupported music type", mtype)
		return
	}
	p.Play(source)
}
