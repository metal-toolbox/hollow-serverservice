//go:build testtools
// +build testtools

package dbtools

import (
	"context"
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
	FixtureDell210700 *models.ComponentFirmwareVersion
	FixtureDell210501 *models.ComponentFirmwareVersion
	FixtureSuperMicro *models.ComponentFirmwareVersion
)

func addFixtures() error {
	ctx := context.TODO()

	FixtureFinType = &models.ServerComponentType{
		Name: "Fins",
		Slug: "fins",
	}

	if err := FixtureFinType.Insert(ctx, testDB, boil.Infer()); err != nil {
		return err
	}

	if err := setupNemo(ctx, testDB); err != nil {
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

	if err := setupFirmwareDell210700(ctx, testDB); err != nil {
		return err
	}

	if err := setupFirmwareDell210501(ctx, testDB); err != nil {
		return err
	}

	if err := setupFirmwareSuperMicro(ctx, testDB); err != nil {
		return err
	}

	// excluding Chuckles here since that server is deleted
	FixtureServers = models.ServerSlice{FixtureNemo, FixtureDory, FixtureMarlin}
	FixtureDeletedServers = models.ServerSlice{FixtureChuckles}

	//nolint:gocritic
	FixtureAllServers = append(FixtureServers, FixtureDeletedServers...)

	return nil
}

func setupNemo(ctx context.Context, db *sqlx.DB) error {
	FixtureNemo = &models.Server{
		Name:         null.StringFrom("Nemo"),
		FacilityCode: null.StringFrom("Sydney"),
	}

	if err := FixtureNemo.Insert(ctx, db, boil.Infer()); err != nil {
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

func setupFirmwareDell210700(ctx context.Context, db *sqlx.DB) error {
	FixtureDell210700 = &models.ComponentFirmwareVersion{
		Vendor:      null.StringFrom("Dell"),
		Model:       null.StringFrom("R615"),
		Filename:    null.StringFrom("foobar"),
		Version:     null.StringFrom("21.07.00"),
		Component:   null.StringFrom("system"),
		Utility:     null.StringFrom("dsu"),
		Sha:         null.StringFrom("foobar"),
		UpstreamURL: null.StringFrom("https://linux.dell.com/repo/hardware/DSU_21.07.00/"),
	}

	if err := FixtureDell210700.Insert(ctx, db, boil.Infer()); err != nil {
		return err
	}

	return nil
}

func setupFirmwareDell210501(ctx context.Context, db *sqlx.DB) error {
	FixtureDell210501 = &models.ComponentFirmwareVersion{
		Vendor:      null.StringFrom("Dell"),
		Model:       null.StringFrom("R615"),
		Filename:    null.StringFrom("foobar"),
		Version:     null.StringFrom("21.05.01"),
		Component:   null.StringFrom("system"),
		Utility:     null.StringFrom("dsu"),
		Sha:         null.StringFrom("foobar"),
		UpstreamURL: null.StringFrom("https://linux.dell.com/repo/hardware/DSU_21.07.00/"),
	}

	if err := FixtureDell210501.Insert(ctx, db, boil.Infer()); err != nil {
		return err
	}

	return nil
}

func setupFirmwareSuperMicro(ctx context.Context, db *sqlx.DB) error {
	FixtureSuperMicro = &models.ComponentFirmwareVersion{
		Vendor:      null.StringFrom("SuperMicro"),
		Model:       null.StringFrom("X11DPH-T"),
		Filename:    null.StringFrom("SMT_X11AST2500_173_11.bin"),
		Version:     null.StringFrom("1.73.11"),
		Component:   null.StringFrom("bmc"),
		Utility:     null.StringFrom("sum"),
		Sha:         null.StringFrom("83d220484495e79a3c20e16c21a0d751a71519ac7058350d8a38e1f55efb0211"),
		UpstreamURL: null.StringFrom("http://install.packet.net/firmware/fup/supermicro/X11DPH-T/bmc/SMT_X11AST2500_173_11.bin"),
	}

	if err := FixtureSuperMicro.Insert(ctx, db, boil.Infer()); err != nil {
		return err
	}

	return nil
}
