package main

import (
	"github.com/alecthomas/kong"
)

type Context struct {
	DryRun bool
}

type Cli struct {
	DryRun bool `description:"Don't actually download or update, just show what it would change"`

	Backup  BackupCmd  `cmd help:"backup DNS from DigitalOcean API"`
	Restore RestoreCmd `cmd help:"restore DNS to Gandi API"`
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
