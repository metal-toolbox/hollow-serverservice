package serverservice

import (
	"github.com/pkg/errors"
	"github.com/tealeg/xlsx"
	"golang.org/x/exp/maps"
)

const (
	// column name
	serialNumName string = "SERIALNUM"
	subItemName   string = "SUB-ITEM"
	subSerialName string = "SUB-SERIAL" // value of the sub-item
	// interested catogories in sub-item
	aocMacAddrName string = "MAC-AOC-ADDRESS"
	bmcMacAddrName string = "MAC-ADDRESS"
	numDefipmiName string = "NUM-DEFIPMI"
	numDefipwdName string = "NUM-DEFPWD"
)

type categoryColNum struct {
	serialNumCol int // column number of the serial number, -1 means no such column
	subItemCol   int // column number of the sub-item, -1 means no such column
	subSerialCol int // column number of the sub-serial, -1 means no such column
}

func newCategoryColNum() *categoryColNum {
	return &categoryColNum{
		serialNumCol: -1,
		subItemCol:   -1,
		subSerialCol: -1,
	}
}

// SerialNumBomInfo includes informations of a serial number writing to database
type SerialNumBomInfo struct {
	// Write to the bom_info table.
	bom Bom
	// Write to the aoc_mac_address_bom table. One serial number can have multiple aoc mac addresses.
	aocMacAddrsBom []AocMacAddressBom
	// Write to the bmc_mac_address_bom table. One serial number can have multiple bmc mac addresses.
	bmcMacAddrsBom []BmcMacAddressBom
}

func newSerialNumBomInfo(serialNum string) *SerialNumBomInfo {
	return &SerialNumBomInfo{
		bom: Bom{SerialNum: serialNum},
	}
}

// ParseXlsxFile is the helper function to parse xlsx to bomInfos.
func ParseXlsxFile(fileBytes []byte) ([]*SerialNumBomInfo, error) {
	file, err := xlsx.OpenBinary(fileBytes)
	if err != nil {
		return nil, errors.New("failed to open the file")
	}

	infoMap := make(map[string]*SerialNumBomInfo)

	for _, sheet := range file.Sheets {
		var categoryCol *categoryColNum
		for _, row := range sheet.Rows {
			if categoryCol == nil {
				categoryCol = newCategoryColNum()

				for i, cell := range row.Cells {
					switch cell.Value {
					case serialNumName:
						categoryCol.serialNumCol = i
					case subItemName:
						categoryCol.subItemCol = i
					case subSerialName:
						categoryCol.subSerialCol = i
					}
				}

				if categoryCol.serialNumCol == -1 || categoryCol.subItemCol == -1 || categoryCol.subSerialCol == -1 {
					return nil, errors.Errorf("missing colomn, serial num %v, sub-item %v, sub-serial %v", categoryCol.serialNumCol, categoryCol.subItemCol, categoryCol.subSerialCol)
				}

				continue
			}

			// There won't be any out of idex issue since any non-existing value will default to empty string.
			cells := row.Cells
			serialNum := cells[categoryCol.serialNumCol].Value

			if len(serialNum) == 0 {
				return nil, errors.New("empty serial number")
			}

			bomInfo, ok := infoMap[serialNum]
			if !ok {
				bomInfo = newSerialNumBomInfo(serialNum)
				infoMap[serialNum] = bomInfo
			}

			switch cells[categoryCol.subItemCol].Value {
			case aocMacAddrName:
				aocMacAddress := cells[categoryCol.subSerialCol].Value
				if len(aocMacAddress) == 0 {
					return nil, errors.New("empty aoc mac address")
				}

				if len(bomInfo.bom.AocMacAddress) > 0 {
					bomInfo.bom.AocMacAddress += ","
				}

				bomInfo.bom.AocMacAddress += aocMacAddress
				bomInfo.aocMacAddrsBom = append(bomInfo.aocMacAddrsBom, AocMacAddressBom{
					SerialNum:     serialNum,
					AocMacAddress: aocMacAddress,
				})
			case bmcMacAddrName:
				bmcMacAddress := cells[categoryCol.subSerialCol].Value
				if len(bmcMacAddress) == 0 {
					return nil, errors.New("empty bmc mac address")
				}

				if len(bomInfo.bom.BmcMacAddress) > 0 {
					bomInfo.bom.BmcMacAddress += ","
				}

				bomInfo.bom.BmcMacAddress += bmcMacAddress
				bomInfo.bmcMacAddrsBom = append(bomInfo.bmcMacAddrsBom, BmcMacAddressBom{
					SerialNum:     serialNum,
					BmcMacAddress: bmcMacAddress,
				})
			case numDefipmiName:
				bomInfo.bom.NumDefipmi = cells[categoryCol.subSerialCol].Value
			case numDefipwdName:
				bomInfo.bom.NumDefpwd = cells[categoryCol.subSerialCol].Value
			}
		}
	}

	return maps.Values(infoMap), nil
}
