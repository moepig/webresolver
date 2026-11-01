FROM golang:1.27.3@sha256:2920bb124f1a2e12a1d39f0cab756ae0068441ee1da7655dec7f5e374944095b AS build

WORKDIR /build

COPY go.mod ./
RUN go mod download

COPY *.go ./

RUN go build -o /app

# ---------------
FROM gcr.io/distroless/base-debian13:nonroot@sha256:9039a88709bcc2bacf5a2faac331d65dd6840f763f225ad9594cc0c67bcd0fba

WORKDIR /

COPY --from=build /app /app

EXPOSE 8080
USER nonroot:nonroot

ENTRYPOINT ["/app"]