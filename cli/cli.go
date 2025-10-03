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
	// encIn := flag.Bool("stdin", false, "read from stdin instead of a file")
	// encOut := flag.Bool("stdout", false, "write to stdout instead of a file")
	// encOutputFormat := flag.String("f", "hex", "output format; values: txt(base64), bin")

	decFlags := flag.NewFlagSet("dec", flag.ExitOnError)
	// decIn := flag.Bool("stdin", false, "read from stdin instead of file")
	// decOut := flag.Bool("stdout", false, "write to stdout instead of a file")

	ciphertext := []byte{}
	// plaintext := []byte{}

	switch args[1] {
	case "enc":
		err := encFlags.Parse(args[2:])
		if err != nil {
			log.Printf("error parsing flags: %v\n", err)
			return 1
		}

		f, err := os.Open(filepath.Clean(encFlags.Arg(0)))
		if err != nil {
			log.Printf("error opening file %q: %v\n", encFlags.Arg(0), err)
			return 1
		}
		defer f.Close()

		data, err := io.ReadAll(f)
		if err != nil {
			log.Printf("error reading file: %v\n", err)
			return 1
		}

		sc := bufio.NewScanner(os.Stdin)
		fmt.Print("password for encrypting >> ")
		var password string
		for sc.Scan() {
			password = strings.TrimSpace(sc.Text())
			if password != "" {
				break
			}
		}
		ciphertext, err = pkg.EncryptWithPassword(data, []byte(password))
		if err != nil {
			log.Printf("error encrypting: %v\n", err)
			return 1
		}

		outputFileName := f.Name() + ".enc"
		outputFile, err := os.Create(outputFileName)
		if err != nil {
			log.Printf("error making outputfile: %v\n", err)
			return 1
		}
		defer outputFile.Close()

		n, err := outputFile.Write(ciphertext)
		if err != nil {
			log.Printf("error writing to file: %v\n", err)
			return 1
		}
		fmt.Printf("Wrote %d bytes to %s\n", n, outputFileName)

	case "dec":
		decFlags.Parse(args[2:])
	}

	return 0
}
