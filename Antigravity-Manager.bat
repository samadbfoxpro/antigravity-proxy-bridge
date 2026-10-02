@echo off
chcp 65001 >nul
setlocal enabledelayedexpansion
title Antigravity Proxy Manager

:: Default Proxy Configuration
set DEFAULT_PORT=10809
set PROXY_HOST=127.0.0.1
set PROXY_PORT=%DEFAULT_PORT%

:MAIN_MENU
cls
set TARGET_PROXY=http://%PROXY_HOST%:%PROXY_PORT%

:: Check Status for IDE
set IDE_STATUS=[DISABLED]
set SETTINGS_FILE=%APPDATA%\Antigravity IDE\User\settings.json
if exist %SETTINGS_FILE% (
    findstr /i http.proxy %SETTINGS_FILE% >nul 2>&1
    if !errorlevel! equ 0 (
        set IDE_STATUS=[ENABLED ]
    )
)

:: Check Status for Antigravity 2.0 (User Env Var)
set V2_STATUS=[DISABLED]
reg query HKCU\Environment /v HTTP_PROXY >nul 2>&1
if !errorlevel! equ 0 (
    set V2_STATUS=[ENABLED ]
)

echo ======================================================================
echo                  Antigravity Proxy Manager
echo ======================================================================
echo.
echo   [ STATUS ]
echo   * Antigravity IDE Proxy : %IDE_STATUS%
echo   * Antigravity 2.0 Proxy : %V2_STATUS%
echo.
echo   [ CONFIG ]
echo   * Target Proxy Address  : %TARGET_PROXY%
echo.
echo ----------------------------------------------------------------------
echo   [1] Toggle / Set Proxy for: Antigravity IDE
echo   [2] Toggle / Set Proxy for: Antigravity 2.0
echo   [3] Set Proxy for BOTH (IDE + 2.0)
echo   [4] Reset / Disable ALL Proxies
echo   --------------------------------------------------------------------
echo   [5] Test Proxy Connection (Ping Google via Proxy)
echo   [6] Change Proxy Port (Current: %PROXY_PORT%)
echo   [0] Exit
echo ======================================================================
set /p CHOICE=Enter your choice [0-6]: 

if %CHOICE%==1 goto SET_IDE
if %CHOICE%==2 goto SET_V2
if %CHOICE%==3 goto SET_BOTH
if %CHOICE%==4 goto RESET_ALL
if %CHOICE%==5 goto TEST_CONN
if %CHOICE%==6 goto CHANGE_PORT
if %CHOICE%==0 exit /b
goto MAIN_MENU

:SET_IDE
echo.
echo [*] Applying proxy to Antigravity IDE...
powershell -NoProfile -Command ^
   = [System.Environment]::ExpandEnvironmentVariables('%SETTINGS_FILE%'); ^
   = Split-Path ; if (-not (Test-Path )) { New-Item -ItemType Directory -Force -Path | Out-Null }; ^
   = @{}; ^
  if (Test-Path ) { try { = (Get-Content -Raw -Encoding UTF8 | ConvertFrom-Json -AsHashtable) } catch { = @{} } }; ^
  ['http.proxy'] = '%TARGET_PROXY%'; ^
  ['http.proxySupport'] = 'override'; ^
  ['http.proxyStrictSSL'] = False; ^
   | ConvertTo-Json -Depth 10 | Set-Content -Encoding UTF8
echo [SUCCESS] Antigravity IDE proxy updated to: %TARGET_PROXY%
timeout /t 2 >nul
goto MAIN_MENU

:SET_V2
echo.
echo [*] Setting proxy for Antigravity 2.0 (User Env Vars)...
setx HTTP_PROXY %TARGET_PROXY% >nul
setx HTTPS_PROXY %TARGET_PROXY% >nul
echo [SUCCESS] Antigravity 2.0 proxy updated to: %TARGET_PROXY%
timeout /t 2 >nul
goto MAIN_MENU

:SET_BOTH
echo.
echo [*] Applying to both IDE and Antigravity 2.0...
call :SET_IDE >nul 2>&1
call :SET_V2 >nul 2>&1
echo [SUCCESS] Both Antigravity IDE and 2.0 are now configured!
timeout /t 2 >nul
goto MAIN_MENU

:RESET_ALL
echo.
echo [*] Resetting IDE settings...
powershell -NoProfile -Command ^
   = [System.Environment]::ExpandEnvironmentVariables('%SETTINGS_FILE%'); ^
  if (Test-Path ) { ^
   try { ^
   = (Get-Content -Raw -Encoding UTF8 | ConvertFrom-Json -AsHashtable); ^
   .Remove('http.proxy'); ^
   .Remove('http.proxySupport'); ^
   .Remove('http.proxyStrictSSL'); ^
   | ConvertTo-Json -Depth 10 | Set-Content -Encoding UTF8; ^
   } catch {} ^
  }

echo [*] Removing Antigravity 2.0 environment variables...
REG delete HKCU\Environment /F /V HTTP_PROXY >nul 2>&1
REG delete HKCU\Environment /F /V HTTPS_PROXY >nul 2>&1
echo [SUCCESS] All proxy settings have been disabled and reset!
timeout /t 2 >nul
goto MAIN_MENU

:TEST_CONN
echo.
echo [*] Testing connection to https://www.google.com through %TARGET_PROXY% ...
powershell -NoProfile -Command ^
  try { ^
   $sw = [System.Diagnostics.Stopwatch]::StartNew(); ^
   $response = Invoke-WebRequest -Uri 'https://www.google.com' -Proxy '%TARGET_PROXY%' -TimeoutSec 7 -UseBasicParsing; ^
   $sw.Stop(); ^
   Write-Host '[SUCCESS] Proxy is working!' -ForegroundColor Green; ^
   Write-Host ('[INFO] Response Status: ' + $response.StatusCode + ' (' + $sw.ElapsedMilliseconds + ' ms)') -ForegroundColor Cyan; ^
  } catch { ^
   Write-Host ('[FAILED] Connection failed: ' + $_.Exception.Message) -ForegroundColor Red; ^
  }
echo.
pause
goto MAIN_MENU

:CHANGE_PORT
echo.
echo Current Port is: %PROXY_PORT%
set /p NEW_PORT=Enter new proxy port (e.g. 10809, 2081, 7890): 
if not %NEW_PORT%==" set PROXY_PORT=%NEW_PORT%
goto MAIN_MENU
