# limited network driver

## Visualize dependencies

### Install *loov/goda*
```shell
go install github.com/loov/goda@latest
```
### Generate graph
```shell
goda graph ./... | dot -Tsvg -o graph.svg
```

## Run test with coverage

```shell
go test -coverprofile=coverage.out $(go list ./... | grep -v "/mocks")
```
