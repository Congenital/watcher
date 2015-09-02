# watcher

###Run Model
>Parallel
>Linear
>>SetRunModel(model int)
>>>func SetParallel()
>>>func SetLinear()

>Handle
>>type Handler func()
>>SetHandle(_signal os.Signal, handle Handler)
>DefaultHandle
>>nothing

>Start Watcher
>>Listen()
