package main

import (
	"log"
	"os"
	"path/filepath"
	"strings"
)

var monthsOutputRoot = filepath.Join(os.Getenv("HOME"), "dev", "rotblauer", "cattracks-split-gz", "output")

func main() {

	tippeFiles := []string{}
	filepath.Walk(monthsOutputRoot, func(path string, info os.FileInfo, err error) error {
		if info.IsDir() {
			// skip
			return nil
		}
		if filepath.Ext(path) != ".gz" {
			return nil
		}
		if len(tippeFiles) == 2 {
			return nil
		}
		if !strings.Contains(path, "2023-") {
			return nil
		}
		tippeFiles = append(tippeFiles, path)
		return nil
	})
	
	// run tippecanoe on two months of json.gz files
	for _, f := range tippeFiles {
		out := strings.Replace(f, ".json.gz", ".mbtiles", 1)
		in := f
		name := filepath.Base(f) + "-layer"
		err := runTippe(out, in, name)
		if err != nil {
			log.Fatal(err)
		}
	}
}
