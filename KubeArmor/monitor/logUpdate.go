// SPDX-License-Identifier: Apache-2.0
// Copyright 2021 Authors of KubeArmor

package monitor

import (
	"fmt"
	"strconv"
	"strings"

	kl "github.com/kubearmor/KubeArmor/KubeArmor/common"
	tp "github.com/kubearmor/KubeArmor/KubeArmor/types"
)

// ========== //
// == Logs == //
// ========== //

// UpdateContainerInfoByContainerID Function
func (mon *SystemMonitor) UpdateContainerInfoByContainerID(log tp.Log) tp.Log {
	Containers := *(mon.Containers)
	ContainersLock := *(mon.ContainersLock)

	ContainersLock.RLock()
	defer ContainersLock.RUnlock()

	if val, ok := Containers[log.ContainerID]; ok {
		// update pod info
		log.NamespaceName = val.NamespaceName
		log.PodName = val.EndPointName
		log.Labels = val.Labels

		// update container info
		log.ContainerName = val.ContainerName
		log.ContainerImage = val.ContainerImage

		// get merged directory
		log.MergedDir = val.MergedDir

		// update policy flag
		log.PolicyEnabled = val.PolicyEnabled

		// update visibility flags
		log.ProcessVisibilityEnabled = val.ProcessVisibilityEnabled
		log.FileVisibilityEnabled = val.FileVisibilityEnabled
		log.NetworkVisibilityEnabled = val.NetworkVisibilityEnabled
		log.CapabilitiesVisibilityEnabled = val.CapabilitiesVisibilityEnabled
	}

	return log
}

// BuildLogBase Function
func (mon *SystemMonitor) BuildLogBase(eventID int32, msg ContextCombined) tp.Log {
	log := tp.Log{}

	timestamp, updatedTime := kl.GetDateTimeNow()

	log.Timestamp = timestamp
	log.UpdatedTime = updatedTime
	log.ContainerID = msg.ContainerID

	if log.ContainerID != "" {
		log = mon.UpdateContainerInfoByContainerID(log)
	} else {
		// update host policy flag
		log.PolicyEnabled = mon.Node.PolicyEnabled

		// update host visibility flags
		log.ProcessVisibilityEnabled = mon.Node.ProcessVisibilityEnabled
		log.FileVisibilityEnabled = mon.Node.FileVisibilityEnabled
		log.NetworkVisibilityEnabled = mon.Node.NetworkVisibilityEnabled
		log.CapabilitiesVisibilityEnabled = mon.Node.CapabilitiesVisibilityEnabled
	}

	log.HostPPID = int32(msg.ContextSys.HostPPID)
	log.HostPID = int32(msg.ContextSys.HostPID)

	log.PPID = int32(msg.ContextSys.PPID)
	log.PID = int32(msg.ContextSys.PID)
	log.UID = int32(msg.ContextSys.UID)

	if msg.ContextSys.EventID == SysExecve || msg.ContextSys.EventID == SysExecveAt {
		log.Source = mon.GetParentExecPath(msg.ContainerID, msg.ContextSys.HostPID)
	} else {
		log.Source = mon.GetCommand(msg.ContainerID, msg.ContextSys.HostPID)
	}

	log.ParentProcessName = mon.GetExecPath(msg.ContainerID, msg.ContextSys.HostPPID)
	log.ProcessName = mon.GetExecPath(msg.ContainerID, msg.ContextSys.HostPID)

	return log
}

// UpdateLogBase Function (SYS_EXECVE, SYS_EXECVEAT)
func (mon *SystemMonitor) UpdateLogBase(eventID int32, log tp.Log) tp.Log {
	if log.ParentProcessName == "" || !strings.HasPrefix(log.ParentProcessName, "/") {
		parentProcessName := mon.GetParentExecPath(log.ContainerID, uint32(log.HostPID))
		if parentProcessName != "" {
			log.ParentProcessName = parentProcessName
		}
	}

	if log.ProcessName == "" || !strings.HasPrefix(log.ProcessName, "/") {
		processName := mon.GetExecPath(log.ContainerID, uint32(log.HostPID))
		if processName != "" {
			log.ProcessName = processName
		}
	}

	if log.Source == "" || !strings.HasPrefix(log.Source, "/") {
		source := mon.GetExecPath(log.ContainerID, uint32(log.HostPPID))
		if source != "" {
			log.Source = source
		}
	}

	return log
}

// UpdateLogs Function
func (mon *SystemMonitor) UpdateLogs() {
	for {
		select {
		case <-StopChan:
			return

		case msg, valid := <-mon.ContextChan:
			if !valid {
				continue
			}

			// generate a log
			log := mon.BuildLogBase(msg.ContextSys.EventID, msg)

			switch msg.ContextSys.EventID {
			case SysOpen:
				if len(msg.ContextArgs) != 2 {
					continue
				}

				var fileName string
				var fileOpenFlags string

				if val, ok := msg.ContextArgs[0].(string); ok {
					fileName = val
				}
				if val, ok := msg.ContextArgs[1].(string); ok {
					fileOpenFlags = val
				}

				log.Operation = "File"
				log.Resource = fileName
				log.Data = "syscall=" + getSyscallName(int32(msg.ContextSys.EventID)) + " flags=" + fileOpenFlags

			case SysOpenAt:
				if len(msg.ContextArgs) != 3 {
					continue
				}

				var fd string
				var fileName string
				var fileOpenFlags string

				if val, ok := msg.ContextArgs[0].(int32); ok {
					fd = strconv.Itoa(int(val))
				}
				if val, ok := msg.ContextArgs[1].(string); ok {
					fileName = val
				}
				if val, ok := msg.ContextArgs[2].(string); ok {
					fileOpenFlags = val
				}

				log.Operation = "File"
				log.Resource = fileName
				log.Data = "syscall=" + getSyscallName(int32(msg.ContextSys.EventID)) + " fd=" + fd + " flags=" + fileOpenFlags

			case SysUnlink:
				if len(msg.ContextArgs) != 2 {
					continue
				}

				var fileName string
				if val, ok := msg.ContextArgs[1].(string); ok {
					fileName = val
				}

				log.Operation = "File"
				log.Resource = fileName
				log.Data = "syscall=" + getSyscallName(int32(msg.ContextSys.EventID))

			case SysUnlinkAt:
				if len(msg.ContextArgs) != 3 {
					continue
				}

				var fileName string
				var fileUnlinkAtFlags string

				if val, ok := msg.ContextArgs[1].(string); ok {
					fileName = val
				}
				if val, ok := msg.ContextArgs[2].(string); ok {
					fileUnlinkAtFlags = val
				}

				log.Operation = "File"
				log.Resource = fileName
				log.Data = "syscall=" + getSyscallName(int32(msg.ContextSys.EventID)) + " flags=" + fileUnlinkAtFlags

			case SysRmdir:
				if len(msg.ContextArgs) != 1 {
					continue
				}

				var fileName string
				if val, ok := msg.ContextArgs[0].(string); ok {
					fileName = val
				}

				log.Operation = "File"
				log.Resource = fileName
				log.Data = "syscall=" + getSyscallName(int32(msg.ContextSys.EventID))

			case SysClose:
				if len(msg.ContextArgs) != 1 {
					continue
				}

				var fd string

				if val, ok := msg.ContextArgs[0].(int32); ok {
					fd = strconv.Itoa(int(val))
				}

				log.Operation = "File"
				log.Resource = ""
				log.Data = "syscall=" + getSyscallName(int32(msg.ContextSys.EventID)) + " fd=" + fd

			case SysSocket: // domain, type, proto
				if len(msg.ContextArgs) != 3 {
					continue
				}

				var sockDomain string
				var sockType string
				var sockProtocol int32

				if val, ok := msg.ContextArgs[0].(string); ok {
					sockDomain = val
				}
				if val, ok := msg.ContextArgs[1].(string); ok {
					sockType = val
				}
				if val, ok := msg.ContextArgs[2].(int32); ok {
					sockProtocol = val
				}

				log.Operation = "Network"
				log.Resource = "domain=" + sockDomain + " type=" + sockType + " protocol=" + getProtocol(sockProtocol)
				log.Data = "syscall=" + getSyscallName(int32(msg.ContextSys.EventID))

			case TCPConnect, TCPConnectv6, TCPAccept, TCPAcceptv6:
				if len(msg.ContextArgs) != 2 {
					continue
				}
				var sockAddr map[string]string
				var protocol string
				if val, ok := msg.ContextArgs[0].(string); ok {
					protocol = val
				}

				if val, ok := msg.ContextArgs[1].(map[string]string); ok {
					sockAddr = val
				}

				log.Operation = "Network"
				log.Resource = "remoteip=" + sockAddr["sin_addr"] + " port=" + sockAddr["sin_port"] + " protocol=" + protocol
				if msg.ContextSys.EventID == TCPConnect || msg.ContextSys.EventID == TCPConnectv6 {
					log.Data = "kprobe=tcp_connect"
				} else {
					log.Data = "kprobe=tcp_accept"
				}
				log.Data = log.Data + " domain=" + sockAddr["sa_family"]

			case SysConnect: // fd, sockaddr
				if len(msg.ContextArgs) != 2 {
					continue
				}

				var fd string
				var sockAddr map[string]string

				if val, ok := msg.ContextArgs[0].(int32); ok {
					fd = strconv.Itoa(int(val))
				}
				if val, ok := msg.ContextArgs[1].(map[string]string); ok {
					sockAddr = val
				}

				log.Operation = "Network"
				log.Resource = ""

				for k, v := range sockAddr {
					if log.Resource == "" {
						log.Resource = k + "=" + v
					} else {
						log.Resource = log.Resource + " " + k + "=" + v
					}
				}

				log.Data = "syscall=" + getSyscallName(int32(msg.ContextSys.EventID)) + " fd=" + fd

			case SysAccept: // fd, sockaddr
				if len(msg.ContextArgs) != 2 {
					continue
				}

				var fd string
				var sockAddr map[string]string

				if val, ok := msg.ContextArgs[0].(int32); ok {
					fd = strconv.Itoa(int(val))
				}
				if val, ok := msg.ContextArgs[1].(map[string]string); ok {
					sockAddr = val
				}

				log.Operation = "Network"
				log.Resource = ""
				log.Data = "syscall=" + getSyscallName(int32(msg.ContextSys.EventID)) + " fd=" + fd

				for k, v := range sockAddr {
					if log.Resource == "" {
						log.Resource = k + "=" + v
					} else {
						log.Resource = log.Resource + " " + k + "=" + v
					}
				}

			case SysBind: // fd, sockaddr
				if len(msg.ContextArgs) != 2 {
					continue
				}

				var fd string
				var sockAddr map[string]string

				if val, ok := msg.ContextArgs[0].(int32); ok {
					fd = strconv.Itoa(int(val))
				}
				if val, ok := msg.ContextArgs[1].(map[string]string); ok {
					sockAddr = val
				}

				log.Operation = "Network"
				log.Resource = ""

				for k, v := range sockAddr {
					if log.Resource == "" {
						log.Resource = k + "=" + v
					} else {
						log.Resource = log.Resource + " " + k + "=" + v
					}
				}

				log.Data = "syscall=" + getSyscallName(int32(msg.ContextSys.EventID)) + " fd=" + fd

			case SysListen: // fd
				if len(msg.ContextArgs) != 2 {
					continue
				}

				var fd string

				if val, ok := msg.ContextArgs[0].(int32); ok {
					fd = strconv.Itoa(int(val))
				}

				log.Operation = "Network"
				log.Resource = ""
				log.Data = "syscall=" + getSyscallName(int32(msg.ContextSys.EventID)) + " fd=" + fd

			default:
				continue
			}

			// get error message
			if msg.ContextSys.Retval < 0 {
				message := getErrorMessage(msg.ContextSys.Retval)
				if message != "" {
					log.Result = message
				} else {
					log.Result = fmt.Sprintf("Unknown (%d)", msg.ContextSys.Retval)
				}
			} else {
				log.Result = "Passed"
			}

			// push the generated log
			if mon.Logger != nil {
				go mon.Logger.PushLog(log)
			}
		}
	}
}
