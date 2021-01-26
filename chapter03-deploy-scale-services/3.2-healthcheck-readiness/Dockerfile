FROM 3dsinteractive/golang:1.14-alpine3.9-librdfkafka1.4.0

# 1. Add all go files and build
WORKDIR /go/src/bitbucket.org/automationworkshop/main
ADD . /go/src/bitbucket.org/automationworkshop/main
RUN go build -mod vendor -i -tags "musl static_all" .

# ================
FROM 3dsinteractive/alpine:3.9

# 2. Use multi stage docker file, and copy executable file from previous stage
COPY --from=0 /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=0 /go/src/bitbucket.org/automationworkshop/main/main /main

# 3. Add entrypoint.sh (The file that will use as entrypoint when run go program)
ADD ./entrypoint.sh /entrypoint.sh

# 4. Create user 1001 (It is my practice to always use user with id 1001)
RUN adduser -u 1001 -D -s /bin/sh -G ping 1001
RUN chown 1001:1001 /entrypoint.sh
RUN chown 1001:1001 /main

# 5. Make entrypoint.sh and main executable
RUN chmod +x /entrypoint.sh
RUN chmod +x /main

# 6. Set default user
USER 1001

# 7. Expose port 8080
EXPOSE 8080

# 8. Start program by run entrypoint.sh
ENTRYPOINT ["/entrypoint.sh"]
