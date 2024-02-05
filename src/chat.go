package src

import (
	"context"
	"fmt"

	pubsub "github.com/libp2p/go-libp2p-pubsub"
	"github.com/libp2p/go-libp2p/core/peer"
)

const defaultuser = "New-user"
const defaultroom = "lobby"

// Structure that represents a PubSub Chat Room
type ChatRoom struct {
	Inbound  chan chatmessage // Channel of incoming messages
	Outbound chan string      // Channel of outgoing messages
	Logs     chan chatlog     // Channel of chat log messages

	RoomName string
	UserName string

	SelfID peer.ID // Represents the host ID of the peer

	psctx    context.Context      // Represents the chat room lifecycle context
	pscancel context.CancelFunc   // Represents the chat room lifecycle cancellation function
	pstopic  *pubsub.Topic        // Represents the PubSub Topic of the ChatRoom
	psub     *pubsub.Subscription // Represents the PubSub Subscription for the topic
}

type chatmessage struct {
	Message    string `json:"message"`
	SenderID   string `json:"senderid"`
	SenderName string `json:"sendername"`
}

// A structure that represents a chat log
type chatlog struct {
	logprefix string
	logmsg    string
}

// A constructor function that generates and returns a new
// ChatRoom for a given P2PHost, username and roomname
func JoinChatRoom(p2phost *P2P, username string, roomname string) (*ChatRoom, error) {
	// Create a PubSub topic with the room name
	topic, err := p2phost.PubSub.Join(fmt.Sprintf("room-peerchat-%s", roomname))
	if err != nil {
		return nil, err
	}

	// Subscribe to the PubSub topic
	sub, err := topic.Subscribe()
	if err != nil {
		return nil, err
	}

	if username == "" {
		username = defaultuser
	}

	if roomname == "" {
		roomname = defaultroom
	}

	// Cancellable context
	pubsubctx, cancel := context.WithCancel(context.Background())

	// ChatRoom object
	chatroom := &ChatRoom{
		Inbound:  make(chan chatmessage),
		Outbound: make(chan string),
		Logs:     make(chan chatlog),

		psctx:    pubsubctx,
		pscancel: cancel,
		pstopic:  topic,
		psub:     sub,

		RoomName: roomname,
		UserName: username,
		SelfID:   p2phost.Host.ID(),
	}

	go chatroom.SubLoop() // Start the subscribe loop
	go chatroom.PubLoop() // Start the publish loop

	return chatroom, nil
}

// Method of ChatRoom that returns a list of all peer IDs connected to chat room topic
func (cr *ChatRoom) PeerList() []peer.ID {
	return cr.pstopic.ListPeers() // slice of peer IDs
}

// Method of ChatRoom that updates the chat room by subscribing to the new topic
func (cr *ChatRoom) Exit() {
	defer cr.pscancel()

	cr.psub.Cancel()   // Cancel the existing subscription
	cr.pstopic.Close() // Close the topic handler
}

func (cr *ChatRoom) UpdateUser(username string) {
	cr.UserName = username
}
