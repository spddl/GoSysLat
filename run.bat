@echo off

SET filename=GoSysLat

:loop
cls

@REM gocritic check -enable="#performance" ./...
@REM gocritic check -enableAll -disable="#experimental,#opinionated,#commentedOutCode" ./...

go build
@REM go build -race -o %filename%.exe
@REM go build -tags debug -o %filename%.exe

IF %ERRORLEVEL% EQU 0 %filename%.exe
@REM IF %ERRORLEVEL% EQU 0 %filename%.exe -ogl -port 4 -time 30s -print
@REM IF %ERRORLEVEL% EQU 0 %filename%.exe -ogl -fullscreen -port 4 -print
@REM IF %ERRORLEVEL% EQU 0 %filename%.exe -d3d9 -fullscreen -port 4 -print
@REM .\GoSysLat.exe
@REM .\GoSysLat.exe -ogl -port 4 -time 30s
@REM .\GoSysLat.exe -ogl -fullscreen -port 4 -time 30s

pause
goto loop