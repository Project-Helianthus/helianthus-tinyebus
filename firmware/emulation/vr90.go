package emulation

import (
	"fmt"
	"strings"
)

const (
	DefaultVR90Address      = byte(0x15)
	DefaultVR90Manufacturer = byte(0xB5)
	DefaultVR90DeviceID     = "B7V00"
	DefaultVR90Software     = uint16(0x0422)
	DefaultVR90Hardware     = uint16(0x5503)
	vr90DeviceIDLength      = 5
)

var (
	defaultVR90ResponseDelayMillis = uint32(8)
	defaultVR90Timing              = TimingConstraints{
		MinResponseDelayMillis: 5,
		MaxResponseDelayMillis: 30,
	}
)

type VR90Profile struct {
	Address             byte
	Manufacturer        byte
	DeviceID            string
	Software            uint16
	Hardware            uint16
	ResponseDelayMillis uint32
	Timing              TimingConstraints
}

func DefaultVR90Profile() VR90Profile {
	return VR90Profile{
		Address:             DefaultVR90Address,
		Manufacturer:        DefaultVR90Manufacturer,
		DeviceID:            DefaultVR90DeviceID,
		Software:            DefaultVR90Software,
		Hardware:            DefaultVR90Hardware,
		ResponseDelayMillis: defaultVR90ResponseDelayMillis,
		Timing:              defaultVR90Timing,
	}
}

func NewVR90Target(profile VR90Profile) (*Target, error) {
	normalized, err := normalizeVR90Profile(profile)
	if err != nil {
		return nil, err
	}
	return &Target{
		Name:          fmt.Sprintf("vr90-minimal-0x%02x", normalized.Address),
		Address:       normalized.Address,
		DefaultTiming: normalized.Timing,
		Rules: []Rule{
			{
				Name:    "identify",
				Matcher: MatchPrimarySecondary(0x07, 0x04),
				Builder: BuildFunc(func(_ Frame) (ResponsePlan, error) {
					return ResponsePlan{
						DelayMillis: normalized.ResponseDelayMillis,
						Data:        normalized.identificationPayload(),
					}, nil
				}),
			},
		},
	}, nil
}

func normalizeVR90Profile(profile VR90Profile) (VR90Profile, error) {
	if profile.Address == 0 {
		profile.Address = DefaultVR90Address
	}
	if profile.Manufacturer == 0 {
		profile.Manufacturer = DefaultVR90Manufacturer
	}
	if strings.TrimSpace(profile.DeviceID) == "" {
		profile.DeviceID = DefaultVR90DeviceID
	}
	if profile.Software == 0 {
		profile.Software = DefaultVR90Software
	}
	if profile.Hardware == 0 {
		profile.Hardware = DefaultVR90Hardware
	}
	if profile.ResponseDelayMillis == 0 {
		profile.ResponseDelayMillis = defaultVR90ResponseDelayMillis
	}
	if !profile.Timing.active() {
		profile.Timing = defaultVR90Timing
	}
	if err := profile.Timing.validate(profile.ResponseDelayMillis); err != nil {
		return VR90Profile{}, err
	}

	trimmedID := strings.TrimSpace(profile.DeviceID)
	if trimmedID == "" {
		return VR90Profile{}, fmt.Errorf("vr90 profile empty device id: %w", ErrInvalidConfiguration)
	}
	if len(trimmedID) > vr90DeviceIDLength {
		trimmedID = trimmedID[:vr90DeviceIDLength]
	}
	profile.DeviceID = trimmedID

	return profile, nil
}

func (profile VR90Profile) identificationPayload() []byte {
	deviceID := fmt.Sprintf("%-*s", vr90DeviceIDLength, profile.DeviceID)
	return []byte{
		profile.Manufacturer,
		deviceID[0],
		deviceID[1],
		deviceID[2],
		deviceID[3],
		deviceID[4],
		byte(profile.Software >> 8),
		byte(profile.Software & 0xFF),
		byte(profile.Hardware >> 8),
		byte(profile.Hardware & 0xFF),
	}
}
