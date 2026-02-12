package emulation

import (
	"errors"
	"testing"
)

func TestHarnessRunSequence_Deterministic(t *testing.T) {
	t.Parallel()

	target := &Target{
		Name:    "deterministic",
		Address: 0x15,
		Rules: []Rule{
			{
				Name:    "identify",
				Matcher: MatchPrimarySecondary(0x07, 0x04),
				Builder: BuildFunc(func(Frame) (ResponsePlan, error) {
					return ResponsePlan{
						DelayMillis: 8,
						Data:        []byte{0xB5},
					}, nil
				}),
			},
		},
	}

	harness := NewHarness(target)
	responses, err := harness.RunSequence([]QueryStep{
		{
			Frame: Frame{
				Source:    0x10,
				Target:    0x15,
				Primary:   0x07,
				Secondary: 0x04,
			},
		},
		{
			AdvanceMillis: 12,
			Frame: Frame{
				Source:    0x11,
				Target:    0x15,
				Primary:   0x07,
				Secondary: 0x04,
			},
		},
		{
			AdvanceMillis: 3,
			Frame: Frame{
				Source:    0x12,
				Target:    0x15,
				Primary:   0x07,
				Secondary: 0x04,
			},
		},
	})
	if err != nil {
		t.Fatalf("RunSequence() error = %v", err)
	}
	if len(responses) != 3 {
		t.Fatalf("len(responses) = %d; want 3", len(responses))
	}

	wantRequested := []uint32{0, 12, 15}
	wantRespondAt := []uint32{8, 20, 23}
	wantSources := []byte{0x10, 0x11, 0x12}
	for idx := range responses {
		response := responses[idx]
		if response.RequestedAtMillis != wantRequested[idx] {
			t.Fatalf("responses[%d].RequestedAtMillis = %d; want %d", idx, response.RequestedAtMillis, wantRequested[idx])
		}
		if response.RespondAtMillis != wantRespondAt[idx] {
			t.Fatalf("responses[%d].RespondAtMillis = %d; want %d", idx, response.RespondAtMillis, wantRespondAt[idx])
		}
		if response.Frame.Source != 0x15 || response.Frame.Target != wantSources[idx] {
			t.Fatalf("responses[%d].Frame source/target = 0x%02x/0x%02x; want 0x15/0x%02x", idx, response.Frame.Source, response.Frame.Target, wantSources[idx])
		}
	}
	if harness.NowMillis() != 15 {
		t.Fatalf("NowMillis() = %d; want 15", harness.NowMillis())
	}

	history := harness.History()
	if len(history) != len(responses) {
		t.Fatalf("len(history) = %d; want %d", len(history), len(responses))
	}
	for idx := range responses {
		if history[idx].RequestedAtMillis != responses[idx].RequestedAtMillis || history[idx].RespondAtMillis != responses[idx].RespondAtMillis {
			t.Fatalf("history[%d] does not match response[%d]", idx, idx)
		}
	}
}

func TestHarnessRunSequence_Errors(t *testing.T) {
	t.Parallel()

	target := &Target{
		Address: 0x15,
		Rules: []Rule{
			{
				Name:    "identify",
				Matcher: MatchPrimarySecondary(0x07, 0x04),
				Builder: BuildFunc(func(Frame) (ResponsePlan, error) {
					return ResponsePlan{
						DelayMillis: 8,
						Data:        []byte{0xB5},
					}, nil
				}),
			},
		},
	}

	cases := []struct {
		name    string
		harness *Harness
		steps   []QueryStep
		want    error
	}{
		{
			name:    "missing target",
			harness: NewHarness(nil),
			steps: []QueryStep{
				{
					Frame: Frame{},
				},
			},
			want: ErrInvalidConfiguration,
		},
		{
			name:    "query mismatch",
			harness: NewHarness(target),
			steps: []QueryStep{
				{
					Frame: Frame{
						Source:    0x10,
						Target:    0x08,
						Primary:   0x07,
						Secondary: 0x04,
					},
				},
			},
			want: ErrRequestTargetMismatch,
		},
	}

	for _, test := range cases {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			_, err := test.harness.RunSequence(test.steps)
			if !errors.Is(err, test.want) {
				t.Fatalf("RunSequence() error = %v; want %v", err, test.want)
			}
		})
	}
}

func TestHarnessHistory_DeepCopiesPayload(t *testing.T) {
	t.Parallel()

	target := &Target{
		Address: 0x15,
		Rules: []Rule{
			{
				Name:    "identify",
				Matcher: MatchPrimarySecondary(0x07, 0x04),
				Builder: BuildFunc(func(Frame) (ResponsePlan, error) {
					return ResponsePlan{
						DelayMillis: 8,
						Data:        []byte{0xB5},
					}, nil
				}),
			},
		},
	}

	harness := NewHarness(target)
	_, err := harness.Query(Frame{
		Source:    0x10,
		Target:    0x15,
		Primary:   0x07,
		Secondary: 0x04,
	})
	if err != nil {
		t.Fatalf("Query() error = %v", err)
	}

	history := harness.History()
	if len(history) != 1 {
		t.Fatalf("len(history) = %d; want 1", len(history))
	}
	history[0].Frame.Data[0] = 0x00

	freshHistory := harness.History()
	if got := freshHistory[0].Frame.Data[0]; got != 0xB5 {
		t.Fatalf("freshHistory[0].Frame.Data[0] = 0x%02x; want 0xb5", got)
	}
}

func TestValidateResponseEnvelope(t *testing.T) {
	t.Parallel()

	responses := []EmulatedResponse{
		{
			RequestedAtMillis: 0,
			RespondAtMillis:   8,
		},
		{
			RequestedAtMillis: 10,
			RespondAtMillis:   21,
		},
	}

	if err := ValidateResponseEnvelope(responses, ResponseEnvelope{
		MinDelayMillis: 5,
		MaxDelayMillis: 15,
	}); err != nil {
		t.Fatalf("ValidateResponseEnvelope() error = %v", err)
	}

	err := ValidateResponseEnvelope(responses, ResponseEnvelope{
		MinDelayMillis: 9,
		MaxDelayMillis: 15,
	})
	if !errors.Is(err, ErrTimingConstraint) {
		t.Fatalf("ValidateResponseEnvelope() error = %v; want %v", err, ErrTimingConstraint)
	}

	err = ValidateResponseEnvelope(responses, ResponseEnvelope{
		MinDelayMillis: 20,
		MaxDelayMillis: 15,
	})
	if !errors.Is(err, ErrInvalidConfiguration) {
		t.Fatalf("ValidateResponseEnvelope() error = %v; want %v", err, ErrInvalidConfiguration)
	}

	err = ValidateResponseEnvelope([]EmulatedResponse{
		{
			RequestedAtMillis: 10,
			RespondAtMillis:   9,
		},
	}, ResponseEnvelope{
		MinDelayMillis: 0,
	})
	if !errors.Is(err, ErrTimingConstraint) {
		t.Fatalf("ValidateResponseEnvelope() error = %v; want %v", err, ErrTimingConstraint)
	}
}
