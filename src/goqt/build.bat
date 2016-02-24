@echo off
cd/d "%~dp0"

call:build_dir qtdrv
call:build_dir tools\rcc

md "%GOROOT%pkgs\bin\%GO_PKG_ARCH%\"
move /y bin\*.exe "%GOROOT%pkgs\bin\%GO_PKG_ARCH%\"
move /y bin\*.dll "%GOROOT%pkgs\bin\%GO_PKG_ARCH%\"

cd ui
go install

goto:eof


:build_dir

cd/d "%~f1"
rd/s/q debug release

qmake "CONFIG+=release"
mingw32-make
cd/d "%~dp0"
goto:eof
