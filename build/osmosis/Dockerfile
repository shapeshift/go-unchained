FROM golang:1.17.5

RUN go install github.com/cespare/reflex@latest

COPY reflex.conf /

CMD ["reflex", "-c", "/reflex.conf"]