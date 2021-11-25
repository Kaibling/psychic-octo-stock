FROM  golang:1.17.3-alpine3.14
RUN apk update && apk add go make
COPY . /wrk
WORKDIR /wrk
RUN cd /wrk &&  make build
EXPOSE 8080
ENTRYPOINT ["make","build-start"]