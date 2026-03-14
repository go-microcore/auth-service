FROM gcr.io/distroless/static:nonroot
ARG TARGETARCH
WORKDIR /
COPY --chmod=755 bin/auth-service_linux_${TARGETARCH} /auth-service
ENTRYPOINT ["/auth-service"]
CMD ["help"]
