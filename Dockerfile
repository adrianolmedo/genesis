FROM golang:1.17-alpine AS build

WORKDIR /src/
COPY . .
RUN go mod download
RUN CGO_ENABLED=0 go install cmd/rest/rest.go

ARG USERNAME
ARG USER_UID
ARG USER_GID=${USER_UID}
RUN addgroup -g ${USER_GID} -S ${USERNAME} \
    && adduser -D -u ${USER_UID} -S ${USERNAME} -s /bin/sh ${USERNAME}

FROM scratch

COPY app.sra .
COPY app.sra.pub .
COPY --from=build /go/bin/rest /bin/rest
USER ${USERNAME}
ENTRYPOINT ["/bin/rest"]