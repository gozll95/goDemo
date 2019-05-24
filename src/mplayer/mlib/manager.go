package mlib

import (
	"errors"
)

type MusicManager struct {
	musics []MusicEntry
}

//创建Music管理对象
func NewMusicManager() *MusicManager {
	a := make([]MusicEntry, 0)
	return &MusicManager{
		musics: a,
	}
}

//长度
func (m *MusicManager) Len() int {
	return len(m.musics)
}

//通过序号得到音乐文件
func (m *MusicManager) Get(index int) (music *MusicEntry, err error) {
	if index < 0 || index > len(m.musics) {
		return nil, errors.New("Index out of range.")
	}

	return &m.musics[index], nil

}

func (m *MusicManager) Find(name string) *MusicEntry {
	if len(m.musics) == 0 || name == "" {
		return nil
	}
	//遍历查找
	for _, v := range m.musics {
		if v.Name == name {
			return &v
		}
	}
	//没找到
	return nil
}

func (m *MusicManager) Add(music *MusicEntry) {
	m.musics = append(m.musics, *music)
}

func (m *MusicManager) Remove(index int) *MusicEntry {
	if index < 0 || index >= len(m.musics) {
		return nil
	}

	removedMusic := &m.musics[index]

	//从数组切片中删除元素
	if index > 0 && index < len(m.musics)-1 { //中间元素
		m.musics = append(m.musics[:index-1], m.musics[index+1:]...)
	} else if index == 0 { //删除仅有的一个元素
		m.musics = make([]MusicEntry, 0)
	} else { //删除最后一个元素
		m.musics = m.musics[:index-1]
	}

	//返回删除的元素
	return removedMusic
}

func (m *MusicManager) RemoveByName(name string) *MusicEntry {
	if len(m.musics) == 0 {
		return nil
	}

	for i, v := range m.musics {
		if v.Name == name {
			return m.Remove(i)
		}
	}

	return nil
}
