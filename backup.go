package main

import (
	"context"
	"fmt"
	"os"

	"github.com/libdns/digitalocean"
)

type BackupCmd struct {
	Zone string `arg required help:"DNS Zone to backup"`
	//Image pullimage.PullImageCmd `cmd default:"1" help:"Pull a specific image"`

	//DryRun bool `description:"Don't actually download or update, just show what it would change"`
	DoToken string `help:"Digital Ocean API key" env:"DO_API_KEY"`

	Format string `default:"json" help:"output format (json, spew, printf)"`
}

func (cmd *BackupCmd) Run(ctx *Context) error {
	provider := digitalocean.Provider{APIToken: cmd.DoToken}

	records, err := provider.GetRecords(context.TODO(), cmd.Zone)
	if err != nil {
		fmt.Printf("ERROR: %s\n", err.Error())
	}

	// testName := "libdns-test"
	// testId := ""
	// for _, record := range records {
	// 	fmt.Printf("%s (.%s): %s, %s, %d\n", record.Name, cmd.Zone, record.Value, record.Type, record.TTL)
	// 	if record.Name == testName {
	// 		testId = record.ID
	// 	}

	// }
	err = printRecords(os.Stdout, records, cmd.Format)
	if err != nil {
		return err
	}

	return nil
}
