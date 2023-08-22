package serverservice

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"reflect"
	"sort"
	"strings"
	"testing"
)

var testSerialNumBomInfo1 = &SerialNumBomInfo{
	bom: Bom{
		SerialNum:     "test-serial-1",
		AocMacAddress: "FakeAOC1,FakeAOC2",
		BmcMacAddress: "FakeMac1,FakeMac2",
		NumDefipmi:    "FakeDEFI1",
		NumDefpwd:     "FakeDEFPWD1",
	},
	aocMacAddrsBom: []AocMacAddressBom{
		{
			AocMacAddress: "FakeAOC1",
			SerialNum:     "test-serial-1",
		},
		{
			AocMacAddress: "FakeAOC2",
			SerialNum:     "test-serial-1",
		},
	},
	bmcMacAddrsBom: []BmcMacAddressBom{
		{
			BmcMacAddress: "FakeMac1",
			SerialNum:     "test-serial-1",
		},
		{
			BmcMacAddress: "FakeMac2",
			SerialNum:     "test-serial-1",
		},
	},
}

var testSerialNumBomInfo2 = &SerialNumBomInfo{
	bom: Bom{
		SerialNum:     "test-serial-2",
		AocMacAddress: "FakeAOC3,FakeAOC4",
		BmcMacAddress: "FakeMac3,FakeMac4",
		NumDefipmi:    "FakeDEFI2",
		NumDefpwd:     "FakeDEFPWD2",
	},
	aocMacAddrsBom: []AocMacAddressBom{
		{
			AocMacAddress: "FakeAOC3",
			SerialNum:     "test-serial-2",
		},
		{
			AocMacAddress: "FakeAOC4",
			SerialNum:     "test-serial-2",
		},
	},
	bmcMacAddrsBom: []BmcMacAddressBom{
		{
			BmcMacAddress: "FakeMac3",
			SerialNum:     "test-serial-2",
		},
		{
			BmcMacAddress: "FakeMac4",
			SerialNum:     "test-serial-2",
		},
	},
}

func sortSerialNumBomInfos(serialNumInfos []*SerialNumBomInfo) {
	sort.Slice(serialNumInfos, func(i, j int) bool {
		return serialNumInfos[i].bom.SerialNum < serialNumInfos[j].bom.SerialNum
	})

	for _, s := range serialNumInfos {
		sortAocMacAddrBoms(s.aocMacAddrsBom)
		sortBmcMacAddrBoms(s.bmcMacAddrsBom)
	}
}

func sortAocMacAddrBoms(aocMacAddrBoms []AocMacAddressBom) {
	sort.Slice(aocMacAddrBoms, func(i, j int) bool {
		return aocMacAddrBoms[i].AocMacAddress < aocMacAddrBoms[j].AocMacAddress
	})
}

func sortBmcMacAddrBoms(bmcMacAddrBoms []BmcMacAddressBom) {
	sort.Slice(bmcMacAddrBoms, func(i, j int) bool {
		return bmcMacAddrBoms[i].BmcMacAddress < bmcMacAddrBoms[j].BmcMacAddress
	})
}

func TestParseXlsxFile(t *testing.T) {
	var testCases = []struct {
		testName                 string
		filePath                 string
		expectedErr              bool
		expectedErrMsg           string
		expectedSerialNumBomInfo []*SerialNumBomInfo
	}{
		{
			testName:                 "file missing serial number",
			filePath:                 "./testdata/test_empty_serial.xlsx",
			expectedErr:              true,
			expectedErrMsg:           "empty serial number",
			expectedSerialNumBomInfo: nil,
		},
		{
			testName:                 "file missing aocMacAddress ",
			filePath:                 "./testdata/test_empty_aocMacAddress.xlsx",
			expectedErr:              true,
			expectedErrMsg:           "empty aoc mac address",
			expectedSerialNumBomInfo: nil,
		},
		{
			testName:                 "file missing bmcMacAddress",
			filePath:                 "./testdata/test_empty_bmcMacAddress.xlsx",
			expectedErr:              true,
			expectedErrMsg:           "empty bmc mac address",
			expectedSerialNumBomInfo: nil,
		},
		{
			testName:                 "valid file for single bom",
			filePath:                 "./testdata/test_valid_one_bom.xlsx",
			expectedErr:              false,
			expectedErrMsg:           "",
			expectedSerialNumBomInfo: []*SerialNumBomInfo{testSerialNumBomInfo1},
		},
		{
			testName:                 "valid file for multiple bom",
			filePath:                 "./testdata/test_valid_multiple_boms.xlsx",
			expectedErr:              false,
			expectedErrMsg:           "",
			expectedSerialNumBomInfo: []*SerialNumBomInfo{testSerialNumBomInfo1, testSerialNumBomInfo2},
		},
		{
			testName:                 "file missing SERIALNUM col",
			filePath:                 "./testdata/test_empty_serial_col.xlsx",
			expectedErr:              true,
			expectedErrMsg:           "missing colomn, serial num -1, sub-item 4, sub-serial 5",
			expectedSerialNumBomInfo: nil,
		},
	}

	for _, tt := range testCases {
		t.Run(tt.testName, func(t *testing.T) {
			file, err := os.Open(tt.filePath)
			if err != nil {
				t.Fatalf("os.Open(%v) failed to open file %v\n", tt.filePath, err)
			}

			// Translate file to bytes since ParseXlsxFile accepts bytes of file as argument,
			// which is the format of a file reading from the HTTP request.
			stat, err := file.Stat()
			if err != nil {
				fmt.Println(err)
				return
			}
			bs := make([]byte, stat.Size())
			_, err = bufio.NewReader(file).Read(bs)
			if err != nil && err != io.EOF {
				fmt.Println(err)
				return
			}

			serialNumInfos, err := ParseXlsxFile(bs)
			if tt.expectedErr {
				if !strings.Contains(err.Error(), tt.expectedErrMsg) {
					t.Fatalf("test %v failed, got %v, expect %v", tt.testName, err, tt.expectedErrMsg)
				}
				if serialNumInfos != nil {
					t.Fatalf("test %v expect nil serialNumInfos, got %v", tt.testName, len(serialNumInfos))
				}
				return
			}
			if err != nil {
				t.Fatalf("test %v failed to parse Xlsx file: %v", tt.testName, err)
			}

			if len(serialNumInfos) != len(tt.expectedSerialNumBomInfo) {
				t.Fatalf("test %v parsed incorrect numbers of serialNumInfos, got %v, expect %v", tt.testName, len(serialNumInfos), tt.expectedSerialNumBomInfo)
			}

			// Sort the serialNumInfos to avoid unexpected orders of the maps.Values
			sortSerialNumBomInfos(serialNumInfos)
			sortSerialNumBomInfos(tt.expectedSerialNumBomInfo)
			for i := range serialNumInfos {
				info := serialNumInfos[i]
				expectedInfo := tt.expectedSerialNumBomInfo[i]
				if info.bom != expectedInfo.bom {
					t.Fatalf("test %v parsed incorrect bom info, got %v, expect %v", tt.testName, info.bom, expectedInfo.bom)
				}
				if !reflect.DeepEqual(info.aocMacAddrsBom, expectedInfo.aocMacAddrsBom) {
					t.Fatalf("test %v parsed incorrect aoc mac addr, got %v, expect %v", tt.testName, info.aocMacAddrsBom, expectedInfo.aocMacAddrsBom)
				}
				if !reflect.DeepEqual(info.bmcMacAddrsBom, expectedInfo.bmcMacAddrsBom) {
					t.Fatalf("test %v parsed incorrect bmc mac addr, got %v, expect %v", tt.testName, info.bmcMacAddrsBom, expectedInfo.bmcMacAddrsBom)
				}
			}
		})
	}
}
