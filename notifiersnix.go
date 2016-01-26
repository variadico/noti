// +build linux freebsd

package main

import "fmt"

// desktopNotify triggers a Notify notification.
func desktopNotify() {
	fmt.Println("desktopNotify")
}

// speechNotify triggers an eSpeak notification.
func speechNotify() {
	fmt.Println("speechNotify")
}
