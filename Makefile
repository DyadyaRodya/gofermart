include .env
LOCAL_BIN:=$(CURDIR)/bin

.PHONY: install_bin
install_bin: # install binary dependencies
	mkdir -p $(LOCAL_BIN)
	GOBIN=$(LOCAL_BIN) go mod tidy
	GOBIN=$(LOCAL_BIN) go install github.com/vektra/mockery/v2@latest

.PHONY: install
install: install_bin

.PHONY:
mockery:
	$(LOCAL_BIN)/mockery --name $(name) --dir $(dir) --output $(dir)/mocks

.PHONY:
mock:
	echo "Nothing to mock"

.PHONY: lint
lint: # run statictest
	go vet -vettool=/usr/bin/statictest ./...

.PHONY: tests
tests: # run unit tests
	go test -race -coverprofile=coverage.out ./...


.PHONY: test-proj
test-proj: # run gophermarttest
	gophermarttest \
                -test.v -test.run=^TestGophermart$$ \
                -gophermart-binary-path=cmd/gophermart/gophermart \
                -gophermart-host=localhost \
                -gophermart-port=8080 \
                -gophermart-database-uri="postgresql://postgres:postgres@postgres/praktikum?sslmode=disable" \
                -accrual-binary-path=cmd/accrual/accrual_linux_amd64 \
                -accrual-host=localhost \
                -accrual-port=$(random unused-port) \
                -accrual-database-uri="postgresql://postgres:postgres@postgres/praktikum?sslmode=disable"

.PHONY: accrual-start
accrual-start: # start accrual server
	RUN_ADDRESS=${ACCRUAL_SYSTEM_ADDRESS_PARAM} DATABASE_URI="${ACCRUAL_DATABASE_URI}" nohup ./cmd/accrual/accrual_linux_amd64 &

.PHONY: accrual-stop
accrual-stop: # stop accrual server
	ps -ef | awk '$$8=="./cmd/accrual/accrual_linux_amd64" {print $$2}' | xargs -r kill
