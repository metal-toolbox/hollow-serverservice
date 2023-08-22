package serverservice

// AocMacAddressBom provide a struct to map the aoc_mac_address table.
type AocMacAddressBom struct {
	AocMacAddress string `json:"aoc_mac_address"`
	SerialNum     string `json:"serial_num"`
}
