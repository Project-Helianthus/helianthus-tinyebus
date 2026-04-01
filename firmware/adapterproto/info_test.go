package adapterproto

import "testing"

func TestParseAdapterVersion(t *testing.T) {
	t.Parallel()

	version, err := ParseAdapterVersion([]byte{0x12, 0x01, 0xAB, 0xCD, 0x1E, 0x34, 0x56, 0x78})
	if err != nil {
		t.Fatalf("ParseAdapterVersion() error = %v", err)
	}
	if version.Version != 0x12 || version.Features != 0x01 {
		t.Fatalf("Version/Features = 0x%02x/0x%02x; want 0x12/0x01", version.Version, version.Features)
	}
	if version.Checksum != 0xABCD || version.Jumpers != 0x1E {
		t.Fatalf("Checksum/Jumpers = 0x%04x/0x%02x; want 0xabcd/0x1e", version.Checksum, version.Jumpers)
	}
	if version.BootloaderVersion != 0x34 || version.BootloaderChecksum != 0x5678 {
		t.Fatalf("Bootloader version/checksum = 0x%02x/0x%04x; want 0x34/0x5678", version.BootloaderVersion, version.BootloaderChecksum)
	}
	if !version.SupportsInfo || !version.HasChecksum || !version.HasBootloader {
		t.Fatalf("Version flags = %+v; want supports/info/checksum/bootloader enabled", version)
	}
	if !version.SupportsInfoID(AdapterInfoVersion) {
		t.Fatalf("SupportsInfoID(version) = false; want true")
	}
	if !version.SupportsInfoID(AdapterInfoResetInfo) {
		t.Fatalf("SupportsInfoID(reset_info) = false; want true")
	}
}

func TestParseAdapterVersionRejectsInvalidLength(t *testing.T) {
	t.Parallel()

	if _, err := ParseAdapterVersion([]byte{0x01}); err == nil {
		t.Fatalf("ParseAdapterVersion() error = nil; want error")
	}
}

func TestParseAdapterResetInfo(t *testing.T) {
	t.Parallel()

	reset, err := ParseAdapterResetInfo([]byte{0x03, 0x07})
	if err != nil {
		t.Fatalf("ParseAdapterResetInfo() error = %v", err)
	}
	if reset.Cause != "watchdog" || reset.CauseCode != 0x03 || reset.RestartCount != 0x07 {
		t.Fatalf("Reset = %+v; want watchdog/0x03/0x07", reset)
	}
}
