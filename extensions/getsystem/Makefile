BOFNAME := getsystem
COMINCLUDE := -I common
LIBINCLUDE := -l secur32 -l psapi -l kernel32
CC_x64 := x86_64-w64-mingw32-gcc
CC_x86 := i686-w64-mingw32-gcc
CC=x86_64-w64-mingw32-clang

all:
	$(CC_x64) -o $(BOFNAME).o $(COMINCLUDE) -Os -c main.c -DBOF 

test:
	$(CC_x64) main.c $(COMINCLUDE) $(LIBINCLUDE) -Os -o $(BOFNAME).exe
