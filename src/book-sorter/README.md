我们准备开发一个排序算法的比较程序，
从命令行指定输入的数据文件和输出的数据文件，并指定对应的排序算法

该程序的用法如下所示:
    USAGE: sorter –i <in> –o <out> –a <qsort|bubblesort>

一个具体的执行过程如下:
    $ ./sorter –I in.dat –o out.dat –a qsort
    The sorting process costs 10us to complete.


当然，如果输入不合法，应该给出对应的提示，接下来我们一步步实现这个程序。

# 主程序
- 获取并解析命令行输入;
- 从对应文件中读取输入数据;
- 调用对应的排序函数;
- 将排序的结果输出到对应的文件中
- 打印排序所花费时间的信息。
