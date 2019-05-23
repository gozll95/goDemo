/*
	len(s)

	如果s是string，len(s)返回字符串中的字节个数
	如何s是[n]T, *[n]T的数组类型，len(s)返回数组的长度n
	如果s是[]T的Slice类型，len(s)返回slice的当前长度
	如果s是map[K]T的map类型，len(s)返回map中的已定义的key的个数
	如果s是chan T类型，那么len(s)返回当前在buffered channel中排队（尚未读取）的元素个数
*/

