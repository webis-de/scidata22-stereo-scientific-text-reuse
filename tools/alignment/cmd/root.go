/*******************************************************************************
 * Last modified: 05.11.20, 23:17
 * Author: Lukas Gienapp
 * Copyright (c) 2020 picapica.org
 */

package cmd

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"os"
)

var (
	logLevel string
	rootCmd  = &cobra.Command{
		Use:   "alignment",
		Short: "The alignment component of Picapica",
		Long:  `Start and configure the alignment server component of Picapica project`,
	}
)

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	// Add subcommands
	rootCmd.AddCommand(alignCmd)
	rootCmd.AddCommand(readCmd)
	// Flags
	rootCmd.PersistentFlags().StringVar(&logLevel, "log_level", "info", "level of log verbosity")
	// Logging setup
	log.SetFormatter(&log.JSONFormatter{})
	log.SetOutput(os.Stdout)
	switch logLevel {
	case "trace":
		log.SetLevel(log.TraceLevel)
	case "debug":
		log.SetLevel(log.DebugLevel)
	case "info":
		log.SetLevel(log.InfoLevel)
	case "warn":
		log.SetLevel(log.WarnLevel)
	case "error":
		log.SetLevel(log.ErrorLevel)
	default:
		log.SetLevel(log.InfoLevel)
	}
}
