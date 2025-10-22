package main

import (
	"log"
	"os"
	"strings"

	"github.com/mildlybrutal/kvStore/internal/cluster"
	"github.com/mildlybrutal/kvStore/internal/node"
	"github.com/spf13/cobra"
)

var (
	id    string
	addr  string
	peers string
)

func main() {
	root := &cobra.Command{
		Use:   "kvnode",
		Short: "Run a kvStore node",
		RunE: func(cmd *cobra.Command, args []string) error {
			if addr == "" || peers == "" {
				return cmd.Usage()
			}
			ps := strings.Split(peers, ",")
			if len(ps) == 0 {
				return cmd.Usage()
			}
			// Expect peers to include this node's addr; split into self + others
			self := addr
			var others []string
			for _, p := range ps {
				if p == "" || p == self {
					continue
				}
				others = append(others, p)
			}
			r := cluster.NewCluster(self, others)

			n := node.NewNode(id, addr, r)
			log.Printf("starting %s at %s; peers: %v", id, addr, ps)
			n.Serve() // blocks
			return nil
		},
	}

	root.Flags().StringVar(&id, "id", "node-1", "node id")
	root.Flags().StringVar(&addr, "addr", ":8001", "listen address for this node")
	root.Flags().StringVar(&peers, "peers", ":8001,:8002,:8003", "comma-separated peer addresses (include this node)")

	if err := root.Execute(); err != nil {
		log.Println(err)
		os.Exit(1)
	}
}
