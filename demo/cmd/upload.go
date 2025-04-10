package cmd

import (
	"fmt"
	"os"
	"time"

	"github.com/20ritiksingh/dmfs/pkg/p2p"
	"github.com/spf13/cobra"
)

var path string
var upload = &cobra.Command{
	Use: "upload",
	Run: func(cmd *cobra.Command, args []string) {
		// fmt.Println("this is upload and given path is ", path)
		if len(Cluster) == 0 {
			fmt.Println("no servers found")
			return
		}
		if path == "" {
			fmt.Println("enter file path")
			return
		}
		file, err := os.ReadFile(path)
		if err != nil {
			fmt.Println("could not read file")
			return
		}
		if Cluster[0].Check(path) {
			fmt.Println("file already exists")
			return
		}
		err = Cluster[0].StoreData(&p2p.Payload{
			Key:  path,
			Data: file,
		})
		if err != nil {
			fmt.Println(err)
			return
		}
		time.Sleep(time.Second)
	},
}

func init() {
	upload.Flags().StringVarP(&path, "file", "f", "", "path of file to upload")
}
