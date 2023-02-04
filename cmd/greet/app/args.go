package app

type Args struct {
	Names      []string `arg:"positional,required" help:"names to greet"`
	Seperately bool     `arg:"-s,--seperately" help:"greet each name seperately" default:"false"`
	Verbose    bool     `arg:"-v,--verbose" help:"verbose output" default:"false"`
	Config     string   `arg:"-c,--config" help:"config file" default:"config.toml"`
}
