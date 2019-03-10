# Spidy
Web Crawler on Go language (Golang).

# How to use
`docker run --rm xaxaxasanyok/spidy -u "https://example.com"`

Usage:
```
Usage of /bin/spidy:
  -d uint
        Maximum depth, 0 - no limit (default 2)
  -r uint
        Requests per second, 0 - no limit (default 5)
  -u string
        Target Uri, required
```

# Build from source
Download source code:
`$ go github.com/3a9LL/spidy`
Go to project directory like $GOPATH/github.com/3a9LL/spidy and build this:
`$ go build -ldflags="-s -w"`

# Example
Lets crawl by url "https://golang.org/doc/articles/go_command.html":
```bash
$ docker run --rm 3abpwasm/spidy -u "https://golang.org/doc/articles/go_command.html"
[+] Spidy started working
http://golang.org
https://golang.org/doc/articles/go_command.html
https://golang.org
https://golang.org/doc/tos.html
https://golang.org/blog/
https://tour.golang.org/
https://golang.org/
https://golang.org/doc/
https://golang.org/pkg/
https://golang.org/project/
https://golang.org/dl/
https://blog.golang.org/
https://golang.org/help/
https://golang.org/LICENSE
https://golang.org/doc/code.html
http://play.golang.org/
https://golang.org/cmd/go/
[+] Spidy finished
```
