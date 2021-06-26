package main

import (
	"fmt"
	"syscall"
	"unsafe"
)

func main() {
	var mod = syscall.NewLazyDLL("user32.dll")
	var proc = mod.NewProc("MessageBoxW")
	var MB_YESNOCANCEL = 0x00000003

	p1, _ := syscall.UTF16PtrFromString("This test is Done.")
	p2, _ := syscall.UTF16PtrFromString("Done Title")

	ret, _, _ := proc.Call(0,
		uintptr(unsafe.Pointer(p1)),
		uintptr(unsafe.Pointer(p2)),
		uintptr(MB_YESNOCANCEL))
	fmt.Printf("Return: %d\n", ret)

}
