
dev:
	wails dev

devWithBrowser:
	wails dev -browser

build:
	wails build

buildForWindows:
	wails build -platform=windows/amd64 -s -m -trimpath -o EchoProxy_windows_amd64.exe

buildForWindowsArm64:
	wails build -platform=windows/arm64 -s -m -trimpath -o EchoProxy_windows_arm64.exe

buildForDarwinArm64:
	wails build -platform=darwin/arm64 -s -m -trimpath -o EchoProxy_darwin_arm64.app

buildForDarwinAmd64:
	wails build -platform=darwin/amd64 -s -m -trimpath -o EchoProxy_darwin_amd64.app

buildForDarwinUniversal:
	wails build -platform=darwin/universal -s -m -trimpath -o EchoProxy_darwin_amd64.app
