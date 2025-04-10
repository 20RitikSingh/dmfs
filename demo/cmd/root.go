package cmd

import (
	"fmt"
	"os"
	"time"

	"github.com/20ritiksingh/dmfs/pkg/p2p"
	"github.com/20ritiksingh/dmfs/server"
	"github.com/20ritiksingh/dmfs/store"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func makeServer(listenAddr string, root string, nodes ...string) *server.FileServer {
	opts := p2p.TCPTransporterOptions{
		ListenAddr: listenAddr,
		Handshake:  func(p2p.Peer) error { return nil },
		Decoder:    p2p.GOBDecoder{},
	}
	tcpTransporter := p2p.NewTCPTransporter(opts)
	fsOpts := server.FileServerOpts{
		RootDir:      root,
		Transporter:  tcpTransporter,
		GeneratePath: store.GenerateCASPath,
		Bootstrap:    nodes,
	}
	fs := server.NewFileServer(fsOpts)
	tcpTransporter.OnPeer = fs.OnPeer
	return fs
}

type DemoServerOpts struct {
	ListenAddr string   `mapstructure:"listenAddr"`
	Root       string   `mapstructure:"root"`
	Nodes      []string `mapstructure:"nodes"`
}

type Config struct {
	Servers []DemoServerOpts `mapstructure:"servers"`
}

var rootCmd = &cobra.Command{
	Use:   "dmfs",
	Short: "dfms is a distributed mesh file system",
	Long:  `DESC goes here`,
	Run: func(cmd *cobra.Command, args []string) {
		// Do Stuff Here
	},
	// PersistentPostRun: func(cmd *cobra.Command, args []string) {
	// 	for _, s := range Cluster {
	// 		s.Stop()
	// 	}
	// },
}

var Cluster = make([]*server.FileServer, 0)
var cfg Config

func initCofig() {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")

	err := viper.ReadInConfig()
	if err != nil {
		fmt.Println("config file not found")
		return
	}
	err = viper.Unmarshal(&cfg)
	if err != nil {
		fmt.Println("failed to unmarshall config")
	}
	for _, s := range cfg.Servers {
		Cluster = append(Cluster, makeServer(s.ListenAddr, s.Root, s.Nodes...))
	}
	for _, s := range Cluster {
		go s.Start()
		time.Sleep(time.Millisecond * 100)
		// fmt.Println("wait")
	}
	// time.Sleep(time.Second * 5)
}

func init() {
	//flags go here

	//cmdz go here
	cobra.OnInitialize(initCofig)
	rootCmd.AddCommand(upload)
	rootCmd.AddCommand(delete)
	rootCmd.AddCommand(read)

}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
