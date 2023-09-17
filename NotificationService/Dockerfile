FROM golang

ADD . /NotificationService/

WORKDIR /NotificationService/

EXPOSE 8181

ENV FBRQ_TOKEN=eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3MjU0NDM5ODgsImlzcyI6ImZhYnJpcXVlIiwibmFtZSI6Imh0dHBzOi8vdC5tZS9zaWRvcmVua29faXZhbiJ9.rhfT9od8RdtPK65W9XwtifvFuB2MGP84S-itmvf8PdI

RUN go get ./...
