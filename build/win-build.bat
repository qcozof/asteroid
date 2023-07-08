
@echo off

::0.set version

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

::1.go install github.com/akavel/rsrc@latest
echo install rsrc
go install github.com/akavel/rsrc@latest

::2.Generates a .syso file with embedded resources
echo rsrc manifest
rsrc -manifest win-admin.manifest -o ../asteroid.syso

cd ..
::Build executable

echo build
go build

del asteroid.syso

move /y "%inputFile%".bak "%inputFile%"

echo finished.
