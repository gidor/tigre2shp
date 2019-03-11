package main

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"

	"github.com/gidor/tigre2shp/config"

	"github.com/gen2brain/dlgs"
	"github.com/gidor/tigre2shp/feature"
	"github.com/gidor/tigre2shp/tigre"
	"github.com/jonas-p/go-shp"
)

type Mapping struct {
	In   uint32
	Out  uint32
	Atts map[string]string
}

func LoadMapings() map[uint32]Mapping {
	cfg, ok := config.Get()
	if !ok {
		fmt.Errorf("Errore leggendo configuration")
	}
	out := make(map[uint32]Mapping)
	mappings, ok := config.GetArray(cfg, "Mappings")
	if !ok {
		fmt.Errorf("Errore leggendo Mappings")
	}
	for _, m := range mappings {
		mm := new(Mapping)
		i, _ := config.GetInt(m, "In")
		mm.In = uint32(i)
		i, _ = config.GetInt(m, "Out")
		mm.Out = uint32(i)
		mm.Atts = make(map[string]string)
		atts, _ := config.GetMap(m, "Atts")
		for k := range atts {
			v, _ := config.GetString(atts, k)
			mm.Atts[k] = v
		}
		out[mm.In] = *mm
	}
	return out
}

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

func (d *DataSet) GetShp(id uint32) (*shpFile, bool) {
	i := int(id)
	for _, item := range d.Files {
		if item.descriptor.Code == i {
			return &item, true
		}
	}
	return nil, false
}

func (d *DataSet) closeShp() {
	for _, item := range d.Files {
		item.handler.Close()
	}
	d.Files = d.Files[0:0]
}

func convertvalue(fd feature.FieldDescription, bs string) (setval interface{}) {
	switch fd.Type {
	case feature.Character:
		setval = bs
	case feature.Numeric:
		if fd.Decimals == 0 {
			val, _ := strconv.ParseInt(bs, 10, 32)
			setval = int(val)
		} else {
			val, _ := strconv.ParseFloat(bs, 64)
			setval = float64(val)
		}
	case feature.Date:
		setval = bs
	}
	return
}

func buildshape(gt int, item tigre.Oggetto) (geom shp.Shape) {
	switch gt {
	case feature.Point:
		g := new(shp.Point)
		g.X = item.Punti[0].X
		g.Y = item.Punti[0].Y
		geom = g
	case feature.Polygon:
	case feature.Line:
		ps := make([]shp.Point, 0)
		for _, p := range item.Punti {
			ps = append(ps, shp.Point{X: p.X, Y: p.Y})
		}
		g := shp.NewPolyLine([][]shp.Point{ps})
		geom = g
	default:
		g := new(shp.Point)
		g.X = item.Punti[0].X
		g.Y = item.Punti[0].Y
		geom = g
	}
	return
}

func process(outdataset *DataSet, oggetti []tigre.Oggetto, mappings map[uint32]Mapping) {
	for _, item := range oggetti {
		mapping, ok := mappings[item.TipoOggetto]
		if !ok {
			fmt.Printf("no mapping for item %d\n", item.TipoOggetto)
			continue
		}
		shpfile, ok := outdataset.GetShp(mapping.Out)
		if !ok {
			fmt.Errorf("no shapefile for item %d", mapping.Out)
			continue
		}
		// write geometry
		index := shpfile.handler.Write(buildshape(shpfile.descriptor.GeometryType, item))
		// write data
		r, _ := regexp.Compile("<([0-9]+)>")
		for i, f := range shpfile.descriptor.Attribute {
			val, ok := mapping.Atts[f.FName]
			if ok {
				// ho mapping per l'attributo
				bs := ""
				ss := r.FindAllString(val, 22)
				for i := 0; i < len(ss); i++ {
					seq, _ := strconv.Atoi(strings.Trim(ss[i], "<>"))
					part, ok := item.Attributi[uint32(seq)]
					if ok {
						bs += part
					}
				}
				if len(bs) > 0 {
					shpfile.handler.WriteAttribute(int(index), i, convertvalue(f, bs))
				}
			} else {
				// non ho mapping  verifico default
				if len(f.Default) > 0 {
					shpfile.handler.WriteAttribute(int(index), i, convertvalue(f, f.Default))
				}
			}
		}
	}

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
	mappingss := LoadMapings()
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
	process(outdataset, ogg, mappingss)

	// fmt.Println(ogg)
	// fmt.Println(dirShp)
	outdataset.closeShp()
	// tigre.Test()
	// shp.Open(dirShp)
	fmt.Println("Finito")
}
