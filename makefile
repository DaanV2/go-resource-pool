assembly:
	go tool compile -S x.go

test:
	go test ./... --cover -coverprofile=reports/coverage.out

coverage-report: test
	go tool cover -html=reports/coverage.out