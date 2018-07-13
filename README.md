# Gift certificate export
Create a config.json file at the root level with basic auth credentials 

```
{
	"host": "www.somehostname.com",
	"credentials": "Basic V2h5IGFyZSB5b3UgZGVjb2RpbmcgdGhpcz8="
}
```

Either compile the main.go file using `go build main.go` then `./main` or run it using `go run main.go`
