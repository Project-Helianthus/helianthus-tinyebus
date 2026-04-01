package adapterproto

import (
	"encoding/binary"
	"fmt"
)

type AdapterInfoID byte

const (
	AdapterInfoVersion      AdapterInfoID = 0x00
	AdapterInfoHardwareID   AdapterInfoID = 0x01
	AdapterInfoHardwareConf AdapterInfoID = 0x02
	AdapterInfoTemperature  AdapterInfoID = 0x03
	AdapterInfoSupplyVolt   AdapterInfoID = 0x04
	AdapterInfoBusVoltage   AdapterInfoID = 0x05
	AdapterInfoResetInfo    AdapterInfoID = 0x06
	AdapterInfoWiFiRSSI     AdapterInfoID = 0x07
)

var adapterInfoIDNames = [8]string{
	"version", "hw_id", "hw_config", "temperature",
	"supply_voltage", "bus_voltage", "reset_info", "wifi_rssi",
}

func (id AdapterInfoID) String() string {
	if int(id) < len(adapterInfoIDNames) {
		return adapterInfoIDNames[id]
	}
	return fmt.Sprintf("unknown(0x%02x)", byte(id))
}

type AdapterVersion struct {
	Version  byte `json:"version"`
	Features byte `json:"features"`

	Checksum uint16 `json:"checksum"`
	Jumpers  byte   `json:"jumpers"`

	BootloaderVersion  byte   `json:"bootloader_version"`
	BootloaderChecksum uint16 `json:"bootloader_checksum"`

	HasChecksum   bool `json:"has_checksum"`
	HasBootloader bool `json:"has_bootloader"`
	SupportsInfo  bool `json:"supports_info"`
	IsWiFi        bool `json:"is_wifi"`
	IsEthernet    bool `json:"is_ethernet"`
	IsHighSpeed   bool `json:"is_high_speed"`
	IsV31         bool `json:"is_v31"`
}

func (v AdapterVersion) VersionResponseLen() int {
	if v.HasBootloader {
		return 8
	}
	if v.HasChecksum {
		return 5
	}
	return 2
}

func (v AdapterVersion) SupportsInfoID(id AdapterInfoID) bool {
	if !v.SupportsInfo {
		return false
	}
	switch id {
	case AdapterInfoVersion, AdapterInfoHardwareID, AdapterInfoHardwareConf,
		AdapterInfoTemperature, AdapterInfoSupplyVolt, AdapterInfoBusVoltage:
		return true
	case AdapterInfoResetInfo:
		return v.HasBootloader
	case AdapterInfoWiFiRSSI:
		return v.HasChecksum && v.IsWiFi
	default:
		return false
	}
}

func ParseAdapterVersion(data []byte) (AdapterVersion, error) {
	switch len(data) {
	case 2, 5, 8:
	default:
		return AdapterVersion{}, fmt.Errorf("adapter version response has invalid length (%d bytes)", len(data))
	}

	version := AdapterVersion{
		Version:  data[0],
		Features: data[1],
	}
	version.SupportsInfo = version.Features&0x01 != 0

	switch len(data) {
	case 8:
		version.HasBootloader = true
		version.BootloaderVersion = data[5]
		version.BootloaderChecksum = binary.BigEndian.Uint16(data[6:8])
		fallthrough
	case 5:
		version.HasChecksum = true
		version.Checksum = binary.BigEndian.Uint16(data[2:4])
		version.Jumpers = data[4]
		version.IsWiFi = version.Jumpers&0x08 != 0
		version.IsEthernet = version.Jumpers&0x04 != 0
		version.IsHighSpeed = version.Jumpers&0x02 != 0
		version.IsV31 = version.Jumpers&0x10 != 0
	}

	return version, nil
}

type AdapterResetInfo struct {
	Cause        string `json:"cause"`
	CauseCode    byte   `json:"cause_code"`
	RestartCount byte   `json:"restart_count"`
}

var resetCauseNames = map[byte]string{
	1: "power_on",
	2: "brown_out",
	3: "watchdog",
	4: "clear",
	5: "external_reset",
	6: "stack_overflow",
	7: "memory_failure",
}

func ParseAdapterResetInfo(data []byte) (AdapterResetInfo, error) {
	if len(data) < 2 {
		return AdapterResetInfo{}, fmt.Errorf("adapter reset info too short (%d bytes)", len(data))
	}

	cause := "unknown"
	if name, ok := resetCauseNames[data[0]]; ok {
		cause = name
	}

	return AdapterResetInfo{
		Cause:        cause,
		CauseCode:    data[0],
		RestartCount: data[1],
	}, nil
}
