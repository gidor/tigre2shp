package main

import (
	"fmt"
	"os"
	"path/filepath"
	"tigre2shp/config"
	"tigre2shp/tigre"

	"github.com/gen2brain/dlgs"
)

func selectDir(msg string) (dir string, err error) {
	ok := false
	for i := 0; i < 3; i++ {
		println(ok)
		println(i)
		dir, ok, err = dlgs.File(msg, "", true)
		if err != nil {
			panic(err)
		}
		if ok {
			return
		}
	}
	if !ok {
		println("3 tenttivi")
		os.Exit(1)
	}
	return
}

func glob(base string) ([]string, error) {
	var files []string
	err := filepath.Walk(base, func(path string, info os.FileInfo, err error) error {
		fmt.Println(path)
		if path == base {
			return nil
		} else if !info.IsDir() {
			files = append(files, path)
		} else {
			return filepath.SkipDir
		}
		return nil
	})
	return files, err
}

func main() {
	conf := config.Get()
	fmt.Println(conf)
	var (
		dirMeta, dirShp string
		err             error
	)
	if len(os.Args) == 1 {
		dirMeta, err = selectDir("Seleziona directory Metafile Tigre")
		if err != nil {
			println(err)
			os.Exit(1)
		}
		dirShp, err = selectDir("Seleziona directory Output")
		if err != nil {
			println(err)
			os.Exit(1)
		}
	} else if len(os.Args) == 2 {
		println("due dir ")
		os.Exit(1)
	} else {
		dirMeta = os.Args[1]
		dirShp = os.Args[2]
	}
	dataset := tigre.NewDataset(dirMeta)
	ogg := dataset.Get()

	fmt.Println(ogg)
	fmt.Println(dirShp)
	// tigre.Test()
	// shp.Open(dirShp)
}
