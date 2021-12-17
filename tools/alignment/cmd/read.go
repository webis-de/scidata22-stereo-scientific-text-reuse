/*******************************************************************************
 * Last modified: 08.12.20, 12:04
 * Author: Lukas Gienapp
 * Copyright (c) 2020 picapica.org
 */

package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"math"
	"os"
	"picapica.org/picapica4/alignment/conf"
	"picapica.org/picapica4/alignment/core/decomposing"
	"picapica.org/picapica4/alignment/core/types"
)

type InputAlignmentTask struct {
	Config string `json:"Config"`
	Source string `json:"Source"`
	Target string `json:"Target"`
}

var (
	returnContext bool
	readCmd       = &cobra.Command{
		Use:   "read",
		Short: "Reads an alignment task from a json file",
		Long:  `Takes a json file with two texts and an alignment configuration as input and computes their local alignment`,
		RunE: func(cmd *cobra.Command, args []string) error {

			var task InputAlignmentTask
			err := json.NewDecoder(os.Stdin).Decode(&task)
			if err != nil {
				log.Fatal(err)
			}

			cfg := &viper.Viper{}
			cfg.SetConfigType("json")
			_ = cfg.ReadConfig(bytes.NewBuffer([]byte(task.Config)))
			aligner := conf.AlignerFromViper(cfg)

			decomposer := decomposing.NewStringDecomposer(" ", ".")
			source_sequence, err := decomposer.Decompose(task.Source)
			if err != nil {
				return err
			}
			target_sequence, err := decomposer.Decompose(task.Target)
			if err != nil {
				return err
			}
			seeds, err := aligner.Align(source_sequence, target_sequence)
			if err != nil {
				return err
			}

			msg, err := json.Marshal(generateResults(seeds, task.Source, task.Target, returnText, returnContext))
			if err != nil {
				return err
			}
			fmt.Println(string(msg))
			return nil
		},
	}
)

func generateResults(seeds []*types.Seed, source string, target string, returnText bool, returnContext bool) []Result {
	results := make([]Result, len(seeds), len(seeds))
	for i, s := range seeds {
		var (
			textSource   string = ""
			beforeSource string = ""
			afterSource  string = ""
			textTarget   string = ""
			beforeTarget string = ""
			afterTarget  string = ""
		)
		if returnText {
			textSource = source[s.SourceSpan.Begin:s.SourceSpan.End]
			textTarget = target[s.TargetSpan.Begin:s.TargetSpan.End]
		}
		if returnContext {
			beforeSource = source[int(math.Max(float64(s.SourceSpan.Begin)-101, 0)):s.SourceSpan.Begin]
			afterSource = source[s.SourceSpan.End:int(math.Min(float64(len(source)), float64(s.SourceSpan.End+101)))]
			beforeTarget = target[int(math.Max(float64(s.TargetSpan.Begin)-101, 0)):s.TargetSpan.Begin]
			afterTarget = target[s.TargetSpan.End:int(math.Min(float64(len(target)), float64(s.TargetSpan.End+101)))]
		}
		results[i] = Result{
			Source: Partial{
				Begin:  s.SourceSpan.Begin,
				End:    s.SourceSpan.End,
				Text:   textSource,
				Before: beforeSource,
				After:  afterSource,
			},
			Target: Partial{
				Begin:  s.TargetSpan.Begin,
				End:    s.TargetSpan.End,
				Text:   textTarget,
				Before: beforeTarget,
				After:  afterTarget,
			},
		}
	}
	return results
}

// Adds flags to the start cmd and writes their values to the viper conf
func init() {
	readCmd.PersistentFlags().BoolVar(&returnText, "text", false, "include text spans of seeds in output")
	readCmd.PersistentFlags().BoolVar(&returnContext, "context", false, "include surrounding text of seeds in output")
}
