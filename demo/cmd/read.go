package cmd

import (
	"fmt"
	"io"
	"time"

	"github.com/spf13/cobra"
)

var read = &cobra.Command{
	Use: "read",
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
		r, err := Cluster[0].ReadData(path)
		if err != nil {
			fmt.Println(err)
			return
		}
		data, err := io.ReadAll(r)
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Println(string(data))
		time.Sleep(time.Second)
	},
}

func init() {
	read.Flags().StringVarP(&path, "file", "f", "", "path of file to upload")
}
