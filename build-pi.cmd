@setlocal EnableDelayedExpansion
@set CGO_ENABLED=0
@set GOOS=linux
@set GOARCH=arm
@set GOARM=6
go build -v -o .test\sensor.arm ./cmd/sensor || exit /b !ERRORLEVEL!
go build -v -o .test\muthur.arm ./cmd/muthur || exit /b !ERRORLEVEL!
go build -v -o .test\dispmatrix.arm ./cmd/dispmatrix || exit /b !ERRORLEVEL!
@rem scp .test/sensor.arm pi@protopi:~/
rclone copy --stats-one-line -P --sftp-host protopi --sftp-user pi --sftp-key-file ~/.ssh/id_rsa --include *.arm .test/ :sftp:
scp docker-compose-rpi.yml pi@protopi:~/docker-compose.yml
scp Dockerfile-arm pi@protopi:~/Dockerfile
@endlocal
