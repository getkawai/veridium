package main

import (
	"github.com/kawai-network/veridium/pkg/nodeos"
)

// CpuInfo represents CPU information (frontend compatible)
type CpuInfo struct {
	Model string   `json:"model"`
	Speed int      `json:"speed"`
	Times CpuTimes `json:"times"`
}

// CpuTimes represents CPU timing information (frontend compatible)
type CpuTimes struct {
	User uint64 `json:"user"`
	Nice uint64 `json:"nice"`
	Sys  uint64 `json:"sys"`
	Idle uint64 `json:"idle"`
	IRQ  uint64 `json:"irq"`
}

// UserInfoResult represents user information (frontend compatible)
type UserInfoResult struct {
	Username string `json:"username"`
	Uid      int    `json:"uid"`
	Gid      int    `json:"gid"`
	Shell    string `json:"shell"`
	Homedir  string `json:"homedir"`
}

// NetworkInterface represents network interface information (frontend compatible)
type NetworkInterface struct {
	Address  string `json:"address"`
	Netmask  string `json:"netmask"`
	Family   string `json:"family"`
	Mac      string `json:"mac"`
	ScopeID  int    `json:"scopeid,omitempty"`
	Internal bool   `json:"internal"`
}

// NodeOsService provides Node.js os module equivalents as a Wails service
type NodeOsService struct{}

// Arch returns the operating system CPU architecture
func (n *NodeOsService) Arch() string {
	return nodeos.Arch()
}

// Platform returns the operating system platform
func (n *NodeOsService) Platform() string {
	return nodeos.Platform()
}

// Type returns the operating system name
func (n *NodeOsService) Type() string {
	return nodeos.Type()
}

// Release returns the operating system release
func (n *NodeOsService) Release() string {
	return nodeos.Release()
}

// Version returns the operating system version
func (n *NodeOsService) Version() string {
	return nodeos.Version()
}

// Hostname returns the hostname of the operating system
func (n *NodeOsService) Hostname() string {
	return nodeos.Hostname()
}

// Homedir returns the home directory of the current user
func (n *NodeOsService) Homedir() string {
	return nodeos.Homedir()
}

// Tmpdir returns the default directory for temporary files
func (n *NodeOsService) Tmpdir() string {
	return nodeos.Tmpdir()
}

// Endianness returns the endianness of the CPU
func (n *NodeOsService) Endianness() string {
	return nodeos.Endianness()
}

// Cpus returns information about each CPU/core installed
func (n *NodeOsService) Cpus() []CpuInfo {
	cpus := nodeos.Cpus()
	result := make([]CpuInfo, len(cpus))

	for i, cpu := range cpus {
		result[i] = CpuInfo{
			Model: cpu.Model,
			Speed: cpu.Speed,
			Times: CpuTimes{
				User: cpu.Times.User,
				Nice: cpu.Times.Nice,
				Sys:  cpu.Times.Sys,
				Idle: cpu.Times.Idle,
				IRQ:  cpu.Times.IRQ,
			},
		}
	}

	return result
}

// Totalmem returns the total amount of system memory in bytes
func (n *NodeOsService) Totalmem() uint64 {
	return nodeos.Totalmem()
}

// Freemem returns the amount of free system memory in bytes
func (n *NodeOsService) Freemem() uint64 {
	return nodeos.Freemem()
}

// Loadavg returns an array containing the 1, 5, and 15 minute load averages
func (n *NodeOsService) Loadavg() []float64 {
	return nodeos.Loadavg()
}

// Uptime returns the system uptime in seconds
func (n *NodeOsService) Uptime() uint64 {
	return nodeos.Uptime()
}

// NetworkInterfaces returns an object containing network interfaces
func (n *NodeOsService) NetworkInterfaces() map[string][]NetworkInterface {
	interfaces := nodeos.NetworkInterfaces()
	result := make(map[string][]NetworkInterface)

	for name, ifaces := range interfaces {
		resultIfaces := make([]NetworkInterface, len(ifaces))
		for i, iface := range ifaces {
			resultIfaces[i] = NetworkInterface{
				Address:  iface.Address,
				Netmask:  iface.Netmask,
				Family:   iface.Family,
				Mac:      iface.Mac,
				ScopeID:  iface.ScopeID,
				Internal: iface.Internal,
			}
		}
		result[name] = resultIfaces
	}

	return result
}

// UserInfo returns information about the currently effective user
func (n *NodeOsService) UserInfo() (UserInfoResult, error) {
	info, err := nodeos.UserInfo()
	if err != nil {
		return UserInfoResult{}, err
	}

	return UserInfoResult{
		Username: info.Username,
		Uid:      info.Uid,
		Gid:      info.Gid,
		Shell:    info.Shell,
		Homedir:  info.Homedir,
	}, nil
}

// GetPriority gets the scheduling priority of a process
func (n *NodeOsService) GetPriority(pid ...int) int {
	if len(pid) > 0 {
		return nodeos.GetPriority(pid[0])
	}
	return nodeos.GetPriority()
}

// SetPriority sets the scheduling priority of a process
func (n *NodeOsService) SetPriority(priority int, pid ...int) error {
	if len(pid) > 0 {
		return nodeos.SetPriority(priority, pid[0])
	}
	return nodeos.SetPriority(priority)
}

// AvailableParallelism returns an estimate of the default amount of parallelism
func (n *NodeOsService) AvailableParallelism() int {
	return nodeos.AvailableParallelism()
}

// Machine returns the machine type (hardware platform)
func (n *NodeOsService) Machine() string {
	return nodeos.Machine()
}

// VersionString returns a string identifying the operating system version
func (n *NodeOsService) VersionString() string {
	return nodeos.VersionString()
}

// Get constants
func (n *NodeOsService) EOL() string {
	return nodeos.EOL
}

func (n *NodeOsService) DevNull() string {
	return nodeos.DevNull
}

// Priority constants
func (n *NodeOsService) PRIORITY_LOWEST() int {
	return nodeos.PRIORITY_LOWEST
}

func (n *NodeOsService) PRIORITY_LOW() int {
	return nodeos.PRIORITY_LOW
}

func (n *NodeOsService) PRIORITY_NORMAL() int {
	return nodeos.PRIORITY_NORMAL
}

func (n *NodeOsService) PRIORITY_HIGH() int {
	return nodeos.PRIORITY_HIGH
}

func (n *NodeOsService) PRIORITY_HIGHEST() int {
	return nodeos.PRIORITY_HIGHEST
}

// Signal constants
func (n *NodeOsService) SIGABRT() string {
	return nodeos.SIGABRT
}

func (n *NodeOsService) SIGALRM() string {
	return nodeos.SIGALRM
}

func (n *NodeOsService) SIGBUS() string {
	return nodeos.SIGBUS
}

func (n *NodeOsService) SIGFPE() string {
	return nodeos.SIGFPE
}

func (n *NodeOsService) SIGHUP() string {
	return nodeos.SIGHUP
}

func (n *NodeOsService) SIGILL() string {
	return nodeos.SIGILL
}

func (n *NodeOsService) SIGINT() string {
	return nodeos.SIGINT
}

func (n *NodeOsService) SIGKILL() string {
	return nodeos.SIGKILL
}

func (n *NodeOsService) SIGPIPE() string {
	return nodeos.SIGPIPE
}

func (n *NodeOsService) SIGQUIT() string {
	return nodeos.SIGQUIT
}

func (n *NodeOsService) SIGSEGV() string {
	return nodeos.SIGSEGV
}

func (n *NodeOsService) SIGTERM() string {
	return nodeos.SIGTERM
}

func (n *NodeOsService) SIGUSR1() string {
	return nodeos.SIGUSR1
}

func (n *NodeOsService) SIGUSR2() string {
	return nodeos.SIGUSR2
}
