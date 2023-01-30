FROM golang:1.20rc3-alpine3.17 as Builder
WORKDIR /go/src
COPY ./src/ ./
RUN go get . && CGO_ENABLED=0 go install

FROM scratch
COPY --from=Builder /go/bin/proxauth /

ENV CONFIG_FILE = "/config/config.yaml"
ENV SERVER_SECRET="changeMe"
ENV PORT="8080"
ENV JWT_EXPIRATION_DURATION="24h"

VOLUME "/config"
EXPOSE 8080

ENTRYPOINT [ "/proxauth" ]
