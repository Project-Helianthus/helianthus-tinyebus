package emulation

import (
	"fmt"
	"strings"
)

const (
	identifyOnlyDeviceIDLength = 5

	DefaultVR71Address      = byte(0x26)
	DefaultVR71Manufacturer = byte(0xB5)
	DefaultVR71DeviceID     = "VR_71"
	DefaultVR71Software     = uint16(0x0100)
	DefaultVR71Hardware     = uint16(0x5904)
)

var (
	defaultIdentifyOnlyResponseDelayMillis = uint32(8)
	defaultIdentifyOnlyTiming              = TimingConstraints{
		MinResponseDelayMillis: 5,
		MaxResponseDelayMillis: 30,
	}
)

type IdentifyOnlyProfile struct {
	Name                string
	Address             byte
	Manufacturer        byte
	DeviceID            string
	Software            uint16
	Hardware            uint16
	ResponseDelayMillis uint32
	Timing              TimingConstraints
}

func PresetVR90IdentifyOnlyProfile() IdentifyOnlyProfile {
	return IdentifyOnlyProfile{
		Name:                fmt.Sprintf("vr90-minimal-0x%02x", DefaultVR90Address),
		Address:             DefaultVR90Address,
		Manufacturer:        DefaultVR90Manufacturer,
		DeviceID:            DefaultVR90DeviceID,
		Software:            DefaultVR90Software,
		Hardware:            DefaultVR90Hardware,
		ResponseDelayMillis: defaultVR90ResponseDelayMillis,
		Timing:              defaultVR90Timing,
	}
}

func PresetVR71IdentifyOnlyProfile() IdentifyOnlyProfile {
	return IdentifyOnlyProfile{
		Name:                fmt.Sprintf("vr71-minimal-0x%02x", DefaultVR71Address),
		Address:             DefaultVR71Address,
		Manufacturer:        DefaultVR71Manufacturer,
		DeviceID:            DefaultVR71DeviceID,
		Software:            DefaultVR71Software,
		Hardware:            DefaultVR71Hardware,
		ResponseDelayMillis: defaultIdentifyOnlyResponseDelayMillis,
		Timing:              defaultIdentifyOnlyTiming,
	}
}

func NewIdentifyOnlyTarget(profile IdentifyOnlyProfile) (*Target, error) {
	normalized, err := normalizeIdentifyOnlyProfile(profile)
	if err != nil {
		return nil, err
	}
	return &Target{
		Name:          normalized.Name,
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

func normalizeIdentifyOnlyProfile(profile IdentifyOnlyProfile) (IdentifyOnlyProfile, error) {
	profile.Name = strings.TrimSpace(profile.Name)
	if profile.Name == "" {
		profile.Name = fmt.Sprintf("identify-only-0x%02x", profile.Address)
	}
	if profile.Address == 0 {
		return IdentifyOnlyProfile{}, fmt.Errorf("identify-only profile empty address: %w", ErrInvalidConfiguration)
	}
	if profile.Manufacturer == 0 {
		return IdentifyOnlyProfile{}, fmt.Errorf("identify-only profile empty manufacturer: %w", ErrInvalidConfiguration)
	}

	trimmedID := strings.TrimSpace(profile.DeviceID)
	if trimmedID == "" {
		return IdentifyOnlyProfile{}, fmt.Errorf("identify-only profile empty device id: %w", ErrInvalidConfiguration)
	}
	if len(trimmedID) > identifyOnlyDeviceIDLength {
		trimmedID = trimmedID[:identifyOnlyDeviceIDLength]
	}
	profile.DeviceID = trimmedID

	if profile.ResponseDelayMillis == 0 {
		profile.ResponseDelayMillis = defaultIdentifyOnlyResponseDelayMillis
	}
	if !profile.Timing.active() {
		profile.Timing = defaultIdentifyOnlyTiming
	}
	if err := profile.Timing.validate(profile.ResponseDelayMillis); err != nil {
		return IdentifyOnlyProfile{}, err
	}

	return profile, nil
}

func (profile IdentifyOnlyProfile) identificationPayload() []byte {
	deviceID := fmt.Sprintf("%-*s", identifyOnlyDeviceIDLength, profile.DeviceID)
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
