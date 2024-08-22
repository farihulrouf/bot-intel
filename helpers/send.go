package helpers

import (
	"context"
	"fmt"

	"github.com/rs/zerolog/log"
	"go.mau.fi/whatsmeow"
	waProto "go.mau.fi/whatsmeow/binary/proto"
	"go.mau.fi/whatsmeow/types"
	"google.golang.org/protobuf/proto"
)

func SendMessageToPhoneNumber(client *whatsmeow.Client, recipient, message string) error {
	// Convert recipient to JID
	jid, err := types.ParseJID(recipient + "@s.whatsapp.net")
	if err != nil {
		log.Error().Err(err).Msg("Invalid recipient JID")
		return fmt.Errorf("invalid recipient JID: %v", err)
	}

	// Create the message
	msg := &waProto.Message{
		Conversation: proto.String(message),
	}

	// Send the message
	_, err = client.SendMessage(context.Background(), jid, msg)
	if err != nil {
		log.Error().Err(err).Msg("Error sending message")
		return fmt.Errorf("error sending message: %v", err)
	}

	log.Info().Msgf("Successfully sent message '%s' to phone number: %s", message, recipient)
	return nil
}
