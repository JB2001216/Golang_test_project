/*
Copyright 2014 Tamás Gulácsi.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

     http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package punchhole

import (
	"os"
	"syscall"
	"unsafe"
)

const (
	//http://source.winehq.org/source/include/winnt.h#L4605
	file_read_data  = 1
	file_write_data = 2

	// METHOD_BUFFERED	0
	method_buffered = 0
	// FILE_ANY_ACCESS   0
	file_any_access = 0
	// FILE_DEVICE_FILE_SYSTEM   0x00000009
	file_device_file_system = 0x00000009
	// FILE_SPECIAL_ACCESS   (FILE_ANY_ACCESS)
	file_special_access = file_any_access
	file_read_access    = file_read_data
	file_write_access   = file_write_data

	// http://source.winehq.org/source/include/winioctl.h
	// #define CTL_CODE 	(  	DeviceType,
	//		Function,
	//		Method,
	//		Access  		 )
	//    ((DeviceType) << 16) | ((Access) << 14) | ((Function) << 2) | (Method)

	// FSCTL_SET_COMPRESSION   CTL_CODE(FILE_DEVICE_FILE_SYSTEM, 16, METHOD_BUFFERED, FILE_READ_DATA | FILE_WRITE_DATA)
	fsctl_set_compression = (file_device_file_system << 16) | ((file_read_access | file_write_access) << 14) | (16 << 2) | method_buffered
	// FSCTL_SET_SPARSE   CTL_CODE(FILE_DEVICE_FILE_SYSTEM, 49, METHOD_BUFFERED, FILE_SPECIAL_ACCESS)
	fsctl_set_sparse = (file_device_file_system << 16) | (file_special_access << 14) | (49 << 2) | method_buffered
	// FSCTL_SET_ZERO_DATA   CTL_CODE(FILE_DEVICE_FILE_SYSTEM, 50, METHOD_BUFFERED, FILE_WRITE_DATA)
	fsctl_set_zero_data = (file_device_file_system << 16) | (file_write_data << 14) | (50 << 2) | method_buffered
)

var (
	modkernel32 = syscall.NewLazyDLL("kernel32.dll")

	procDeviceIOControl = modkernel32.NewProc("DeviceIoControl")
)

func init() {
	PunchHole = punchHoleWindows
}

// http://msdn.microsoft.com/en-us/library/windows/desktop/aa364411%28v=vs.85%29.aspx
// typedef struct _FILE_ZERO_DATA_INFORMATION {
//  LARGE_INTEGER FileOffset;
//  LARGE_INTEGER BeyondFinalZero;
//} FILE_ZERO_DATA_INFORMATION, *PFILE_ZERO_DATA_INFORMATION;
type fileZeroDataInformation struct {
	FileOffset, BeyondFinalZero int64
}

// puncHoleWindows punches a hole into the given file starting at offset,
// measuring "size" bytes
// (http://msdn.microsoft.com/en-us/library/windows/desktop/aa364597%28v=vs.85%29.aspx)
func punchHoleWindows(file *os.File, offset, size int64) (err error) {
	lpInBuffer := fileZeroDataInformation{
		FileOffset:      offset,
		BeyondFinalZero: offset + size}
	lpBytesReturned := make([]byte, 8)
	// BOOL
	// WINAPI
	// DeviceIoControl( (HANDLE) hDevice,              // handle to a file
	//                  FSCTL_SET_ZERO_DATA,           // dwIoControlCode
	//                  (LPVOID) lpInBuffer,           // input buffer
	//                  (DWORD) nInBufferSize,         // size of input buffer
	//                  NULL,                          // lpOutBuffer
	//                  0,                             // nOutBufferSize
	//                  (LPDWORD) lpBytesReturned,     // number of bytes returned
	//                  (LPOVERLAPPED) lpOverlapped ); // OVERLAPPED structure
	r1, _, e1 := syscall.Syscall9(procDeviceIOControl.Addr(), 8,
		file.Fd(),
		uintptr(fsctl_set_zero_data),
		uintptr(unsafe.Pointer(&lpInBuffer)),
		8,
		0,
		0,
		uintptr(unsafe.Pointer(&lpBytesReturned[0])),
		0,
		0)
	if r1 == 0 {
		if e1 != 0 {
			err = error(e1)
		} else {
			err = syscall.EINVAL
		}
	}
	return
}
