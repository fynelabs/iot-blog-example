# iot-blog-example

This repository is to demonstrate the content of https://fynelabs.com/2022/12/22/using-go-and-fyne-for-your-next-embedded-software-development/ .

# Install

On Linux, install and run the server with the following command:
```
$ go install github.com/fynelabs/iot-blog-example/cmd/iot-server@latest
go: downloading github.com/fynelabs/iot-blog-example v0.0.0-20221222182935-ac6c9f928fcb
$ iot-server
2022/12/22 13:56:19 Server running on port 8080
```

Now, your developer workstation, you can get the graphical user install after installing Fyne: https://developer.fyne.io/started/ with the following command on Linux:
```
$ sudo fyne get github.com/fynelabs/iot-blog-example/cmd/iot-ui
```
On other OS:
```
$ sudo fyne get github.com/fynelabs/iot-blog-example/cmd/iot-ui
```

Now, you can just run iot-ui directly from your desktop environment.

It is possible to compile this application for other OS easily using fyne-cross.
