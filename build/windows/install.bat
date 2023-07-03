@echo off
color f0

set serviceName=asteroid.exe
set nssm=nssm.exe
set currentPath=%~dp0
set appPath=%currentPath%..\..\
set currentBatName=%~nx0

title install %serviceName% as service
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

echo uninstall old service ...
start uninstall.bat exitNow

echo wait 5 seconds ...
ping 127.0.0.1 -n 5 > nul
echo.

echo install
%nssm% install %serviceName% %appPath%%serviceName%.exe
echo.

echo set Application path
%nssm% set %serviceName% Application %appPath%%serviceName%.exe
echo.

echo set Application directory
%nssm% set %serviceName% AppDirectory %appPath%
echo.

echo set Application restart delay
%nssm% set %serviceName% AppRestartDelay 30000
echo.

echo start service
%nssm% start %serviceName%
echo.

pause