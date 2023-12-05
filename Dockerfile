FROM --platform=$BUILDPLATFORM golang:alpine AS builder
WORKDIR /build
COPY deploy_no_license.sh /usr/local/share/deploy_no_license.sh
COPY hls_tools_copy.tar /usr/local/share/hls_tools.tar
COPY go.mod ./
RUN go mod download
COPY . ./
ARG TARGETOS
ARG TARGETARCH
ENV GOOS $TARGETOS
ENV GOARCH $TARGETARCH
RUN go build -o streamvalidator.exe ./cmd

FROM --platform=linux/amd64 centos:centos7
COPY --from=builder ["/build/streamvalidator.exe", "/"]
COPY --from=builder ["/usr/local/share/hls_tools.tar", "/"]
RUN tar xvf hls_tools.tar
COPY --from=builder ["/usr/local/share/deploy_no_license.sh", "/deploy.sh"]
RUN chmod +x deploy.sh
RUN ./deploy.sh
ENV LD_LIBRARY_PATH="${LD_LIBRARY_PATH}:/usr/local/lib"
EXPOSE 3001
RUN chmod +x /streamvalidator.exe
ENTRYPOINT [ "/streamvalidator.exe" ]