@echo off
cd/d "%~dp0"

call ..\..\..\init.bat

cd..
for /f %%g in ('dir/s/b .git') do (
    cd/d "%%~dpg"
    go install -v
)

pause