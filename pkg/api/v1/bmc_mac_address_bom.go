package serverservice

// BmcMacAddressBom provide a struct to map the bmc_mac_address table.
type BmcMacAddressBom struct {
	BmcMacAddress string `json:"bmc_mac_address"`
	SerialNum     string `json:"serial_num"`
}
