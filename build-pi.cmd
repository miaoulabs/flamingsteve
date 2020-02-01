@setlocal EnableDelayedExpansion
@rem @set PKG_CONFIG_PATH=C:\tools\raspberry\arm-linux-gnueabihf\sysroot\usr\lib\arm-linux-gnueabihf\pkgconfig
@rem @set CGO_CFLAGS=--sysroot=C:\tools\gcc-arm-none-eabi\sysroot -IC:\tools\gcc-arm-none-eabi\sysroot\opt\vc\include -IC:\tools\gcc-arm-none-eabi\sysroot\opt\vc\include\interface\vmcs_host
@rem @set CC=arm-linux-gnueabihf-gcc
@rem @set CGO_LDFLAGS=-lbcm_host -LC:\tools\gcc-arm-none-eabi\sysroot\opt\vc\lib
@rem @set PATH=C:\tools\gcc-arm-none-eabi\bin;%PATH%
@set GOOS=linux
@set GOARCH=arm
@set GOARM=6
@set CGO_ENABLED=1
go build -v -o .test\sensor.arm ./cmd/sensor || exit /b !ERRORLEVEL!
scp .test/sensor.arm pi@protopi:~/
@endlocal
