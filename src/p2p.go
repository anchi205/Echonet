package src

import (
	"context"

	"github.com/libp2p/go-libp2p"
	dht "github.com/libp2p/go-libp2p-kad-dht"
	pubsub "github.com/libp2p/go-libp2p-pubsub"
	host "github.com/libp2p/go-libp2p/core/host"
	discovery "github.com/libp2p/go-libp2p/p2p/discovery/routing"
	"github.com/multiformats/go-multiaddr"
	"github.com/sirupsen/logrus"
)

// Service connection ID
const service = "anchi205/Echonet"

// Structure representing P2P Host
type P2P struct {
	// Represents the host context layer
	// Context object used for handling the lifecycle of the P2P host
	Ctx context.Context

	// Instance of the libp2p.Host interface
	Host host.Host

	// Instance of the Kademlia DHT for peer discovery
	KadDHT *dht.IpfsDHT

	// Represents the peer discovery service
	Discovery *discovery.RoutingDiscovery

	// GossipSub-based pubsub router for handling publish/subscribe messaging.
	PubSub *pubsub.PubSub
}

/*
A constructor function that generates and returns a P2PHost for a given context object.
Constructs a libp2p host with a multiaddr on 0.0.0.0/0 IPV4 address and configure it
with NATPortMap to open a port in the firewall using UPnP. A GossipSub pubsub router
is initialized for transport and a Kademlia DHT for peer discovery
*/
func NewP2P() *P2P {
	ctx := context.Background() // Setup a background context

	// Create a new multiaddr object
	sourcemultiaddr, _ := multiaddr.NewMultiaddr("/ip4/0.0.0.0/tcp/0")

	// Construct a new LibP2P host with the multiaddr and the NAT Port Map
	// Initializes a new libp2p.Host (libhost) using libp2p.New.
	libhost, err := libp2p.New(
		libp2p.ListenAddrs(sourcemultiaddr), // listens on the specified multiaddress.
		libp2p.NATPortMap(),                 // Enables NAT port mapping
	)

	if err != nil {
		logrus.WithFields(logrus.Fields{
			"error": err.Error(),
		}).Fatalln("P2P Host Creation Failed!")
	}

	// Create a new PubSub service which uses a GossipSub router
	// Initializes a GossipSub-based pubsub router (gossip)
	pubsubhandler, err := pubsub.NewGossipSub(ctx, libhost)

	if err != nil {
		logrus.WithFields(logrus.Fields{
			"error": err.Error(),
		}).Fatalln("GossipSub Router Creation Failed!")
	}

	// Bind the LibP2P host to a Kademlia DHT peer
	// Initializes a Kademlia DHT (kaddht) for peer discovery
	// DHT is configured in server mode (dht.ModeServer).
	kaddht, err := dht.New(ctx, libhost, dht.Mode(dht.ModeServer))

	if err != nil {
		logrus.WithFields(logrus.Fields{
			"error": err.Error(),
		}).Fatalln("Kademlia DHT Creation Failed!")
	}

	// Create a peer discovery service using the Kad DHT
	routingdiscovery := discovery.NewRoutingDiscovery(kaddht)

	// Pointer to a new P2PHost instance with the created host, DHT, and pubsub router.
	return &P2P{
		Ctx:       ctx,
		Host:      libhost,
		KadDHT:    kaddht,
		PubSub:    pubsubhandler,
		Discovery: routingdiscovery,
	}
}
