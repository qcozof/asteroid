@echo off

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

echo finished.