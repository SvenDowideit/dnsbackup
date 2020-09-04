package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/libdns/gandi"
	"github.com/libdns/libdns"

	"github.com/onaci/docker-ona/outputs"
)

type RestoreCmd struct {
	Zone     string `arg required help:"DNS Zone to restore"`
	Filename string `arg required help:"json formated DNS file to restore"`

	GandiToken string `help:"Digital Ocean API key" env:"GANDIV5_API_KEY"`
	Format     string `default:"| {{.Name}} | {{.Value}} | {{.Type}} | {{.TTL}} |" help:"output format (json, spew, printf)"`
}

func (cmd *RestoreCmd) Run(ctx *Context) error {
	// Open our jsonFile
	jsonFile, err := os.Open("ona.im.json")
	// if we os.Open returns an error then handle it
	if err != nil {
		fmt.Println(err)
		return err
	}
	// defer the closing of our jsonFile so that we can parse it later on
	defer jsonFile.Close()
	byteValue, _ := ioutil.ReadAll(jsonFile)

	var newRecords []libdns.Record
	if err := json.Unmarshal(byteValue, &newRecords); err != nil {
		fmt.Println(err)
		return err
	}

	fmt.Printf("\n\nrecords read from %s\n", cmd.Filename)
	err = printRecords(os.Stdout, newRecords, cmd.Format)
	if err != nil {
		return err
	}

	provider := gandi.Provider{APIToken: cmd.GandiToken}
	records, err := provider.GetRecords(context.TODO(), cmd.Zone)
	if err != nil {
		fmt.Printf("ERROR: %s\n", err.Error())
		return err
	}

	fmt.Printf("\n\nrecords read from %s using API\n", cmd.Zone)
	err = printRecords(os.Stdout, records, cmd.Format)
	if err != nil {
		return err
	}

	// need to update the newRecord ID's to use the destination ID, or empty string if it doesn't exist
	for _, nr := range newRecords {
		nr.ID = ""
		for _, v := range records {
			if v.Name == nr.Name {
				nr.ID = v.ID
				break
			}
		}
		if !ctx.DryRun {
			err = addOrUpdateRecord(&provider, cmd.Zone, nr)
			if err != nil {
				fmt.Printf("ERROR: %s\n", err.Error())
			}
		}
	}

	return nil
}

func printRecords(out *os.File, records []libdns.Record, format string) error {
	res := make([]interface{}, len(records))
	for i, v := range records {
		res[i] = v
	}
	err := outputs.FormatArray(out, res, "", format)
	if err != nil {
		return err
	}
	return nil
}

func addOrUpdateRecord(provider *gandi.Provider, zone string, r libdns.Record) error {
	if r.ID != "" {
		fmt.Printf("Replacing entry for %s\n", r.Name)
		_, err := provider.SetRecords(context.TODO(), zone, []libdns.Record{r})
		if err != nil {
			fmt.Printf("ERROR: %s\n", err.Error())
			return err
		}
	} else {
		fmt.Printf("Creating new entry for %s\n", r.Name)
		_, err := provider.AppendRecords(context.TODO(), zone, []libdns.Record{r})
		if err != nil {
			fmt.Printf("ERROR: %s\n", err.Error())
			return err
		}
	}
	return nil
}
