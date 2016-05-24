package main

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/attic-labs/noms/clients/go/flags"
	"github.com/attic-labs/noms/clients/go/util"
	"github.com/attic-labs/noms/d"
)

func main() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "usage: %s <url> <dataset>\n", os.Args[0])
		flag.PrintDefaults()
	}

	flags.RegisterDatabaseFlags()
	flag.Parse()

	if len(flag.Args()) != 2 {
		util.CheckError(errors.New("expected url and dataset flags"))
	}

	spec, err := flags.ParseDatasetSpec(flag.Arg(1))
	util.CheckError(err)
	ds, err := spec.Dataset()
	util.CheckError(err)

	url := flag.Arg(0)
	if url == "" {
		flag.Usage()
	}

	res, err := http.Get(url)
	if err != nil {
		log.Fatalf("Error fetching %s: %+v\n", url, err)
	} else if res.StatusCode != 200 {
		log.Fatalf("Error fetching %s: %s\n", url, res.Status)
	}
	defer res.Body.Close()

	var jsonObject interface{}
	err = json.NewDecoder(res.Body).Decode(&jsonObject)
	if err != nil {
		log.Fatalln("Error decoding JSON: ", err)
	}

	_, err = ds.Commit(util.NomsValueFromDecodedJSON(jsonObject))
	d.Exp.NoError(err)
	ds.Database().Close()
}
