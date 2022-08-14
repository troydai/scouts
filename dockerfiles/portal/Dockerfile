FROM alpine:3.14

RUN apk update && apk upgrade
RUN apk add bind-tools curl bash bash-doc bash-completion apache2-utils

USER 1000

ENTRYPOINT ["tail", "-f", "/dev/null"]
