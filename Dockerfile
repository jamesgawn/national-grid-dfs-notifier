FROM golang:alpine as build
WORKDIR $GOPATH/src/github.com/jamesgawn/ng-dfs-notifier

# Copy the rest of the project and build
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -a -o /ng-dfs-notifier ./main.go

# Reset to scratch to drop all of the above layers and only copy over the final binary
FROM scratch
ENV HOME=/home
COPY --from=build /ng-dfs-notifier /bin/ng-dfs-notifier

ENTRYPOINT ["/bin/ng-dfs-notifier"]