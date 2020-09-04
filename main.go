package main

import (
	"context"
	"fmt"
	"time"

	"github.com/alecthomas/kong"
	"github.com/libdns/digitalocean"
	"github.com/libdns/libdns"
)

type Context struct {
	DryRun bool
}

type Cli struct {
	DryRun bool `description:"Don't actually download or update, just show what it would change"`

	Backup  BackupCmd  `cmd help:"backup DNS from DigitalOcean API"`
	Restore RestoreCmd `cmd help:"restore DNS to Gandi API"`
}

type BackupCmd struct {
	Zone string `arg required help:"DNS Zone to backup"`
	//Image pullimage.PullImageCmd `cmd default:"1" help:"Pull a specific image"`

	//DryRun bool `description:"Don't actually download or update, just show what it would change"`
	DoToken string `help:"Digital Ocean API key" env:"DO_API_KEY"`
}

type RestoreCmd struct {
	Zone string `arg required help:"DNS Zone to restore"`
	//Image pullimage.PullImageCmd `cmd default:"1" help:"Pull a specific image"`

	GandiToken string `help:"Digital Ocean API key" env:"GANDIV5_API_KEY"`
}

func (cmd *BackupCmd) Run(ctx *Context) error {
	provider := digitalocean.Provider{APIToken: cmd.DoToken}

	records, err := provider.GetRecords(context.TODO(), cmd.Zone)
	if err != nil {
		fmt.Printf("ERROR: %s\n", err.Error())
	}

	testName := "libdns-test"
	testId := ""
	for _, record := range records {
		fmt.Printf("%s (.%s): %s, %s, %d\n", record.Name, cmd.Zone, record.Value, record.Type, record.TTL)
		if record.Name == testName {
			testId = record.ID
		}

	}

	if testId != "" {
		// fmt.Printf("Delete entry for %s (id:%s)\n", testName, testId)
		// _, err = provider.DeleteRecords(context.TODO(), zone, []libdns.Record{libdns.Record{
		// 	ID: testId,
		// }})
		// if err != nil {
		// 	fmt.Printf("ERROR: %s\n", err.Error())
		// }
		// Set only works if we have a record.ID
		fmt.Printf("Replacing entry for %s\n", testName)
		_, err = provider.SetRecords(context.TODO(), cmd.Zone, []libdns.Record{libdns.Record{
			Type:  "TXT",
			Name:  testName,
			Value: fmt.Sprintf("Replacement test entry created by libdns %s", time.Now()),
			TTL:   time.Duration(30) * time.Second,
			ID:    testId,
		}})
		if err != nil {
			fmt.Printf("ERROR: %s\n", err.Error())
		}
	} else {
		fmt.Printf("Creating new entry for %s\n", testName)
		_, err = provider.AppendRecords(context.TODO(), cmd.Zone, []libdns.Record{libdns.Record{
			Type:  "TXT",
			Name:  testName,
			Value: fmt.Sprintf("This is a test entry created by libdns %s", time.Now()),
			TTL:   time.Duration(30) * time.Second,
		}})
		if err != nil {
			fmt.Printf("ERROR: %s\n", err.Error())
		}
	}
	return nil
}

func (cmd *RestoreCmd) Run(ctx *Context) error {
	return nil
}

func main() {
	cli := &Cli{}
	ctx := kong.Parse(cli,
		kong.Name("dnsbackup"),
		kong.Description("backup and restore dns using libdns"),
		kong.UsageOnError(),
		kong.ConfigureHelp(kong.HelpOptions{
			Compact: true,
		}),
		kong.Vars{})
	//options.LoadAllEnvs()
	err := ctx.Run(&Context{DryRun: cli.DryRun})
	ctx.FatalIfErrorf(err)
}
