/*******************************************************************************
 * Last modified: 08.12.20, 12:04
 * Author: Lukas Gienapp
 * Copyright (c) 2020 picapica.org
 */

package cmd

import (
	"encoding/json"
	"fmt"
	"github.com/spf13/cobra"
	"picapica.org/picapica4/alignment/conf"
	"picapica.org/picapica4/alignment/core/combinating"
	"picapica.org/picapica4/alignment/core/decomposing"
	"picapica.org/picapica4/alignment/core/extending"
	"picapica.org/picapica4/alignment/core/normalizing"
	"picapica.org/picapica4/alignment/core/reducing"
	"picapica.org/picapica4/alignment/core/seeding"
	"picapica.org/picapica4/alignment/core/similarity"
)

type Partial struct {
	Begin  uint32 `json:"begin"`
	End    uint32 `json:"end"`
	Text   string `json:"text,omitempty"`
	Before string `json:"before,omitempty"`
	After  string `json:"after,omitempty"`
}

type Result struct {
	Source Partial `json:"source"`
	Target Partial `json:"target"`
}

var (
	n          uint32
	overlap    uint32
	theta      uint32
	returnText bool
	alignCmd   = &cobra.Command{
		Use:   "align",
		Short: "Aligns two given texts",
		Long:  `Takes two input strings and computes their local alignment`,
		RunE: func(cmd *cobra.Command, args []string) error {
			source := args[0]
			target := args[1]
			aligner := conf.AlignerFromPreset(conf.Preset{
				Mode:       "",
				Similarity: similarity.NewIdentitySimilarity(),
				Seeder:     seeding.NewHashSeeder(n, overlap, normalizing.NewIdentityNormalizer(), combinating.NewStringValueCombinator(" ")),
				Extender:   extending.NewRangeExtender(theta),
				Reducer:    reducing.NewIdentityReducer(),
			})
			decomposer := decomposing.NewStringDecomposer(" ", ".")

			source_sequence, err := decomposer.Decompose(source)
			if err != nil {
				return err
			}

			target_sequence, err := decomposer.Decompose(args[1])
			if err != nil {
				return err
			}

			seeds, err := aligner.Align(source_sequence, target_sequence)
			if err != nil {
				return err
			}

			results := make([]Result, len(seeds), len(seeds))
			for i, s := range seeds {
				if returnText {
					results[i] = Result{
						Source: Partial{
							Begin: s.SourceSpan.Begin - 1,
							End:   s.SourceSpan.End,
							Text:  source[s.SourceSpan.Begin-1 : s.SourceSpan.End],
						},
						Target: Partial{
							Begin: s.TargetSpan.Begin - 1,
							End:   s.TargetSpan.End,
							Text:  target[s.TargetSpan.Begin-1 : s.TargetSpan.End],
						},
					}
				} else {
					results[i] = Result{
						Source: Partial{
							Begin: s.SourceSpan.Begin - 1,
							End:   s.SourceSpan.End,
							Text:  "",
						},
						Target: Partial{
							Begin: s.TargetSpan.Begin - 1,
							End:   s.TargetSpan.End,
							Text:  "",
						},
					}
				}
			}
			msg, err := json.Marshal(results)
			if err != nil {
				return err
			}
			fmt.Println(string(msg))
			return nil
		},
	}
)

// Adds flags to the start cmd and writes their values to the viper conf
func init() {
	alignCmd.PersistentFlags().Uint32VarP(&n, "ngram", "n", 8, "ngram length")
	alignCmd.PersistentFlags().Uint32VarP(&overlap, "overlap", "o", 7, "ngram overlap")
	alignCmd.PersistentFlags().Uint32VarP(&theta, "theta", "t", 300, "extension range")
	alignCmd.PersistentFlags().BoolVar(&returnText, "return_text", false, "include text spans of seeds in output")
}
