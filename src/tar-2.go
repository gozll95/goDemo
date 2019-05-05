package main

import (
	"archive/tar"
	"compress/gzip"
	"flag"
	"fmt"
	"io"
	"os"
)

func main() {
	// get the arguments from the command line
	numPtr := flag.Int("n", 4, "an integer")
	flag.Parse()

	sourceFile := flag.Arg(0)

	if sourceFile == "" {
		fmt.Println("Dude, you didn't pass in a tar file!")
		os.Exit(1)
	}

	fmt.Println("arg 1: ", flag.Arg(0))
	fmt.Println("arg n: ", *numPtr)

	processFile(sourceFile, *numPtr)
}

func processFile(srcFile string, num int) {
	f, err := os.Open(srcFile)
	if err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}

	defer f.Close()

	gzf, err := gzip.NewReader(f)
	if err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}

	tarReader := tar.NewReader(gzf)

	i := 0
	for {
		header, err := tarReader.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			fmt.Println(err)
			os.Exit(-1)
		}

		name := header.Name
		fmt.Println("name is ", name)

		switch header.Typeflag {
		case tar.TypeDir:
			continue
		case tar.TypeReg:
			fmt.Println("(", i, ")", "Name: ", name)
			if i == num {
				fmt.Println(" --- ")
				//fmt.Println(tarReader)
				io.Copy(os.Stdout, tarReader)
				fmt.Println(" --- ")
				os.Exit(0)
			}
		default:
			fmt.Printf("%s : %c %s %s\n",
				"Yikes! Unable to figure out type",
				header.Typeflag,
				"in file",
				name,
			)
		}
		i++
	}
}
