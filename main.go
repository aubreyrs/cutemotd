package main

import (
	"fmt"
	"os"
	"os/user"
	"strings"
	"time"

	"github.com/shirou/gopsutil/v3/host"
)

const (
	reset    = "\033[0m"
	mauve    = "\033[38;2;203;166;247m"
	pink     = "\033[38;2;245;194;231m"
	lavender = "\033[38;2;180;190;254m"
)

const (
	transBlue  = "\033[38;2;91;206;250m"
	transPink  = "\033[38;2;245;169;184m"
	transWhite = "\033[38;2;255;255;255m"
)

const catAscii = `
%s  ,   ,       %s
%s { \w/ }     %s %s★ %s ★%s
%s  '>!<'      %s %s☆ %s ☆%s
%s  (/^\)      %s %s★ %s ★%s
%s  '   '      %s

`

type motdInfo struct {
	username string
	ip       string
	uptime   string
}

func getUsername() string {
	if len(os.Args) > 1 && os.Args[1] == "test" {
		return "kitty"
	}

	currentUser, err := user.Current()
	if err != nil {
		return "unknown"
	}
	return currentUser.Username
}

func getConnectingIP() string {
	if len(os.Args) > 1 && os.Args[1] == "test" {
		return "127.0.0.1"
	}

	sshClient := os.Getenv("SSH_CLIENT")
	if sshClient == "" {
		return "local session"
	}
	return strings.Split(sshClient, " ")[0]
}

func formatUptime(uptimeSeconds uint64) string {
	duration := time.Duration(uptimeSeconds) * time.Second
	days := int(duration.Hours() / 24)
	hours := int(duration.Hours()) % 24
	minutes := int(duration.Minutes()) % 60

	if days > 0 {
		return fmt.Sprintf("%dd %dh %dm", days, hours, minutes)
	} else if hours > 0 {
		return fmt.Sprintf("%dh %dm", hours, minutes)
	}
	return fmt.Sprintf("%dm", minutes)
}

func getUptime() string {
	if len(os.Args) > 1 && os.Args[1] == "test" {
		return "1d 2h 3m"
	}

	info, err := host.Info()
	if err != nil {
		return "unknown"
	}
	return formatUptime(info.Uptime)
}

func formatMOTD(info motdInfo) string {
	welcomeMsg := fmt.Sprintf("welcome back %s :3", info.username)
	ipMsg := fmt.Sprintf("you're connecting from %s :o", info.ip)
	uptimeMsg := fmt.Sprintf("uptime: %s", info.uptime)

	return fmt.Sprintf(catAscii,
		transBlue, reset,
		transPink, reset, mauve, welcomeMsg, reset,
		transWhite, reset, pink, ipMsg, reset,
		transPink, reset, lavender, uptimeMsg, reset,
		transBlue, reset,
	)
}

func main() {
	info := motdInfo{
		username: getUsername(),
		ip:       getConnectingIP(),
		uptime:   getUptime(),
	}

	fmt.Print(formatMOTD(info))
}
