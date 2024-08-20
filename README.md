# Sliver Dispatcher
> By @nigerald and @bri5ee
Gotta pwn all the blue teamers in RvB. Mass dispatch commands to everyone or select users. Only supports sessions right now.

# Usage
Start the client
```
./sliver-dispatch -config /path/to/sliver_config
```

Yolo and tab autocomplete

# Client Commands
| Command | Notes |
|---------|-------|
| debug | enables debug mode |
| get_configs | get agent builds and profiles |
| get_sessions | gets active sessions and selected sessions |
| run_all_linux | run a capability on all linux sessions |
| run_all_selected | run a capability on all selected sessions. Capabilities that are incompatible with an OS skip the selected session |
| run_all_windows | run a capability on all windows sessions |
| select_sessions | select sessions |
| help | get more detailed help info on any of these commands |

# Dispatch

| OS      | Capability | Syntax | Notes |
|---------|------------|-----------|-------|
| Windows | getsystem | getsystem [/path/to/shc] [process_to_spawn] | spawns a process as SYSTEM and injects shellcode into it via the getsystem bof |
| Windows | noseferatu | nosfereatu [/path/to/nosferatu.bin]| inject the nosferatu dll (prepended with srdi) via sliver's built-in injection technique |
| Windows | shinject | shinject [/path/to/shellcode] [process name]| inject shellcode to the target process using syscalls_shinject bof |
| Either | execute | execute [executable] [args] | run an executable via starting it as a process |
| Either | script | script [/path/to/script] | runs a powershell or bash script |
| Either | upload | upload [source_path] [dest_path]| upload a file |