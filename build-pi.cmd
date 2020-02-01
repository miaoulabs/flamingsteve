@setlocal EnableDelayedExpansion
@set CGO_ENABLED=0
@set GOOS=linux
@set GOARCH=arm
@set GOARM=6
go build -v -o .test\sensor.arm ./cmd/sensor || exit /b !ERRORLEVEL!
scp .test/sensor.arm pi@protopi:~/
@endlocal
