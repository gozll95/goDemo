Simple Media Player(SMP)

# 功能:
- 音乐库功能,使用者可以查看、添加和删除里面的音乐曲目
- 播放音乐
- 支持MP3和WAV,但是同时也能随时扩展以支持更多的音乐类型
- 退出程序

# 将接受以下命令:
- 音乐库管理命令:lib,包括list/add/remove命令
- 播放管理: play命令,play后带歌曲名参数
- 退出程序: q命令

## 音乐库

### 音乐库管理模块
每首音乐都包含以下信息:
- 唯一的ID
- 音乐名
- 艺术家名
- 音乐位置
- 音乐文件类型(MP3和WAV等)

type Music struct{
    Id string
    Name string
    Artist string
    Source string
    Type string
}


 $ go run mplayer.go
    Enter following commands to control the player:
    lib list -- View the existing music lib
    lib add <name><artist><source><type> -- Add a music to the music lib
    lib remove <name> -- Remove the specified music from the lib
    play <name> -- Play the specified music
    Enter command-> lib add HugeStone MJ ~/MusicLib/hs.mp3 MP3
    Enter command-> play HugeStone
    Playing MP3 music ~/MusicLib/hs.mp3
    ..........
    Finished playing ~/MusicLib/hs.mp3
    Enter command-> lib list
    1 : HugeStone MJ ~/MusicLib/hs.mp3 MP3
    Enter command-> lib view
    Enter command-> q


# 遗留问题
1. 多任务
当前，我们这个程序还只是单任务程序，即同时只能执行一个任务，比如音乐正在播放时， 用户不能做其他任何事情。作为一个运行在现代多任务操作系统上的应用程序，这种做法肯定是 无法被用户接受的。音乐播放过程不应导致用户界面无法响应，因此播放应该在一个单独的线程中，并能够与主程序相互通信。而且像一般的媒体播放器一样，在播放音乐的同时，我们甚至也 要支持一些视觉效果的播放，即至少需要这么几个线程:用户界面、音乐播放和视频播放。
考虑到这个需求，我们自然而然地想到了使用Go语言的看家本领goroutine，比如将上面的播 放进行稍微修改后即可将Play()函数作为一个独立的goroutine运行。

 2. 控制播放
因为当前这个设计是单任务的，所以播放过程无法接受外部的输入。然而作为一个成熟的播 放器，我们至少需要支持暂停和停止等功能，甚至包括设置当前播放位置等。假设我们已经将播 放过程放到一个独立的goroutine中，那么现在就是如何对这个goroutine进行控制的问题，这可以 使用Go语言的channel功能来实现。
