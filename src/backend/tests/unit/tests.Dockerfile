FROM golang:alpine

WORKDIR /app

COPY . /app

CMD ["go", "test", "-count=1", "./internal/..."]
# CMD ["ls", "allure-results"]
