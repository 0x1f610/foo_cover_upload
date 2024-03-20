@echo off

set GOARCH=amd64
:: linux builds
set GOOS=linux
go build -o .build/foo_cover_upload_linux_amd64
echo Linux OK
:: windows builds
set GOOS=windows
go build -o .build/foo_cover_upload_win_amd64.exe
echo Windows OK

pause