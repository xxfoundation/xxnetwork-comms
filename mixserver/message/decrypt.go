package message

import (
	"gitlab.com/privategrity/crypto/cyclic"
	"gitlab.com/privategrity/server/cryptops/precomputation"
	"gitlab.com/privategrity/server/node"
	"gitlab.com/privategrity/server/services"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"time"

	pb "gitlab.com/privategrity/comms/mixmessages"

	jww "github.com/spf13/jwalterweatherman"
)

func SendPrecompDecrypt(nextServer string, round node.Round, input *pb.PrecompDecryptMessage) (*pb.Ack, error) {
	// Create Dispatcher for Decrypt
	dcPrecompDecrypt := services.DispatchCryptop(node.Grp, precomputation.Decrypt{}, nil, nil, round)

	// Convert input message to equivalent SlotDecrypt
	slotDecrypt := &precomputation.SlotDecrypt{
		Slot:                         input.Slot,
		EncryptedMessageKeys:         cyclic.NewIntFromBytes(input.EncryptedMessageKeys),
		EncryptedRecipientIDKeys:     cyclic.NewIntFromBytes(input.EncryptedRecipientIDKeys),
		PartialMessageCypherText:     cyclic.NewIntFromBytes(input.PartialMessageCypherText),
		PartialRecipientIDCypherText: cyclic.NewIntFromBytes(input.PartialRecipientIDCypherText),
	}
	// Type assert SlotDecrypt to Slot
	var slot services.Slot = slotDecrypt

	// Pass slot as input to Decrypt
	dcPrecompDecrypt.InChannel <- &slot

	// Get output from Decrypt
	output := <-dcPrecompDecrypt.OutChannel
	// Type assert Slot to SlotDecrypt
	out := (*output).(*precomputation.SlotDecrypt)

	// Attempt to connect to nextServer
	conn, err := grpc.Dial(nextServer, grpc.WithInsecure())
	// Check for an error
	if err != nil {
		jww.ERROR.Printf("Failed to connect to server at %v\n", nextServer)
	}

	// Prepare to send a message
	c := pb.NewMixMessageServiceClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), 300*time.Millisecond)
	defer cancel()

	// Send the PrecompDecrypt message using the Decrypt output
	result, err := c.PrecompDecrypt(ctx, &pb.PrecompDecryptMessage{
		Slot:                         out.Slot,
		EncryptedMessageKeys:         out.EncryptedMessageKeys.Bytes(),
		EncryptedRecipientIDKeys:     out.EncryptedRecipientIDKeys.Bytes(),
		PartialMessageCypherText:     out.PartialMessageCypherText.Bytes(),
		PartialRecipientIDCypherText: out.PartialRecipientIDCypherText.Bytes(),
	})
	// Make sure there are no errors with sending the message
	if err != nil {
		jww.ERROR.Printf("PrecompDecrypt: Error received: %s", err)
	}

	return result, err
}
