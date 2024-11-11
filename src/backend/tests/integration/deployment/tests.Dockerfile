FROM golang:alpine

WORKDIR /app

COPY . /app

CMD ["go", "test", "./tests/integration", "-count=1"]
# CMD ["ls"]
