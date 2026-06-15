package main

import (
	"fmt"
	"os"
	"os/exec"
	"syscall"
	"strconv"
)

func main() {
	switch os.Args[1] {
	case "run":
		run()
	case "child":
		child()
	default:
		fmt.Println("usage: gocker run <command>")
	}
}

func run() {
	cmd := exec.Command("/proc/self/exe", append([]string{"child"}, os.Args[2:]...)...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Cloneflags: syscall.CLONE_NEWUTS |
			syscall.CLONE_NEWPID |
			syscall.CLONE_NEWNS,
	}
	if err := cmd.Run(); err != nil {
		fmt.Println("error:", err)
	}
}

func child() {
     cgroups()
	 syscall.Sethostname([]byte("container"))
    

    syscall.Chroot("/tmp/gocker/alpine")
    os.Chdir("/")
    
    syscall.Mount("proc", "/proc", "proc", 0, "")

    cmd := exec.Command(os.Args[2], os.Args[3:]...)
    cmd.Stdin = os.Stdin
    cmd.Stdout = os.Stdout
    cmd.Stderr = os.Stderr
    cmd.Run()

    syscall.Unmount("/proc", 0)
}
func cgroups() {
    cgPath := "/sys/fs/cgroup/gocker"

    
    os.MkdirAll(cgPath, 0755)

    
    os.WriteFile(cgPath+"/memory.max", []byte("536870912"), 0700)

    
    os.WriteFile(cgPath+"/cpu.weight", []byte("50"), 0700)

    
    pid := strconv.Itoa(os.Getpid())
    os.WriteFile(cgPath+"/cgroup.procs", []byte(pid), 0700)
}