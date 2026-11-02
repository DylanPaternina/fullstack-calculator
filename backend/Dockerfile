FROM golang:1.27-alpine AS build
WORKDIR /src
COPY go.mod ./
COPY cmd ./cmd
COPY internal ./internal
RUN CGO_ENABLED=0 go build -trimpath -ldflags="-s -w" -o /out/calculator ./cmd/server

FROM scratch
COPY --from=build /out/calculator /calculator
USER 65532:65532
ENV ADDR=0.0.0.0:8080
EXPOSE 8080
ENTRYPOINT ["/calculator"]
