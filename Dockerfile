FROM gcr.io/distroless/static@sha256:fac888659ca3eb59f7d5dcb0d62540cc5c53615e2671062b36c815d000da8ef4

COPY arrebato /bin/arrebato
ENTRYPOINT ["/bin/arrebato"]
