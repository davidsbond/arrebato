FROM gcr.io/distroless/static@sha256:2ad95019a0cbf07e0f917134f97dd859aaccc09258eb94edcb91674b3c1f448f

COPY arrebato /usr/bin/arrebato
COPY LICENSE /usr/bin/LICENSE
COPY licenses /usr/bin/licenses

ENTRYPOINT ["/usr/bin/arrebato"]
