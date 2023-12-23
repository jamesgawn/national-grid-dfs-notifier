FROM golang:latest as build
WORKDIR $GOPATH/src/github.com/jamesgawn/ng-dfs-notifier

# Copy the rest of the project and build
COPY . .
RUN CGO_ENABLED=1 go build -o /ng-dfs-notifier ./main.go

# Reset to scratch to drop all of the above layers and only copy over the final binary
FROM golang:latest
ENV HOME=/home
COPY --from=build /ng-dfs-notifier /ng-dfs-notifier
COPY --from=build /etc/ssl/certs /etc/ssl/certs
RUN mkdir storage

VOLUME ["/storage"]

ENTRYPOINT ["/ng-dfs-notifier"]