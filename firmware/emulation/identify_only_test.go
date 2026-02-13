package emulation

import (
	"bytes"
	"errors"
	"testing"
)

func TestNewIdentifyOnlyTarget_PayloadAndTiming(t *testing.T) {
	t.Parallel()

	target, err := NewIdentifyOnlyTarget(IdentifyOnlyProfile{
		Name:                "  custom-identify  ",
		Address:             0x30,
		Manufacturer:        0xAA,
		DeviceID:            "VR71LONG",
		Software:            0x1234,
		Hardware:            0x5678,
		ResponseDelayMillis: 9,
		Timing: TimingConstraints{
			MinResponseDelayMillis: 5,
			MaxResponseDelayMillis: 15,
		},
	})
	if err != nil {
		t.Fatalf("NewIdentifyOnlyTarget() error = %v", err)
	}

	if target.Name != "custom-identify" {
		t.Fatalf("Name = %q; want %q", target.Name, "custom-identify")
	}
	if target.Address != 0x30 {
		t.Fatalf("Address = 0x%02x; want 0x30", target.Address)
	}

	response, err := target.Emulate(RequestEvent{
		AtMillis: 20,
		Frame: Frame{
			Source:    0x10,
			Target:    0x30,
			Primary:   0x07,
			Secondary: 0x04,
		},
	})
	if err != nil {
		t.Fatalf("Emulate() error = %v", err)
	}
	if response.RespondAtMillis != 29 {
		t.Fatalf("RespondAtMillis = %d; want 29", response.RespondAtMillis)
	}

	wantData := []byte{
		0xAA, 'V', 'R', '7', '1', 'L',
		0x12, 0x34,
		0x56, 0x78,
	}
	if !bytes.Equal(response.Frame.Data, wantData) {
		t.Fatalf("Frame data = %x; want %x", response.Frame.Data, wantData)
	}
}

func TestNewIdentifyOnlyTarget_Errors(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name    string
		profile IdentifyOnlyProfile
		want    error
	}{
		{
			name: "empty address",
			profile: IdentifyOnlyProfile{
				Manufacturer: 0xB5,
				DeviceID:     "VR_71",
			},
			want: ErrInvalidConfiguration,
		},
		{
			name: "empty manufacturer",
			profile: IdentifyOnlyProfile{
				Address:  0x26,
				DeviceID: "VR_71",
			},
			want: ErrInvalidConfiguration,
		},
		{
			name: "empty device id",
			profile: IdentifyOnlyProfile{
				Address:      0x26,
				Manufacturer: 0xB5,
				DeviceID:     "  ",
			},
			want: ErrInvalidConfiguration,
		},
		{
			name: "timing violation",
			profile: IdentifyOnlyProfile{
				Address:             0x26,
				Manufacturer:        0xB5,
				DeviceID:            "VR_71",
				ResponseDelayMillis: 2,
				Timing: TimingConstraints{
					MinResponseDelayMillis: 5,
					MaxResponseDelayMillis: 15,
				},
			},
			want: ErrTimingConstraint,
		},
	}

	for _, test := range cases {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			_, err := NewIdentifyOnlyTarget(test.profile)
			if !errors.Is(err, test.want) {
				t.Fatalf("NewIdentifyOnlyTarget() error = %v; want %v", err, test.want)
			}
		})
	}
}

func TestIdentifyOnlyProfile_Presets(t *testing.T) {
	t.Parallel()

	vr90 := PresetVR90IdentifyOnlyProfile()
	if vr90.Address != DefaultVR90Address ||
		vr90.Manufacturer != DefaultVR90Manufacturer ||
		vr90.DeviceID != DefaultVR90DeviceID ||
		vr90.Software != DefaultVR90Software ||
		vr90.Hardware != DefaultVR90Hardware {
		t.Fatalf("VR90 preset mismatch: %+v", vr90)
	}

	vr71 := PresetVR71IdentifyOnlyProfile()
	if vr71.Address != DefaultVR71Address ||
		vr71.Manufacturer != DefaultVR71Manufacturer ||
		vr71.DeviceID != DefaultVR71DeviceID ||
		vr71.Software != DefaultVR71Software ||
		vr71.Hardware != DefaultVR71Hardware {
		t.Fatalf("VR_71 preset mismatch: %+v", vr71)
	}
}

func TestIdentifyOnlyTarget_TimingAndHistorySafety(t *testing.T) {
	t.Parallel()

	profile := PresetVR71IdentifyOnlyProfile()
	target, err := NewIdentifyOnlyTarget(profile)
	if err != nil {
		t.Fatalf("NewIdentifyOnlyTarget() error = %v", err)
	}

	harness := NewHarness(target)
	frame := Frame{
		Source:    0x10,
		Target:    profile.Address,
		Primary:   0x07,
		Secondary: 0x04,
	}

	first, err := harness.Query(frame)
	if err != nil {
		t.Fatalf("Query(first) error = %v", err)
	}
	if first.RespondAtMillis != 8 {
		t.Fatalf("first.RespondAtMillis = %d; want 8", first.RespondAtMillis)
	}

	snapshot := harness.History()
	if len(snapshot) != 1 {
		t.Fatalf("len(snapshot) = %d; want 1", len(snapshot))
	}
	snapshot[0].Frame.Data[0] = 0x00

	harness.AdvanceMillis(3)
	second, err := harness.Query(frame)
	if err != nil {
		t.Fatalf("Query(second) error = %v", err)
	}

	if second.RequestedAtMillis != 3 {
		t.Fatalf("second.RequestedAtMillis = %d; want 3", second.RequestedAtMillis)
	}
	if second.RespondAtMillis != 11 {
		t.Fatalf("second.RespondAtMillis = %d; want 11", second.RespondAtMillis)
	}

	wantData := profile.identificationPayload()
	if !bytes.Equal(second.Frame.Data, wantData) {
		t.Fatalf("second.Frame.Data = %x; want %x", second.Frame.Data, wantData)
	}

	freshHistory := harness.History()
	if len(freshHistory) != 2 {
		t.Fatalf("len(freshHistory) = %d; want 2", len(freshHistory))
	}
	if !bytes.Equal(freshHistory[0].Frame.Data, wantData) {
		t.Fatalf("freshHistory[0].Frame.Data = %x; want %x", freshHistory[0].Frame.Data, wantData)
	}

	if err := ValidateResponseEnvelope(freshHistory, ResponseEnvelope{
		MinDelayMillis: profile.Timing.MinResponseDelayMillis,
		MaxDelayMillis: profile.Timing.MaxResponseDelayMillis,
	}); err != nil {
		t.Fatalf("ValidateResponseEnvelope() error = %v", err)
	}
}

func TestSmokeVR71IdentifyOnlyProfile(t *testing.T) {
	profile := PresetVR71IdentifyOnlyProfile()
	target, err := NewIdentifyOnlyTarget(profile)
	if err != nil {
		t.Fatalf("NewIdentifyOnlyTarget() error = %v", err)
	}

	harness := NewHarness(target)
	responses, err := harness.RunSequence([]QueryStep{
		{
			Frame: Frame{
				Source:    0x10,
				Target:    profile.Address,
				Primary:   0x07,
				Secondary: 0x04,
			},
		},
	})
	if err != nil {
		t.Fatalf("RunSequence() error = %v", err)
	}
	if len(responses) != 1 {
		t.Fatalf("len(responses) = %d; want 1", len(responses))
	}
	if gotID := string(responses[0].Frame.Data[1:6]); gotID != DefaultVR71DeviceID {
		t.Fatalf("DeviceID = %q; want %q", gotID, DefaultVR71DeviceID)
	}

	if err := ValidateResponseEnvelope(responses, ResponseEnvelope{
		MinDelayMillis: profile.Timing.MinResponseDelayMillis,
		MaxDelayMillis: profile.Timing.MaxResponseDelayMillis,
	}); err != nil {
		t.Fatalf("ValidateResponseEnvelope() error = %v", err)
	}
}
