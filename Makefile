client:
	go run cmd/client/main.go

# Запустить все авто тесты
test:
	go test -v ./...

# Запустить все авто тесты с подсчетом процента покрытия
testPC:
	go test -v -coverprofile=coverage.out ./internal/...
	go tool cover -func=coverage.out | grep "total:"
	rm coverage.out

# Создать отчёт по покрытию кода тестами
testCoverage:
	go test -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html
	rm coverage.out
	@echo "Coverage report generated: coverage.html"

# генерируем код
gen:
	go generate ./...

# Остановить выполнение PostgreSQL,
# чтобы запустить в докере
stopPG:
	sudo systemctl stop postgresql

# Создание файла размером 1 ГБ
create1:
	mkdir -p ./test_file/
	dd if=/dev/zero of=./test_file/1GB_file bs=1G count=1

# Создание файла размером 10 ГБ
create10:
	mkdir -p ./test_file/
	dd if=/dev/zero of=./test_file/10GB_file bs=1G count=10

# Создание файла размером 34 ГБ
create34:
	mkdir -p ./test_file/
	dd if=/dev/zero of=./test_file/34GB_file bs=1G count=34

# Создание файла размером 100 ГБ
create100:
	mkdir -p ./test_file/
	dd if=/dev/zero of=./test_file/100GB_file bs=1G count=100

createCertificates:
	cd ./certificate/ &&	openssl req -x509 -newkey rsa:4096 -keyout key.pem -out cert.pem -days 365 -nodes \
	-subj "/C=RU/ST=Moscow/L=Moscow/O=GophKeeper/OU=Server/CN=localhost" \
	-addext "subjectAltName=DNS:localhost,DNS:127.0.0.1,IP:127.0.0.1"

