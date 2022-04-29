// (c) 2021, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"sync"
	"syscall"
	"time"

	"github.com/chain4travel/caminogo/utils/logging"
	"github.com/chain4travel/magellan/api"
	"github.com/chain4travel/magellan/balance"
	"github.com/chain4travel/magellan/cfg"
	"github.com/chain4travel/magellan/db"
	"github.com/chain4travel/magellan/models"
	"github.com/chain4travel/magellan/replay"
	oreliusRpc "github.com/chain4travel/magellan/rpc"
	"github.com/chain4travel/magellan/services/rewards"
	"github.com/chain4travel/magellan/servicesctrl"
	"github.com/chain4travel/magellan/stream"
	"github.com/chain4travel/magellan/stream/consumers"
	"github.com/chain4travel/magellan/utils"
	"github.com/go-sql-driver/mysql"
	"github.com/gorilla/rpc/v2"
	"github.com/gorilla/rpc/v2/json2"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/spf13/cobra"
)

const (
	rootCmdUse  = "magelland [command]"
	rootCmdDesc = "Daemons for Magellan."

	apiCmdUse  = "api"
	apiCmdDesc = "Runs the API daemon"

	streamCmdUse  = "stream"
	streamCmdDesc = "Runs stream commands"

	streamReplayCmdUse  = "replay"
	streamReplayCmdDesc = "Runs the replay"

	streamIndexerCmdUse  = "indexer"
	streamIndexerCmdDesc = "Runs the stream indexer daemon"

	envCmdUse  = "env"
	envCmdDesc = "Displays information about the Magellan environment"

	defaultReplayQueueSize    = int(2000)
	defaultReplayQueueThreads = int(4)
)

var wgGlobal = &sync.WaitGroup{}

func main() {
	// Check for multiple commands
	argFlags := []string{}
	argCmds := [][]string{{}}
	cmdPos := 0
	for _, a := range os.Args[1:] {
		if strings.HasPrefix(a, "-") {
			argFlags = append(argFlags, a)
		} else if a == "," {
			argCmds = append(argCmds, []string{})
			cmdPos++
		} else {
			argCmds[cmdPos] = append(argCmds[cmdPos], a)
		}
	}
	for _, a := range argCmds {
		if cmd, err := getCmd(); err != nil {
			log.Fatalln("Failed to run:", err.Error())
		} else {
			cmd.SetArgs(append(argFlags, a...))
			wgGlobal.Add(1)
			go func() {
				cmd.Execute()
				wgGlobal.Done()
			}()
		}
	}
	wgGlobal.Wait()
}

// Execute runs the root command for magellan
func getCmd() (*cobra.Command, error) {
	rand.Seed(time.Now().UnixNano())

	var (
		runErr             error
		config             = &cfg.Config{}
		serviceControl     = &servicesctrl.Control{}
		configFile         = func() *string { s := ""; return &s }()
		replayqueuesize    = func() *int { i := defaultReplayQueueSize; return &i }()
		replayqueuethreads = func() *int { i := defaultReplayQueueThreads; return &i }()
		cmd                = &cobra.Command{Use: rootCmdUse, Short: rootCmdDesc, Long: rootCmdDesc,
			PersistentPreRun: func(cmd *cobra.Command, args []string) {
				c, err := cfg.NewFromFile(*configFile)
				if err != nil {
					log.Fatalln("Failed to read config file", *configFile, ":", err.Error())
				}
				lf := logging.NewFactory(c.Logging)
				alog, err := lf.Make("magellan")
				if err != nil {
					log.Fatalln("Failed to create log", c.Logging.Directory, ":", err.Error())
				}

				mysqllogger := &MysqlLogger{
					Log: alog,
				}
				_ = mysql.SetLogger(mysqllogger)

				models.SetBech32HRP(c.NetworkID)

				serviceControl.Log = alog
				serviceControl.Services = c.Services
				serviceControl.ServicesCfg = *c
				serviceControl.Chains = c.Chains
				serviceControl.Persist = db.NewPersist()
				serviceControl.Features = c.Features
				persist := db.NewPersist()
				serviceControl.BalanceManager = balance.NewManager(persist, serviceControl)
				err = serviceControl.Init(c.NetworkID)
				if err != nil {
					log.Fatalln("Failed to create service control", ":", err.Error())
				}

				*config = *c

				if config.MetricsListenAddr != "" {
					sm := http.NewServeMux()
					sm.Handle("/metrics", promhttp.Handler())
					go func() {
						err = http.ListenAndServe(config.MetricsListenAddr, sm)
						if err != nil {
							log.Fatalln("Failed to start metrics listener", err.Error())
						}
					}()
					alog.Info("Starting metrics handler on %s", config.MetricsListenAddr)
				}
				if config.AdminListenAddr != "" {
					rpcServer := rpc.NewServer()
					codec := json2.NewCodec()
					rpcServer.RegisterCodec(codec, "application/json")
					rpcServer.RegisterCodec(codec, "application/json;charset=UTF-8")
					api := oreliusRpc.NewAPI(alog)
					if err := rpcServer.RegisterService(api, "api"); err != nil {
						log.Fatalln("Failed to start admin listener", err.Error())
					}
					sm := http.NewServeMux()
					sm.Handle("/api", rpcServer)
					go func() {
						err = http.ListenAndServe(config.AdminListenAddr, sm)
						if err != nil {
							log.Fatalln("Failed to start metrics listener", err.Error())
						}
					}()
				}
			},
		}
	)

	// Add flags and commands
	cmd.PersistentFlags().StringVarP(configFile, "config", "c", "config.json", "config file")
	cmd.PersistentFlags().IntVarP(replayqueuesize, "replayqueuesize", "", defaultReplayQueueSize, fmt.Sprintf("replay queue size default %d", defaultReplayQueueSize))
	cmd.PersistentFlags().IntVarP(replayqueuethreads, "replayqueuethreads", "", defaultReplayQueueThreads, fmt.Sprintf("replay queue size threads default %d", defaultReplayQueueThreads))

	cmd.AddCommand(
		createReplayCmds(serviceControl, config, &runErr, replayqueuesize, replayqueuethreads),
		createStreamCmds(serviceControl, config, &runErr),
		createAPICmds(serviceControl, config, &runErr),
		createEnvCmds(config, &runErr))

	return cmd, nil
}

func createAPICmds(sc *servicesctrl.Control, config *cfg.Config, runErr *error) *cobra.Command {
	return &cobra.Command{
		Use:   apiCmdUse,
		Short: apiCmdDesc,
		Long:  apiCmdDesc,
		Run: func(cmd *cobra.Command, args []string) {
			lc, err := api.NewServer(sc, *config)
			if err != nil {
				*runErr = err
				return
			}
			runListenCloser(lc)
		},
	}
}

func createReplayCmds(sc *servicesctrl.Control, config *cfg.Config, runErr *error, replayqueuesize *int, replayqueuethreads *int) *cobra.Command {
	replayCmd := &cobra.Command{
		Use:   streamReplayCmdUse,
		Short: streamReplayCmdDesc,
		Long:  streamReplayCmdDesc,
		Run: func(cmd *cobra.Command, args []string) {
			replay := replay.NewDB(sc, config, *replayqueuesize, *replayqueuethreads)
			err := replay.Start()
			if err != nil {
				*runErr = err
			}
		},
	}

	return replayCmd
}

func createStreamCmds(sc *servicesctrl.Control, config *cfg.Config, runErr *error) *cobra.Command {
	streamCmd := &cobra.Command{
		Use:   streamCmdUse,
		Short: streamCmdDesc,
		Long:  streamCmdDesc,
		Run: func(cmd *cobra.Command, args []string) {
			_ = cmd.Help()
			os.Exit(0)
		},
	}

	// Add sub commands
	streamCmd.AddCommand(&cobra.Command{
		Use:   streamIndexerCmdUse,
		Short: streamIndexerCmdDesc,
		Long:  streamIndexerCmdDesc,
		Run: func(cmd *cobra.Command, arg []string) {
			runStreamProcessorManagers(
				sc,
				config,
				runErr,
				producerFactories(sc, config),
				[]consumers.ConsumerFactory{
					consumers.IndexerConsumer,
				},
				[]stream.ProcessorFactoryChainDB{
					consumers.IndexerDB,
					consumers.IndexerConsensusDB,
				},
				[]stream.ProcessorFactoryInstDB{
					consumers.IndexerCChainDB(),
				},
			)(cmd, arg)
		},
	})

	return streamCmd
}

func producerFactories(sc *servicesctrl.Control, cfg *cfg.Config) []utils.ListenCloser {
	var factories []utils.ListenCloser
	factories = append(factories, stream.NewProducerCChain(sc, *cfg))
	for _, v := range cfg.Chains {
		switch v.VMType {
		case consumers.IndexerAVMName:
			p, err := stream.NewProducerChain(sc, *cfg, v.ID, stream.EventTypeDecisions, stream.IndexTypeTransactions, stream.IndexXChain)
			if err != nil {
				panic(err)
			}
			factories = append(factories, p)
			p, err = stream.NewProducerChain(sc, *cfg, v.ID, stream.EventTypeConsensus, stream.IndexTypeVertices, stream.IndexXChain)
			if err != nil {
				panic(err)
			}
			factories = append(factories, p)
		case consumers.IndexerPVMName:
			p, err := stream.NewProducerChain(sc, *cfg, v.ID, stream.EventTypeDecisions, stream.IndexTypeBlocks, stream.IndexPChain)
			if err != nil {
				panic(err)
			}
			factories = append(factories, p)
		}
	}

	if sc.IsCChainIndex {
		p, err := stream.NewProducerChain(sc, *cfg, cfg.CchainID, stream.EventTypeDecisions, stream.IndexTypeBlocks, stream.IndexCChain)
		if err != nil {
			panic(err)
		}
		factories = append(factories, p)
	}

	return factories
}

func createEnvCmds(config *cfg.Config, runErr *error) *cobra.Command {
	return &cobra.Command{
		Use:   envCmdUse,
		Short: envCmdDesc,
		Long:  envCmdDesc,
		Run: func(_ *cobra.Command, _ []string) {
			configBytes, err := json.MarshalIndent(config, "", "    ")
			if err != nil {
				*runErr = err
				return
			}

			fmt.Println(string(configBytes))
		},
	}
}

// runListenCloser runs the ListenCloser until signaled to stop
func runListenCloser(lc utils.ListenCloser) {
	// Start listening in the background
	go func() {
		if err := lc.Listen(); err != nil {
			log.Fatalln("Daemon listen error:", err.Error())
		}
	}()

	// Wait for exit signal
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	<-sigCh

	// Stop server
	if err := lc.Close(); err != nil {
		log.Fatalln("Daemon shutdown error:", err.Error())
	}
}

// runStreamProcessorManagers returns a cobra command that instantiates and runs
// a set of stream process managers
func runStreamProcessorManagers(
	sc *servicesctrl.Control,
	config *cfg.Config,
	runError *error,
	listenCloseFactories []utils.ListenCloser,
	consumerFactories []consumers.ConsumerFactory,
	factoriesChainDB []stream.ProcessorFactoryChainDB,
	factoriesInstDB []stream.ProcessorFactoryInstDB,
) func(_ *cobra.Command, _ []string) {
	return func(_ *cobra.Command, _ []string) {
		wg := &sync.WaitGroup{}

		bm, _ := sc.BalanceManager.(*balance.Manager)
		err := bm.Start(sc.IsAccumulateBalanceIndexer)
		if err != nil {
			*runError = err
			return
		}
		defer func() {
			bm.Close()
		}()

		rh := &rewards.Handler{}
		err = rh.Start(sc)
		if err != nil {
			*runError = err
			return
		}
		defer func() {
			rh.Close()
		}()

		err = consumers.Bootstrap(sc, config.NetworkID, config.Chains, consumerFactories)
		if err != nil {
			*runError = err
			return
		}

		sc.BalanceManager.Exec()

		runningControl := utils.NewRunning()

		err = consumers.IndexerFactories(sc, config, factoriesChainDB, factoriesInstDB, wg, runningControl)
		if err != nil {
			*runError = err
			return
		}

		for _, listenCloseFactory := range listenCloseFactories {
			wg.Add(1)
			go func(lc utils.ListenCloser) {
				wg.Done()
				if err := lc.Listen(); err != nil {
					log.Fatalln("Daemon listen error:", err.Error())
				}
			}(listenCloseFactory)
		}

		// Wait for exit signal
		sigCh := make(chan os.Signal, 1)
		signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
		<-sigCh

		for _, listenCloseFactory := range listenCloseFactories {
			if err := listenCloseFactory.Close(); err != nil {
				log.Println("Daemon shutdown error:", err.Error())
			}
		}

		runningControl.Close()

		wg.Wait()
	}
}

type MysqlLogger struct {
	Log logging.Logger
}

func (m *MysqlLogger) Print(v ...interface{}) {
	s := fmt.Sprint(v...)
	m.Log.Warn("mysql %s", s)
}
