package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/gen2brain/dlgs"
	"github.com/gidor/tigre2shp/feature"
	"github.com/gidor/tigre2shp/tigre"
	"github.com/jonas-p/go-shp"
)

type shpFile struct {
	descriptor feature.FeatureDescription
	filename   string
	handler    *shp.Writer
}

type DataSet struct {
	Path  string
	Files []shpFile
}

// NewDataset
func NewDataset(path string) *DataSet {
	out := new(DataSet)
	out.Path = path
	out.Files = make([]shpFile, 0)
	return out
}

// SetFeatures
func (d *DataSet) SetFeatures(features []feature.FeatureDescription) {
	for _, descr := range features {
		d.addShp(descr)
	}
}

func (d *DataSet) addShp(desc feature.FeatureDescription) {
	var (
		gt      shp.ShapeType
		shpItem shpFile
	)

	switch desc.GeometryType {
	case feature.Point:
		gt = shp.POINT
	case feature.Polygon:
		gt = shp.POLYGON
	case feature.Line:
		gt = shp.POLYLINE
	default:
		gt = shp.POINT
	}
	shpItem.descriptor = desc
	shpItem.filename = filepath.Join(d.Path, desc.Name+".shp")
	shpItem.handler, _ = shp.Create(shpItem.filename, gt)
	// fields to write
	fields := make([]shp.Field, 0)
	for _, att := range desc.Attribute {
		switch att.Type {
		case feature.Character:
			fields = append(fields, shp.StringField(att.FName, uint8(att.Len)))
		case feature.Numeric:
			if att.Decimals == 0 {
				fields = append(fields, shp.NumberField(att.FName, uint8(att.Len)))
			} else {
				fields = append(fields, shp.FloatField(att.FName, uint8(att.Len), uint8(att.Len)))
			}
		case feature.Date:
			fields = append(fields, shp.DateField(att.FName))
		case feature.Logical:
			fields = append(fields, shp.NumberField(att.FName, 1))
		default:
			fields = append(fields, shp.StringField(att.FName, 20))
		}
	}
	shpItem.handler.SetFields(fields)
	d.Files = append(d.Files, shpItem)
}

func (d *DataSet) closeShp() {
	for _, item := range d.Files {
		item.handler.Close()
	}
	d.Files = d.Files[0:0]
}

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
	// conf, _ := config.Get()
	// fmt.Println(conf)
	var (
		dirMeta, dirShp string
		err             error
	)
	defs := feature.Load()
	// fmt.Println(defs)
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
	outdataset := NewDataset(dirShp)
	outdataset.SetFeatures(defs)
	dataset := tigre.NewDataset(dirMeta)
	ogg := dataset.Get()

	fmt.Println(ogg)
	fmt.Println(dirShp)
	outdataset.closeShp()
	// tigre.Test()
	// shp.Open(dirShp)
}
