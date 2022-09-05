FROM golang:latest as builder
COPY ./src /src
WORKDIR /src
RUN go build

FROM ubuntu:latest

RUN useradd user
WORKDIR /home/user
RUN chown user:user /home/user
USER user

COPY --from=builder /src/scoin /home/user/scoin

EXPOSE 8090
ENTRYPOINT ["/home/user/scoin", "-init"]
