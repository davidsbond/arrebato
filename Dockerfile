FROM gcr.io/distroless/static@sha256:21d3f84a4f37c36199fd07ad5544dcafecc17776e3f3628baf9a57c8c0181b3f

COPY arrebato /usr/bin/arrebato
COPY LICENSE /usr/bin/LICENSE
COPY licenses /usr/bin/licenses

ENTRYPOINT ["/usr/bin/arrebato"]
