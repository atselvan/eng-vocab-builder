FROM golang:alpine3.16

ARG BINARY_NAME="eng-vocab-builder"

RUN mkdir /appl

ENV PATH="${PATH}:/appl"

COPY bin/${BINARY_NAME}-linux /appl/${BINARY_NAME}

COPY docker-context/config.env /appl

WORKDIR /appl

RUN chmod +x *

CMD ["eng-vocab-builder"]
