package app

type Args struct {
	Host     string `arg:"-H,--host" help:"host address"`
	Port     int    `arg:"-p,--port" help:"port"`
	Username string `arg:"-u,--user" help:"username"`
	Verbose  bool   `arg:"-v,--verbose" help:"verbose output"`
	Config   string `arg:"-c,--config" help:"config file"`
}
