FROM golang:1.17-alpine AS build

WORKDIR /src/
COPY . .
RUN go mod download
RUN CGO_ENABLED=0 go install cmd/rest/rest.go

FROM scratch

COPY app.sra .
COPY app.sra.pub .
COPY --from=build /go/bin/rest /bin/rest
ENTRYPOINT ["/bin/rest"]
