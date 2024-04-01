@echo off
setlocal

rem Define array elements
set "PROTO_NAMES=auth conversation errinfo friend group msg msggateway push rtc sdkws third user statistics"

rem Loop through each element in the array
for %%i in (%PROTO_NAMES%) do (
  protoc --go_out=plugins=grpc:./%%i --go_opt=module=github.com/openimsdk/protocol/%%i %%i/%%i.proto
)

endlocal
