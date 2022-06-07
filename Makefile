build: main.go
	go build -o ./bin/demoexam-checker main.go
	GOOS=windows go build -o ./bin/demoexam-checker.exe main.go

