package main

import (
	"encoding/gob"
	"net"

	"github.com/20ritiksingh/dmfs/demo/cmd"
)

// func makeServer(listenAddr string, root string, nodes ...string) *server.FileServer {
// 	opts := p2p.TCPTransporterOptions{
// 		ListenAddr: listenAddr,
// 		Handshake:  func(p2p.Peer) error { return nil },
// 		Decoder:    p2p.GOBDecoder{},
// 	}
// 	tcpTransporter := p2p.NewTCPTransporter(opts)
// 	fsOpts := server.FileServerOpts{
// 		RootDir:      root,
// 		Transporter:  tcpTransporter,
// 		GeneratePath: store.GenerateCASPath,
// 		Bootstrap:    nodes,
// 	}
// 	fs := server.NewFileServer(fsOpts)
// 	tcpTransporter.OnPeer = fs.OnPeer
// 	return fs
// }

func main() {
	gob.Register(&net.TCPAddr{})
	// log.Println("Its Alive!!")
	// fs1 := makeServer(":8080", "fs1", ":4000", ":8090")
	// go func() {
	// 	if err := fs1.Start(); err != nil {
	// 		log.Fatalf("error starting file server 1 %s", err)
	// 	}
	// }()

	// fs2 := makeServer(":8090", "fs2", ":8080")
	// time.Sleep(time.Second * 5)
	// go func() {
	// 	if err := fs2.Start(); err != nil {
	// 		log.Fatalln("error starting server 2", err)
	// 	}
	// }()
	// time.Sleep(time.Second * 2)
	// if err := fs2.StoreData(&p2p.Payload{
	// 	Key:  "this is key",
	// 	Data: []byte("this is data "),
	// 	From: &net.TCPAddr{IP: net.ParseIP("127.0.0.1"), Port: 9998},
	// }); err != nil {
	// 	log.Println("error storing data, ", err)
	// }
	// time.Sleep(time.Second * 2)
	// r, err := fs1.ReadData("this is key")
	// if err != nil {
	// 	log.Println(err)
	// } else {
	// 	data := new(bytes.Buffer)
	// 	io.Copy(data, r)
	// 	log.Println("data read: ", data.Bytes())
	// }
	// time.Sleep(time.Second * 10)
	// if err := fs2.DeleteData(&p2p.Payload{
	// 	Key:  "this is key",
	// 	Data: []byte(""),
	// 	From: &net.TCPAddr{IP: net.ParseIP("127.0.0.1"), Port: 9998},
	// }); err != nil {
	// 	log.Println("error deleting ", err)
	// }
	// time.Sleep(time.Second * 20)
	// fs1.Store.Clear()
	// fs2.Store.Clear()
	cmd.Execute()
	// fs1.Stop()
	// fs2.Stop()
	// select {}
}
