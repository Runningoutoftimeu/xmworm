package injector


func Run(payload[] byte, injecttype string) {
	switch injecttype {
	case "shellcode":
		shellcodeLoad(payload)
	case "pe":
		// do some magical thing
	case "dll":
		// do another magical illegal programming technic here
	}
}