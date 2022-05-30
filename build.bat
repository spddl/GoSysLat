cls
@REM broken in 1.18.0
@REM gocritic check -enableAll -disable="#experimental,#opinionated,#commentedOutCode" ./...
@REM gocritic check -enableAll ./...

@REM go build
::go build -ldflags="-s -w -H windowsgui"
go build -ldflags="-s -w"
pause