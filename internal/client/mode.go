package client

import (
	"encoding/json"
	"errors"
)

type Mode uint8

const (
	DirectMode Mode = iota
	RuleMode
	GlobalMode
)

func (m Mode) MarshalJSON() ([]byte, error) {
	return json.Marshal(m.String())
}

func (m *Mode) UnmarshalJSON(data []byte) error {
	var mode string
	err := json.Unmarshal(data, &mode)
	if err != nil {
		return err
	}
	switch mode {
	case "direct":
		*m = DirectMode
	case "rule":
		*m = RuleMode
	case "global":
		*m = GlobalMode
	default:
		return errors.New("unknown mode")
	}
	return nil
}

func (m Mode) String() string {
	switch m {
	case DirectMode:
		return "direct"
	case RuleMode:
		return "rule"
	case GlobalMode:
		return "global"
	default:
		return "unknown"
	}
}
