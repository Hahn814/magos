FROM golang:1-bookworm

RUN useradd -ms /bin/bash legion
USER legion

WORKDIR /usr/src/app
COPY . .

RUN mkdir -p /home/legion/go/bin 
    # && go build -o /home/legion/go/bin/legion ./cmd/legion/

ENV PATH=/home/legion/go/bin:$PATH