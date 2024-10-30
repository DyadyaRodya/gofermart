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
	$(LOCAL_BIN)/mockery --name $(name) --dir $(dir) --output $(output)/mocks

.PHONY:
mock:
	make mockery name=OrderAccrualGateway dir=./internal/interactors/interfaces output=./internal/interactors
	make mockery name=LuhnService dir=./internal/interactors/interfaces output=./internal/interactors
	make mockery name=PasswordService dir=./internal/interactors/interfaces output=./internal/interactors
	make mockery name=LoginService dir=./internal/interactors/interfaces output=./internal/interactors
	make mockery name=RepoSession dir=./internal/interactors/interfaces output=./internal/interactors
	make mockery name=Repository dir=./internal/interactors/interfaces output=./internal/interactors
	make mockery name=UUIDGenerator dir=./internal/interactors/interfaces output=./internal/interactors

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
                -gophermart-database-uri="${DATABASE_URI}" \
                -accrual-binary-path=cmd/accrual/accrual_linux_amd64 \
                -accrual-host=localhost \
                -accrual-port=8100 \
                -accrual-database-uri="${ACCRUAL_DATABASE_URI}"

.PHONY: accrual-start
accrual-start: # start accrual server
	RUN_ADDRESS=${ACCRUAL_SYSTEM_ADDRESS_PARAM} DATABASE_URI="${ACCRUAL_DATABASE_URI}" nohup ./cmd/accrual/accrual_linux_amd64 &

.PHONY: accrual-stop
accrual-stop: # stop accrual server
	ps -ef | awk '$$8=="./cmd/accrual/accrual_linux_amd64" {print $$2}' | xargs -r kill
