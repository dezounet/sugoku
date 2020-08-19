FROM alpine:3.12

COPY sugokud-linux-amd64 /opt/sugokud

# musl (alpine) and glibc so are compatible
RUN mkdir /lib64 && ln -s /lib/libc.musl-x86_64.so.1 /lib64/ld-linux-x86-64.so.2

RUN adduser -D -H -s /dev/null -S sugoku
USER sugoku

EXPOSE 8080

ENTRYPOINT [ "/opt/sugokud", "-p", "8080", "-s", "3" ]
