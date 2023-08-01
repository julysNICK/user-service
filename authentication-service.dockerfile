FROM alpine:lastest

RUN mkdir /app

COPY user /app

CMD ["/app/user"]