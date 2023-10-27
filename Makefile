test:
	go install
	go test -count=1 -v ./cli/exec
