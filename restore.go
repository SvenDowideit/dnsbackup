package main

type RestoreCmd struct {
	Zone string `arg required help:"DNS Zone to restore"`
	//Image pullimage.PullImageCmd `cmd default:"1" help:"Pull a specific image"`

	GandiToken string `help:"Digital Ocean API key" env:"GANDIV5_API_KEY"`
}

func (cmd *RestoreCmd) Run(ctx *Context) error {
	return nil
}
