FROM scratch

COPY ./bin/api-pipeline /go/bin/api-pipeline

ENTRYPOINT ["/go/bin/api-pipeline"]