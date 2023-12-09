# National Grid Demand Flexability Service Notifer
A small project to automatically send a notification to Telegram if there's a new demand flexibility request.

*Build*
```bash
make build
```

*Test*
```bash
make test
```

*Create docker image*
```bash
docker build --pull --rm -f "Dockerfile" -t ng-dfs-notifier:latest "."
```