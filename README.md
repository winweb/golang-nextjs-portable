# golang-nextjs-portable

**golang-nextjs-portable** is a small Go program to showcase the `embed` package
for bundling a static HTML export of a Next.js app.

👉 Read the companion
[article](https://v0x.nl/articles/portable-apps-go-nextjs) that walks
through this project.

<img src="https://v0x.nl/assets/articles/golang-nextjs-portable-og.png">

## Requirements

- Go 1.16.4
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

Run Stress test

```bat
$ C:\xampp\apache\bin\ab.exe -k -c 10000 -n 100000 http://localhost:8080/add
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
