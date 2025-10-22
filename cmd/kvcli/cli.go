package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/mildlybrutal/kvStore/internal/cluster"
	"github.com/mildlybrutal/kvStore/internal/node"
	"github.com/spf13/cobra"
)

var peers string

func buildClient(peersCSV string) *node.Node {
	ps := strings.Split(peersCSV, ",")
	if len(ps) == 0 || ps[0] == "" {
		log.Fatal("no peers provided (use --peers=:8001,:8002,:8003)")
	}
	self := ps[0]
	var others []string
	if len(ps) > 1 {
		others = ps[1:]
	}
	r := cluster.NewCluster(self, others)
	return node.NewNode("client", ":0", r) // not serving
}

func main() {
	rootCmd := &cobra.Command{
		Use:   "kvcli",
		Short: "CLI for the kvStore cluster",
	}

	rootCmd.PersistentFlags().StringVar(&peers, "peers", ":8001,:8002,:8003", "comma-separated peer addresses")

	putCmd := &cobra.Command{
		Use:   "put <key> <value>",
		Short: "Put a key/value",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			key, val := args[0], args[1]
			c := buildClient(peers)
			return c.Put(key, val)
		},
	}

	getCmd := &cobra.Command{
		Use:   "get <key>",
		Short: "Get a value by key",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			key := args[0]
			c := buildClient(peers)
			val, err := c.Get(key)
			if err != nil {
				return err
			}
			fmt.Println(val)
			return nil
		},
	}
	rootCmd.AddCommand(putCmd, getCmd)
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
