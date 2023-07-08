@echo off

echo -------begin...-------

set appName="asteroid"

::-----------------------------------
for /f "delims=" %%i in ('git describe --abbrev=0 --tags') do set "latest_tag=%%i"

set "search=${version}"
set "replace=%latest_tag%"

set "inputFile=description.txt"
set "outputFile=output.txt"

copy "%inputFile%" "%inputFile%".bak

(for /f "delims=" %%i in ('type "%inputFile%"') do (
    set "line=%%i"
    call set "line=%%line:%search%=%replace%%%"
    echo(!line!
)) > "%outputFile%"

move /y "%outputFile%" "%inputFile%"
::-----------------------------------

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

move /y "%inputFile%".bak "%inputFile%"

echo finished

