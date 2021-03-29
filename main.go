package main

import (
	"flag"
	"log"
	"os"
	"pano2cube/cache"
	"pano2cube/saver"
	"pano2cube/worker"
	"time"
	"path/filepath"
    "strings"
)

var (
	tileNames = []string{
		worker.TileUp,
		worker.TileDown,
		worker.TileFront,
		worker.TileRight,
		worker.TileBack,
		worker.TileLeft,
	}

	tileSize          = 2048
	originalImagePath = ""
	outPutDir         = "./build"

	tileSizeCmd          = flag.Int("s", tileSize, "Size in px of final tile")
	originalImagePathCmd = flag.String("i", "", "Path to input equirectangular panorama")
	outPutDirCmd         = flag.String("o", outPutDir, "Path to output directory")
)

func fileNameWithoutExtension(fileName string) string {
	return strings.TrimSuffix(fileName, filepath.Ext(fileName))
}

func main() {

	flag.Parse()
	tileSize = *tileSizeCmd
	originalImagePath = *originalImagePathCmd
	outPutDir = *outPutDirCmd
    process(originalImagePath, tileSize, outPutDir)
}

func process(originalImagePath string, tileSize int, outPutDir string) {

	if originalImagePath == "" {
		flag.PrintDefaults()
		os.Exit(2)
	}

	_, err := os.Stat(originalImagePath)
	if err != nil {
		if os.IsNotExist(err) {
			log.Fatalf("%v not found", originalImagePath)
			os.Exit(2)
		}
	}

    basename := filepath.Base(originalImagePath)
    prefix := fileNameWithoutExtension(basename)

	done := make(chan worker.TileResult)
	timeStart := time.Now()
	cacheResult := cache.CacheAnglesHandler(tileSize)
	for _, tileName := range tileNames {
		tile := worker.Tile{TileName: tileName, TileSize: tileSize}
		go worker.Worker(tile, cacheResult, originalImagePath, done)
	}

	for range tileNames {
		tileResult := <-done
		//err = saver.SaveTile(tileResult, outPutDir)
        err = saver.SaveTileSlices(tileResult, prefix, outPutDir)

		if err != nil {
			log.Fatal(err.Error())
			os.Exit(2)
		}

		log.Printf("Process for tile %v --> finished", tileResult.Tile.TileName)
	}

	timeFinish := time.Now()
	duration := timeFinish.Sub(timeStart)
	log.Printf("Time to render: %v seconds", duration.Seconds())
}
