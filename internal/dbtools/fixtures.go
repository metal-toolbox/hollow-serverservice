//go:build testtools
// +build testtools

package dbtools

import (
	"context"
	"testing"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/volatiletech/null/v8"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/types"

	"go.hollow.sh/serverservice/internal/models"
)

//nolint:revive
var (
	// Namespaces used in Attributes and VersionedAttributes
	FixtureNamespaceMetadata    = "hollow.metadata"
	FixtureNamespaceOtherdata   = "hollow.other_data"
	FixtureNamespaceVersioned   = "hollow.versioned"
	FixtureNamespaceVersionedV2 = "hollow.versioned.v2"

	// Server Component Types
	FixtureFinType *models.ServerComponentType

	FixtureNemo                 *models.Server
	FixtureNemoMetadata         *models.Attribute
	FixtureNemoOtherdata        *models.Attribute
	FixtureNemoLeftFin          *models.ServerComponent
	FixtureNemoRightFin         *models.ServerComponent
	FixtureNemoLeftFinVersioned *models.VersionedAttribute
	FixtureNemoVersionedNew     *models.VersionedAttribute
	FixtureNemoVersionedOld     *models.VersionedAttribute
	FixtureNemoVersionedV2      *models.VersionedAttribute
	FixtureNemoBMCSecret        *models.ServerCredential

	FixtureDory          *models.Server
	FixtureDoryMetadata  *models.Attribute
	FixtureDoryOtherdata *models.Attribute
	FixtureDoryLeftFin   *models.ServerComponent
	FixtureDoryRightFin  *models.ServerComponent

	FixtureMarlin          *models.Server
	FixtureMarlinMetadata  *models.Attribute
	FixtureMarlinOtherdata *models.Attribute
	FixtureMarlinLeftFin   *models.ServerComponent
	FixtureMarlinRightFin  *models.ServerComponent

	// FixtureChuckles represents the fish that was deleted
	// https://pixar.fandom.com/wiki/Chuckles_(Finding_Nemo)
	FixtureChuckles          *models.Server
	FixtureChucklesMetadata  *models.Attribute
	FixtureChucklesOtherdata *models.Attribute
	FixtureChucklesLeftFin   *models.ServerComponent

	FixtureServers        models.ServerSlice
	FixtureDeletedServers models.ServerSlice
	FixtureAllServers     models.ServerSlice

	// ComponentFirmwareVersion fixtures
	FixtureDellR640BMC      *models.ComponentFirmwareVersion
	FixtureDellR640BIOS     *models.ComponentFirmwareVersion
	FixtureDellR640CPLD     *models.ComponentFirmwareVersion
	FixtureDellR6515BMC     *models.ComponentFirmwareVersion
	FixtureDellR6515BIOS    *models.ComponentFirmwareVersion
	FixtureSuperMicro       *models.ComponentFirmwareVersion
	FixtureServerComponents models.ServerComponentSlice

	// ComponentFirmwareSet fixtures
	FixtureFirmwareUUIDsR6515        []string
	FixtureFirmwareSetR6515          *models.ComponentFirmwareSet
	FixtureFirmwareSetR6515Attribute *models.AttributesFirmwareSet

	FixtureFirmwareUUIDsR640        []string
	FixtureFirmwareSetR640          *models.ComponentFirmwareSet
	FixtureFirmwareSetR640Attribute *models.AttributesFirmwareSet
)

func addFixtures(t *testing.T) error {
	ctx := context.TODO()

	FixtureFinType = &models.ServerComponentType{
		Name: "Fins",
		Slug: "fins",
	}

	if err := FixtureFinType.Insert(ctx, testDB, boil.Infer()); err != nil {
		return err
	}

	if err := setupNemo(ctx, testDB, t); err != nil {
		return err
	}

	if err := setupDory(ctx, testDB); err != nil {
		return err
	}

	if err := setupMarlin(ctx, testDB); err != nil {
		return err
	}

	if err := setupChuckles(ctx, testDB); err != nil {
		return err
	}

	if err := setupFirmwareDellR640(ctx, testDB); err != nil {
		return err
	}

	if err := setupFirmwareDellR6515(ctx, testDB); err != nil {
		return err
	}

	if err := setupFirmwareSuperMicro(ctx, testDB); err != nil {
		return err
	}

	if err := setupFirmwareSetR6515(ctx, testDB); err != nil {
		return err
	}

	if err := setupFirmwareSetR640(ctx, testDB); err != nil {
		return err
	}

	// excluding Chuckles here since that server is deleted
	FixtureServers = models.ServerSlice{FixtureNemo, FixtureDory, FixtureMarlin}
	FixtureDeletedServers = models.ServerSlice{FixtureChuckles}

	FixtureServerComponents = models.ServerComponentSlice{FixtureDoryLeftFin, FixtureDoryRightFin}

	//nolint:gocritic
	FixtureAllServers = append(FixtureServers, FixtureDeletedServers...)

	return nil
}

func setupNemo(ctx context.Context, db *sqlx.DB, t *testing.T) error {
	FixtureNemo = &models.Server{
		Name:         null.StringFrom("Nemo"),
		FacilityCode: null.StringFrom("Sydney"),
	}

	if err := FixtureNemo.Insert(ctx, db, boil.Infer()); err != nil {
		return err
	}

	bmc, err := models.ServerCredentialTypes(models.ServerCredentialTypeWhere.Slug.EQ("bmc")).One(ctx, db)
	if err != nil {
		return err
	}

	keeper := TestSecretKeeper(t)

	value, err := Encrypt(ctx, keeper, "super-secret-bmc-password")
	if err != nil {
		return err
	}

	FixtureNemoBMCSecret = &models.ServerCredential{
		ServerCredentialTypeID: bmc.ID,
		Password:               value,
	}

	if err := FixtureNemo.AddServerCredentials(ctx, db, true, FixtureNemoBMCSecret); err != nil {
		return err
	}

	FixtureNemoMetadata = &models.Attribute{
		Namespace: FixtureNamespaceMetadata,
		Data:      types.JSON([]byte(`{"age":6,"location":"Fishbowl"}`)),
	}

	FixtureNemoOtherdata = &models.Attribute{
		Namespace: FixtureNamespaceOtherdata,
		Data:      types.JSON([]byte(`{"enabled": true, "type": "clown", "lastUpdated": 1624960800, "nested": {"tag": "finding-nemo", "number": 1}}`)),
	}

	if err := FixtureNemo.AddAttributes(ctx, db, true, FixtureNemoMetadata, FixtureNemoOtherdata); err != nil {
		return err
	}

	FixtureNemoLeftFin = &models.ServerComponent{
		ServerComponentTypeID: FixtureFinType.ID,
		Name:                  null.StringFrom("Normal Fin"),
		Model:                 null.StringFrom("Normal Fin"),
		Serial:                null.StringFrom("Left"),
	}

	FixtureNemoRightFin = &models.ServerComponent{
		ServerComponentTypeID: FixtureFinType.ID,
		Name:                  null.StringFrom("My Lucky Fin"),
		Vendor:                null.StringFrom("Barracuda"),
		Model:                 null.StringFrom("A Lucky Fin"),
		Serial:                null.StringFrom("Right"),
	}

	if err := FixtureNemo.AddServerComponents(ctx, db, true, FixtureNemoLeftFin, FixtureNemoRightFin); err != nil {
		return err
	}

	FixtureNemoLeftFinVersioned = &models.VersionedAttribute{
		Namespace: FixtureNamespaceVersioned,
		Data:      types.JSON([]byte(`{"something": "cool"}`)),
	}

	if err := FixtureNemoLeftFin.AddVersionedAttributes(ctx, db, true, FixtureNemoLeftFinVersioned); err != nil {
		return err
	}

	FixtureNemoVersionedV2 = &models.VersionedAttribute{
		Namespace: FixtureNamespaceVersionedV2,
		Data:      types.JSON([]byte(`{"something": "cool"}`)),
	}

	if err := FixtureNemo.AddVersionedAttributes(ctx, db, true, FixtureNemoVersionedV2); err != nil {
		return err
	}

	FixtureNemoVersionedOld = &models.VersionedAttribute{
		Namespace: FixtureNamespaceVersioned,
		Data:      types.JSON([]byte(`{"name": "old"}`)),
	}

	FixtureNemoVersionedNew = &models.VersionedAttribute{
		Namespace: FixtureNamespaceVersioned,
		Data:      types.JSON([]byte(`{"name": "new"}`)),
	}

	// Insert old and new in a separate transaction to ensure the new one has a later timestamp and is indeed new
	if err := FixtureNemo.AddVersionedAttributes(ctx, db, true, FixtureNemoVersionedOld); err != nil {
		return err
	}

	return FixtureNemo.AddVersionedAttributes(ctx, db, true, FixtureNemoVersionedNew)
}

func setupDory(ctx context.Context, db *sqlx.DB) error {
	FixtureDory = &models.Server{
		Name:         null.StringFrom("Dory"),
		FacilityCode: null.StringFrom("Ocean"),
	}

	if err := FixtureDory.Insert(ctx, db, boil.Infer()); err != nil {
		return err
	}

	FixtureDoryMetadata = &models.Attribute{
		Namespace: FixtureNamespaceMetadata,
		Data:      types.JSON([]byte(`{"age":12,"location":"East Australian Current"}`)),
	}

	FixtureDoryOtherdata = &models.Attribute{
		Namespace: FixtureNamespaceOtherdata,
		Data:      types.JSON([]byte(`{"enabled": true, "type": "blue-tang", "lastUpdated": 1624960400, "nested": {"tag": "finding-nemo", "number": 2}}`)),
	}

	if err := FixtureDory.AddAttributes(ctx, db, true, FixtureDoryMetadata, FixtureDoryOtherdata); err != nil {
		return err
	}

	FixtureDoryLeftFin = &models.ServerComponent{
		ServerComponentTypeID: FixtureFinType.ID,
		Name:                  null.StringFrom("Normal Fin"),
		Model:                 null.StringFrom("Normal Fin"),
		Serial:                null.StringFrom("Left"),
	}

	FixtureDoryRightFin = &models.ServerComponent{
		ServerComponentTypeID: FixtureFinType.ID,
		Name:                  null.StringFrom("Normal Fin"),
		Serial:                null.StringFrom("Right"),
	}

	return FixtureDory.AddServerComponents(ctx, db, true, FixtureDoryLeftFin, FixtureDoryRightFin)
}

func setupMarlin(ctx context.Context, db *sqlx.DB) error {
	FixtureMarlin = &models.Server{
		Name:         null.StringFrom("Marlin"),
		FacilityCode: null.StringFrom("Ocean"),
	}

	if err := FixtureMarlin.Insert(ctx, db, boil.Infer()); err != nil {
		return err
	}

	FixtureMarlinMetadata = &models.Attribute{
		Namespace: FixtureNamespaceMetadata,
		Data:      types.JSON([]byte(`{"age":10,"location":"East Australian Current"}`)),
	}

	FixtureMarlinOtherdata = &models.Attribute{
		Namespace: FixtureNamespaceOtherdata,
		Data:      types.JSON([]byte(`{"enabled": false, "type": "clown", "lastUpdated": 1624960000, "nested": {"tag": "finding-nemo", "number": 3}}`)),
	}

	if err := FixtureMarlin.AddAttributes(ctx, db, true, FixtureMarlinMetadata, FixtureMarlinOtherdata); err != nil {
		return err
	}

	FixtureMarlinLeftFin = &models.ServerComponent{
		ServerComponentTypeID: FixtureFinType.ID,
		Name:                  null.StringFrom("Normal Fin"),
		Model:                 null.StringFrom("Normal Fin"),
		Serial:                null.StringFrom("Left"),
	}

	FixtureMarlinRightFin = &models.ServerComponent{
		ServerComponentTypeID: FixtureFinType.ID,
		Name:                  null.StringFrom("Normal Fin"),
		Serial:                null.StringFrom("Right"),
	}

	return FixtureMarlin.AddServerComponents(ctx, db, true, FixtureMarlinLeftFin, FixtureMarlinRightFin)
}

func setupChuckles(ctx context.Context, db *sqlx.DB) error {
	FixtureChuckles = &models.Server{
		Name:         null.StringFrom("Chuckles"),
		FacilityCode: null.StringFrom("Aquarium"),
		DeletedAt:    null.TimeFrom(time.Date(2003, 5, 30, 0, 0, 0, 0, time.UTC)),
	}

	if err := FixtureChuckles.Insert(ctx, db, boil.Infer()); err != nil {
		return err
	}

	FixtureChucklesMetadata = &models.Attribute{
		Namespace: FixtureNamespaceMetadata,
		Data:      types.JSON([]byte(`{"age":1,"location":"Old shipwreck"}`)),
	}

	FixtureChucklesOtherdata = &models.Attribute{
		Namespace: FixtureNamespaceOtherdata,
		Data:      types.JSON([]byte(`{"enabled": false, "type": "goldfish", "lastUpdated": 1624960000, "nested": {"tag": "finding-nemo", "number": 4}}`)),
	}

	if err := FixtureChuckles.AddAttributes(ctx, db, true, FixtureChucklesMetadata, FixtureChucklesOtherdata); err != nil {
		return err
	}

	FixtureChucklesLeftFin = &models.ServerComponent{
		ServerComponentTypeID: FixtureFinType.ID,
		Model:                 null.StringFrom("Belly"),
		Serial:                null.StringFrom("Up"),
	}

	return FixtureChuckles.AddServerComponents(ctx, db, true, FixtureChucklesLeftFin)
}

func setupFirmwareDellR640(ctx context.Context, db *sqlx.DB) error {
	FixtureFirmwareUUIDsR640 = []string{}

	FixtureDellR640BMC = &models.ComponentFirmwareVersion{
		Vendor:        "Dell",
		Model:         "R640",
		Filename:      "iDRAC-with-Lifecycle-Controller_Firmware_P8HC9_WN64_5.10.00.00_A00.EXE",
		Version:       "5.10.00.00",
		Component:     "bmc",
		Checksum:      "98db2fe5bca0745151d678ddeb26679464ccb13ca3f1a3d289b77e211344402f",
		UpstreamURL:   "https://vendor.com/firmwares/iDRAC-with-Lifecycle-Controller_Firmware_P8HC9_WN64_5.10.00.00_A00.EXE",
		RepositoryURL: "https://example-firmware-bucket.s3.amazonaws.com/firmware/dell/r640/bmc/iDRAC-with-Lifecycle-Controller_Firmware_P8HC9_WN64_5.10.00.00_A00.EXE",
	}

	if err := FixtureDellR640BMC.Insert(ctx, db, boil.Infer()); err != nil {
		return err
	}

	FixtureFirmwareUUIDsR640 = append(FixtureFirmwareUUIDsR640, FixtureDellR640BMC.ID)

	FixtureDellR640BIOS = &models.ComponentFirmwareVersion{
		Vendor:        "Dell",
		Model:         "R640",
		Filename:      "bios.exe",
		Version:       "2.4.4",
		Component:     "bios",
		Checksum:      "78ad2fe5bca0745151d678ddeb26679464ccb13ca3f1a3d289b77e211344402f",
		UpstreamURL:   "https://vendor.com/firmwares/bios-2.4.4.EXE",
		RepositoryURL: "https://example-firmware-bucket.s3.amazonaws.com/firmware/dell/r640/bios/bios-2.4.4.EXE",
	}

	if err := FixtureDellR640BIOS.Insert(ctx, db, boil.Infer()); err != nil {
		return err
	}

	FixtureFirmwareUUIDsR640 = append(FixtureFirmwareUUIDsR640, FixtureDellR640BIOS.ID)

	// This fixture is not included in FixtureFirmwareUUIDsR640 slice
	// since its part of the test TestIntegrationServerComponentFirmwareSetUpdate
	// where its added into the firmware set.
	FixtureDellR640CPLD = &models.ComponentFirmwareVersion{
		Vendor:        "Dell",
		Model:         "R640",
		Filename:      "cpld.exe",
		Version:       "1.0.1",
		Component:     "cpld",
		Checksum:      "676d2fe5bca0745151d678ddeb26679464ccb13ca3f1a3d289b77e211344402f",
		UpstreamURL:   "https://vendor.com/firmwares/cpld.exe",
		RepositoryURL: "https://example-firmware-bucket.s3.amazonaws.com/firmware/dell/r640/cpld/cpld.EXE",
	}

	if err := FixtureDellR640CPLD.Insert(ctx, db, boil.Infer()); err != nil {
		return err
	}

	return nil
}

func setupFirmwareDellR6515(ctx context.Context, db *sqlx.DB) error {
	FixtureFirmwareUUIDsR6515 = []string{}

	FixtureDellR6515BIOS = &models.ComponentFirmwareVersion{
		Vendor:        "Dell",
		Model:         "R6515",
		Filename:      "BIOS_C4FT0_WN64_2.6.6.EXE",
		Version:       "2.6.6",
		Component:     "bios",
		Checksum:      "1ddcb3c3d0fc5925ef03a3dde768e9e245c579039dd958fc0f3a9c6368b6c5f4",
		UpstreamURL:   "https://vendor.com/firmwares/BIOS_C4FT0_WN64_2.6.6.EXE",
		RepositoryURL: "https://example-firmware-bucket.s3.amazonaws.com/firmware/dell/r6515/bios/BIOS_C4FT0_WN64_2.6.6.EXE",
	}

	if err := FixtureDellR6515BIOS.Insert(ctx, db, boil.Infer()); err != nil {
		return err
	}

	FixtureFirmwareUUIDsR6515 = append(FixtureFirmwareUUIDsR6515, FixtureDellR6515BIOS.ID)

	FixtureDellR6515BMC = &models.ComponentFirmwareVersion{
		Vendor:        "Dell",
		Model:         "R6515",
		Filename:      "BMC-5.20.20.20.EXE",
		Version:       "5.20.20.20",
		Component:     "bmc",
		Checksum:      "abccb3c3d0fc5925ef03a3dde768e9e245c579039dd958fc0f3a9c6368b6c5f4",
		UpstreamURL:   "https://vendor.com/firmwares/BMC-5.20.20.20.EXE",
		RepositoryURL: "https://example-firmware-bucket.s3.amazonaws.com/firmware/dell/r6515/bmc/BMC-5.20.20.20.EXE",
	}

	if err := FixtureDellR6515BMC.Insert(ctx, db, boil.Infer()); err != nil {
		return err
	}

	FixtureFirmwareUUIDsR6515 = append(FixtureFirmwareUUIDsR6515, FixtureDellR6515BMC.ID)

	return nil
}

func setupFirmwareSuperMicro(ctx context.Context, db *sqlx.DB) error {
	FixtureSuperMicro = &models.ComponentFirmwareVersion{
		Vendor:        "SuperMicro",
		Model:         "X11DPH-T",
		Filename:      "SMT_X11AST2500_173_11.bin",
		Version:       "1.73.11",
		Component:     "bmc",
		Checksum:      "83d220484495e79a3c20e16c21a0d751a71519ac7058350d8a38e1f55efb0211",
		UpstreamURL:   "https://vendor.com/firmwares/SMT_X11AST2500_173_11.bin",
		RepositoryURL: "https://example-firmware-bucket.s3.amazonaws.com/firmware/supermicro/X11DPH-T/bmc/SMT_X11AST2500_173_11.bin",
	}

	if err := FixtureSuperMicro.Insert(ctx, db, boil.Infer()); err != nil {
		return err
	}

	return nil
}

func setupFirmwareSetR6515(ctx context.Context, db *sqlx.DB) error {
	// setup firmware fixtures if they haven't been
	if len(FixtureFirmwareUUIDsR6515) == 0 {
		if err := setupFirmwareDellR6515(ctx, db); err != nil {
			return err
		}
	}

	FixtureFirmwareSetR6515 = &models.ComponentFirmwareSet{Name: "r6515"}

	if err := FixtureFirmwareSetR6515.Insert(ctx, db, boil.Infer()); err != nil {
		return err
	}

	FixtureFirmwareSetR6515Attribute = &models.AttributesFirmwareSet{
		FirmwareSetID: null.StringFrom(FixtureFirmwareSetR6515.ID),
		Namespace:     "sh.hollow.firmware_set.labels",
		Data:          types.JSON([]byte(`{"vendor": "dell", "model": "r6515"}`)),
	}

	if err := FixtureFirmwareSetR6515.AddFirmwareSetAttributesFirmwareSets(ctx, db, true, FixtureFirmwareSetR6515Attribute); err != nil {
		return err
	}

	for _, firmwareID := range FixtureFirmwareUUIDsR6515 {
		m := &models.ComponentFirmwareSetMap{
			FirmwareSetID: FixtureFirmwareSetR6515.ID,
			FirmwareID:    firmwareID,
		}

		if err := m.Insert(ctx, db, boil.Infer()); err != nil {
			return err
		}
	}

	return nil
}

func setupFirmwareSetR640(ctx context.Context, db *sqlx.DB) error {
	// setup firmware fixtures if they haven't been
	if len(FixtureFirmwareUUIDsR640) == 0 {
		if err := setupFirmwareDellR640(ctx, db); err != nil {
			return err
		}
	}

	FixtureFirmwareSetR640 = &models.ComponentFirmwareSet{Name: "r640"}

	if err := FixtureFirmwareSetR640.Insert(ctx, db, boil.Infer()); err != nil {
		return err
	}

	FixtureFirmwareSetR640Attribute = &models.AttributesFirmwareSet{
		FirmwareSetID: null.StringFrom(FixtureFirmwareSetR640.ID),
		Namespace:     "sh.hollow.firmware_set.labels",
		Data:          types.JSON([]byte(`{"vendor": "dell", "model": "r640"}`)),
	}

	if err := FixtureFirmwareSetR640.AddFirmwareSetAttributesFirmwareSets(ctx, db, true, FixtureFirmwareSetR640Attribute); err != nil {
		return err
	}

	for _, firmwareID := range FixtureFirmwareUUIDsR640 {
		m := &models.ComponentFirmwareSetMap{
			FirmwareSetID: FixtureFirmwareSetR640.ID,
			FirmwareID:    firmwareID,
		}

		if err := m.Insert(ctx, db, boil.Infer()); err != nil {
			return err
		}
	}

	return nil
}
