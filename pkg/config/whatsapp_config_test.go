package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestDefaultConfig_WhatsAppAckReactionDisabled(t *testing.T) {
	cfg := DefaultConfig()
	if cfg.Channels.WhatsApp.AckReaction {
		t.Fatal("AckReaction should be false by default")
	}
}

func TestDefaultConfig_WhatsAppDMOnlyDisabled(t *testing.T) {
	cfg := DefaultConfig()
	if cfg.Channels.WhatsApp.DMOnly {
		t.Fatal("DMOnly should be false by default")
	}
}

func TestLoadConfig_WhatsAppAckReactionEnabled(t *testing.T) {
	dir := t.TempDir()
	configPath := filepath.Join(dir, "config.json")
	data := `{
		"channels": {
			"whatsapp": {
				"ack_reaction": true
			}
		},
		"model_list": [{"model_name":"test","model":"openai/gpt-4","api_key":"sk-test"}]
	}`
	if err := os.WriteFile(configPath, []byte(data), 0o600); err != nil {
		t.Fatalf("setup: %v", err)
	}
	cfg, err := LoadConfig(configPath)
	if err != nil {
		t.Fatalf("LoadConfig() error: %v", err)
	}
	if !cfg.Channels.WhatsApp.AckReaction {
		t.Fatal("AckReaction should be true when set in config")
	}
}

func TestLoadConfig_WhatsAppAckReactionDefaultsFalse(t *testing.T) {
	dir := t.TempDir()
	configPath := filepath.Join(dir, "config.json")
	data := `{
		"channels": {
			"whatsapp": {
				"enabled": true
			}
		},
		"model_list": [{"model_name":"test","model":"openai/gpt-4","api_key":"sk-test"}]
	}`
	if err := os.WriteFile(configPath, []byte(data), 0o600); err != nil {
		t.Fatalf("setup: %v", err)
	}
	cfg, err := LoadConfig(configPath)
	if err != nil {
		t.Fatalf("LoadConfig() error: %v", err)
	}
	if cfg.Channels.WhatsApp.AckReaction {
		t.Fatal("AckReaction should remain false when not set in config")
	}
}

func TestLoadConfig_WhatsAppDMOnlyEnabled(t *testing.T) {
	dir := t.TempDir()
	configPath := filepath.Join(dir, "config.json")
	data := `{
		"channels": {
			"whatsapp": {
				"dm_only": true
			}
		},
		"model_list": [{"model_name":"test","model":"openai/gpt-4","api_key":"sk-test"}]
	}`
	if err := os.WriteFile(configPath, []byte(data), 0o600); err != nil {
		t.Fatalf("setup: %v", err)
	}
	cfg, err := LoadConfig(configPath)
	if err != nil {
		t.Fatalf("LoadConfig() error: %v", err)
	}
	if !cfg.Channels.WhatsApp.DMOnly {
		t.Fatal("DMOnly should be true when set in config")
	}
}

func TestDefaultConfig_WhatsAppSelfChatPrefix(t *testing.T) {
	cfg := DefaultConfig()
	if cfg.Channels.WhatsApp.SelfChatPrefix != "[picoclaw]" {
		t.Fatalf("SelfChatPrefix = %q, want %q", cfg.Channels.WhatsApp.SelfChatPrefix, "[picoclaw]")
	}
}

func TestLoadConfig_WhatsAppSelfChatPrefixCustom(t *testing.T) {
	dir := t.TempDir()
	configPath := filepath.Join(dir, "config.json")
	data := `{
		"channels": {
			"whatsapp": {
				"self_chat_prefix": "[bot]"
			}
		},
		"model_list": [{"model_name":"test","model":"openai/gpt-4","api_key":"sk-test"}]
	}`
	if err := os.WriteFile(configPath, []byte(data), 0o600); err != nil {
		t.Fatalf("setup: %v", err)
	}
	cfg, err := LoadConfig(configPath)
	if err != nil {
		t.Fatalf("LoadConfig() error: %v", err)
	}
	if cfg.Channels.WhatsApp.SelfChatPrefix != "[bot]" {
		t.Fatalf("SelfChatPrefix = %q, want %q", cfg.Channels.WhatsApp.SelfChatPrefix, "[bot]")
	}
}

func TestLoadConfig_WhatsAppSelfChatPrefixDisabled(t *testing.T) {
	dir := t.TempDir()
	configPath := filepath.Join(dir, "config.json")
	data := `{
		"channels": {
			"whatsapp": {
				"self_chat_prefix": ""
			}
		},
		"model_list": [{"model_name":"test","model":"openai/gpt-4","api_key":"sk-test"}]
	}`
	if err := os.WriteFile(configPath, []byte(data), 0o600); err != nil {
		t.Fatalf("setup: %v", err)
	}
	cfg, err := LoadConfig(configPath)
	if err != nil {
		t.Fatalf("LoadConfig() error: %v", err)
	}
	if cfg.Channels.WhatsApp.SelfChatPrefix != "" {
		t.Fatalf("SelfChatPrefix = %q, want empty", cfg.Channels.WhatsApp.SelfChatPrefix)
	}
}
