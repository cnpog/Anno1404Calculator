package main

import (
	"anno1404/pkg/helper"
	"encoding/binary"
	"fmt"
	"log"

	"golang.org/x/sys/windows"
)

func main() {
	var base uintptr = 0x00ED2BFC
	var beggars uintptr = 0x980
	// var peasants uintptr = 0x9A0
	// var citizens uintptr = 0x9C0
	// var patricians uintptr = 0x9E0
	// var noblemen uintptr = 0xA00
	// var nomads uintptr = 0x920
	// var envoys uintptr = 0x940
	procID, err := helper.GetProcessIDByProcessName("Addon.exe")
	if err != nil {
		log.Fatal(err)
	}
	modBaseAddr, err := helper.GetModuleEntryByProcID(procID)
	if err != nil {
		log.Fatal(err)
	}
	handle, err := windows.OpenProcess(windows.PROCESS_VM_READ, false, procID)
	if err != nil {
		log.Fatal(err)
	}

	guidAddr, err := helper.FindDMAAddy(handle, modBaseAddr+base, beggars)
	if err != nil {
		log.Fatal(err)
	}
	var buf [4]byte
	if err := windows.ReadProcessMemory(handle, guidAddr, &buf[0], 4, nil); err != nil {
		log.Fatal(err)
	}
	fmt.Println(binary.LittleEndian.Uint32(buf[:]))

}
