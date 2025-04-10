package server

import (
	"bytes"
	"errors"
	"io"
	"log"
	"sync"

	"github.com/20ritiksingh/dmfs/pkg/p2p"
	"github.com/20ritiksingh/dmfs/store"
)

type FileServerOpts struct {
	Transporter  p2p.Transporter
	GeneratePath store.GeneratePath
	RootDir      string
	Bootstrap    []string
}
type FileServer struct {
	FileServerOpts
	Store  *store.Store
	quitch chan struct{}

	peerMapMutex sync.Mutex
	peers        map[string]p2p.Peer
}

func NewFileServer(opts FileServerOpts) *FileServer {
	return &FileServer{
		FileServerOpts: opts,
		quitch:         make(chan struct{}),
		peers:          make(map[string]p2p.Peer),
		Store: store.NewStore(store.StoreOptions{
			RootDir:      opts.RootDir,
			GeneratePath: opts.GeneratePath,
		}),
	}
}

func (fs *FileServer) Check(key string) bool {
	return fs.Store.Has(key)
}

func (fs *FileServer) ReadData(key string) (io.Reader, error) {
	if fs.Store.Has(key) {
		return fs.Store.Read(key)
	}
	return nil, errors.New("file not found")
}

func (fs *FileServer) DeleteData(payload *p2p.Payload) error {
	// log.Printf("Deleting %+v", payload)
	if !fs.Store.Has(payload.Key) {
		return nil
	}
	if err := fs.Store.Delete(payload.Key); err != nil {
		return err
	}
	return fs.broadcast(payload)
}

func (fs *FileServer) broadcast(payload *p2p.Payload) error {
	// log.Printf("boradcasting %+v to %d peers", payload, len(fs.peers))
	peers := []io.Writer{}
	for _, peer := range fs.peers {
		peers = append(peers, peer)
	}
	if len(peers) == 0 {
		return nil
	}
	multiWriter := io.MultiWriter(peers...)
	if err := p2p.GOBEncode(multiWriter, payload); err != nil {
		log.Printf("error broadcasting to peers: %s", err)
		return err
	}
	return nil
}

func (fs *FileServer) StoreData(payload *p2p.Payload) error {
	fs.Store.Write(payload.Key, bytes.NewReader(payload.Data))
	return fs.broadcast(payload)
}

func (fs *FileServer) OnPeer(peer p2p.Peer) error {
	fs.peerMapMutex.Lock()
	defer fs.peerMapMutex.Unlock()
	fs.peers[peer.RemoteAddr().String()] = peer
	log.Println("new peer added: ", peer.RemoteAddr())
	return nil
}

func (fs *FileServer) loop() {
	defer func() {
		log.Println("file server has stopped")
		fs.Transporter.Close()
	}()
	for {
		select {
		case payload := <-fs.Transporter.Consume():
			// log.Printf("payload recived %+v by %s", payload, fs.RootDir)
			if len(payload.Data) == 0 {
				// log.Printf("this is delete signal")
				if fs.Store.Has(payload.Key) {
					fs.DeleteData(&payload)
				}
				continue
			}
			if fs.Store.Has(payload.Key) {
				// log.Println("already exists")
				continue
			}
			fs.StoreData(&payload)
		case <-fs.quitch:
			return
		}
	}
}

func (fs *FileServer) bootstrapNetwork() {
	for _, addr := range fs.Bootstrap {
		go func() {
			if err := fs.Transporter.Dial(addr); err != nil {
				log.Printf("failed to connect to peer %s while bootstraping", addr)
			}
		}()
	}
}

func (fs *FileServer) Start() error {
	fs.bootstrapNetwork()
	if err := fs.Transporter.ListenAndAccept(); err != nil {
		log.Println(err)
		return err
	}
	fs.loop()
	return nil
}
func (fs *FileServer) Stop() {
	close(fs.quitch)
}
