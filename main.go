package main

import (
	"crypto/md5"
	"embed"
	"flag"
	"fmt"
	"os"

	"github.com/0x1f610/foo_cover_upload/api"
	"github.com/0x1f610/foo_cover_upload/upload"
	"github.com/0x1f610/foo_cover_upload/utils"
)

var (
	// Modes
	isModeServe bool
	isModeGen   bool
	isModeHelp  bool

	// For uploading
	url            string
	authentication string

	// For hosting
	host          string
	address       string
	passHash      string
	redisHost     string
	redisPassword string

	//go:embed html/*
	FS embed.FS
)

func main() {
	flag.BoolVar(&isModeServe, "serve", false, "Launch a server.")
	flag.BoolVar(&isModeGen, "gen", false, "Generate a key to use for the image server.")
	flag.BoolVar(&isModeHelp, "help", false, "Displays this help message.")

	flag.StringVar(&url, "url", "", "Address of the server to upload to, including port.")
	flag.StringVar(&authentication, "auth", "", "Image server password, if required.")

	flag.StringVar(&host, "host", "", "The host of your server. Use your domain that can be used by Discord to access images.")
	flag.StringVar(&address, "addr", "127.0.0.1:2131", "Listening address of the image upload server.")
	flag.StringVar(&passHash, "hash", "", "Sets the MD5 hash of your password on the server.")
	flag.StringVar(&redisHost, "redis-host", "127.0.0.1:6379", "Host of the redis instance.")
	flag.StringVar(&redisPassword, "redis-password", "", "Password of the redis instance.")

	flag.Usage = func() {
		w := flag.CommandLine.Output()

		fmt.Fprintf(w, "\n===============================\n\n")
		fmt.Fprintf(w, "foo_cover_upload HELP PAGE\n")
		fmt.Fprintf(w, "\n===============================\n\n")

		fmt.Fprintf(w, "This program houses both the client and the server.\n")
		fmt.Fprintf(w, "By default, the program will use client mode, which is intended for uploading.\n")
		fmt.Fprintf(w, "To change mode to server, you must pass the -serve flag to the program along with the necessary parameters.\n")
		fmt.Fprintf(w, "\n===============================\n\n")

		fmt.Fprintf(w, "AVAILABLE MODES:\n")
		flag.VisitAll(func(f *flag.Flag) {
			if f.Name == "serve" || f.Name == "gen" || f.Name == "help" {
				fmt.Printf("-%s\n\t%s\n", f.Name, f.Usage)
			}
		})
		fmt.Fprintf(w, "\n===============================\n\n")

		fmt.Fprintf(w, "AVAILABLE PARAMETERS:\n")
		flag.VisitAll(func(f *flag.Flag) {
			if f.Name != "serve" && f.Name != "gen" && f.Name != "help" {
				fmt.Printf("-%s\n\t%s\n", f.Name, f.Usage)
			}
		})
	}

	flag.Parse()

	// If no arguments, print help page
	if len(os.Args) == 1 {
		flag.Usage()
	} else {
		if isModeServe {
			// Check all necessary parameters
			if host == "" {
				fmt.Println("ERROR: You must enter the host of the image server.")
				fmt.Println("Exitting.")
			} else {
				// Set up upload server
				api.Run(FS, host, address, passHash, redisHost, redisPassword)
			}
		} else if isModeGen {
			// Generate key
			newKey := utils.GenerateString(24)
			newKeyHash := md5.Sum([]byte(newKey))

			fmt.Printf("Key: %s\nHash: %x\n\n", newKey, newKeyHash)
			fmt.Println("Please use the hash with the -hash flag in the server.")
			fmt.Println("and the password with the -auth flag in the client.")
			fmt.Println("They will not be shown to you again.")
		} else if isModeHelp {
			// Print help page
			flag.Usage()
		} else {
			// If no mode is used, it implies the user wants to upload an image
			if url != "" {
				upload.Run(url, authentication)
			} else {
				fmt.Println("ERROR: No image server address defined.")
				fmt.Println("Exitting.")
			}
		}
	}
}
