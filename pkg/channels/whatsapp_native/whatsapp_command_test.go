//go:build whatsapp_native

package whatsapp

import (
	"context"
	"fmt"
	"testing"
	"time"

	"go.mau.fi/whatsmeow"
	"go.mau.fi/whatsmeow/proto/waE2E"
	"go.mau.fi/whatsmeow/store"
	"go.mau.fi/whatsmeow/types"
	"go.mau.fi/whatsmeow/types/events"
	"google.golang.org/protobuf/proto"

	"github.com/sipeed/picoclaw/pkg/bus"
	"github.com/sipeed/picoclaw/pkg/channels"
	"github.com/sipeed/picoclaw/pkg/config"
)

func TestHandleIncoming_DoesNotConsumeGenericCommandsLocally(t *testing.T) {
	messageBus := bus.NewMessageBus()
	ch := &WhatsAppNativeChannel{
		BaseChannel: channels.NewBaseChannel("whatsapp_native", config.WhatsAppConfig{}, messageBus, nil),
		runCtx:      context.Background(),
	}

	evt := &events.Message{
		Info: types.MessageInfo{
			MessageSource: types.MessageSource{
				Sender: types.NewJID("1001", types.DefaultUserServer),
				Chat:   types.NewJID("1001", types.DefaultUserServer),
			},
			ID:       "mid1",
			PushName: "Alice",
		},
		Message: &waE2E.Message{
			Conversation: proto.String("/new"),
		},
	}

	ch.handleIncoming(evt)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	select {
	case <-ctx.Done():
		t.Fatal("timeout waiting for message to be forwarded")
		return
	case inbound, ok := <-messageBus.InboundChan():
		if !ok {
			t.Fatal("expected inbound message to be forwarded")
		}
		if inbound.Channel != "whatsapp_native" {
			t.Fatalf("channel=%q", inbound.Channel)
		}
		if inbound.Content != "/new" {
			t.Fatalf("content=%q", inbound.Content)
		}
	}
}

// mockLIDStore implements store.LIDStore for testing LID→phone resolution.
type mockLIDStore struct {
	pnForLID map[types.JID]types.JID
}

func (m *mockLIDStore) PutManyLIDMappings(_ context.Context, _ []store.LIDMapping) error {
	return nil
}
func (m *mockLIDStore) PutLIDMapping(_ context.Context, _, _ types.JID) error { return nil }
func (m *mockLIDStore) GetPNForLID(_ context.Context, lid types.JID) (types.JID, error) {
	if pn, ok := m.pnForLID[lid]; ok {
		return pn, nil
	}
	return types.EmptyJID, fmt.Errorf("not found")
}
func (m *mockLIDStore) GetLIDForPN(_ context.Context, _ types.JID) (types.JID, error) {
	return types.EmptyJID, fmt.Errorf("not implemented")
}
func (m *mockLIDStore) GetManyLIDsForPNs(_ context.Context, _ []types.JID) (map[types.JID]types.JID, error) {
	return nil, nil
}

func TestHandleIncoming_LIDResolvedToPhoneForAllowList(t *testing.T) {
	lidJID := types.NewJID("76802773020849", types.HiddenUserServer)
	phoneJID := types.NewJID("46730956607", types.DefaultUserServer)

	messageBus := bus.NewMessageBus()
	ch := &WhatsAppNativeChannel{
		BaseChannel: channels.NewBaseChannel(
			"whatsapp_native",
			config.WhatsAppConfig{},
			messageBus,
			// Allow list uses the plain phone number — this should match
			// after LID resolution.
			[]string{"46730956607"},
		),
		runCtx: context.Background(),
		client: whatsmeow.NewClient(&store.Device{
			LIDs: &mockLIDStore{
				pnForLID: map[types.JID]types.JID{lidJID: phoneJID},
			},
		}, nil),
	}

	evt := &events.Message{
		Info: types.MessageInfo{
			MessageSource: types.MessageSource{
				Sender: lidJID,
				Chat:   lidJID,
			},
			ID:       "mid-lid-1",
			PushName: "Arthur",
		},
		Message: &waE2E.Message{
			Conversation: proto.String("Hello from LID"),
		},
	}

	ch.handleIncoming(evt)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	select {
	case <-ctx.Done():
		t.Fatal("message was dropped — LID sender should have been resolved and allowed")
	case inbound := <-messageBus.InboundChan():
		if inbound.Content != "Hello from LID" {
			t.Fatalf("unexpected content: %q", inbound.Content)
		}
	}
}

func TestHandleIncoming_LIDNotInAllowList_Dropped(t *testing.T) {
	lidJID := types.NewJID("76802773020849", types.HiddenUserServer)
	phoneJID := types.NewJID("46730956607", types.DefaultUserServer)

	messageBus := bus.NewMessageBus()
	ch := &WhatsAppNativeChannel{
		BaseChannel: channels.NewBaseChannel(
			"whatsapp_native",
			config.WhatsAppConfig{},
			messageBus,
			// Allow list has a different number — should NOT match.
			[]string{"11111111111"},
		),
		runCtx: context.Background(),
		client: whatsmeow.NewClient(&store.Device{
			LIDs: &mockLIDStore{
				pnForLID: map[types.JID]types.JID{lidJID: phoneJID},
			},
		}, nil),
	}

	evt := &events.Message{
		Info: types.MessageInfo{
			MessageSource: types.MessageSource{
				Sender: lidJID,
				Chat:   lidJID,
			},
			ID:       "mid-lid-2",
			PushName: "Arthur",
		},
		Message: &waE2E.Message{
			Conversation: proto.String("Should be dropped"),
		},
	}

	ch.handleIncoming(evt)

	ctx, cancel := context.WithTimeout(context.Background(), 200*time.Millisecond)
	defer cancel()

	select {
	case <-ctx.Done():
		// Expected — message was correctly dropped.
	case <-messageBus.InboundChan():
		t.Fatal("message should have been dropped — sender not in allow list")
	}
}
