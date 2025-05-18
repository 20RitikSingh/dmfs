package chunk

import (
	"errors"
	"io"
	"log"
	"os"

	encryption "github.com/20ritiksingh/dmfs/pkg/encryption"
)

// ChunkHandler reads a file in chunks and encrypts each chunk
// using the provided Codec. It sends the encrypted chunks to a channel.
// The chunks can be decrypted and written to a file using the Construct method.
// The ChunkHandler is designed to handle large files efficiently by processing them in smaller chunks.
// FileHasher is used to generate a unique key for the file.
type ChunkHandler struct {
	fileKey   string
	ChunkSize int
	FilePath  string
	ChunkChan chan Chunk
	Codec     encryption.Codec
}

// ChunkReaderConfig is a configuration struct for ChunkHandler.
// It contains the file path, chunk size, codec for encryption/decryption,
// and the channel buffer size for chunk communication.
// The ChunkSize defines the size of each chunk to be read from the file.
// The Codec is used for encrypting and decrypting the chunks.
// The ChanBuffer defines the size of the channel buffer for chunk communication.
// The Path is the file path of the file to be read.
type ChunkReaderConfig struct {
	FileKey    string
	Path       string
	ChunkSize  int
	Codec      encryption.Codec
	ChanBuffer int
}

// NewChunkReader creates a new ChunkHandler with the provided configuration.
func NewChunkReader(cfg ChunkReaderConfig) (*ChunkHandler, error) {
	if cfg.ChunkSize <= 0 {
		return nil, errors.New("chunk size must be greater than 0")
	}
	if cfg.Path == "" {
		return nil, errors.New("file path cannot be empty")
	}
	if cfg.Codec == nil {
		return nil, errors.New("codec cannot be nil")
	}
	if cfg.ChanBuffer <= 0 {
		return nil, errors.New("channel buffer size must be greater than 0")
	}
	if _, err := os.Stat(cfg.Path); os.IsNotExist(err) {
		return nil, errors.New("file does not exist")
	}
	return &ChunkHandler{
		fileKey:   cfg.FileKey,
		FilePath:  cfg.Path,
		ChunkSize: cfg.ChunkSize,
		ChunkChan: make(chan Chunk, cfg.ChanBuffer),
		Codec:     cfg.Codec,
	}, nil

}

// Read reads the file in chunks and encrypts each chunk
// using the provided Codec. It sends the encrypted chunks to the ChunkChan channel.
func (cr ChunkHandler) Read() error {
	defer close(cr.ChunkChan)
	buf := make([]byte, cr.ChunkSize)
	file, err := os.Open(cr.FilePath)
	if err != nil {
		log.Println("error opening file: ", err)
		return err
	}
	defer file.Close()

	//file encryption has to be done chunk wise to prevent loading large files in RAM
	//encrytion logic goes here
	i := 0
	for {
		n, err := file.Read(buf)
		if n > 0 {
			data := make([]byte, n)
			copy(data, buf[:n])
			cr.Codec.Encrypt(data)
			// Send the encrypted chunk to the channel
			cr.ChunkChan <- Chunk{
				Key:   "",
				Index: i,
				Data:  data,
			}
			i++
		}
		if err != nil && err != io.EOF {
			log.Println("error reading file: ", err)
			return err
		}
		if err == io.EOF {
			break
		}
	}
	return nil
}

// Construct takes a file and writes the chunks after decryption
// to the provided file path. It uses the Codec to decrypt the chunks.
func (cr ChunkHandler) ReconstructFile(path string) error {
	w, err := os.Create(path)
	if err != nil {
		log.Println("error creating file: ", err)
		return err
	}
	defer w.Close()

	for chunk := range cr.ChunkChan {
		cr.Codec.Decrypt(chunk.Data)
		if _, err := w.Write(chunk.Data); err != nil {
			log.Println("error writing file: ", err)
			return err
		}
		if err := w.Sync(); err != nil {
			log.Println("error syncing file: ", err)
			return err
		}
	}
	return nil
}

// Close closes the ChunkHandler and releases any resources.
// It is important to call this method to ensure that all resources are released properly.
func (cr ChunkHandler) Close() error {
	if cr.ChunkChan == nil {
		return nil
	}
	close(cr.ChunkChan)
	return nil
}

type Chunk struct {
	Key   string
	Index int
	Data  []byte
}
