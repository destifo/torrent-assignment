package bencodeParser

import (
	"bytes"
	"crypto/sha1"
	"fmt"
	"io"
	"log"

	"github.com/destifo/torrent-assignment/models/Torrent"
	"github.com/jackpal/bencode-go"
)

type bencodeInfo struct {
    Pieces      string `bencode:"pieces"`
    PieceLength int    `bencode:"piece length"`
    Length      int    `bencode:"length"`
    Name        string `bencode:"name"`
}

type bencodeTorrent struct {
    Announce string      `bencode:"announce"`
    Info     bencodeInfo `bencode:"info"`
}

func (i *bencodeInfo) hash() ([20]byte, error) {
	var buf bytes.Buffer
	err := bencode.Marshal(&buf, *i)
	if err != nil {
		return [20]byte{}, err
	}
	h := sha1.Sum(buf.Bytes())
	return h, nil
}

func (i *bencodeInfo) splitPieceHashes() ([][20]byte, error) {
	hashLen := 20 // Length of SHA-1 hash
	buf := []byte(i.Pieces)
	if len(buf)%hashLen != 0 {
		err := fmt.Errorf("Received malformed pieces of length %d", len(buf))
		return nil, err
	}
	numHashes := len(buf) / hashLen
	hashes := make([][20]byte, numHashes)

	for i := 0; i < numHashes; i++ {
		copy(hashes[i][:], buf[i*hashLen:(i+1)*hashLen])
	}
	return hashes, nil
}

func (bto *bencodeTorrent) ToTorrentFile() (torrentFile.TorrentFile, error) {
    infoHash, err := bto.Info.hash()
    if err != nil {
        log.Fatalf("error during hashing the bencode info map", err)
    }
    
    piecesHahes, err := bto.Info.splitPieceHashes()
    if err != nil {
        log.Fatalf("error converting pieces into hashes", err)
    }

    torrent := torrentFile.TorrentFile {
        Announce: bto.Announce,
        InfoHash: infoHash,
        PieceHashes: piecesHahes,
        PieceLength: bto.Info.PieceLength,
        Length: bto.Info.Length,
        Name: bto.Info.Name,
    }

    return torrent, nil;

}

// takes torrent file contents from io, unmarshalles it to bencodeTorrent struct
func Open(r io.Reader) (*bencodeTorrent, error) {
    bto := bencodeTorrent{}
    err := bencode.Unmarshal(r, &bto)
    if err != nil {
        return nil, err
    }
    return &bto, nil
}