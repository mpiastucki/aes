package cli

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/happymanju/aes/pkg"
)

func Run(args []string) int {
	encFlags := flag.NewFlagSet("enc", flag.ExitOnError)
	encIn := flag.Bool("stdin", false, "read from stdin instead of a file")
	encOut := flag.Bool("stdout", false, "write to stdout instead of a file")
	encOutputFormat := flag.String("f", "hex", "output format; values: txt(base64), bin")

	decFlags := flag.NewFlagSet("dec", flag.ExitOnError)
	decIn := flag.Bool("stdin", false, "read from stdin instead of file")
	decOut := flag.Bool("stdout", false, "write to stdout instead of a file")

	var ciphertext []byte
	plaintext := []byte{}

	switch args[1] {
	case "enc":
		encFlags.Parse(args[2:])
		sc := bufio.NewScanner(os.Stdin)
		var out *bufio.Writer

		if *encIn {
			fmt.Println("enter plaintext to encrypt")
			fmt.Print(">> ")
			for sc.Scan() {

				if err := sc.Err(); err != nil {
					log.Printf("error reading from stdin: %v\n", err)
				}
				data := []byte(sc.Text())
				plaintext = append(plaintext, data...)
			}
		} else {
			fp := filepath.Clean(args[len(args)-1])
			f, err := os.Open(fp)
			if err != nil {
				log.Printf("error opening target file: %v\n", err)
				return 1
			}
			defer f.Close()

			data, err := io.ReadAll(f)
			if err != nil {
				log.Printf("error reading entire file: %v\n", err)
			}
			password := ""
			fmt.Println("password for encryption:")
			fmt.Print(">>")

			for sc.Scan() {
				t := sc.Text()
				if t == "" {
					fmt.Print(">>")
					continue
				}
				password = strings.TrimSpace(sc.Text())
				break
			}
			ciphertext, err = pkg.EncryptWithPassword(data, []byte(password))
			if err != nil {
				log.Printf("error encrypting with password: %v\n", err)
				return 1
			}
		}
		if *encOut {
			out = bufio.NewWriter(os.Stdout)
		}

	case "dec":
		decFlags.Parse(args[2:])
	}

	return 0
}
