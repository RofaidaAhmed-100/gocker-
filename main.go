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
			syscall.CLONE_NEWNS |
			syscall.CLONE_NEWNET,
	}
	cmd.Start()
    networkSetup(cmd.Process.Pid) 
    cmd.Wait()
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
func networkSetup(pid int) {
	
    os.MkdirAll("/var/run/netns", 0755)
    exec.Command("ln", "-sfT",
        fmt.Sprintf("/proc/%d/ns/net", pid),
        "/var/run/netns/gocker").Run()

    exec.Command("ip", "link", "add", "name", "gocker0", "type", "bridge").Run()
    exec.Command("ip", "addr", "add", "10.0.0.1/24", "dev", "gocker0").Run()
    exec.Command("ip", "link", "set", "gocker0", "up").Run()

    exec.Command("ip", "link", "add", "veth0", "type", "veth", "peer", "name", "veth1").Run()
    exec.Command("ip", "link", "set", "veth0", "master", "gocker0").Run()
    exec.Command("ip", "link", "set", "veth0", "up").Run()
    exec.Command("ip", "link", "set", "veth1", "netns", "gocker").Run()

    exec.Command("ip", "netns", "exec", "gocker",
        "ip", "link", "set", "lo", "up").Run()
    exec.Command("ip", "netns", "exec", "gocker",
        "ip", "link", "set", "veth1", "up").Run()
    exec.Command("ip", "netns", "exec", "gocker",
        "ip", "addr", "add", "10.0.0.2/24", "dev", "veth1").Run()
    exec.Command("ip", "netns", "exec", "gocker",
        "ip", "route", "add", "default", "via", "10.0.0.1").Run()

    
    defer exec.Command("ip", "netns", "delete", "gocker").Run()
}
func networkContainerSetup() {
    exec.Command("ip", "link", "set", "lo", "up").Run()
    exec.Command("ip", "link", "set", "veth1", "up").Run()
    exec.Command("ip", "addr", "add", "10.0.0.2/24", "dev", "veth1").Run()
    exec.Command("ip", "route", "add", "default", "via", "10.0.0.1").Run()
}