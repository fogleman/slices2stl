package main

import (
	"fmt"
	"image"
	_ "image/jpeg"
	_ "image/png"
	"os"

	"github.com/fogleman/fauxgl"
	"github.com/fogleman/image3d"
	"github.com/fogleman/mc"
)

const (
	xyScale     = 1
	zScale      = 1
	zStretch    = 2
	mcThreshold = 0.5
)

func loadImage(path string) (image.Image, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	im, _, err := image.Decode(file)
	return im, err
}

func trianglesToMesh(triangles []mc.Triangle) *fauxgl.Mesh {
	triangles2 := make([]*fauxgl.Triangle, len(triangles))
	for i, t := range triangles {
		p1 := fauxgl.Vector{t.V1.X, t.V1.Y, t.V1.Z * zStretch}
		p2 := fauxgl.Vector{t.V2.X, t.V2.Y, t.V2.Z * zStretch}
		p3 := fauxgl.Vector{t.V3.X, t.V3.Y, t.V3.Z * zStretch}
		triangles2[i] = fauxgl.NewTriangleForPoints(p1, p2, p3)
	}
	return fauxgl.NewTriangleMesh(triangles2)
}

type Evaluator struct {
	Slices *image3d.Image3D
}

func (e *Evaluator) Evaluate(x, y, z float64) float64 {
	c := e.Slices.At(x/xyScale, y/xyScale, z/zScale)
	return float64(c.R) / 65535
}

func main() {
	var images []image.Image
	for _, path := range os.Args[1:] {
		fmt.Println(path)
		im, err := loadImage(path)
		if err != nil {
			panic(err)
		}
		images = append(images, im)
	}
	w := float64(images[0].Bounds().Size().X) * xyScale
	h := float64(images[0].Bounds().Size().Y) * xyScale
	d := float64(len(images)-1) * zScale
	im := image3d.NewImage3D(images)
	e := &Evaluator{im}
	triangles := mc.MarchingCubes(e, 0, 0, -1, w, h, d+1, 1, 1, 1, mcThreshold)
	mesh := trianglesToMesh(triangles)
	mesh.SaveSTL("out.stl")
}
