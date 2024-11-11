FROM golang:alpine

WORKDIR /app

COPY . /app

CMD ["go", "test", "./internal/..."]
# CMD ["ls", "allure-results"]
