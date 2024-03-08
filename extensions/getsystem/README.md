# getsystem-bof

This is a BOF that can be used within Sliver (and probably Cobaltstrike) to execute shellcode as SYSTEM when you are administrator or have the SeDebugPrivilege.

```
getsystem-bof notepad.exe PATH_TO_SHELLCODE
```

This is a well-known technique to get SYSTEM privileges via PPID spoofing. The BOF uses the boku7 shellcode injector and aquires the SeDebugPrivilege to get access to process handles of SYSTEM processes. It finds the PID of a SYSTEM process (here winlogon.exe) to use as parent of a new sacrificial process and injects the supplied shellocde into it. As the child process of a SYSTEM process inherits the parent's privileges, the shellcode will be running as SYSTEM as well.

Credits:
* https://github.com/boku7/spawn
* https://gist.github.com/G0ldenGunSec/8ca0e853dd5637af2881697f8de6aecc
* https://github.com/trustedsec/CS-Situational-Awareness-BOF
* https://github.com/py7hagoras/GetSystem 