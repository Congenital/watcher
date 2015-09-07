package watcher

import (
	"github.com/Congenital/log/v0.2/log"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

/*

信号    取值    默认动作    含义（发出信号的原因）
SIGHUP  1   Term    终端的挂断或进程死亡
SIGINT  2   Term    来自键盘的中断信号
SIGQUIT     3   Core    来自键盘的离开信号
SIGILL  4   Core    非法指令
SIGABRT     6   Core    来自abort的异常信号
SIGFPE  8   Core    浮点例外
SIGKILL     9   Term    杀死
SIGSEGV     11  Core    段非法错误(内存引用无效)
SIGPIPE     13  Term    管道损坏：向一个没有读进程的管道写数据
SIGALRM     14  Term    来自alarm的计时器到时信号
SIGTERM     15  Term    终止
SIGUSR1     30,10,16    Term    用户自定义信号1
SIGUSR2     31,12,17    Term    用户自定义信号2
SIGCHLD     20,17,18    Ign     子进程停止或终止
SIGCONT     19,18,25    Cont    如果停止，继续执行
SIGSTOP     17,19,23    Stop    非来自终端的停止信号
SIGTSTP     18,20,24    Stop    来自终端的停止信号
SIGTTIN     21,21,26    Stop    后台进程读终端
SIGTTOU     22,22,27    Stop    后台进程写终端

SIGBUS  10,7,10     Core    总线错误（内存访问错误）
SIGPOLL         Term    Pollable事件发生(Sys V)，与SIGIO同义
SIGPROF     27,27,29    Term    统计分布图用计时器到时
SIGSYS  12,-,12     Core    非法系统调用(SVr4)
SIGTRAP     5   Core    跟踪/断点自陷
SIGURG  16,23,21    Ign     socket紧急信号(4.2BSD)
SIGVTALRM   26,26,28    Term    虚拟计时器到时(4.2BSD)
SIGXCPU     24,24,30    Core    超过CPU时限(4.2BSD)
SIGXFSZ     25,25,31    Core    超过文件长度限制(4.2BSD)

SIGIOT  6   Core    IOT自陷，与SIGABRT同义
SIGEMT  7,-,7       Term
SIGSTKFLT   -,16,-  Term    协处理器堆栈错误(不使用)
SIGIO   23,29,22    Term    描述符上可以进行I/O操作
SIGCLD  -,-,18  Ign     与SIGCHLD同义
SIGPWR  29,30,19    Term    电力故障(System V)
SIGINFO     29,-,-      与SIGPWR同义
SIGLOST     -,-,-   Term    文件锁丢失
SIGWINCH    28,28,20    Ign     窗口大小改变(4.3BSD, Sun)
SIGUNUSED   -,31,-  Term    未使用信号(will be SIGSYS)

*/

const (
	Parallel = 0
	Linear   = 1
)

type Handler func()
type DefaultHandler func(os.Signal)

var signal_buff_size = 0

var signal_channel chan os.Signal
var signal_buff_size_lock = &sync.RWMutex{}

var exit_channel = make(chan int)

var handler_list = make(map[os.Signal]Handler)
var list_lock = &sync.RWMutex{}

var run_model int = Linear
var run_model_lock = &sync.RWMutex{}

var defaultHandler DefaultHandler = DefaultHandle
var defaultHandlerLock = &sync.RWMutex{}

func init() {
	log.Info("init")
}

func Listen() {
	log.Info("Listen")
	signal_channel = make(chan os.Signal, signal_buff_size)

	signal.Notify(signal_channel,
		syscall.SIGHUP,  /*终端的挂断或进程死亡*/
		syscall.SIGINT,  /*来自键盘的中断信号*/
		syscall.SIGQUIT, /*来自键盘的离开信号*/
		syscall.SIGILL,  /*非法指令*/
		syscall.SIGABRT, /*来自abort的异常信号*/
		syscall.SIGFPE,  /*浮点例外*/
		syscall.SIGKILL, /*杀死*/
		syscall.SIGSEGV, /*段非法错误(内存引用无效)*/
		syscall.SIGPIPE, /*管道损坏：向一个没有读进程的管道写数据*/
		syscall.SIGALRM, /*来自alarm的计时器到时信号*/
		syscall.SIGTERM, /*终止*/
		syscall.SIGUSR1, /*用户自定义信号1*/
		syscall.SIGUSR2, /*用户自定义信号2*/
		syscall.SIGCHLD, /*子进程停止或终止*/
		syscall.SIGCONT, /*如果停止，继续执行*/
		syscall.SIGSTOP, /*非来自终端的停止信号*/
		syscall.SIGTSTP, /*来自终端的停止信号*/
		syscall.SIGTTIN, /*后台进程读终端*/
		syscall.SIGTTOU, /*后台进程写终端*/

		syscall.SIGBUS,    /*总线错误（内存访问错误）*/
		syscall.SIGPOLL,   /*Pollable事件发生(Sys V)，与SIGIO同义*/
		syscall.SIGPROF,   /*统计分布图用计时器到时*/
		syscall.SIGSYS,    /*非法系统调用(SVr4)*/
		syscall.SIGTRAP,   /*跟踪/断点自陷*/
		syscall.SIGURG,    /*socket紧急信号(4.2BSD)*/
		syscall.SIGVTALRM, /*虚拟计时器到时(4.2BSD)*/
		syscall.SIGXCPU,   /*超过CPU时限(4.2BSD)*/
		syscall.SIGXFSZ,   /*超过文件长度限制(4.2BSD)*/

		syscall.SIGIOT,    /*IOT自陷，与SIGABRT同义*/
		syscall.SIGSTKFLT, /*协处理器堆栈错误(不使用)*/
		syscall.SIGIO,     /*描述符上可以进行I/O操作*/
		syscall.SIGTERM,   /*Term*/
		syscall.SIGCLD,    /*与SIGCHLD同义*/
		syscall.SIGPWR,    /*电力故障(System V)*/
		syscall.SIGWINCH,  /*窗口大小改变(4.3BSD, Sun)*/
		syscall.SIGUNUSED, /*未使用信号(will be SIGSYS)*/
	)

	go Switch()
}

func Switch() {
	log.Info("Switch")

	for {
		signal_value := <-signal_channel
		log.Info(signal_value)
		Handle(signal_value)
	}
}

func Handle(_signal os.Signal) {
	log.Info("Handle")

	list_lock.RLock()
	handle, ok := handler_list[_signal]
	list_lock.RUnlock()

	if ok && handle != nil {
		if GetRunModel() == Linear {
			handle()
		} else {
			go handle()
		}
		return
	}

	if GetRunModel() == Linear {
		if defaultHandler != nil {
			defaultHandler(_signal)
		} else {
			DefaultHandle(_signal)
		}
	} else {
		if defaultHandler != nil {
			go defaultHandler(_signal)
		} else {
			go DefaultHandle(_signal)
		}
	}

}

func SetHandle(_signal os.Signal, handle Handler) {
	log.Info("SetHandle")

	list_lock.Lock()
	defer list_lock.Unlock()

	handler_list[_signal] = handle
}

func DefaultHandle(_signal os.Signal) {
	log.Info("DefaultHandle")

	log.Info(_signal)
}

func SetDefaultHandle(handle DefaultHandler) {
	defaultHandlerLock.Lock()
	defer defaultHandlerLock.Unlock()

	defaultHandler = handle
}

func ClearDefaultHandle() {
	defaultHandlerLock.Lock()
	defer defaultHandlerLock.Unlock()

	defaultHandler = DefaultHandle
}

func Exit(code int) {
	log.Info("Exit")

	signal.Stop(signal_channel)
	exit_channel <- code
}

func GetExit() chan int {
	return exit_channel
}

func GetExitCode() int {
	log.Info("GetExit")

	return <-exit_channel
}

func Stop() {
	log.Info("Stop")

	signal.Stop(signal_channel)
}

func SetRunModel(model int) {
	run_model_lock.Lock()
	defer run_model_lock.Unlock()

	run_model = model
}

func GetRunModel() int {
	run_model_lock.RLock()
	defer run_model_lock.RUnlock()

	return run_model
}

func SetParallel() {
	SetRunModel(Parallel)
}

func SetLinear() {
	SetRunModel(Linear)
}

func GetBuffSize() int {
	signal_buff_size_lock.RLock()
	defer signal_buff_size_lock.RUnlock()

	return signal_buff_size
}

func SetBuffSize(size int) {
	signal_buff_size_lock.Lock()
	defer signal_buff_size_lock.Unlock()

	signal_buff_size = size
}

func SendSignal(pid int, _signal os.Signal) {
	process, err := os.FindProcess(pid)
	if err != nil {
		log.Error(err)
		return
	}

	process.Signal(_signal)
}

func ReStart(pid int) {
	SendSignal(pid, syscall.SIGUSR1)
}

func ShutDown(pid int) {
	SendSignal(pid, syscall.SIGUSR2)
}
