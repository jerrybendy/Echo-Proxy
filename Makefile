
dev:
	wails dev

devWithBrowser:
	wails dev -browser

build:
	wails build

buildFrontend:
	cd frontend
	npm run build

buildForWindows:
	wails build -platform=windows/amd64 -s -m -trimpath -o EchoProxy.exe
	cd build/bin
	zip EchoProxy_windows_amd64.zip EchoProxy.exe
	rm -rf EchoProxy.exe

buildForDarwinUniversal:
	wails build -platform=darwin/universal -s -m -trimpath
	cd build/bin
	zip -r EchoProxy_darwin_universal.zip "./Echo Proxy.app"
	rm -rf "Echo Proxy.app"
