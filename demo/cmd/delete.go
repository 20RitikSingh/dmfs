package cmd

import (
	"fmt"
	"path/filepath"
	"time"

	"github.com/20ritiksingh/dmfs/pkg/p2p"
	"github.com/spf13/cobra"
)

var key string
var delete = &cobra.Command{
	Use: "delete",
	Run: func(cmd *cobra.Command, args []string) {
		// fmt.Println("this is upload and given path is ", path)
		if len(Cluster) == 0 {
			fmt.Println("no servers found")
			return
		}
		if key == "" {
			fmt.Println("enter file path")
			return
		}

		err := Cluster[0].DeleteData(&p2p.Payload{
			Key:  filepath.Base(key),
			Data: []byte(""),
		})
		if err != nil {
			fmt.Println(err)
			return
		}
		time.Sleep(time.Second)
	},
}

func init() {
	delete.Flags().StringVarP(&key, "file", "f", "", "path of file to upload")
}
