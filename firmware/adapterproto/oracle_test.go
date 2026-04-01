package adapterproto

import (
	"encoding/json"
	"testing"
)

func TestDefaultOracleReportIsDeterministic(t *testing.T) {
	t.Parallel()

	first := DefaultOracleReport()
	second := DefaultOracleReport()

	firstJSON, err := json.Marshal(first)
	if err != nil {
		t.Fatalf("json.Marshal(first) error = %v", err)
	}
	secondJSON, err := json.Marshal(second)
	if err != nil {
		t.Fatalf("json.Marshal(second) error = %v", err)
	}

	if string(firstJSON) != string(secondJSON) {
		t.Fatalf("DefaultOracleReport() was not deterministic")
	}
}

func TestDefaultOracleReportIncludesInfoAndScanSamples(t *testing.T) {
	t.Parallel()

	report := DefaultOracleReport()

	if len(report.INFO.Queries) != 8 {
		t.Fatalf("report.INFO.Queries len = %d; want 8", len(report.INFO.Queries))
	}
	if !report.INFO.Queries[0].Modeled || report.INFO.Queries[0].Version == nil {
		t.Fatalf("report.INFO.Queries[0] = %+v; want modeled version sample", report.INFO.Queries[0])
	}
	if !report.INFO.Queries[6].Modeled || report.INFO.Queries[6].ResetInfo == nil {
		t.Fatalf("report.INFO.Queries[6] = %+v; want modeled reset-info sample", report.INFO.Queries[6])
	}
	if report.INFO.Queries[7].Supported != report.INFO.Version.Parsed.SupportsInfoID(AdapterInfoWiFiRSSI) {
		t.Fatalf("report.INFO.Queries[7].Supported = %v; want %v", report.INFO.Queries[7].Supported, report.INFO.Version.Parsed.SupportsInfoID(AdapterInfoWiFiRSSI))
	}
	if report.Scan.Initial.State != 3 || report.Scan.Initial.StateName != "ready" {
		t.Fatalf("report.Scan.Initial = %+v; want ready snapshot", report.Scan.Initial)
	}
	if len(report.Scan.DeadlineExamples) < 3 {
		t.Fatalf("report.Scan.DeadlineExamples len = %d; want >= 3", len(report.Scan.DeadlineExamples))
	}
	if report.Scan.DeadlineExamples[0].EffectiveDelay != scanMinimumDelay {
		t.Fatalf("report.Scan.DeadlineExamples[0].EffectiveDelay = 0x%02x; want 0x%02x", report.Scan.DeadlineExamples[0].EffectiveDelay, scanMinimumDelay)
	}
	if len(report.Scan.DispatchSamples) != 7 {
		t.Fatalf("report.Scan.DispatchSamples len = %d; want 7", len(report.Scan.DispatchSamples))
	}
	if report.Scan.StatusSnapshot.Builder != "build_status_snapshot_frame" {
		t.Fatalf("report.Scan.StatusSnapshot.Builder = %q; want build_status_snapshot_frame", report.Scan.StatusSnapshot.Builder)
	}
	if report.Scan.StatusVariant.Builder != "build_status_variant_frame" {
		t.Fatalf("report.Scan.StatusVariant.Builder = %q; want build_status_variant_frame", report.Scan.StatusVariant.Builder)
	}
	if report.Scan.StatusSnapshot.PrefixHex[4] != "35" || report.Scan.StatusVariant.PrefixHex[4] != "37" {
		t.Fatalf("scan status prefixes = %#v / %#v; want 0x35 and 0x37", report.Scan.StatusSnapshot.PrefixHex, report.Scan.StatusVariant.PrefixHex)
	}
}
