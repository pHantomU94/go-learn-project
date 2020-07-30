package main

import (
	"log"
	"os"
	"os/exec"
	"syscall"
)

// NewNamespace 创建新的命名空间隔离的sh
func NewNamespace() {
	cmd := exec.Command("bash")
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Cloneflags: syscall.CLONE_NEWUTS |
			syscall.CLONE_NEWIPC |
			syscall.CLONE_NEWPID |
			syscall.CLONE_NEWUSER |
			syscall.CLONE_NEWNS |
			syscall.CLONE_NEWNET,
		UidMappings: []syscall.SysProcIDMap{
			{
				// 容器的UID
				ContainerID: 1,
				// 宿主机的UID
				HostID: 0,
				Size: 1,
			},
		},
		GidMappings: []syscall.SysProcIDMap{
			{
				ContainerID: 1,
				HostID: 0,
				Size: 1,
			},
		},
	}
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		log.Fatal(err)
	}
}

func main() {
	NewNamespace()
}