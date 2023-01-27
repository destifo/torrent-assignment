package main

import (
	"crypto/rand"
	"log"
	"os"

	bencodeParser "github.com/destifo/torrent-assignment/packages/parsers"
)

func main() {

	// reading file
	torrentFilePath := "resources/archlinux-2019.12.01-x86_64.iso.torrent"
	encodedTorrentFile, err := os.Open(torrentFilePath)
	if err != nil {
		log.Fatalf("error reading file", err)
	}

	// parsing torrent file
	bto, err := bencodeParser.Open(encodedTorrentFile)
	if err != nil {
		log.Fatalf("Error parsing bencode file to BencodeTorrent object")
	}

	torrentFile, err := bto.ToTorrentFile()
	if err != nil {
		log.Fatalf("error converting to torrent struct: %s", err)
	}

    var peerID [20]byte
	_, err = rand.Read(peerID[:])
	if err != nil {
		log.Fatalf("error creating ")
	}
    trackerUrl, err := torrentFile.BuildTrackerURL()
}
