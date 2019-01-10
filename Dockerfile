FROM golang AS builder

RUN mkdir /src

COPY go.mod go.sum /src/

WORKDIR /src

RUN go mod download

COPY . /src/

RUN CGO_ENABLED=0 GOOS=linux go build -o migrator

FROM scratch

COPY --from=builder /src/migrator /src/

ENTRYPOINT ["/src/migrator"]
