package src

import "encoding/json"

func (cr *ChatRoom) SubLoop() {
	for {
		select {
		case <-cr.psctx.Done():
			return

		default:
			message, err := cr.psub.Next(cr.psctx) // Read a message from the subscription
			if err != nil {
				close(cr.Inbound) // Close the messages queue (subscription has closed)
				cr.Logs <- chatlog{logprefix: "suberr", logmsg: "subscription has closed"}
				return
			}

			if message.ReceivedFrom == cr.SelfID { // Check if message is from self
				continue
			}

			cm := &chatmessage{}

			err = json.Unmarshal(message.Data, cm) // Unmarshal the message data into a ChatMessage
			if err != nil {
				cr.Logs <- chatlog{logprefix: "suberr", logmsg: "could not unmarshal JSON"}
				continue
			}

			cr.Inbound <- *cm // Send the ChatMessage into the message queue
		}
	}
}

/*
SubLoop:
   Continously reads from the subscription until either the subscription or pubsub context closes.
   The recieved message is parsed sent into the inbound channel
*/
