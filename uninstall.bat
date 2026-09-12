@echo off
setlocal

echo ========================================
echo       AGY++ Windows Uninstaller
echo ========================================
echo.

powershell -ExecutionPolicy Bypass -NoProfile -File "%~dp0scripts\install\install.ps1" -Uninstall

if %ERRORLEVEL% NEQ 0 (
    echo [ERROR] Uninstall failed.
    pause
    exit /b %ERRORLEVEL%
)

echo.
echo [SUCCESS] agy++ has been uninstalled.
echo.
pause
