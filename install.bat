@echo off
setlocal

echo ========================================
echo       AGY++ Windows Installer
echo ========================================
echo.

powershell -ExecutionPolicy Bypass -NoProfile -File "%~dp0scripts\install\install.ps1"

if %ERRORLEVEL% NEQ 0 (
    echo [ERROR] Installation failed.
    pause
    exit /b %ERRORLEVEL%
)

echo.
echo [SUCCESS] agy++ is installed!
echo You can now open a new terminal and type:
echo   agy++
echo   agy++ -y
echo   agy++ diagnostics
echo.
pause
