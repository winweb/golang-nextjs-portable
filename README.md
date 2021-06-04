# golang-nextjs-portable

**golang-nextjs-portable** is a small Go program to showcase the `embed` package
for bundling a static HTML export of a Next.js app.

👉 Read the companion
[article](https://v0x.nl/articles/portable-apps-go-nextjs) that walks
through this project.

### *** After I fork these project from David Stotijn.

### *** I changed net/http to [Fiber V2](https://github.com/gofiber/fiber) and use [GORM V2](https://github.com/go-gorm/gorm) connect the [SQLite3](https://github.com/mattn/go-sqlite3) database for stress test.

<img src="https://v0x.nl/assets/articles/golang-nextjs-portable-og.png">

## Requirements

- Go 1.16
- Yarn

## Installing

Clone or download the repository:

```sh
$ git clone git@github.com:dstotijn/golang-nextjs-portable.git
```

## Usage

From the repository root directory, generate the static HTML export of the Next.js
app, and build the Go binary:

```sh
$ cd nextjs
$ yarn install
$ yarn run export
$ cd ..
$ go build main.go
```

On Windows PowerShell

```bat
$ cd nextjs\
$ set-ExecutionPolicy RemoteSigned -Scope CurrentUser
$ Get-ExecutionPolicy
$ yarn install
$ yarn run export-win
$ cd ..
$ go build main.go
```

Run Go unit test

```sh
$ go test -v
```

Run Stress test on Windows 10

```bat
$ C:\xampp\apache\bin\ab.exe -k -p test.json -T application/json -c 10000 -n 70000 http://localhost:8080/add
```

Then run the binary:

```sh
$ ./golang-nextjs-portable

2021/04/27 14:55:38 Starting HTTP server at http://localhost:8080 ...
```

On Windows PowerShell

```bat
$ ./main.exe

2021/05/22 20:29:22 Starting HTTP server at http://localhost:8080 ...
```

On Docker

```sh
docker build .
```

## License

[MIT](/LICENSE)

---

© 2021 David Stotijn — [Twitter](https://twitter.com/dstotijn), [Email](mailto:dstotijn@gmail.com), [Homepage](https://v0x.nl)
