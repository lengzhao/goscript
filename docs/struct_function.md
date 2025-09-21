为我梳理代码中逻辑，给出详细的流程文档：
1. struct方法调用的编译逻辑
2. 普通方法调用的编译逻辑
3. opCall的执行逻辑

改造：
1. struct.Function，比如有struct Rectangle，对应的方法为Area，则自动注册方法Rectangle.Area( r Rectagle)
2. 如果是可以修改struct的，指针传递的，则注册的是Rectangle.Set(r *Rectagle, value int)
自动将struct对象转成函数参数的第一个参数
这样的方式，就能够统一普通函数调用和struct的方法调用
