@echo off
REM Build script for Windows
set BINARY_NAME=boilerblade.exe
set CMD_DIR=cmd\cli
set BUILD_DIR=bin

echo Building %BINARY_NAME%...
if not exist %BUILD_DIR% mkdir %BUILD_DIR%
go build -o %BUILD_DIR%\%BINARY_NAME% .\%CMD_DIR%

if %ERRORLEVEL% EQU 0 (
    echo ✓ Binary built successfully: %BUILD_DIR%\%BINARY_NAME%
    echo.
    echo To use the binary, add to PATH:
    echo   set PATH=%%PATH%%;%CD%\%BUILD_DIR%
    echo.
    echo Or copy to a directory in your PATH:
    echo   copy %BUILD_DIR%\%BINARY_NAME% C:\Windows\System32\
) else (
    echo ✗ Build failed
    exit /b 1
)
