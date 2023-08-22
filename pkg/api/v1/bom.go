package serverservice

// Bom provide a struct to map the bom_info table.
type Bom struct {
	SerialNum     string `json:"serial_num"`
	AocMacAddress string `json:"aoc_mac_address"`
	BmcMacAddress string `json:"bmc_mac_address"`
	NumDefipmi    string `json:"num_defi_pmi"`
	NumDefpwd     string `json:"num_def_pwd"`
	Metro         string `json:"metro"`
}
