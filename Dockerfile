FROM golang:1.17-alpine AS build

WORKDIR /src/
COPY . .
RUN go mod download
RUN CGO_ENABLED=0 make install

FROM scratch
COPY --from=build /go/bin/rest /bin/rest
ENTRYPOINT ["/bin/rest"]
