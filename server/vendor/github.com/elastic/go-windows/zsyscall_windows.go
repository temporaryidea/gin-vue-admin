// Licensed to Elasticsearch B.V. under one or more contributor
// license agreements. See the NOTICE file distributed with
// this work for additional information regarding copyright
// ownership. Elasticsearch B.V. licenses this file to you under
// the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing,
// software distributed under the License is distributed on an
// "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
// KIND, either express or implied.  See the License for the
// specific language governing permissions and limitations
// under the License.

// Code generated by 'go generate'; DO NOT EDIT.

//lint:file-ignore SA1019 Generated code will not be updated to use SyscallN as per https://github.com/golang/go/issues/57914.

package windows

import (
	"syscall"
	"unsafe"

	"golang.org/x/sys/windows"
)

var _ unsafe.Pointer

// Do the interface allocations only once for common
// Errno values.
const (
	errnoERROR_IO_PENDING = 997
)

var (
	errERROR_IO_PENDING error = syscall.Errno(errnoERROR_IO_PENDING)
	errERROR_EINVAL     error = syscall.EINVAL
)

// errnoErr returns common boxed Errno values, to prevent
// allocations at runtime.
func errnoErr(e syscall.Errno) error {
	switch e {
	case 0:
		return errERROR_EINVAL
	case errnoERROR_IO_PENDING:
		return errERROR_IO_PENDING
	}
	// TODO: add more here, after collecting data on the common
	// error values see on Windows. (perhaps when running
	// all.bat?)
	return e
}

var (
	modkernel32 = windows.NewLazySystemDLL("kernel32.dll")
	modntdll    = windows.NewLazySystemDLL("ntdll.dll")
	modpsapi    = windows.NewLazySystemDLL("psapi.dll")
	modversion  = windows.NewLazySystemDLL("version.dll")

	procGetNativeSystemInfo       = modkernel32.NewProc("GetNativeSystemInfo")
	procGetProcessHandleCount     = modkernel32.NewProc("GetProcessHandleCount")
	procGetSystemTimes            = modkernel32.NewProc("GetSystemTimes")
	procGetTickCount64            = modkernel32.NewProc("GetTickCount64")
	procGlobalMemoryStatusEx      = modkernel32.NewProc("GlobalMemoryStatusEx")
	procReadProcessMemory         = modkernel32.NewProc("ReadProcessMemory")
	procNtQueryInformationProcess = modntdll.NewProc("NtQueryInformationProcess")
	procEnumProcesses             = modpsapi.NewProc("EnumProcesses")
	procGetProcessImageFileNameA  = modpsapi.NewProc("GetProcessImageFileNameA")
	procGetProcessMemoryInfo      = modpsapi.NewProc("GetProcessMemoryInfo")
	procGetFileVersionInfoSizeW   = modversion.NewProc("GetFileVersionInfoSizeW")
	procGetFileVersionInfoW       = modversion.NewProc("GetFileVersionInfoW")
	procVerQueryValueW            = modversion.NewProc("VerQueryValueW")
)

func _GetNativeSystemInfo(systemInfo *SystemInfo) {
	syscall.Syscall(procGetNativeSystemInfo.Addr(), 1, uintptr(unsafe.Pointer(systemInfo)), 0, 0)
	return
}

func _GetProcessHandleCount(handle syscall.Handle, pdwHandleCount *uint32) (err error) {
	r1, _, e1 := syscall.Syscall(procGetProcessHandleCount.Addr(), 2, uintptr(handle), uintptr(unsafe.Pointer(pdwHandleCount)), 0)
	if r1 == 0 {
		err = errnoErr(e1)
	}
	return
}

func _GetSystemTimes(idleTime *syscall.Filetime, kernelTime *syscall.Filetime, userTime *syscall.Filetime) (err error) {
	r1, _, e1 := syscall.Syscall(procGetSystemTimes.Addr(), 3, uintptr(unsafe.Pointer(idleTime)), uintptr(unsafe.Pointer(kernelTime)), uintptr(unsafe.Pointer(userTime)))
	if r1 == 0 {
		err = errnoErr(e1)
	}
	return
}

func _GetTickCount64() (millis uint64, err error) {
	r0, _, e1 := syscall.Syscall(procGetTickCount64.Addr(), 0, 0, 0, 0)
	millis = uint64(r0)
	if millis == 0 {
		err = errnoErr(e1)
	}
	return
}

func _GlobalMemoryStatusEx(buffer *MemoryStatusEx) (err error) {
	r1, _, e1 := syscall.Syscall(procGlobalMemoryStatusEx.Addr(), 1, uintptr(unsafe.Pointer(buffer)), 0, 0)
	if r1 == 0 {
		err = errnoErr(e1)
	}
	return
}

func _ReadProcessMemory(handle syscall.Handle, baseAddress uintptr, buffer uintptr, size uintptr, numRead *uintptr) (err error) {
	r1, _, e1 := syscall.Syscall6(procReadProcessMemory.Addr(), 5, uintptr(handle), uintptr(baseAddress), uintptr(buffer), uintptr(size), uintptr(unsafe.Pointer(numRead)), 0)
	if r1 == 0 {
		err = errnoErr(e1)
	}
	return
}

func _NtQueryInformationProcess(handle syscall.Handle, infoClass uint32, info uintptr, infoLen uint32, returnLen *uint32) (ntStatus uint32) {
	r0, _, _ := syscall.Syscall6(procNtQueryInformationProcess.Addr(), 5, uintptr(handle), uintptr(infoClass), uintptr(info), uintptr(infoLen), uintptr(unsafe.Pointer(returnLen)), 0)
	ntStatus = uint32(r0)
	return
}

func _EnumProcesses(lpidProcess *uint32, cb uint32, lpcbNeeded *uint32) (err error) {
	r1, _, e1 := syscall.Syscall(procEnumProcesses.Addr(), 3, uintptr(unsafe.Pointer(lpidProcess)), uintptr(cb), uintptr(unsafe.Pointer(lpcbNeeded)))
	if r1 == 0 {
		err = errnoErr(e1)
	}
	return
}

func _GetProcessImageFileNameA(handle syscall.Handle, imageFileName *byte, nSize uint32) (len uint32, err error) {
	r0, _, e1 := syscall.Syscall(procGetProcessImageFileNameA.Addr(), 3, uintptr(handle), uintptr(unsafe.Pointer(imageFileName)), uintptr(nSize))
	len = uint32(r0)
	if len == 0 {
		err = errnoErr(e1)
	}
	return
}

func _GetProcessMemoryInfo(handle syscall.Handle, psmemCounters *ProcessMemoryCountersEx, cb uint32) (err error) {
	r1, _, e1 := syscall.Syscall(procGetProcessMemoryInfo.Addr(), 3, uintptr(handle), uintptr(unsafe.Pointer(psmemCounters)), uintptr(cb))
	if r1 == 0 {
		err = errnoErr(e1)
	}
	return
}

func _GetFileVersionInfoSize(filename string, handle uintptr) (size uint32, err error) {
	var _p0 *uint16
	_p0, err = syscall.UTF16PtrFromString(filename)
	if err != nil {
		return
	}
	return __GetFileVersionInfoSize(_p0, handle)
}

func __GetFileVersionInfoSize(filename *uint16, handle uintptr) (size uint32, err error) {
	r0, _, e1 := syscall.Syscall(procGetFileVersionInfoSizeW.Addr(), 2, uintptr(unsafe.Pointer(filename)), uintptr(handle), 0)
	size = uint32(r0)
	if size == 0 {
		err = errnoErr(e1)
	}
	return
}

func _GetFileVersionInfo(filename string, reserved uint32, dataLen uint32, data *byte) (success bool, err error) {
	var _p0 *uint16
	_p0, err = syscall.UTF16PtrFromString(filename)
	if err != nil {
		return
	}
	return __GetFileVersionInfo(_p0, reserved, dataLen, data)
}

func __GetFileVersionInfo(filename *uint16, reserved uint32, dataLen uint32, data *byte) (success bool, err error) {
	r0, _, e1 := syscall.Syscall6(procGetFileVersionInfoW.Addr(), 4, uintptr(unsafe.Pointer(filename)), uintptr(reserved), uintptr(dataLen), uintptr(unsafe.Pointer(data)), 0, 0)
	success = r0 != 0
	if !success {
		err = errnoErr(e1)
	}
	return
}

func _VerQueryValueW(data *byte, subBlock string, pBuffer *uintptr, len *uint32) (success bool, err error) {
	var _p0 *uint16
	_p0, err = syscall.UTF16PtrFromString(subBlock)
	if err != nil {
		return
	}
	return __VerQueryValueW(data, _p0, pBuffer, len)
}

func __VerQueryValueW(data *byte, subBlock *uint16, pBuffer *uintptr, len *uint32) (success bool, err error) {
	r0, _, e1 := syscall.Syscall6(procVerQueryValueW.Addr(), 4, uintptr(unsafe.Pointer(data)), uintptr(unsafe.Pointer(subBlock)), uintptr(unsafe.Pointer(pBuffer)), uintptr(unsafe.Pointer(len)), 0, 0)
	success = r0 != 0
	if !success {
		err = errnoErr(e1)
	}
	return
}
