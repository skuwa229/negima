
FROM golang:1.15.1-alpine3.12 as builder

RUN apk --update upgrade
RUN apk --no-cache --no-progress add make git gcc musl-dev

WORKDIR /build
COPY . .
RUN go build .

FROM node:10-alpine
RUN apk update && apk add --no-cache --virtual ca-certificates
COPY --from=builder /build/negima /usr/bin/negima

LABEL version="1.0.0"
LABEL repository="https://github.com/skuwa229/negima"
LABEL homepage="https://github.com/skuwa229/negima"
LABEL maintainer="Shota Kuwahara"

LABEL com.github.actions.name="Negima"
LABEL com.github.actions.description="Negima"
LABEL com.github.actions.icon="check"
LABEL com.github.actions.color="green"

ENV JEST_CMD ./node_modules/.bin/jest
COPY entrypoint.sh /entrypoint.sh

ENTRYPOINT ["/entrypoint.sh"]
CMD [""]
