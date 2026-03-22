package whatsapp

import (
	"context"
	"testing"
	"time"

	"github.com/sipeed/picoclaw/pkg/bus"
	"github.com/sipeed/picoclaw/pkg/channels"
	"github.com/sipeed/picoclaw/pkg/config"
)

func TestHandleIncomingMessage_DMOnly_DropsGroupMessages(t *testing.T) {
	messageBus := bus.NewMessageBus()
	ch := &WhatsAppChannel{
		BaseChannel: channels.NewBaseChannel("whatsapp", config.WhatsAppConfig{DMOnly: true}, messageBus, nil),
		config:      config.WhatsAppConfig{DMOnly: true},
		ctx:         context.Background(),
	}

	// Group message: chat ID differs from sender ID
	ch.handleIncomingMessage(map[string]any{
		"type":    "message",
		"id":      "mid1",
		"from":    "user1",
		"chat":    "group-chat-123",
		"content": "hello group",
	})

	select {
	case msg := <-messageBus.InboundChan():
		t.Fatalf("expected group message to be dropped, got: %+v", msg)
	case <-time.After(50 * time.Millisecond):
		// Good — nothing published
	}
}

func TestHandleIncomingMessage_DMOnly_AllowsDirectMessages(t *testing.T) {
	messageBus := bus.NewMessageBus()
	ch := &WhatsAppChannel{
		BaseChannel: channels.NewBaseChannel("whatsapp", config.WhatsAppConfig{DMOnly: true}, messageBus, nil),
		config:      config.WhatsAppConfig{DMOnly: true},
		ctx:         context.Background(),
	}

	// DM: chat ID equals sender ID
	ch.handleIncomingMessage(map[string]any{
		"type":    "message",
		"id":      "mid2",
		"from":    "user1",
		"chat":    "user1",
		"content": "hello dm",
	})

	select {
	case msg := <-messageBus.InboundChan():
		if msg.Content != "hello dm" {
			t.Fatalf("content=%q, want %q", msg.Content, "hello dm")
		}
	case <-time.After(time.Second):
		t.Fatal("expected DM to be forwarded")
	}
}

func TestHandleIncomingMessage_DMOnlyFalse_AllowsGroupMessages(t *testing.T) {
	messageBus := bus.NewMessageBus()
	ch := &WhatsAppChannel{
		BaseChannel: channels.NewBaseChannel("whatsapp", config.WhatsAppConfig{DMOnly: false}, messageBus, nil),
		config:      config.WhatsAppConfig{DMOnly: false},
		ctx:         context.Background(),
	}

	// Group message should go through when dm_only is false
	ch.handleIncomingMessage(map[string]any{
		"type":    "message",
		"id":      "mid3",
		"from":    "user1",
		"chat":    "group-chat-456",
		"content": "hello group",
	})

	select {
	case msg := <-messageBus.InboundChan():
		if msg.Content != "hello group" {
			t.Fatalf("content=%q, want %q", msg.Content, "hello group")
		}
	case <-time.After(time.Second):
		t.Fatal("expected group message to be forwarded when dm_only is false")
	}
}
