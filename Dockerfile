FROM golang:alpine as build
WORKDIR $GOPATH/src/github.com/jamesgawn/ng-dfs-notifier

# Copy the rest of the project and build
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -a -o /ng-dfs-notifier ./main.go

# Reset to scratch to drop all of the above layers and only copy over the final binary
FROM scratch
ENV HOME=/home
COPY --from=build /ng-dfs-notifier /ng-dfs-notifier
COPY --from=build /etc/ssl/certs /etc/ssl/certs

VOLUME ["/storage"]

ENTRYPOINT ["/ng-dfs-notifier"]