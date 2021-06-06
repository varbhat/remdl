package main

import (
	"flag"
	"fmt"
	"os"
)

func checkDir(dir string) {
	fi, err := os.Stat(dir)

	if os.IsNotExist(err) {
		er := os.MkdirAll(dir, 0755)
		if er != nil {
			Err.Fatalln("Error Creating Directory")
		}
		Info.Println("Created Directory ", dir)
		return
	}

	if !fi.IsDir() {
		Err.Fatalln("Non-Directory File Present")
	}
}

func init() {

	flag.Usage = func() {
		fmt.Fprintf(flag.CommandLine.Output(), "remdl is Self-Hosted HTTP Downloader\n\n")
		fmt.Fprintf(flag.CommandLine.Output(), "Usage of %s:\n", os.Args[0])

		flag.VisitAll(func(f *flag.Flag) {
			fmt.Fprintf(flag.CommandLine.Output(), " -%-5v   %v\n", f.Name, f.Usage)
		})
		fmt.Fprintf(flag.CommandLine.Output(), " -%-5v   %v\n", "help", "<opt>  Print this Help")
	}

	flag.StringVar(&Cred.ListenAddress, "addr", ":7000", `<addr> Listen Address (Default: ":7000")`)
	flag.StringVar(&Cred.TLSKeyPath, "key", "", "<path> Path to TLS Key (Required for HTTPS)")
	flag.StringVar(&Cred.TLSCertPath, "cert", "", "<path> Path to TLS Certificate (Required for HTTPS)")
	flag.StringVar(&Cred.DirPath, "dir", "remdldir", `<path> remdl Directory (Default: "remdldir")`)
	flag.StringVar(&Cred.Username, "user", "remdluser", `<user> Username (Default: "remdluser")`)
	flag.StringVar(&Cred.Password, "pass", "remdlpassword", `<pass> Password (Default: "remdlpassword")`)
	flag.Parse()

	if len(flag.Args()) != 0 {
		fmt.Fprintf(flag.CommandLine.Output(), "Invalid Flags Provided: %s\n\n", flag.Args())
		flag.Usage()
		return
	}

	Cred.Locked = true

	if len(os.Getenv("REMDL_USERNAME")) > 0 {
		Cred.Username = os.Getenv("REMDL_USERNAME")
	}
	if len(os.Getenv("REMDL_PASSWORD")) > 0 {
		Cred.Username = os.Getenv("REMDL_PASSWORD")
	}

	// Create Required SubDirectories if not present
	checkDir(Cred.DirPath)

}
