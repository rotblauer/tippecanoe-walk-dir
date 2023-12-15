package main

import (
	"flag"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/fsnotify/fsnotify"
)

var _source = filepath.Join(os.Getenv("HOME"), "dev", "rotblauer", "cattracks-split-cats-uniqcell-gz", "output")
var flagSourceDir = flag.String("source", _source, "Source directory containing .json.gz files")
var flagOutputRootFilepath = flag.String("output", filepath.Join(_source, "mbtiles"), "Output root dir")
var flagForce = flag.Bool("force", false, "Force overwrite of existing files (otherwise skip if .mbtiles is newer than .json.gz)")
var flagEnableFSWatch = flag.Bool("enable-fs-watch", false, "Enable watching of source directory for changes")

// walkDirRunTippe run tippecanoe to create .mbtiles for each .json.gz file that doesn't yet have one.
// Without --force it may be run idempotently, skipping all cases where file.mbiles is newer than file.json.gz.
// The layer name is the base filename of the .json.gz file, suffixed with "-layer", eg 'ia.json.gz-layer'.
func walkDirRunTippe(dir string, changedPath string) {
	log.Println("source dir:", dir)
	filepath.Walk(dir, func(path string, jsonGZInfo os.FileInfo, err error) error {
		if err != nil {
			log.Println("error walking path:", path, err)
			return nil
		}
		if jsonGZInfo.IsDir() {
			return nil
		}
		if filepath.Ext(path) != ".gz" {
			return nil
		}

		tilesPath := strings.Replace(path, ".json.gz", ".mbtiles", 1)
		out := filepath.Join(*flagOutputRootFilepath, filepath.Base(tilesPath))

		// --force?
		if !*flagForce {
			if tilesInfo, err := os.Stat(out); err == nil {
				// if the modtime of the .json.gz is older than the .mbtiles, skip
				if jsonGZInfo.ModTime().Before(tilesInfo.ModTime()) {
					log.Printf("skipping %s, already exists and is fresh: %s\n", path, tilesPath)
					return nil
				}
			}
		}

		log.Println("found file:", path)

		f := path
		log.Println("running tippecanoe on:", f)
		in := f
		name := filepath.Base(f) + "-layer"
		if err := runTippe(out, in, name); err != nil {
			log.Println(err)
		}
		return nil
	})
}

func main() {
	flag.Parse()
	os.MkdirAll(*flagOutputRootFilepath, 0755)
	walkDirRunTippe(*flagSourceDir, "-")

	if !*flagEnableFSWatch {
		return
	}

	// Create new watcher.
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}
	defer watcher.Close()

	// Start listening for events.
	go func() {
		for {
			select {
			case event, ok := <-watcher.Events:
				if !ok {
					return
				}
				log.Println("event:", event)
				if event.Has(fsnotify.Write) {
					log.Println("modified file:", event.Name)
					walkDirRunTippe(*flagSourceDir, event.Name)
				}
			case err, ok := <-watcher.Errors:
				if !ok {
					return
				}
				log.Println("error:", err)
			}
		}
	}()

	// Add a path.
	err = watcher.Add(*flagSourceDir)
	if err != nil {
		log.Fatal(err)
	}

	// Block main goroutine forever.
	<-make(chan struct{})
}
