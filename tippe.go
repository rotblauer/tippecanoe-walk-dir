package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"time"
)

type prefixedWriter struct {
	*log.Logger
}

func (pw prefixedWriter) Write(p []byte) (n int, err error) {
	pw.Logger.Println(string(p))
	return len(p), nil
}

func runTippe(out, in string, tilesetname string) error {
	tippCmd, tippargs, tipperr := getTippyProcess(out, in, tilesetname)
	if tipperr != nil {
		return tipperr
	}

	start := time.Now()
	defer func() {
		log.Printf("Finished tippecanoe on %s: %s\n", out, time.Since(start))
	}()

	// log.Println("> [", tilesetname, "]", tippCmd, tippargs)
	tippmycanoe := exec.Command(tippCmd, tippargs...)
	tippmycanoe.Stdout = os.Stdout
	tippmycanoe.Stderr = os.Stderr

	prefix := fmt.Sprintf("[%s] ", tilesetname)
	stdout := prefixedWriter{log.New(os.Stdout, prefix, log.LstdFlags|log.Lmsgprefix)}
	tippmycanoe.Stdout = stdout
	stderr := prefixedWriter{log.New(os.Stderr, prefix, log.LstdFlags|log.Lmsgprefix)}
	tippmycanoe.Stderr = stderr

	err := tippmycanoe.Start()
	if err != nil {
		log.Println("Error starting Cmd", err)
		os.Exit(1)
	}

	if err := tippmycanoe.Wait(); err != nil {
		return err
	}
	return nil
}

func getTippyProcess(out string, in string, tilesetname string) (tippCmd string, tippargs []string, err error) {
	tippCmd = "/usr/local/bin/tippecanoe"
	tippargs = []string{
		// "--maximum-tile-bytes", "500000", // num bytes/tile,default: 500kb=500000
		"--maximum-tile-bytes", fmt.Sprintf("%d", 350_000), // num bytes/tile,default: 500kb=500000
		"--cluster-densest-as-needed",
		"--cluster-distance=1",
		"--calculate-feature-density",
		"--include", "UnixTime",
		"--include", "Activity",
		"--include", "Elevation",
		"--include", "Speed",
		"--include", "Accuracy",
		"-EUnixTime:max",
		"-EElevation:max",
		"-ESpeed:max", // mean",
		"-EAccuracy:mean",
		"-EPressure:mean",
		"--single-precision",
		"-r1", // == --drop-rate
		"--minimum-zoom", "3",
		"--maximum-zoom", "18",
		"--json-progress", "--progress-interval", "30",
		"-l", tilesetname, // TODO: what's difference layer vs name?
		"-n", tilesetname,
		"-o", out,
		"--force",
		"--read-parallel", in,
	}

	// 'in' should be an existing file
	_, err = os.Stat(in)
	if err != nil {
		return
	}

	// Use alternate tippecanoe path if 'bash -c which tippecanoe' returns something without error and different than default
	if b, e := exec.Command("bash -c", "which", "tippecanoe").Output(); e == nil && string(b) != tippCmd {
		tippCmd = string(b)
	}
	return
}

func runTileJoin(in1, in2, out string) error {
	tjCmd := "/usr/local/bin/tile-join"
	tjargs := []string{
		"-o", out,
		"--force", // remove out.mbtiles if it already exists
		in1, in2,
	}
	log.Println("> [tile-join]", tjCmd, tjargs)
	tj := exec.Command(tjCmd, tjargs...)
	tj.Stdout = os.Stdout
	tj.Stderr = os.Stderr

	err := tj.Start()
	if err != nil {
		log.Println("Error starting Cmd", err)
		os.Exit(1)
	}
	if err := tj.Wait(); err != nil {
		return err
	}
	return nil
}
