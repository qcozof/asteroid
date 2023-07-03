@echo off
set serviceName=asteroid.exe
set currentPath=%~dp0
set currentBatName=%~nx0
set nssm=nssm.exe

title upgrade %serviceName% service
echo current path: %currentPath%
echo app path: %appPath%

echo %currentPath%>>1aa.txt

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

::taskkill /f /im %serviceName% .exe

echo restart service
%nssm% restart %serviceName%
echo.

if "%1"=="exitNow" (
 ::exit 0
) else (
 pause
)

