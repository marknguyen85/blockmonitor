# Start from a Debian image with the latest version of Go installed
# and a workspace (GOPATH) configured at /go.
FROM golang

RUN  mkdir -p /go/src \
  && mkdir -p /go/bin \
  && mkdir -p /go/pkg
ENV GOPATH=/go
ENV PATH=$GOPATH/bin:$PATH   

# now copy your app to the proper build path
RUN mkdir -p $GOPATH/src/telealert
ADD . $GOPATH/src/telealert

# should be able to build now
WORKDIR $GOPATH/src/telealert

RUN go build .
CMD ["/go/src/telealert/telegramalert"]



