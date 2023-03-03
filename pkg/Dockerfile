FROM gcr.io/google-appengine/golang

WORKDIR /bin
COPY /usr/local/bin/wkhtmltopdf /bin/wkhtmltopdf

RUN go build -o /bin/api ./cmd/appengine

ENTRYPOINT [ "/bin/api" ]
