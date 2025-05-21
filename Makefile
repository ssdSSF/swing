build:
	env GOOS=linux GOARCH=arm go build -o ./bin/swing-cli main.go

deploy: build
	scp ./bin/swing-cli $(PI_AT_HOST):~/

test-env:
	echo $(PI_AT_HOST)