package main

import (
	"context"
	"crypto/rand"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	libp2pquic "github.com/costinm/go-libp2p-quic-transport"
	ic "github.com/libp2p/go-libp2p-core/crypto"
	peer "github.com/libp2p/go-libp2p-core/peer"
	ma "github.com/multiformats/go-multiaddr"
)

func main() {
	var addr, pub string
	if len(os.Args) != 3 {
		fmt.Printf("Usage: %s <multiaddr> <peer id>", os.Args[0])
		addr = "5555"
		pub = "12D3KooWBbkYafqbHDtmCpp47aj8P16YVfUGtyBeBB1txENTYU7x"
	} else {
		addr = os.Args[1]
		pub = os.Args[2]
	}
	if err := run(addr, pub); err != nil {
		log.Fatalf(err.Error())
	}
}

func run(raddr string, p string) error {
	peerID, err := peer.Decode(p)
	if err != nil {
		return err
	}
	addr, err := ma.NewMultiaddr(raddr)
	if err != nil {
		return err
	}
	priv, _, err := ic.GenerateECDSAKeyPair(rand.Reader)
	if err != nil {
		return err
	}

	t, err := libp2pquic.NewTransport(priv, nil, nil)
	if err != nil {
		return err
	}

	log.Printf("Dialing %s\n", addr.String())
	conn, err := t.Dial(context.Background(), addr, peerID)
	if err != nil {
		return err
	}
	defer conn.Close()
	str, err := conn.OpenStream()
	if err != nil {
		return err
	}
	const msg = "Hello world!"
	log.Println(conn.RemotePeer().String(), conn.RemoteMultiaddr().String(),
		conn.RemotePublicKey())
	log.Printf("Sending: %s\n", msg)
	if _, err := str.Write([]byte(msg)); err != nil {
		return err
	}
	if err := str.Close(); err != nil {
		return err
	}
	data, err := ioutil.ReadAll(str)
	if err != nil {
		return err
	}
	log.Printf("Received: %s\n", data)
	return nil
}
