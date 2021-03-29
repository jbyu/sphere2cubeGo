package saver

import (
	"os"
	"image"
	"image/jpeg"
	"path/filepath"
	"pano2cube/worker"
	"github.com/disintegration/imaging"
    "fmt"
)

const (
    SliceSize = 512
)

func SaveTile(tileResult worker.TileResult, outPutDir string) error {

	err := os.Mkdir(outPutDir, os.FileMode(os.ModePerm))

	if err != nil {
		if !os.IsExist(err) {
			return err
		}
	}

	finalPath := filepath.Join(outPutDir, tileResult.Tile.TileName+".jpg")

	f, err := os.Create(finalPath)

	if err != nil {
		return err
	}

	defer f.Close()

	err = jpeg.Encode(f, tileResult.Image, &jpeg.Options{100})

	if err != nil {
		return err
	}

	return nil
}


func SaveTileSlices(tileResult worker.TileResult, prefix string, outPutDir string) error {

	err := os.Mkdir(outPutDir, os.FileMode(os.ModePerm))

	if err != nil {
		if !os.IsExist(err) {
			return err
		}
	}

    rgbImg := tileResult.Image.(*image.NRGBA)
    size := tileResult.Tile.TileSize / 1024
    
    // level of detail: 2k -> 1k -> 512
    for k := size; k >= 0; k-- {
        bound := rgbImg.Bounds()
        width := bound.Dx()
        height := bound.Dy()
        count := 0
    
        for i := 0; i < height; i+= SliceSize {
            for j := 0; j < width; j+= SliceSize {
                filename := fmt.Sprintf("%s_%dK_%s_%d.jpg", prefix, k, tileResult.Tile.TileName, count)
                finalPath := filepath.Join(outPutDir, filename)
                sub_image := rgbImg.SubImage(image.Rect(j, i, j + SliceSize, i + SliceSize))
                err = imaging.Save(sub_image, finalPath)
                if err != nil {
                    return err
                }
                count++
            }
        }
        //downsampling
        rgbImg = imaging.Resize(rgbImg, width/2, 0, imaging.Linear)
    }
    
	return nil
}