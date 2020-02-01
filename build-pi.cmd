@setlocal EnableDelayedExpansion
@set GOOS=linux
@set GOARCH=arm
@set GOARM=6
@set CGO_ENABLED=1
go build -v -o .test\sensor.arm ./cmd/sensor || exit /b !ERRORLEVEL!
scp .test/sensor.arm pi@protopi:~/
@endlocal
