package container

import (
	"github.com/sirupsen/logrus"
	"os"
	"syscall"
)

/*
这里的init函数是再容器内部执行的，
也就是说，代码执行这里后，
容器所在的进程其实已经创建出来了，
这是本容器执行的第一个进程。
使用mount先去挂载proc文件系统，
以便后面通过ps等系统命令去查看当前进程资源情况。
 */
func RunContainerInitProcess(command string, args []string) error {

	logrus.Infof("command %s", command)
	defaultMountFlags := syscall.MS_NOEXEC | syscall.MS_NOSUID | syscall.MS_NODEV
	syscall.Mount("proc", "/proc", "proc", uintptr(defaultMountFlags), "")
	argv := []string{command}
	if err := syscall.Exec(command, argv, os.Environ()); err != nil {
		logrus.Errorf(err.Error())
	}
	return nil
}
/*
syscall.Exec

使用Docker创建起来第一个容器之后，PID为1的进程，（容器内的第一个程序），是指定的前台程序。
容器创建之后，执行的第一个进程并不是用户的进程，而是初始化的进程

如果通过ps命令查看就会发现，容器的第一个进程变成自己的init，这和预想的不一样。
但PID为1的进程是不能被kill掉的，如果该进程被kill，容器也就推出了。

syscall.Exec这个方法，
最终调用了Kernel的int execve(const char *filename，char *const argv[], char *const envp[]);
这个系统函数。
它的作用是执行当前的filename对应的程序。
它会覆盖当前进程的镜像、数据和堆栈等信息，包括PID，这些最终都会被将要运行的进程覆盖掉。

也就是说。调用这个方法，将用户指定的进程运行起来，把最初的init进程给替换掉，这样当进入到容器内部的时候，
就会发现容器内的第一个程序就是我们指定的进程了

这也是目前Docker使用的容器引擎runC的实现方式之一。
 */