FROM gcr.io/distroless/static@sha256:d6fa9db9548b5772860fecddb11d84f9ebd7e0321c0cb3c02870402680cc315f

COPY arrebato /bin/arrebato
ENTRYPOINT ["/bin/arrebato"]
