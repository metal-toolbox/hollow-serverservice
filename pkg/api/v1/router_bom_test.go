package serverservice_test

import (
	"context"
	"reflect"
	"strings"
	"testing"

	serverservice "go.hollow.sh/serverservice/pkg/api/v1"
)

func TestIntegrationBomUpload(t *testing.T) {
	testCases := []struct {
		testName                   string
		uploadBoms                 []serverservice.Bom
		expectedUploadErrorMsg     string
		expectedUploadErr          bool
		aocMacAddress              string
		expectedAocMacAddressError bool
	}{
		{
			testName: "upload 1 bom and get by aoc mac address",
			uploadBoms: []serverservice.Bom{
				{
					SerialNum:     "fakeSerialNum1",
					AocMacAddress: "fakeAocMacAddress1,fakeAocMacAddress2",
					BmcMacAddress: "fakeBmcMacAddress1,fakeBmcMacAddress2",
					NumDefiPmi:    "fakeNumDefipmi1",
					NumDefPWD:     "fakeNumDefpwd1",
					Metro:         "fakeMetro1",
				},
			},
			expectedUploadErr:          false,
			expectedUploadErrorMsg:     "",
			aocMacAddress:              "fakeAocMacAddress1",
			expectedAocMacAddressError: false,
		},
		{
			testName: "upload 2 boms and get by aoc mac address",
			uploadBoms: []serverservice.Bom{
				{
					SerialNum:     "fakeSerialNum1",
					AocMacAddress: "fakeAocMacAddress1,fakeAocMacAddress2",
					BmcMacAddress: "fakeBmcMacAddress1,fakeBmcMacAddress2",
					NumDefiPmi:    "fakeNumDefipmi1",
					NumDefPWD:     "fakeNumDefpwd1",
					Metro:         "fakeMetro1",
				},
				{
					SerialNum:     "fakeSerialNum2",
					AocMacAddress: "fakeAocMacAddress3,fakeAocMacAddress4",
					BmcMacAddress: "fakeBmcMacAddress3,fakeBmcMacAddress4",
					NumDefiPmi:    "fakeNumDefipmi2",
					NumDefPWD:     "fakeNumDefpwd2",
					Metro:         "fakeMetro2",
				},
			},
			expectedUploadErr:          false,
			expectedUploadErrorMsg:     "",
			aocMacAddress:              "fakeAocMacAddress3",
			expectedAocMacAddressError: false,
		},
		{
			testName: "upload duplicate serial number",
			uploadBoms: []serverservice.Bom{
				{
					SerialNum:     "fakeSerialNum1",
					AocMacAddress: "fakeAocMacAddress1,fakeAocMacAddress2",
					BmcMacAddress: "fakeBmcMacAddress1,fakeBmcMacAddress2",
					NumDefiPmi:    "fakeNumDefipmi1",
					NumDefPWD:     "fakeNumDefpwd1",
					Metro:         "fakeMetro1",
				},
				{
					SerialNum:     "fakeSerialNum1",
					AocMacAddress: "fakeAocMacAddress3,fakeAocMacAddress4",
					BmcMacAddress: "fakeBmcMacAddress3,fakeBmcMacAddress4",
					NumDefiPmi:    "fakeNumDefipmi2",
					NumDefPWD:     "fakeNumDefpwd2",
					Metro:         "fakeMetro2",
				},
			},
			expectedUploadErr:          true,
			expectedUploadErrorMsg:     "unable to insert into bom_info: pq: duplicate key value violates unique constraint",
			aocMacAddress:              "fakeAocMacAddress3",
			expectedAocMacAddressError: false,
		},
		{
			testName: "upload duplicate AocMacAddress",
			uploadBoms: []serverservice.Bom{
				{
					SerialNum:     "fakeSerialNum1",
					AocMacAddress: "fakeAocMacAddress1,fakeAocMacAddress2",
					BmcMacAddress: "fakeBmcMacAddress1,fakeBmcMacAddress2",
					NumDefiPmi:    "fakeNumDefipmi1",
					NumDefPWD:     "fakeNumDefpwd1",
					Metro:         "fakeMetro1",
				},
				{
					SerialNum:     "fakeSerialNum2",
					AocMacAddress: "fakeAocMacAddress1,fakeAocMacAddress3",
					BmcMacAddress: "fakeBmcMacAddress3,fakeBmcMacAddress4",
					NumDefiPmi:    "fakeNumDefipmi2",
					NumDefPWD:     "fakeNumDefpwd2",
					Metro:         "fakeMetro2",
				},
			},
			expectedUploadErr:          true,
			expectedUploadErrorMsg:     "unable to insert into aoc_mac_address: pq: duplicate key value violates unique constraint",
			aocMacAddress:              "fakeAocMacAddress3",
			expectedAocMacAddressError: false,
		},
		{
			testName: "upload duplicate BmcMacAddress",
			uploadBoms: []serverservice.Bom{
				{
					SerialNum:     "fakeSerialNum1",
					AocMacAddress: "fakeAocMacAddress1,fakeAocMacAddress2",
					BmcMacAddress: "fakeBmcMacAddress1,fakeBmcMacAddress2",
					NumDefiPmi:    "fakeNumDefipmi1",
					NumDefPWD:     "fakeNumDefpwd1",
					Metro:         "fakeMetro1",
				},
				{
					SerialNum:     "fakeSerialNum2",
					AocMacAddress: "fakeAocMacAddress3,fakeAocMacAddress4",
					BmcMacAddress: "fakeBmcMacAddress1,fakeBmcMacAddress3",
					NumDefiPmi:    "fakeNumDefipmi2",
					NumDefPWD:     "fakeNumDefpwd2",
					Metro:         "fakeMetro2",
				},
			},
			expectedUploadErr:          true,
			expectedUploadErrorMsg:     "unable to insert into bmc_mac_address: pq: duplicate key value violates unique constraint",
			aocMacAddress:              "fakeBmcMacAddress3",
			expectedAocMacAddressError: false,
		},
		{
			testName: "upload empty serial number",
			uploadBoms: []serverservice.Bom{
				{
					SerialNum:     "",
					AocMacAddress: "fakeAocMacAddress1,fakeAocMacAddress2",
					BmcMacAddress: "fakeBmcMacAddress1,fakeBmcMacAddress2",
					NumDefiPmi:    "fakeNumDefipmi1",
					NumDefPWD:     "fakeNumDefpwd1",
					Metro:         "fakeMetro1",
				},
			},
			expectedUploadErr:          true,
			expectedUploadErrorMsg:     "the primary key serial-num can not be blank",
			aocMacAddress:              "fakeAocMacAddress3",
			expectedAocMacAddressError: false,
		},
		{
			testName: "upload empty AocMacAddress",
			uploadBoms: []serverservice.Bom{
				{
					SerialNum:     "fakeSerialNum1",
					AocMacAddress: "",
					BmcMacAddress: "fakeBmcMacAddress1,fakeBmcMacAddress2",
					NumDefiPmi:    "fakeNumDefipmi1",
					NumDefPWD:     "fakeNumDefpwd1",
					Metro:         "fakeMetro1",
				},
			},
			expectedUploadErr:          true,
			expectedUploadErrorMsg:     "the primary key aoc-mac-address can not be blank",
			aocMacAddress:              "fakeAocMacAddress3",
			expectedAocMacAddressError: false,
		},
		{
			testName: "upload empty BmcMacAddress",
			uploadBoms: []serverservice.Bom{
				{
					SerialNum:     "fakeSerialNum1",
					AocMacAddress: "fakeAocMacAddress1,fakeAocMacAddress2",
					BmcMacAddress: "",
					NumDefiPmi:    "fakeNumDefipmi1",
					NumDefPWD:     "fakeNumDefpwd1",
					Metro:         "fakeMetro1",
				},
			},
			expectedUploadErr:          true,
			expectedUploadErrorMsg:     "the primary key bmc-mac-address can not be blank",
			aocMacAddress:              "fakeAocMacAddress3",
			expectedAocMacAddressError: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.testName, func(t *testing.T) {
			s := serverTest(t)

			authToken := validToken(adminScopes)
			s.Client.SetToken(authToken)

			_, err := s.Client.BillOfMaterialsBatchUpload(context.TODO(), tc.uploadBoms)
			if tc.expectedUploadErr {
				if err == nil {
					t.Fatalf("BillOfMaterialsBatchUpload(%v) expect error, got nil", tc.uploadBoms)
				}
				if !strings.Contains(err.Error(), tc.expectedUploadErrorMsg) {
					t.Fatalf("BillOfMaterialsBatchUpload(%v) expect error %v, got %v", tc.uploadBoms, tc.expectedUploadErrorMsg, err)
				}
				return
			}
			if err != nil {
				t.Fatalf("BillOfMaterialsBatchUpload(%v) failed to upload, err %v", tc.uploadBoms, err)
				return
			}

			_, _, err = s.Client.GetBomInfoByAOCMacAddr(context.TODO(), tc.aocMacAddress)
			if tc.expectedAocMacAddressError {
				if err == nil {
					t.Fatalf("GetBomInfoByAOCMacAddr(%v) expect error, got nil", tc.aocMacAddress)
				}
				return
			}

			if err != nil {
				t.Fatalf("GetBomInfoByAOCMacAddr(%v) failed to get bom, err %v", tc.aocMacAddress, err)
			}
		})
	}
}

func TestUploadInOneTransaction(t *testing.T) {
	uploadBoms :=
		[]serverservice.Bom{
			{
				SerialNum:     "fakeSerialNum1",
				AocMacAddress: "fakeAocMacAddress1,fakeAocMacAddress2",
				BmcMacAddress: "fakeBmcMacAddress1,fakeBmcMacAddress2",
				NumDefiPmi:    "fakeNumDefipmi1",
				NumDefPWD:     "fakeNumDefpwd1",
				Metro:         "fakeMetro1",
			},
			{
				SerialNum:     "fakeSerialNum2",
				AocMacAddress: "fakeAocMacAddress1,fakeAocMacAddress3",
				BmcMacAddress: "fakeBmcMacAddress3,fakeBmcMacAddress4",
				NumDefiPmi:    "fakeNumDefipmi2",
				NumDefPWD:     "fakeNumDefpwd2",
				Metro:         "fakeMetro2",
			},
		}
	expectedUploadErrorMsg := "unable to insert into aoc_mac_address: pq: duplicate key value violates unique constraint"
	expectedGetMsg := "no rows in result set"

	s := serverTest(t)
	authToken := validToken(adminScopes)
	s.Client.SetToken(authToken)

	_, err := s.Client.BillOfMaterialsBatchUpload(context.TODO(), uploadBoms)
	if !strings.Contains(err.Error(), expectedUploadErrorMsg) {
		t.Fatalf("BillOfMaterialsBatchUpload(%v) expect error %v, got %v", uploadBoms, expectedUploadErrorMsg, err)
	}

	bom, _, err := s.Client.GetBomInfoByAOCMacAddr(context.TODO(), uploadBoms[0].AocMacAddress)
	if err == nil || !strings.Contains(err.Error(), expectedGetMsg) {
		t.Fatalf("GetBomInfoByAOCMacAddr(%v) got bom %v err %v, expect nil, %v", uploadBoms[0].AocMacAddress, bom, err, expectedGetMsg)
	}
}

func TestIntegrationGetBomByAocMacAddr(t *testing.T) {
	s := serverTest(t)

	uploadBoms := []serverservice.Bom{
		{
			SerialNum:     "fakeSerialNum1",
			AocMacAddress: "fakeAocMacAddress1,fakeAocMacAddress2",
			BmcMacAddress: "fakeBmcMacAddress1,fakeBmcMacAddress2",
			NumDefiPmi:    "fakeNumDefipmi1",
			NumDefPWD:     "fakeNumDefpwd1",
			Metro:         "fakeMetro1",
		},
		{
			SerialNum:     "fakeSerialNum2",
			AocMacAddress: "fakeAocMacAddress3,fakeAocMacAddress4",
			BmcMacAddress: "fakeBmcMacAddress3,fakeBmcMacAddress4",
			NumDefiPmi:    "fakeNumDefipmi2",
			NumDefPWD:     "fakeNumDefpwd2",
			Metro:         "fakeMetro2",
		},
	}
	authToken := validToken(adminScopes)
	s.Client.SetToken(authToken)

	_, err := s.Client.BillOfMaterialsBatchUpload(context.TODO(), uploadBoms)
	if err != nil {
		t.Fatalf("s.Client.BillOfMaterialsBatchUpload(%v) failed to upload, err %v", uploadBoms, err)
		return
	}

	var testCases = []struct {
		testName                      string
		aocMacAddress                 string
		expectedBom                   serverservice.Bom
		expectedAocMacAddressError    bool
		expectedAocMacAddressErrorMsg string
	}{
		{
			testName:                      "get first bom by first aoc mac address",
			aocMacAddress:                 "fakeAocMacAddress1",
			expectedBom:                   uploadBoms[0],
			expectedAocMacAddressError:    false,
			expectedAocMacAddressErrorMsg: "",
		},
		{
			testName:                      "get first bom by second aoc mac address",
			aocMacAddress:                 "fakeAocMacAddress2",
			expectedBom:                   uploadBoms[0],
			expectedAocMacAddressError:    false,
			expectedAocMacAddressErrorMsg: "",
		},
		{
			testName:                      "get second bom by first aoc mac address",
			aocMacAddress:                 "fakeAocMacAddress3",
			expectedBom:                   uploadBoms[1],
			expectedAocMacAddressError:    false,
			expectedAocMacAddressErrorMsg: "",
		},
		{
			testName:                      "get second bom by second aoc mac address",
			aocMacAddress:                 "fakeAocMacAddress3",
			expectedBom:                   uploadBoms[1],
			expectedAocMacAddressError:    false,
			expectedAocMacAddressErrorMsg: "",
		},
		{
			testName:                      "non-exist aoc mac address",
			aocMacAddress:                 "random",
			expectedBom:                   uploadBoms[1],
			expectedAocMacAddressError:    true,
			expectedAocMacAddressErrorMsg: "sql: no rows in result set",
		},
		{
			testName:                      "empty aoc mac address",
			aocMacAddress:                 "",
			expectedBom:                   uploadBoms[1],
			expectedAocMacAddressError:    true,
			expectedAocMacAddressErrorMsg: "route not found",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.testName, func(t *testing.T) {
			bom, _, err := s.Client.GetBomInfoByAOCMacAddr(context.TODO(), tc.aocMacAddress)
			if tc.expectedAocMacAddressError {
				if err == nil {
					t.Fatalf("GetBomInfoByAOCMacAddr(%v) expect error, got nil", tc.aocMacAddress)
				}
				if !strings.Contains(err.Error(), tc.expectedAocMacAddressErrorMsg) {
					t.Fatalf("GetBomInfoByAOCMacAddr(%v) expect error %v, got %v", tc.aocMacAddress, tc.expectedAocMacAddressErrorMsg, err)
				}
				return
			}
			if err != nil {
				t.Fatalf("GetBomInfoByAOCMacAddr(%v) failed to upload, err %v", tc.aocMacAddress, err)
				return
			}

			if !reflect.DeepEqual(bom, &tc.expectedBom) {
				t.Fatalf("got incorrect bom %v, expect %v", bom, tc.expectedBom)
			}
		})
	}
}

func TestIntegrationGetBomByBmcMacAddr(t *testing.T) {
	s := serverTest(t)

	uploadBoms := []serverservice.Bom{
		{
			SerialNum:     "fakeSerialNum1",
			AocMacAddress: "fakeAocMacAddress1,fakeAocMacAddress2",
			BmcMacAddress: "fakeBmcMacAddress1,fakeBmcMacAddress2",
			NumDefiPmi:    "fakeNumDefipmi1",
			NumDefPWD:     "fakeNumDefpwd1",
			Metro:         "fakeMetro1",
		},
		{
			SerialNum:     "fakeSerialNum2",
			AocMacAddress: "fakeAocMacAddress3,fakeAocMacAddress4",
			BmcMacAddress: "fakeBmcMacAddress3,fakeBmcMacAddress4",
			NumDefiPmi:    "fakeNumDefipmi2",
			NumDefPWD:     "fakeNumDefpwd2",
			Metro:         "fakeMetro2",
		},
	}
	authToken := validToken(adminScopes)
	s.Client.SetToken(authToken)

	_, err := s.Client.BillOfMaterialsBatchUpload(context.TODO(), uploadBoms)
	if err != nil {
		t.Fatalf("s.Client.BillOfMaterialsBatchUpload(%v) failed to upload, err %v", uploadBoms, err)
		return
	}

	var testCases = []struct {
		testName                      string
		bmcMacAddress                 string
		expectedBom                   serverservice.Bom
		expectedBmcMacAddressError    bool
		expectedBmcMacAddressErrorMsg string
	}{
		{
			testName:                      "get first bom by first bmc mac address",
			bmcMacAddress:                 "fakeBmcMacAddress1",
			expectedBom:                   uploadBoms[0],
			expectedBmcMacAddressError:    false,
			expectedBmcMacAddressErrorMsg: "",
		},
		{
			testName:                      "get first bom by second bmc mac address",
			bmcMacAddress:                 "fakeBmcMacAddress2",
			expectedBom:                   uploadBoms[0],
			expectedBmcMacAddressError:    false,
			expectedBmcMacAddressErrorMsg: "",
		},
		{
			testName:                      "get second bom by first bmc mac address",
			bmcMacAddress:                 "fakeBmcMacAddress3",
			expectedBom:                   uploadBoms[1],
			expectedBmcMacAddressError:    false,
			expectedBmcMacAddressErrorMsg: "",
		},
		{
			testName:                      "get second bom by second bmc mac address",
			bmcMacAddress:                 "fakeBmcMacAddress3",
			expectedBom:                   uploadBoms[1],
			expectedBmcMacAddressError:    false,
			expectedBmcMacAddressErrorMsg: "",
		},
		{
			testName:                      "non-exist bmc mac address",
			bmcMacAddress:                 "random",
			expectedBom:                   uploadBoms[1],
			expectedBmcMacAddressError:    true,
			expectedBmcMacAddressErrorMsg: "sql: no rows in result set",
		},
		{
			testName:                      "empty bmc mac address",
			bmcMacAddress:                 "",
			expectedBom:                   uploadBoms[1],
			expectedBmcMacAddressError:    true,
			expectedBmcMacAddressErrorMsg: "route not found",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.testName, func(t *testing.T) {
			bom, _, err := s.Client.GetBomInfoByBMCMacAddr(context.TODO(), tc.bmcMacAddress)
			if tc.expectedBmcMacAddressError {
				if err == nil {
					t.Fatalf("GetBomInfoByBMCMacAddr(%v) expect error, got nil", tc.bmcMacAddress)
				}
				if !strings.Contains(err.Error(), tc.expectedBmcMacAddressErrorMsg) {
					t.Fatalf("GetBomInfoByBMCMacAddr(%v) expect error %v, got %v", tc.bmcMacAddress, tc.expectedBmcMacAddressErrorMsg, err)
				}
				return
			}
			if err != nil {
				t.Fatalf("GetBomInfoByBMCMacAddr(%v) failed to upload, err %v", tc.bmcMacAddress, err)
				return
			}

			if !reflect.DeepEqual(bom, &tc.expectedBom) {
				t.Fatalf("got incorrect bom %v, expect %v", bom, tc.expectedBom)
			}
		})
	}
}
