FROM alpine
COPY stack2slack /bin/stack2slack
ENTRYPOINT ["/bin/stack2slack"]
