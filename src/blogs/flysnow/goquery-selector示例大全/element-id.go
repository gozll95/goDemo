/*
id选择器以#开头，紧跟着元素id的值，使用语法为dom.Find(#id),后面的例子我会简写为Find(#id),大家知道这是代表goquery选择器的即可。

如果有相同的ID，但是它们又分别属于不同的HTML元素怎么办？有好办法，和Element结合起来。比如我们筛选元素为div,并且id是div1的元素，就可以使用Find(div#div1)这样的筛选器进行筛选。

所以这类筛选器的语法为Find(element#id)，这是常用的组合方法，比如后面讲的过滤器也可以采用这种方式组合使用。
*/
