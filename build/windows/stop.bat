@echo off
color f5

set serviceName=asteroid.exe
set currentPath=%~dp0
set currentBatName=%~nx0
set nssm=nssm.exe

title stop %serviceName% service
echo current path: %currentPath%
echo app path: %appPath%
echo.

rem check admin permission
>nul 2>&1 "%SYSTEMROOT%\system32\cacls.exe" "%SYSTEMROOT%\system32\config\system"
if "%errorlevel%"=="0" (
    echo Run as administrator
    echo.
    cd /d %currentPath%
) else (
    echo Please right-click %currentBatName% and select 'Run as administrator'.
    echo.

    pause
    exit
)

echo stop service
%nssm% stop %serviceName%
echo.

pause