compile:
	# compile for Linux
	GOOS=linux GOARCH=amd64 go build -o ./bin/gspot_linux_amd64 ./cmd/gspot

	# compile for macOS
	GOOS=darwin GOARCH=amd64 go build -o ./bin/gspot_darwin_amd64 ./cmd/gspot
	
	# compile for Apple M1
	GOOS=darwin GOARCH=arm64 go build -o ./bin/gspot_darwin_arm64 ./cmd/gspot

	# compile for Windows
	GOOS=windows GOARCH=amd64 go build -o ./bin/gspot_win_amd64.exe ./cmd/gspot