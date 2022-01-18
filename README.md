<h1 align="center">remdl</h1> 
<p align="center">Remote Self Hosted HTTP Downloader</p>
<hr>

# Introduction
* `remdl` is HTTP Downloader intended to be hosted in remote machines.
* It features HTML Web Client (No JS and Not even CSS!) to add,remove,list downloads.
* You can also download files which are downloaded in remdl to your local machine.
* Most of code in `remdl` comes from [grab](https://github.com/cavaliergopher/grab) and you can consider remdl to be fork of [grab](https://github.com/cavaliergopher/grab) with many of its features removed (remdl takes only URL , Optional Filename as input ,and most other features of grab like Checksum Verification,Basic Auth,etc. are removed) and a minimal HTML Web Client which is some raw HTML written to response body just to let remdl be controlled in web browser.
* remdl was written as personal program just to satisfy the occassional downloading needs of author (not intended to be used by anyone else). There are very high chances that this program may not be for you.

# Installation
You can download binary for your OS from [Releases](https://github.com/varbhat/remdl/releases/latest) . Also , if you have [Go](https://golang.org/) installed , you can install `remdl` by typing this in terminal.

```bash
go install github.com/varbhat/remdl@latest
```

# Features
* Ultra Minimalistic User Interface
* Concurrently Downloads Files
* Authentication to access Web Client
* Optional Authentication to download files which are downloaded in remdl
* Delete Downloaded File (Requires Authentication)
* Change Credentials
* Add,Cancel and List Downloads
* Explicit Filename for Download
* Download Directory as tar or zip

# Usage
`remdl` is Self-Hosted HTTP Downloader

```bash
Usage of ./remdl:
 -addr    <addr> Listen Address (Default: ":7000")
 -cert    <path> Path to TLS Certificate (Required for HTTPS)
 -dir     <path> remdl Directory (Default: "remdldir")
 -key     <path> Path to TLS Key (Required for HTTPS)
 -pass    <pass> Password (Default: "remdlpassword")
 -user    <user> Username (Default: "remdluser")
 -help    <opt>  Print this Help
```

# License
BSD-3-Clause License
