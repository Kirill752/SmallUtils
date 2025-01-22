package main

import (
	"fmt"
	"io"
	"os"
	"strconv"
)

func printDir(out io.Writer, path string, printFiles bool, prefix string) error {
	files, err := os.ReadDir(path)
	if err != nil {
		return err
	}
	lastDir := len(files) - 1
	for ; lastDir > -1; lastDir-- {
		if files[lastDir].IsDir() {
			break
		}
	}
	for i, file := range files {
		if file.IsDir() {
			if i == len(files)-1 {
				fmt.Fprintf(out, "%s└───%s\n", prefix, file.Name())
			} else {
				if prefix == "" && i == lastDir {
					fmt.Fprintf(out, "%s└───%s\n", prefix, file.Name())
				} else {
					fmt.Fprintf(out, "%s├───%s\n", prefix, file.Name())
				}
			}
			newPref := "│\t"
			if i == len(files)-1 || prefix == "" && i == lastDir {
				newPref = "\t"
			}
			err = printDir(out, path+"/"+file.Name(), printFiles, prefix+newPref)
		}
	}
	return err
}

func printAll(out io.Writer, path string, printFiles bool, prefix string) error {
	files, err := os.ReadDir(path)
	if err != nil {
		return err
	}
	for i, file := range files {
		inf, _ := file.Info()
		size := "empty"
		if inf.Size() > 0 {
			size = strconv.Itoa(int(inf.Size())) + "b"
		}
		if i == len(files)-1 {
			if !file.IsDir() {
				fmt.Fprintf(out, "%s└───%s (%s)\n", prefix, file.Name(), size)
			} else {
				fmt.Fprintf(out, "%s└───%s\n", prefix, file.Name())
			}
		} else {
			if !file.IsDir() {

				fmt.Fprintf(out, "%s├───%s (%s)\n", prefix, file.Name(), size)
			} else {
				fmt.Fprintf(out, "%s├───%s\n", prefix, file.Name())
			}
		}
		if file.IsDir() {
			newPref := "│\t"
			if i == len(files)-1 {
				newPref = "\t"
			}
			err = printAll(out, path+"/"+file.Name(), printFiles, prefix+newPref)
		}
	}
	return err
}

func dirTree(out io.Writer, path string, printFiles bool) error {
	var err error
	if printFiles {
		err = printAll(out, path, printFiles, "")
	} else {
		err = printDir(out, path, printFiles, "")
	}
	return err
}
