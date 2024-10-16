# winsvr_adm
Windows Server administration - get system statuses

Compiling : go build -o bin/winsvr_adm

running : ./bin/winsvr_adm

open browser : http://serverip:8080

APIs :
 1.   http://serverip:8080/status
 2.   http://serverip:8080/cpu
 3.   http://serverip:8080/storage
 4.   http://serverip:8080/net
 5.   http://serverip:8080/process

Future APIs :
1. list process with percentage
2. kill, restart process
3. server reboot
