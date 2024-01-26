package serverservice

import (
	"strings"

	"github.com/pkg/errors"
	"github.com/volatiletech/null/v8"

	"github.com/metal-toolbox/fleetdb/internal/models"
)

// Bom provides a struct to map the bom_info table.
// Naming conversion is strange here just in order to make it consistent
// with generated BomInfo.
type Bom struct {
	SerialNum     string `json:"serial_num"`      // physical serial number listed outside of a server
	AocMacAddress string `json:"aoc_mac_address"` // Aoc is alternative name of the fiber channel card MAC address
	BmcMacAddress string `json:"bmc_mac_address"`
	NumDefiPmi    string `json:"num_defi_pmi"`
	NumDefPWD     string `json:"num_def_pwd"` // DefPWD is the IPMI Password in the portal
	Metro         string `json:"metro"`
}

// AocMacAddressBom provides a struct to map the aoc_mac_address table.
type AocMacAddressBom struct {
	AocMacAddress string `json:"aoc_mac_address"`
	SerialNum     string `json:"serial_num"`
}

// BmcMacAddressBom provides a struct to map the bmc_mac_address table.
type BmcMacAddressBom struct {
	BmcMacAddress string `json:"bmc_mac_address"`
	SerialNum     string `json:"serial_num"`
}

// toDBModel converts Bom to BomInfo.
func (b *Bom) toDBModel() (*models.BomInfo, error) {
	if b.SerialNum == "" {
		return nil, errors.Errorf("the primary key serial-num can not be blank")
	}

	dbB := &models.BomInfo{
		SerialNum:     b.SerialNum,
		AocMacAddress: null.StringFrom(b.AocMacAddress),
		BMCMacAddress: null.StringFrom(b.BmcMacAddress),
		NumDefiPmi:    null.StringFrom(b.NumDefiPmi),
		NumDefPWD:     null.StringFrom(b.NumDefPWD),
		Metro:         null.StringFrom(b.Metro),
	}

	return dbB, nil
}

// toDBModel converts BomInfo to Bom.
func (b *Bom) fromDBModel(bomInfo *models.BomInfo) error {
	b.SerialNum = bomInfo.SerialNum
	b.AocMacAddress = bomInfo.AocMacAddress.String
	b.BmcMacAddress = bomInfo.BMCMacAddress.String
	b.NumDefiPmi = bomInfo.NumDefiPmi.String
	b.NumDefPWD = bomInfo.NumDefPWD.String
	b.Metro = bomInfo.Metro.String

	return nil
}

// toAocMacAddressDBModels converts Bom to one or multiple AocMacAddress.
func (b *Bom) toAocMacAddressDBModels() ([]*models.AocMacAddress, error) {
	if b.AocMacAddress == "" {
		return nil, errors.Errorf("the primary key aoc-mac-address can not be blank")
	}

	dbAs := []*models.AocMacAddress{}

	AocMacAddrs := strings.Split(b.AocMacAddress, ",")
	for _, aocMacAddr := range AocMacAddrs {
		dbA := &models.AocMacAddress{
			SerialNum:     b.SerialNum,
			AocMacAddress: aocMacAddr,
		}
		dbAs = append(dbAs, dbA)
	}

	return dbAs, nil
}

// toBmcMacAddressDBModels converts Bom to one or multiple BmcMacAddress.
func (b *Bom) toBmcMacAddressDBModels() ([]*models.BMCMacAddress, error) {
	if b.BmcMacAddress == "" {
		return nil, errors.Errorf("the primary key bmc-mac-address can not be blank")
	}

	dbBs := []*models.BMCMacAddress{}

	BmcMacAddrs := strings.Split(b.BmcMacAddress, ",")
	for _, bmcMacAddr := range BmcMacAddrs {
		dbB := &models.BMCMacAddress{
			SerialNum:     b.SerialNum,
			BMCMacAddress: bmcMacAddr,
		}
		dbBs = append(dbBs, dbB)
	}

	return dbBs, nil
}
