FROM golang:alpine

WORKDIR /app

COPY . /app

CMD ["go", "test", "./tests/e2e", "-count=1"]
# CMD ["ls"]
