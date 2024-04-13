package helper

import (
	"encoding/binary"
	"fmt"
	"log"
	"unsafe"

	"golang.org/x/sys/windows"
)

func FindDMAAddy(h windows.Handle, baseAddr uintptr, offsets ...uintptr) (uintptr, error) {
	var addr = baseAddr
	var buf [8]byte
	for _, offset := range offsets {
		if err := windows.ReadProcessMemory(h, addr, &buf[0], uintptr(unsafe.Sizeof(addr)), nil); err != nil {
			return 0, err
		}
		addr = uintptr(binary.LittleEndian.Uint64(buf[:]))
		if addr == 0 {
			return 0, fmt.Errorf("Couldn't resolve address")
		}
		addr += offset
	}
	return addr, nil
}

func GetProcessIDByProcessName(processName string) (uint32, error) {
	handle, err := windows.CreateToolhelp32Snapshot(windows.TH32CS_SNAPPROCESS, 0)
	if err != nil {
		return 0, err
	}
	var pe32 windows.ProcessEntry32
	pe32.Size = uint32(unsafe.Sizeof(pe32))
	_ = windows.Process32First(handle, &pe32)
	for {
		if err := windows.Process32Next(handle, &pe32); err != nil {
			return 0, err
		}
		if windows.UTF16ToString(pe32.ExeFile[:]) == processName {
			break
		}
	}
	return pe32.ProcessID, nil
}

func GetModuleEntryByProcID(procID uint32) (uintptr, error) {
	handle, err := windows.CreateToolhelp32Snapshot(windows.TH32CS_SNAPMODULE, procID)
	if err != nil {
		return 0, err
	}
	var me32 windows.ModuleEntry32
	me32.Size = uint32(unsafe.Sizeof(me32))
	err = windows.Module32First(handle, &me32)
	if err != nil {
		log.Fatal(err)
	}
	return me32.ModBaseAddr, nil
}
