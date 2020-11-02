FROM golang:1.14-alpine

ARG DOCK_PKG_DIR=/go/src/hyc/demo/fractals
WORKDIR $DOCK_PKG_DIR

COPY go.mod .
COPY go.sum .
RUN go mod download

COPY . $DOCK_PKG_DIR

RUN CGO_ENABLED=0 GOOS=linux go build -o fractals

CMD ["./fractals"]
