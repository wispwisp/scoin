FROM ubuntu:latest

RUN apt update && apt install -y wget

ARG go_archive=go1.19.linux-amd64.tar.gz
ARG go_archive_sha=464b6b66591f6cf055bc5df90a9750bf5fbc9d038722bb84a9d56a2bea974be6
ARG URL=https://go.dev/dl

RUN wget ${URL}/${go_archive} && \
  echo "${go_archive_sha} ${go_archive}" | sha256sum --check --status && \
  tar -xvf ${go_archive} && \
  mv go /usr/local && \
  rm ${go_archive}

ENV PATH="/usr/local/go/bin:${PATH}"

RUN useradd user
WORKDIR /home/user

RUN chown user:user /home/user
USER user

