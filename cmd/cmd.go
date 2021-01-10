package cmd

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"os/signal"
	"path/filepath"
	"runtime"
	"runtime/pprof"
	"sync"
	"syscall"
	"time"

	"github.com/choria-io/go-choria/filter"
	"github.com/choria-io/go-choria/protocol"

	"github.com/choria-io/go-choria/build"
	"github.com/choria-io/go-choria/choria"
	"github.com/choria-io/go-choria/config"
	log "github.com/sirupsen/logrus"
	"gopkg.in/alecthomas/kingpin.v2"
)

type application struct {
	app      *kingpin.Application
	command  string
	commands []runableCmd
}

var (
	cli        = application{}
	debug      = false
	configFile = ""
	c          *choria.Framework
	cfg        *config.Config
	ctx        context.Context
	cancel     func()
	wg         *sync.WaitGroup
	mu         = &sync.Mutex{}
	cpuProfile string
	bi         *build.Info
	err        error
)

type stdFilterOptions struct {
	collective string
	factF      []string
	agentsF    []string
	classF     []string
	identityF  []string
	combinedF  []string
	compoundF  string
	dm         string
	dt         int
}

func ParseCLI() (err error) {
	bi = &build.Info{}

	cli.app = kingpin.New("choria", "Choria Orchestration System")
	cli.app.Version(bi.Version())
	cli.app.Author("R.I.Pienaar <rip@devco.net>")

	cli.app.Flag("debug", "Enable debug logging").Short('d').BoolVar(&debug)
	cli.app.Flag("config", "Config file to use").PlaceHolder("FILE").ExistingFileVar(&configFile)
	cli.app.Flag("profile", "Enable CPU profiling and write to the supplied file").Hidden().StringVar(&cpuProfile)

	for _, cmd := range cli.commands {
		err = cmd.Setup()
	}

	cli.command = kingpin.MustParse(cli.app.Parse(os.Args[1:]))

	for _, cmd := range cli.commands {
		if cmd.FullCommand() == cli.command {
			err = cmd.Configure()
			if err != nil {
				return fmt.Errorf("%s failed to configure: %s", cmd.FullCommand(), err)
			}
		}
	}

	return
}

func commonConfigure() error {
	if debug {
		log.SetOutput(os.Stdout)
		log.SetLevel(log.DebugLevel)
		log.Debug("Logging at debug level due to CLI override")
	}

	if configFile == "" {
		configFile = choria.UserConfig()
	}

	cfg, err = config.NewConfig(configFile)
	if err != nil {
		return fmt.Errorf("could not parse configuration: %s", err)
	}

	cfg.ApplyBuildSettings(bi)

	if os.Getenv("INSECURE_ANON_TLS") == "true" {
		cfg.Choria.ClientAnonTLS = true
		cfg.DisableTLSVerify = true
		cfg.DisableSecurityProviderVerify = true
		log.Warn("Using anonymous TLS via environment override")
	}

	if os.Getenv("INSECURE_DISABLE_TLS") == "true" {
		cfg.DisableTLS = true
		log.Warn("Disabling TLS via environment override")
	}

	if os.Getenv("INSECURE_YES_REALLY") == "true" {
		protocol.Secure = "false"
		cfg.DisableTLS = true
		log.Warn("Disabling protocol security via environment override")
	}

	return nil
}

func Run() (err error) {
	wg = &sync.WaitGroup{}
	ran := false

	ctx, cancel = context.WithCancel(context.Background())
	defer cancel()

	go interruptWatcher()

	if cpuProfile != "" {
		cpf, err := os.Create(cpuProfile)
		if err != nil {
			return fmt.Errorf("could not setup profiling: %s", err)
		}
		defer cpf.Close()

		err = pprof.StartCPUProfile(cpf)
		if err != nil {
			return fmt.Errorf("could not setup profiling: %s", err)
		}

		defer pprof.StopCPUProfile()
	}

	if cfg != nil {
		// we do this here so that the command setup has a chance to fiddle the config for
		// things like disabling full verification of the security system during enrollment
		c, err = choria.NewWithConfig(cfg)
		if err != nil {
			return fmt.Errorf("could not initialize Choria: %s", err)
		}

		c.SetupLogging(debug)
	}

	for _, cmd := range cli.commands {
		if cmd.FullCommand() == cli.command {
			ran = true

			wg.Add(1)
			err = cmd.Run(wg)
		}
	}

	if !ran {
		err = fmt.Errorf("could not run the CLI: Invalid command %s", cli.command)
	}

	if err != nil {
		cancel()
	}

	wg.Wait()

	return
}

func forcequit() {
	grace := 2 * time.Second

	if cfg != nil {
		if cfg.SoftShutdownTimeout > 0 {
			grace = time.Duration(cfg.SoftShutdownTimeout) * time.Second
		}
	}

	<-time.NewTimer(grace).C

	dumpGoRoutines()

	log.Errorf("Forced shutdown triggered after %v", grace)

	os.Exit(1)
}

func interruptWatcher() {
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	for {
		select {
		case sig := <-sigs:
			switch sig {
			case syscall.SIGINT, syscall.SIGTERM:
				go forcequit()

				log.Infof("Shutting down on %s", sig)
				cancel()

			case syscall.SIGQUIT:
				dumpGoRoutines()
			}
		case <-ctx.Done():
			return
		}
	}
}

func dumpGoRoutines() {
	mu.Lock()
	defer mu.Unlock()

	now := time.Now().UnixNano()
	pid := os.Getpid()

	tdoutname := filepath.Join(os.TempDir(), fmt.Sprintf("choria-threaddump-%d-%d.txt", pid, now))
	memoutname := filepath.Join(os.TempDir(), fmt.Sprintf("choria-memoryprofile-%d-%d.mprof", pid, now))

	buf := make([]byte, 1<<20)
	stacklen := runtime.Stack(buf, true)

	err := ioutil.WriteFile(tdoutname, buf[:stacklen], 0644)
	if err != nil {
		log.Errorf("Could not produce thread dump: %s", err)
		return
	}

	log.Warnf("Produced thread dump to %s", tdoutname)

	mf, err := os.Create(memoutname)
	if err != nil {
		log.Errorf("Could not produce memory profile: %s", err)
		return
	}
	defer mf.Close()

	err = pprof.WriteHeapProfile(mf)
	if err != nil {
		log.Errorf("Could not produce memory profile: %s", err)
		return
	}

	log.Warnf("Produced memory profile to %s", memoutname)
}

// digs in the application.commands structure looking for a entry with
// the given command string
func cmdWithFullCommand(command string) (cmd runableCmd, ok bool) {
	for _, cmd := range cli.commands {
		if cmd.FullCommand() == command {
			return cmd, true
		}
	}

	return cmd, false
}

func addStdDiscovery(app *kingpin.CmdClause, opts *stdFilterOptions) {
	app.Flag("dm", "Sets a discovery method (broadcast, choria)").EnumVar(&opts.dm, "broadcast", "choria", "mc")
	app.Flag("discovery-timeout", "Timeout for doing discovery").PlaceHolder("SECONDS").IntVar(&opts.dt)
}

func addStdFilter(app *kingpin.CmdClause, opts *stdFilterOptions) {
	app.Flag("wf", "Match hosts with a certain fact").Short('F').StringsVar(&opts.factF)
	app.Flag("wc", "Match hosts with a certain configuration management class").Short('C').StringsVar(&opts.classF)
	app.Flag("wa", "Match hosts with a certain Choria agent").Short('A').StringsVar(&opts.agentsF)
	app.Flag("wi", "Match hosts with a certain Choria identity").Short('I').StringsVar(&opts.identityF)
	app.Flag("with", "Combined classes and facts filter").Short('W').PlaceHolder("FILTER").StringsVar(&opts.combinedF)
	app.Flag("select", "Match hosts using a expr compound filter").Short('S').PlaceHolder("EXPR").StringVar(&opts.compoundF)
	app.Flag("target", "Target a specific sub collective").Short('T').StringVar(&opts.collective)
}

func (o *stdFilterOptions) setDefaults(collective string, dm string, dt int) {
	if o.dm == "" {
		o.dm = dm
	}

	if o.collective == "" {
		o.collective = collective
	}

	if o.dt == 0 {
		o.dt = dt
	}

	if len(o.compoundF) > 0 {
		o.dm = "broadcast"
	}
}

func (o *stdFilterOptions) newFilter(agent string) (*protocol.Filter, error) {
	f, err := filter.NewFilter(
		filter.FactFilter(o.factF...),
		filter.AgentFilter(o.agentsF...),
		filter.ClassFilter(o.classF...),
		filter.IdentityFilter(o.identityF...),
		filter.CombinedFilter(o.combinedF...),
		filter.CompoundFilter(o.compoundF),
	)
	if err != nil {
		return nil, err
	}

	if agent != "" {
		f.AddAgentFilter(agent)
	}

	return f, nil
}
