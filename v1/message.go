package irkit

import "fmt"

// SignalInfo represetns IRKit signal infomation.
type SignalInfo struct {
	// Hostname is IRKit device hostname.
	Hostname string `json:"hostname"`

	// Deviceid is IRKit device identifier.
	Deviceid string `json:"deviceid"`

	// Message is IRKit signal message.
	Message Message `json:"message"`
}

// Message represents IRKit signal message.
type Message struct {
	// Format is format of signal. "raw" only.
	Format string `json:"format"`

	// Freq is IRKit sub-carrier frequency. 38 or 40 only. [kHz]
	Freq int `json:"freq"`

	// Data is IRkit signal consists of ON/OFF of sub carrier frequency.
	// IRKit measures On to Off, Off to On interval using a 2MHz counter.
	// data value is an array of those intervals
	Data []int `json:"data"`
}

func (m *Message) validate() error {
	if m.Format != "raw" {
		return fmt.Errorf("format must be raw")
	}

	if m.Freq != 38 && m.Freq != 40 {
		return fmt.Errorf("freq must 38 or 40")
	}

	if len(m.Data) == 0 {
		return fmt.Errorf("empty data")
	}

	return nil
}
