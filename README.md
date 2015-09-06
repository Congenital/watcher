# watcher(一个golang的signal监听回调模块)

###Run Model (运行模式)
>Parallel (并行)
>Linear (线性)
>>SetRunModel(model int) (设置运行模式)
>>>func SetParallel() (并行模式)
>>>func SetLinear() (线性模式)

>Handle (处理接口)
>>type Handler func() (处理接口实际原型 func())
>>SetHandle(_signal os.Signal, handle Handler) (设置对应信号处理函数)
>DefaultHandle() (预设处理函数)
>>nothing (什么都不做)

>Buff (处理缓冲区)
>>GetBuffSize() int 获得缓冲区大小
>>SetBuffSize(size int) 设置缓冲区大小

>Start Watcher (启动监听)
>>Listen() (开始监听)

>Exit (退出监听)
>>Exit(code int) (退出监听)
>>GetExit() chan int (获得退出Code Channel)
>>GetExitCode() int (获得退出Code)

>Send (发送信号)
>>SendSignal(_pid int, _signal os.Signal) 向指定pid进程发送_signal信号
