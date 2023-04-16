@echo off

echo -------begin...-------

set appName="asteroid"

SET CGO_ENABLED=0
SET GOOS=linux
SET GOARCH=amd64
echo now the CGO_ENABLED:
go env CGO_ENABLED

echo now the GOOS:
go env GOOS

echo ---build...---
echo now the GOARCH:
go env GOARCH
 
cd ../
go build -o %appName% main.go 

echo ---build end.---

SET CGO_ENABLED=1
SET GOOS=windows
SET GOARCH=amd64

echo now the CGO_ENABLED:
go env CGO_ENABLED

echo now the GOOS:
go env GOOS

echo now the GOARCH:
go env GOARCH

echo finished