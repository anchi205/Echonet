package src

import (
	"encoding/json"
	"fmt"
)

// Publishes a chatmessage to the PubSub topic until the pubsub context closes
func (cr *ChatRoom) PubLoop() {
	for {
		select {
		case <-cr.psctx.Done():
			return

		case message := <-cr.Outbound:
			if err := cr.publishMessage(message); err != nil {
				cr.logError("puberr", "could not publish message to topic", err)
			}
		}
	}
}

func (cr *ChatRoom) publishMessage(message string) error {
	// Create a ChatMessage
	chatMsg := cr.createChatMessage(message)

	// Marshal the ChatMessage into JSON
	messageBytes, err := json.Marshal(chatMsg)
	if err != nil {
		return fmt.Errorf("could not marshal JSON: %w", err)
	}

	// Publish message to the topic
	if err := cr.pstopic.Publish(cr.psctx, messageBytes); err != nil {
		return fmt.Errorf("could not publish to topic: %w", err)
	}

	return nil
}

// Creates a ChatMessage with the given message, sender ID, and sender name
func (cr *ChatRoom) createChatMessage(message string) chatmessage {
	return chatmessage{
		Message:    message,
		SenderID:   cr.SelfID.String(),
		SenderName: cr.UserName,
	}
}

// Logs an error with the specified log prefix and message
func (cr *ChatRoom) logError(logPrefix, logMsg string, err error) {
	cr.Logs <- chatlog{logprefix: logPrefix, logmsg: fmt.Sprintf("%s: %v", logMsg, err)}
}
