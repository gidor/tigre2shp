package config

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"path/filepath"

	_ "github.com/mattn/go-sqlite3"
)

var confdb *sql.DB

func dispose() {
	if confdb != nil {
		confdb.Close()
		confdb = nil
	}
}

func dbinit() {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		log.Fatal(err)
	}
	if confdb == nil {
		db, err := sql.Open("sqlite3", filepath.Join(dir, "config.sqlite"))
		if err != nil {
			log.Fatal(err)
		}
		confdb = db
	}
}

// Feature configurata
type Feature struct {
	Fcode     int64
	Tablename string
	Geometry  string
}

// Attributo dele feature
type Attributo struct {
	Campo        string
	Formato      string
	Lun          int64
	Obbligatorio int8
	Defaultval   string
}

// Features estrae le feature definite
func Features() []Feature {
	res := make([]Feature, 0)
	dbinit()
	// sqlStmt := ` select fcode, tablename, geometry from catalogo where table_type='featureclass' order by fcode `
	rows, err := confdb.Query("select fcode, tablename, geometry from catalogo where table_type='featureclass' order by fcode")
	if err != nil {
		log.Fatal(err)
	}
	for rows.Next() {
		var f Feature
		var (
			fc     sql.NullInt64
			tn, ge sql.NullString
		)
		err = rows.Scan(&fc, &tn, &ge)
		if fc.Valid {
			f.Fcode = fc.Int64
		} else {
			f.Fcode = -1
		}

		if tn.Valid {
			f.Tablename = tn.String
		} else {
			f.Tablename = ""
		}

		if ge.Valid {
			f.Geometry = ge.String
		} else {
			f.Geometry = "Point"
		}

		if err != nil {
			log.Fatal(err)
		}
		res = append(res, f)
	}
	err = rows.Err()
	if err != nil {
		log.Fatal(err)
	}
	return res
}

//Attributi estrai gli attributi di una classe
// :param fcode il feature codeda cercare
func Attributi(fcode int64) []Attributo {
	res := make([]Attributo, 0)
	dbinit()
	sqlStmt := `
	SELECT 
		upper(f.campo) campo,
		f.formato,
		f.lun, 
		f.obbligatorio,
		f.defaultval
	from
		catalogo c
		INNER JOIN campi f on c.fcode = f.fcode  and c.tablename = f.tablename
	where c.table_type ='featureclass' and c.fcode = $1
	order by c.fcode, f.id
	 `
	rows, err := confdb.Query(sqlStmt, fcode)
	if err != nil {
		log.Fatal(err)
	}
	for rows.Next() {
		var a Attributo
		var (
			campo, formato, def sql.NullString
			lun, obb            sql.NullInt64
		)
		err = rows.Scan(&campo, &formato, &lun, &obb, &def)
		if campo.Valid {
			a.Campo = campo.String
		} else {
			a.Campo = ""
		}
		if formato.Valid {
			a.Formato = formato.String
		} else {
			a.Formato = ""
		}

		if lun.Valid {
			a.Lun = lun.Int64
		} else {
			a.Lun = 0
		}

		if obb.Valid {
			a.Obbligatorio = int8(obb.Int64)
		} else {
			a.Obbligatorio = 0
		}
		if def.Valid {
			a.Defaultval = def.String
		} else {
			a.Defaultval = ""
		}

		if err != nil {
			fmt.Println(err)
			log.Fatal(err)
		}
		res = append(res, a)
	}
	err = rows.Err()
	if err != nil {
		log.Fatal(err)
	}
	return res
}

// AttMapping --
type AttMapping struct {
	Attributo string
	Valore    string
}

//Mappings
func AttMappings(fcin uint32, fcout uint32) []AttMapping {
	res := make([]AttMapping, 0)
	dbinit()
	sqlStmt := "SELECT upper(attributo) attributo, valore from mappings where fcin=$1 and fcout=$2"
	rows, err := confdb.Query(sqlStmt, fcin, fcout)
	if err != nil {
		log.Fatal(err)
	}
	for rows.Next() {
		var a AttMapping
		var (
			attributo, valore sql.NullString
		)
		err = rows.Scan(&attributo, &valore)
		if attributo.Valid {
			a.Attributo = attributo.String
		} else {
			a.Attributo = ""
		}
		if valore.Valid {
			a.Valore = valore.String
		} else {
			a.Valore = ""
		}
		if err != nil {
			fmt.Println(err)
			log.Fatal(err)
		}
		res = append(res, a)
	}
	err = rows.Err()
	if err != nil {
		log.Fatal(err)
	}
	return res
}

// FMapping --
type FMapping struct {
	FcIn  int64
	FcOut int64
}

//FMappings
func FMappings() []FMapping {
	res := make([]FMapping, 0)
	dbinit()
	sqlStmt := "SELECT fcin, fcout from mappings GROUP by fcin, fcout order by fcin"
	rows, err := confdb.Query(sqlStmt)
	if err != nil {
		log.Fatal(err)
	}
	for rows.Next() {
		var a FMapping
		var (
			fcin, fcout sql.NullInt64
		)
		err = rows.Scan(&fcin, &fcout)
		if fcin.Valid {
			a.FcIn = fcin.Int64
		} else {
			a.FcIn = 0
		}
		if fcout.Valid {
			a.FcOut = fcout.Int64
		} else {
			a.FcOut = 0
		}
		if err != nil {
			fmt.Println(err)
			log.Fatal(err)
		}
		res = append(res, a)
	}
	err = rows.Err()
	if err != nil {
		log.Fatal(err)
	}
	return res
}

// DefaultValue --
type DefaultValue struct {
	Attributo string
	Valore    string
}

//DefaultValues
func DefaultValues() []DefaultValue {
	res := make([]DefaultValue, 0)
	dbinit()
	sqlStmt := "SELECT upper(attributo) attributo, valore from conf_defaults"
	rows, err := confdb.Query(sqlStmt)
	if err != nil {
		log.Fatal(err)
	}
	for rows.Next() {
		var a DefaultValue
		var (
			attributo, valore sql.NullString
		)
		err = rows.Scan(&attributo, &valore)
		if attributo.Valid {
			a.Attributo = attributo.String
		} else {
			a.Attributo = ""
		}
		if valore.Valid {
			a.Valore = valore.String
		} else {
			a.Valore = ""
		}
		if err != nil {
			fmt.Println(err)
			log.Fatal(err)
		}
		res = append(res, a)
	}
	err = rows.Err()
	if err != nil {
		log.Fatal(err)
	}
	return res
}
