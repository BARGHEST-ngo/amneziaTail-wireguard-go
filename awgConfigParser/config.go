package awg

import (
	"fmt"
	"strings"
)

type Awgconf struct {
	Enabled                 bool
	JunkPacketCount         int
	JunkPacketMinSize       int
	JunkPacketMaxSize       int
	InitPacketJunkSize      int         //s1
	RespPacketJunkSize      int         //2
	CookiePacketJunkSize    int         //s3
	TransPacketJunkSize     int         //s4
	InitPacketMagicHeader   MagicHeader //h1
	RespPacketMagicHeader   MagicHeader //h2
	CookiePacketMagicHeader MagicHeader //h3
	TransPacketMagicHeader  MagicHeader //h4
}

type MagicHeader struct {
	min uint32
	max uint32
}

// Min returns the minimum value of the magic header
func (m MagicHeader) Min() uint32 {
	return m.min
}

// Max returns the maximum value of the magic header
func (m MagicHeader) Max() uint32 {
	return m.max
}

const (
	messageTypeInit   = 1
	messageTypeResp   = 2
	messageTypeCookie = 3
	messageTypeTrans  = 4

	messageInitSize     = 148
	messageRespSize     = 92
	messageCookieSize   = 64
	messageTransMinSize = 32
)

func NewConfig() *Awgconf {
	return &Awgconf{}
}

func (c *Awgconf) IsEnabled() bool {
	if c == nil {
		return false
	}
	if !c.Enabled {
		return false
	}
	hasJunkPackets := c.JunkPacketCount > 0
	hasHeaderJunk := c.InitPacketJunkSize > 0 ||
		c.RespPacketJunkSize > 0 ||
		c.CookiePacketJunkSize > 0 ||
		c.TransPacketJunkSize > 0
	hasMagicHeaders := c.InitPacketMagicHeader.min > 4 ||
		c.RespPacketMagicHeader.min > 4 ||
		c.CookiePacketMagicHeader.min > 4 ||
		c.TransPacketMagicHeader.min > 4
	return hasJunkPackets || hasHeaderJunk || hasMagicHeaders
}

func (m MagicHeader) Validate(name string) error {
	//disabled
	if m.min == 0 && m.max == 0 {
		return nil
	}
	//min
	if m.min > 0 && m.min <= 4 {
		return fmt.Errorf("%s: Magic header must be > 4 or 0 (disabled) but got min=%d", name, m.min)
	}
	//order
	if m.max < m.min {
		return fmt.Errorf("%s: max must be >= min (min=%d, max=%d)", name, m.min, m.max)
	}
	return nil
}

// enabled=true
// jc=5
// jmin=20
// jmax=200
// s1=10
// s2=10
// s3=10
// s4=5
// h1=157
// h2=158
// h3=159
// h4=160
func (c *Awgconf) ParseConfg(configStr string) (*Awgconf, error) {
	cfg := NewConfig()
	lines := strings.Split(configStr, "\n")
	for _, raw := range lines {
		line := strings.TrimSpace(raw)
		if len(line) == 0 || strings.HasPrefix(line, "#") {
			continue
		}
		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			return nil, fmt.Errorf("invalid config line")
		}
		key := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])

		if key == "enabled" {
			cfg.Enabled = (value == "true")
			break
		}
	}

	if !cfg.Enabled {
		return cfg, nil
	}

	for _, raw := range lines {
		line := strings.TrimSpace(raw)
		if len(line) == 0 || strings.HasPrefix(line, "#") {
			continue
		}
		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			return nil, fmt.Errorf("invalid config line")
		}
		key := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])
		switch key {
		case "enabled":
			cfg.Enabled = (value == "true")
		case "jc":
			fmt.Sscanf(value, "%d", &cfg.JunkPacketCount)
		case "jmin":
			fmt.Sscanf(value, "%d", &cfg.JunkPacketMinSize)
		case "jmax":
			fmt.Sscanf(value, "%d", &cfg.JunkPacketMaxSize)
		case "s1":
			fmt.Sscanf(value, "%d", &cfg.InitPacketJunkSize)
		case "s2":
			fmt.Sscanf(value, "%d", &cfg.RespPacketJunkSize)
		case "s3":
			fmt.Sscanf(value, "%d", &cfg.CookiePacketJunkSize)
		case "s4":
			fmt.Sscanf(value, "%d", &cfg.TransPacketJunkSize)
		case "h1":
			var val uint32
			fmt.Sscanf(value, "%d", &val)
			cfg.InitPacketMagicHeader = MagicHeader{min: val, max: val}
		case "h2":
			var val uint32
			fmt.Sscanf(value, "%d", &val)
			cfg.RespPacketMagicHeader = MagicHeader{min: val, max: val}
		case "h3":
			var val uint32
			fmt.Sscanf(value, "%d", &val)
			cfg.CookiePacketMagicHeader = MagicHeader{min: val, max: val}
		case "h4":
			var val uint32
			fmt.Sscanf(value, "%d", &val)
			cfg.TransPacketMagicHeader = MagicHeader{min: val, max: val}
		}
	}
	return cfg, nil
}
