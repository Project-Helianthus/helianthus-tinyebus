package emulation

import (
	"bytes"
	"errors"
	"testing"
)

func TestTargetEmulate_MirrorResponseFrame(t *testing.T) {
	t.Parallel()

	responsePayload := []byte{0xB5, 0x42}
	target := &Target{
		Name:    "test-target",
		Address: 0x15,
		Rules: []Rule{
			{
				Name:    "identify",
				Matcher: MatchPrimarySecondary(0x07, 0x04),
				Builder: BuildFunc(func(Frame) (ResponsePlan, error) {
					return ResponsePlan{
						DelayMillis: 9,
						Data:        responsePayload,
					}, nil
				}),
			},
		},
	}

	request := Frame{
		Source:    0x10,
		Target:    0x15,
		Primary:   0x07,
		Secondary: 0x04,
	}
	got, err := target.Emulate(RequestEvent{AtMillis: 50, Frame: request})
	if err != nil {
		t.Fatalf("Emulate() error = %v", err)
	}
	if got.Rule != "identify" {
		t.Fatalf("Rule = %q; want identify", got.Rule)
	}
	if got.RequestedAtMillis != 50 {
		t.Fatalf("RequestedAtMillis = %d; want 50", got.RequestedAtMillis)
	}
	if got.RespondAtMillis != 59 {
		t.Fatalf("RespondAtMillis = %d; want 59", got.RespondAtMillis)
	}
	if got.Frame.Source != 0x15 || got.Frame.Target != 0x10 {
		t.Fatalf("Frame source/target = 0x%02x/0x%02x; want 0x15/0x10", got.Frame.Source, got.Frame.Target)
	}
	if got.Frame.Primary != 0x07 || got.Frame.Secondary != 0x04 {
		t.Fatalf("Frame PB/SB = 0x%02x/0x%02x; want 0x07/0x04", got.Frame.Primary, got.Frame.Secondary)
	}
	if !bytes.Equal(got.Frame.Data, responsePayload) {
		t.Fatalf("Frame data = %x; want %x", got.Frame.Data, responsePayload)
	}

	responsePayload[0] = 0x00
	if got.Frame.Data[0] != 0xB5 {
		t.Fatalf("Frame data was not copied, got %x", got.Frame.Data)
	}
}

func TestTargetEmulate_Errors(t *testing.T) {
	t.Parallel()

	validRequest := Frame{
		Source:    0x10,
		Target:    0x15,
		Primary:   0x07,
		Secondary: 0x04,
	}

	cases := []struct {
		name   string
		target *Target
		event  RequestEvent
		want   error
	}{
		{
			name:   "nil target",
			target: nil,
			event: RequestEvent{
				Frame: validRequest,
			},
			want: ErrInvalidConfiguration,
		},
		{
			name: "target mismatch",
			target: &Target{
				Address: 0x08,
			},
			event: RequestEvent{
				Frame: validRequest,
			},
			want: ErrRequestTargetMismatch,
		},
		{
			name: "missing matcher",
			target: &Target{
				Address: 0x15,
				Rules: []Rule{
					{Name: "broken"},
				},
			},
			event: RequestEvent{
				Frame: validRequest,
			},
			want: ErrInvalidConfiguration,
		},
		{
			name: "missing builder",
			target: &Target{
				Address: 0x15,
				Rules: []Rule{
					{
						Name:    "broken",
						Matcher: MatchPrimarySecondary(0x07, 0x04),
					},
				},
			},
			event: RequestEvent{
				Frame: validRequest,
			},
			want: ErrInvalidConfiguration,
		},
		{
			name: "no matching rule",
			target: &Target{
				Address: 0x15,
				Rules: []Rule{
					{
						Name:    "other",
						Matcher: MatchPrimarySecondary(0x01, 0x02),
						Builder: BuildFunc(func(Frame) (ResponsePlan, error) {
							return ResponsePlan{}, nil
						}),
					},
				},
			},
			event: RequestEvent{
				Frame: validRequest,
			},
			want: ErrNoMatchingRule,
		},
		{
			name: "timing constraint violation",
			target: &Target{
				Address: 0x15,
				DefaultTiming: TimingConstraints{
					MinResponseDelayMillis: 10,
					MaxResponseDelayMillis: 15,
				},
				Rules: []Rule{
					{
						Name:    "identify",
						Matcher: MatchPrimarySecondary(0x07, 0x04),
						Builder: BuildFunc(func(Frame) (ResponsePlan, error) {
							return ResponsePlan{
								DelayMillis: 2,
								Data:        []byte{0x00},
							}, nil
						}),
					},
				},
			},
			event: RequestEvent{
				Frame: validRequest,
			},
			want: ErrTimingConstraint,
		},
	}

	for _, test := range cases {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			_, err := test.target.Emulate(test.event)
			if !errors.Is(err, test.want) {
				t.Fatalf("Emulate() error = %v; want %v", err, test.want)
			}
		})
	}
}

func TestMatchPrimarySecondaryWithPrefix(t *testing.T) {
	t.Parallel()

	matcher := MatchPrimarySecondaryWithPrefix(0x07, 0x04, []byte{0x01, 0x02})

	tests := []struct {
		name  string
		frame Frame
		want  bool
	}{
		{
			name: "matches exact prefix",
			frame: Frame{
				Primary:   0x07,
				Secondary: 0x04,
				Data:      []byte{0x01, 0x02, 0x03},
			},
			want: true,
		},
		{
			name: "fails on wrong prefix",
			frame: Frame{
				Primary:   0x07,
				Secondary: 0x04,
				Data:      []byte{0x01, 0xFF, 0x03},
			},
			want: false,
		},
		{
			name: "fails on short data",
			frame: Frame{
				Primary:   0x07,
				Secondary: 0x04,
				Data:      []byte{0x01},
			},
			want: false,
		},
		{
			name: "fails on different pb sb",
			frame: Frame{
				Primary:   0x07,
				Secondary: 0x05,
				Data:      []byte{0x01, 0x02},
			},
			want: false,
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			if got := matcher.Match(test.frame); got != test.want {
				t.Fatalf("Match() = %v; want %v", got, test.want)
			}
		})
	}
}
