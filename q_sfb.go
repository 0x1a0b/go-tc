package tc

import (
	"bytes"
	"encoding/binary"
	"fmt"

	"github.com/mdlayher/netlink"
)

const (
	tcaSfbUnspec = iota
	tcaSfbParms
)

// Sfb contains attributes of the SBF discipline
type Sfb struct {
	Parms *SfbQopt
}

// unmarshalSfb parses the Sfb-encoded data and stores the result in the value pointed to by info.
func unmarshalSfb(data []byte, info *Sfb) error {
	ad, err := netlink.NewAttributeDecoder(data)
	if err != nil {
		return err
	}
	ad.ByteOrder = nativeEndian
	for ad.Next() {
		switch ad.Type() {
		case tcaSfbParms:
			opt := &SfbQopt{}
			if err := extractSfbQopt(ad.Bytes(), opt); err != nil {
				return err
			}
			info.Parms = opt
		default:
			return fmt.Errorf("extractSbfOptions()\t%d\n\t%v", ad.Type(), ad.Bytes())
		}
	}
	return nil
}

// marshalSbf returns the binary encoding of Sfb
func marshalSfb(info *Sfb) ([]byte, error) {
	options := []tcOption{}

	if info == nil {
		return []byte{}, fmt.Errorf("Sfb options are missing")
	}

	if info.Parms != nil {
		data, err := validateSfbQopt(info.Parms)
		if err != nil {
			return []byte{}, err
		}
		options = append(options, tcOption{Interpretation: vtBytes, Type: tcaSfbParms, Data: data})
	}

	// TODO: improve logic and check combinations
	return marshalAttributes(options)
}

// SfbQopt from include/uapi/linux/pkt_sched.h
type SfbQopt struct {
	RehashInterval uint32 // in ms
	WarmupTime     uint32 //  in ms
	Max            uint32
	BinSize        uint32
	Increment      uint32
	Decrement      uint32
	Limit          uint32
	PenaltyRate    uint32
	PenaltyBurst   uint32
}

func extractSfbQopt(data []byte, info *SfbQopt) error {
	b := bytes.NewReader(data)
	return binary.Read(b, nativeEndian, info)
}

func validateSfbQopt(info *SfbQopt) ([]byte, error) {
	var buf bytes.Buffer
	err := binary.Write(&buf, nativeEndian, *info)
	return buf.Bytes(), err
}