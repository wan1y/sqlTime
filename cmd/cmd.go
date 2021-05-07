package cmd

import (
	"context"
	"fmt"
	"github.com/spf13/cobra"
	"gorm.io/gorm/logger"
	"io/ioutil"
	"os"
	"sqlTime/server"
	"strconv"
	"syscall"
)

var (
	threads           int
	dsns              []string
	dsnFileNames      []string
	sqlFile           string
	logLevel          string
)

var timeCommand = &cobra.Command{
	Use:   "time",
	Short: "Compare SQL execution time",
	Args:  cobra.MinimumNArgs(0),
	Run: func(cmd *cobra.Command, args []string) {
		s := server.New(threads, logLevel, sqlFile)
		if len(dsns) != len(dsnFileNames) {
			panic("dsns and dsnfilenames must be one-to-one correspondence")
		}
		s.SetDsnAndFileNames(dsns, dsnFileNames)
		s.CompareTime()
	},
}

var RootCmd = &cobra.Command{
	Use:   "sqlTime",
	Short: "sql execution time test tool",
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
	},
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

func checkPID(path string) {
	_, err := os.Stat(path)
	pid := os.Getpid()
	if os.IsNotExist(err) {
		logger.Default.Info(context.Background(), fmt.Sprintf("write %d to %s", pid, path))
		if err := ioutil.WriteFile(path, []byte(strconv.Itoa(pid)), 0644); err != nil {
			panic(err)
		}
		return
	}

	data, err := ioutil.ReadFile(path)
	if err != nil {
		panic(err)
	}
	opid, err := strconv.Atoi(string(data))
	if err != nil {
		panic(err)
	}
	if opid == pid {
		return
	}

	p, err := os.FindProcess(opid)
	if err != nil {
		panic(err)
	}
	if err := p.Signal(syscall.Signal(0)); err != nil {
		logger.Default.Info(context.Background(), err.Error())
	} else {
		logger.Default.Info(context.Background(), fmt.Sprintf("%d is alive, exited!", opid))
		os.Exit(1)
	}

	logger.Default.Info(context.Background(), fmt.Sprintf("write %d to %s", pid, path))
	if err := ioutil.WriteFile(path, []byte(strconv.Itoa(pid)), 0644); err != nil {
		panic(err)
	}
	return
}

func init() {
	RootCmd.PersistentFlags().IntVar(&threads, "threads", 1, "concurrent threads")
	RootCmd.PersistentFlags().StringVar(&logLevel, "log-level", "error", "log level: info, warn, error, silent")
	RootCmd.PersistentFlags().StringVar(&sqlFile, "sqlfile", "test.sql", "sql file")
	RootCmd.PersistentFlags().StringSliceVar(&dsnFileNames, "filenames", nil, "file name to save the result")
	RootCmd.PersistentFlags().StringSliceVar(&dsns, "dsns", nil, "set DSN of test dbs,take the first one as the standard")
	RootCmd.AddCommand(timeCommand)
}
