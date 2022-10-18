// Code generated by SQLBoiler 4.11.0 (https://github.com/volatiletech/sqlboiler). DO NOT EDIT.
// This file is meant to be re-generated in place and/or deleted at any time.

package models

import (
	"bytes"
	"context"
	"reflect"
	"testing"

	"github.com/volatiletech/randomize"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/queries"
	"github.com/volatiletech/strmangle"
)

func testComponentFirmwareSetsUpsert(t *testing.T) {
	t.Parallel()

	if len(componentFirmwareSetAllColumns) == len(componentFirmwareSetPrimaryKeyColumns) {
		t.Skip("Skipping table with only primary key columns")
	}

	seed := randomize.NewSeed()
	var err error
	// Attempt the INSERT side of an UPSERT
	o := ComponentFirmwareSet{}
	if err = randomize.Struct(seed, &o, componentFirmwareSetDBTypes, true); err != nil {
		t.Errorf("Unable to randomize ComponentFirmwareSet struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Upsert(ctx, tx, false, nil, boil.Infer(), boil.Infer()); err != nil {
		t.Errorf("Unable to upsert ComponentFirmwareSet: %s", err)
	}

	count, err := ComponentFirmwareSets().Count(ctx, tx)
	if err != nil {
		t.Error(err)
	}
	if count != 1 {
		t.Error("want one record, got:", count)
	}

	// Attempt the UPDATE side of an UPSERT
	if err = randomize.Struct(seed, &o, componentFirmwareSetDBTypes, false, componentFirmwareSetPrimaryKeyColumns...); err != nil {
		t.Errorf("Unable to randomize ComponentFirmwareSet struct: %s", err)
	}

	if err = o.Upsert(ctx, tx, true, nil, boil.Infer(), boil.Infer()); err != nil {
		t.Errorf("Unable to upsert ComponentFirmwareSet: %s", err)
	}

	count, err = ComponentFirmwareSets().Count(ctx, tx)
	if err != nil {
		t.Error(err)
	}
	if count != 1 {
		t.Error("want one record, got:", count)
	}
}

var (
	// Relationships sometimes use the reflection helper queries.Equal/queries.Assign
	// so force a package dependency in case they don't.
	_ = queries.Equal
)

func testComponentFirmwareSets(t *testing.T) {
	t.Parallel()

	query := ComponentFirmwareSets()

	if query.Query == nil {
		t.Error("expected a query, got nothing")
	}
}

func testComponentFirmwareSetsDelete(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &ComponentFirmwareSet{}
	if err = randomize.Struct(seed, o, componentFirmwareSetDBTypes, true, componentFirmwareSetColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize ComponentFirmwareSet struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	if rowsAff, err := o.Delete(ctx, tx); err != nil {
		t.Error(err)
	} else if rowsAff != 1 {
		t.Error("should only have deleted one row, but affected:", rowsAff)
	}

	count, err := ComponentFirmwareSets().Count(ctx, tx)
	if err != nil {
		t.Error(err)
	}

	if count != 0 {
		t.Error("want zero records, got:", count)
	}
}

func testComponentFirmwareSetsQueryDeleteAll(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &ComponentFirmwareSet{}
	if err = randomize.Struct(seed, o, componentFirmwareSetDBTypes, true, componentFirmwareSetColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize ComponentFirmwareSet struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	if rowsAff, err := ComponentFirmwareSets().DeleteAll(ctx, tx); err != nil {
		t.Error(err)
	} else if rowsAff != 1 {
		t.Error("should only have deleted one row, but affected:", rowsAff)
	}

	count, err := ComponentFirmwareSets().Count(ctx, tx)
	if err != nil {
		t.Error(err)
	}

	if count != 0 {
		t.Error("want zero records, got:", count)
	}
}

func testComponentFirmwareSetsSliceDeleteAll(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &ComponentFirmwareSet{}
	if err = randomize.Struct(seed, o, componentFirmwareSetDBTypes, true, componentFirmwareSetColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize ComponentFirmwareSet struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	slice := ComponentFirmwareSetSlice{o}

	if rowsAff, err := slice.DeleteAll(ctx, tx); err != nil {
		t.Error(err)
	} else if rowsAff != 1 {
		t.Error("should only have deleted one row, but affected:", rowsAff)
	}

	count, err := ComponentFirmwareSets().Count(ctx, tx)
	if err != nil {
		t.Error(err)
	}

	if count != 0 {
		t.Error("want zero records, got:", count)
	}
}

func testComponentFirmwareSetsExists(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &ComponentFirmwareSet{}
	if err = randomize.Struct(seed, o, componentFirmwareSetDBTypes, true, componentFirmwareSetColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize ComponentFirmwareSet struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	e, err := ComponentFirmwareSetExists(ctx, tx, o.ID)
	if err != nil {
		t.Errorf("Unable to check if ComponentFirmwareSet exists: %s", err)
	}
	if !e {
		t.Errorf("Expected ComponentFirmwareSetExists to return true, but got false.")
	}
}

func testComponentFirmwareSetsFind(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &ComponentFirmwareSet{}
	if err = randomize.Struct(seed, o, componentFirmwareSetDBTypes, true, componentFirmwareSetColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize ComponentFirmwareSet struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	componentFirmwareSetFound, err := FindComponentFirmwareSet(ctx, tx, o.ID)
	if err != nil {
		t.Error(err)
	}

	if componentFirmwareSetFound == nil {
		t.Error("want a record, got nil")
	}
}

func testComponentFirmwareSetsBind(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &ComponentFirmwareSet{}
	if err = randomize.Struct(seed, o, componentFirmwareSetDBTypes, true, componentFirmwareSetColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize ComponentFirmwareSet struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	if err = ComponentFirmwareSets().Bind(ctx, tx, o); err != nil {
		t.Error(err)
	}
}

func testComponentFirmwareSetsOne(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &ComponentFirmwareSet{}
	if err = randomize.Struct(seed, o, componentFirmwareSetDBTypes, true, componentFirmwareSetColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize ComponentFirmwareSet struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	if x, err := ComponentFirmwareSets().One(ctx, tx); err != nil {
		t.Error(err)
	} else if x == nil {
		t.Error("expected to get a non nil record")
	}
}

func testComponentFirmwareSetsAll(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	componentFirmwareSetOne := &ComponentFirmwareSet{}
	componentFirmwareSetTwo := &ComponentFirmwareSet{}
	if err = randomize.Struct(seed, componentFirmwareSetOne, componentFirmwareSetDBTypes, false, componentFirmwareSetColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize ComponentFirmwareSet struct: %s", err)
	}
	if err = randomize.Struct(seed, componentFirmwareSetTwo, componentFirmwareSetDBTypes, false, componentFirmwareSetColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize ComponentFirmwareSet struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = componentFirmwareSetOne.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}
	if err = componentFirmwareSetTwo.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	slice, err := ComponentFirmwareSets().All(ctx, tx)
	if err != nil {
		t.Error(err)
	}

	if len(slice) != 2 {
		t.Error("want 2 records, got:", len(slice))
	}
}

func testComponentFirmwareSetsCount(t *testing.T) {
	t.Parallel()

	var err error
	seed := randomize.NewSeed()
	componentFirmwareSetOne := &ComponentFirmwareSet{}
	componentFirmwareSetTwo := &ComponentFirmwareSet{}
	if err = randomize.Struct(seed, componentFirmwareSetOne, componentFirmwareSetDBTypes, false, componentFirmwareSetColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize ComponentFirmwareSet struct: %s", err)
	}
	if err = randomize.Struct(seed, componentFirmwareSetTwo, componentFirmwareSetDBTypes, false, componentFirmwareSetColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize ComponentFirmwareSet struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = componentFirmwareSetOne.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}
	if err = componentFirmwareSetTwo.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	count, err := ComponentFirmwareSets().Count(ctx, tx)
	if err != nil {
		t.Error(err)
	}

	if count != 2 {
		t.Error("want 2 records, got:", count)
	}
}

func componentFirmwareSetBeforeInsertHook(ctx context.Context, e boil.ContextExecutor, o *ComponentFirmwareSet) error {
	*o = ComponentFirmwareSet{}
	return nil
}

func componentFirmwareSetAfterInsertHook(ctx context.Context, e boil.ContextExecutor, o *ComponentFirmwareSet) error {
	*o = ComponentFirmwareSet{}
	return nil
}

func componentFirmwareSetAfterSelectHook(ctx context.Context, e boil.ContextExecutor, o *ComponentFirmwareSet) error {
	*o = ComponentFirmwareSet{}
	return nil
}

func componentFirmwareSetBeforeUpdateHook(ctx context.Context, e boil.ContextExecutor, o *ComponentFirmwareSet) error {
	*o = ComponentFirmwareSet{}
	return nil
}

func componentFirmwareSetAfterUpdateHook(ctx context.Context, e boil.ContextExecutor, o *ComponentFirmwareSet) error {
	*o = ComponentFirmwareSet{}
	return nil
}

func componentFirmwareSetBeforeDeleteHook(ctx context.Context, e boil.ContextExecutor, o *ComponentFirmwareSet) error {
	*o = ComponentFirmwareSet{}
	return nil
}

func componentFirmwareSetAfterDeleteHook(ctx context.Context, e boil.ContextExecutor, o *ComponentFirmwareSet) error {
	*o = ComponentFirmwareSet{}
	return nil
}

func componentFirmwareSetBeforeUpsertHook(ctx context.Context, e boil.ContextExecutor, o *ComponentFirmwareSet) error {
	*o = ComponentFirmwareSet{}
	return nil
}

func componentFirmwareSetAfterUpsertHook(ctx context.Context, e boil.ContextExecutor, o *ComponentFirmwareSet) error {
	*o = ComponentFirmwareSet{}
	return nil
}

func testComponentFirmwareSetsHooks(t *testing.T) {
	t.Parallel()

	var err error

	ctx := context.Background()
	empty := &ComponentFirmwareSet{}
	o := &ComponentFirmwareSet{}

	seed := randomize.NewSeed()
	if err = randomize.Struct(seed, o, componentFirmwareSetDBTypes, false); err != nil {
		t.Errorf("Unable to randomize ComponentFirmwareSet object: %s", err)
	}

	AddComponentFirmwareSetHook(boil.BeforeInsertHook, componentFirmwareSetBeforeInsertHook)
	if err = o.doBeforeInsertHooks(ctx, nil); err != nil {
		t.Errorf("Unable to execute doBeforeInsertHooks: %s", err)
	}
	if !reflect.DeepEqual(o, empty) {
		t.Errorf("Expected BeforeInsertHook function to empty object, but got: %#v", o)
	}
	componentFirmwareSetBeforeInsertHooks = []ComponentFirmwareSetHook{}

	AddComponentFirmwareSetHook(boil.AfterInsertHook, componentFirmwareSetAfterInsertHook)
	if err = o.doAfterInsertHooks(ctx, nil); err != nil {
		t.Errorf("Unable to execute doAfterInsertHooks: %s", err)
	}
	if !reflect.DeepEqual(o, empty) {
		t.Errorf("Expected AfterInsertHook function to empty object, but got: %#v", o)
	}
	componentFirmwareSetAfterInsertHooks = []ComponentFirmwareSetHook{}

	AddComponentFirmwareSetHook(boil.AfterSelectHook, componentFirmwareSetAfterSelectHook)
	if err = o.doAfterSelectHooks(ctx, nil); err != nil {
		t.Errorf("Unable to execute doAfterSelectHooks: %s", err)
	}
	if !reflect.DeepEqual(o, empty) {
		t.Errorf("Expected AfterSelectHook function to empty object, but got: %#v", o)
	}
	componentFirmwareSetAfterSelectHooks = []ComponentFirmwareSetHook{}

	AddComponentFirmwareSetHook(boil.BeforeUpdateHook, componentFirmwareSetBeforeUpdateHook)
	if err = o.doBeforeUpdateHooks(ctx, nil); err != nil {
		t.Errorf("Unable to execute doBeforeUpdateHooks: %s", err)
	}
	if !reflect.DeepEqual(o, empty) {
		t.Errorf("Expected BeforeUpdateHook function to empty object, but got: %#v", o)
	}
	componentFirmwareSetBeforeUpdateHooks = []ComponentFirmwareSetHook{}

	AddComponentFirmwareSetHook(boil.AfterUpdateHook, componentFirmwareSetAfterUpdateHook)
	if err = o.doAfterUpdateHooks(ctx, nil); err != nil {
		t.Errorf("Unable to execute doAfterUpdateHooks: %s", err)
	}
	if !reflect.DeepEqual(o, empty) {
		t.Errorf("Expected AfterUpdateHook function to empty object, but got: %#v", o)
	}
	componentFirmwareSetAfterUpdateHooks = []ComponentFirmwareSetHook{}

	AddComponentFirmwareSetHook(boil.BeforeDeleteHook, componentFirmwareSetBeforeDeleteHook)
	if err = o.doBeforeDeleteHooks(ctx, nil); err != nil {
		t.Errorf("Unable to execute doBeforeDeleteHooks: %s", err)
	}
	if !reflect.DeepEqual(o, empty) {
		t.Errorf("Expected BeforeDeleteHook function to empty object, but got: %#v", o)
	}
	componentFirmwareSetBeforeDeleteHooks = []ComponentFirmwareSetHook{}

	AddComponentFirmwareSetHook(boil.AfterDeleteHook, componentFirmwareSetAfterDeleteHook)
	if err = o.doAfterDeleteHooks(ctx, nil); err != nil {
		t.Errorf("Unable to execute doAfterDeleteHooks: %s", err)
	}
	if !reflect.DeepEqual(o, empty) {
		t.Errorf("Expected AfterDeleteHook function to empty object, but got: %#v", o)
	}
	componentFirmwareSetAfterDeleteHooks = []ComponentFirmwareSetHook{}

	AddComponentFirmwareSetHook(boil.BeforeUpsertHook, componentFirmwareSetBeforeUpsertHook)
	if err = o.doBeforeUpsertHooks(ctx, nil); err != nil {
		t.Errorf("Unable to execute doBeforeUpsertHooks: %s", err)
	}
	if !reflect.DeepEqual(o, empty) {
		t.Errorf("Expected BeforeUpsertHook function to empty object, but got: %#v", o)
	}
	componentFirmwareSetBeforeUpsertHooks = []ComponentFirmwareSetHook{}

	AddComponentFirmwareSetHook(boil.AfterUpsertHook, componentFirmwareSetAfterUpsertHook)
	if err = o.doAfterUpsertHooks(ctx, nil); err != nil {
		t.Errorf("Unable to execute doAfterUpsertHooks: %s", err)
	}
	if !reflect.DeepEqual(o, empty) {
		t.Errorf("Expected AfterUpsertHook function to empty object, but got: %#v", o)
	}
	componentFirmwareSetAfterUpsertHooks = []ComponentFirmwareSetHook{}
}

func testComponentFirmwareSetsInsert(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &ComponentFirmwareSet{}
	if err = randomize.Struct(seed, o, componentFirmwareSetDBTypes, true, componentFirmwareSetColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize ComponentFirmwareSet struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	count, err := ComponentFirmwareSets().Count(ctx, tx)
	if err != nil {
		t.Error(err)
	}

	if count != 1 {
		t.Error("want one record, got:", count)
	}
}

func testComponentFirmwareSetsInsertWhitelist(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &ComponentFirmwareSet{}
	if err = randomize.Struct(seed, o, componentFirmwareSetDBTypes, true); err != nil {
		t.Errorf("Unable to randomize ComponentFirmwareSet struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Whitelist(componentFirmwareSetColumnsWithoutDefault...)); err != nil {
		t.Error(err)
	}

	count, err := ComponentFirmwareSets().Count(ctx, tx)
	if err != nil {
		t.Error(err)
	}

	if count != 1 {
		t.Error("want one record, got:", count)
	}
}

func testComponentFirmwareSetToManyAttributes(t *testing.T) {
	var err error
	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()

	var a ComponentFirmwareSet
	var b, c Attribute

	seed := randomize.NewSeed()
	if err = randomize.Struct(seed, &a, componentFirmwareSetDBTypes, true, componentFirmwareSetColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize ComponentFirmwareSet struct: %s", err)
	}

	if err := a.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Fatal(err)
	}

	if err = randomize.Struct(seed, &b, attributeDBTypes, false, attributeColumnsWithDefault...); err != nil {
		t.Fatal(err)
	}
	if err = randomize.Struct(seed, &c, attributeDBTypes, false, attributeColumnsWithDefault...); err != nil {
		t.Fatal(err)
	}

	queries.Assign(&b.ComponentFirmwareSetID, a.ID)
	queries.Assign(&c.ComponentFirmwareSetID, a.ID)
	if err = b.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Fatal(err)
	}
	if err = c.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Fatal(err)
	}

	check, err := a.Attributes().All(ctx, tx)
	if err != nil {
		t.Fatal(err)
	}

	bFound, cFound := false, false
	for _, v := range check {
		if queries.Equal(v.ComponentFirmwareSetID, b.ComponentFirmwareSetID) {
			bFound = true
		}
		if queries.Equal(v.ComponentFirmwareSetID, c.ComponentFirmwareSetID) {
			cFound = true
		}
	}

	if !bFound {
		t.Error("expected to find b")
	}
	if !cFound {
		t.Error("expected to find c")
	}

	slice := ComponentFirmwareSetSlice{&a}
	if err = a.L.LoadAttributes(ctx, tx, false, (*[]*ComponentFirmwareSet)(&slice), nil); err != nil {
		t.Fatal(err)
	}
	if got := len(a.R.Attributes); got != 2 {
		t.Error("number of eager loaded records wrong, got:", got)
	}

	a.R.Attributes = nil
	if err = a.L.LoadAttributes(ctx, tx, true, &a, nil); err != nil {
		t.Fatal(err)
	}
	if got := len(a.R.Attributes); got != 2 {
		t.Error("number of eager loaded records wrong, got:", got)
	}

	if t.Failed() {
		t.Logf("%#v", check)
	}
}

func testComponentFirmwareSetToManyFirmwareSetComponentFirmwareSetMaps(t *testing.T) {
	var err error
	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()

	var a ComponentFirmwareSet
	var b, c ComponentFirmwareSetMap

	seed := randomize.NewSeed()
	if err = randomize.Struct(seed, &a, componentFirmwareSetDBTypes, true, componentFirmwareSetColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize ComponentFirmwareSet struct: %s", err)
	}

	if err := a.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Fatal(err)
	}

	if err = randomize.Struct(seed, &b, componentFirmwareSetMapDBTypes, false, componentFirmwareSetMapColumnsWithDefault...); err != nil {
		t.Fatal(err)
	}
	if err = randomize.Struct(seed, &c, componentFirmwareSetMapDBTypes, false, componentFirmwareSetMapColumnsWithDefault...); err != nil {
		t.Fatal(err)
	}

	b.FirmwareSetID = a.ID
	c.FirmwareSetID = a.ID

	if err = b.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Fatal(err)
	}
	if err = c.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Fatal(err)
	}

	check, err := a.FirmwareSetComponentFirmwareSetMaps().All(ctx, tx)
	if err != nil {
		t.Fatal(err)
	}

	bFound, cFound := false, false
	for _, v := range check {
		if v.FirmwareSetID == b.FirmwareSetID {
			bFound = true
		}
		if v.FirmwareSetID == c.FirmwareSetID {
			cFound = true
		}
	}

	if !bFound {
		t.Error("expected to find b")
	}
	if !cFound {
		t.Error("expected to find c")
	}

	slice := ComponentFirmwareSetSlice{&a}
	if err = a.L.LoadFirmwareSetComponentFirmwareSetMaps(ctx, tx, false, (*[]*ComponentFirmwareSet)(&slice), nil); err != nil {
		t.Fatal(err)
	}
	if got := len(a.R.FirmwareSetComponentFirmwareSetMaps); got != 2 {
		t.Error("number of eager loaded records wrong, got:", got)
	}

	a.R.FirmwareSetComponentFirmwareSetMaps = nil
	if err = a.L.LoadFirmwareSetComponentFirmwareSetMaps(ctx, tx, true, &a, nil); err != nil {
		t.Fatal(err)
	}
	if got := len(a.R.FirmwareSetComponentFirmwareSetMaps); got != 2 {
		t.Error("number of eager loaded records wrong, got:", got)
	}

	if t.Failed() {
		t.Logf("%#v", check)
	}
}

func testComponentFirmwareSetToManyAddOpAttributes(t *testing.T) {
	var err error

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()

	var a ComponentFirmwareSet
	var b, c, d, e Attribute

	seed := randomize.NewSeed()
	if err = randomize.Struct(seed, &a, componentFirmwareSetDBTypes, false, strmangle.SetComplement(componentFirmwareSetPrimaryKeyColumns, componentFirmwareSetColumnsWithoutDefault)...); err != nil {
		t.Fatal(err)
	}
	foreigners := []*Attribute{&b, &c, &d, &e}
	for _, x := range foreigners {
		if err = randomize.Struct(seed, x, attributeDBTypes, false, strmangle.SetComplement(attributePrimaryKeyColumns, attributeColumnsWithoutDefault)...); err != nil {
			t.Fatal(err)
		}
	}

	if err := a.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Fatal(err)
	}
	if err = b.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Fatal(err)
	}
	if err = c.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Fatal(err)
	}

	foreignersSplitByInsertion := [][]*Attribute{
		{&b, &c},
		{&d, &e},
	}

	for i, x := range foreignersSplitByInsertion {
		err = a.AddAttributes(ctx, tx, i != 0, x...)
		if err != nil {
			t.Fatal(err)
		}

		first := x[0]
		second := x[1]

		if !queries.Equal(a.ID, first.ComponentFirmwareSetID) {
			t.Error("foreign key was wrong value", a.ID, first.ComponentFirmwareSetID)
		}
		if !queries.Equal(a.ID, second.ComponentFirmwareSetID) {
			t.Error("foreign key was wrong value", a.ID, second.ComponentFirmwareSetID)
		}

		if first.R.ComponentFirmwareSet != &a {
			t.Error("relationship was not added properly to the foreign slice")
		}
		if second.R.ComponentFirmwareSet != &a {
			t.Error("relationship was not added properly to the foreign slice")
		}

		if a.R.Attributes[i*2] != first {
			t.Error("relationship struct slice not set to correct value")
		}
		if a.R.Attributes[i*2+1] != second {
			t.Error("relationship struct slice not set to correct value")
		}

		count, err := a.Attributes().Count(ctx, tx)
		if err != nil {
			t.Fatal(err)
		}
		if want := int64((i + 1) * 2); count != want {
			t.Error("want", want, "got", count)
		}
	}
}

func testComponentFirmwareSetToManySetOpAttributes(t *testing.T) {
	var err error

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()

	var a ComponentFirmwareSet
	var b, c, d, e Attribute

	seed := randomize.NewSeed()
	if err = randomize.Struct(seed, &a, componentFirmwareSetDBTypes, false, strmangle.SetComplement(componentFirmwareSetPrimaryKeyColumns, componentFirmwareSetColumnsWithoutDefault)...); err != nil {
		t.Fatal(err)
	}
	foreigners := []*Attribute{&b, &c, &d, &e}
	for _, x := range foreigners {
		if err = randomize.Struct(seed, x, attributeDBTypes, false, strmangle.SetComplement(attributePrimaryKeyColumns, attributeColumnsWithoutDefault)...); err != nil {
			t.Fatal(err)
		}
	}

	if err = a.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Fatal(err)
	}
	if err = b.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Fatal(err)
	}
	if err = c.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Fatal(err)
	}

	err = a.SetAttributes(ctx, tx, false, &b, &c)
	if err != nil {
		t.Fatal(err)
	}

	count, err := a.Attributes().Count(ctx, tx)
	if err != nil {
		t.Fatal(err)
	}
	if count != 2 {
		t.Error("count was wrong:", count)
	}

	err = a.SetAttributes(ctx, tx, true, &d, &e)
	if err != nil {
		t.Fatal(err)
	}

	count, err = a.Attributes().Count(ctx, tx)
	if err != nil {
		t.Fatal(err)
	}
	if count != 2 {
		t.Error("count was wrong:", count)
	}

	if !queries.IsValuerNil(b.ComponentFirmwareSetID) {
		t.Error("want b's foreign key value to be nil")
	}
	if !queries.IsValuerNil(c.ComponentFirmwareSetID) {
		t.Error("want c's foreign key value to be nil")
	}
	if !queries.Equal(a.ID, d.ComponentFirmwareSetID) {
		t.Error("foreign key was wrong value", a.ID, d.ComponentFirmwareSetID)
	}
	if !queries.Equal(a.ID, e.ComponentFirmwareSetID) {
		t.Error("foreign key was wrong value", a.ID, e.ComponentFirmwareSetID)
	}

	if b.R.ComponentFirmwareSet != nil {
		t.Error("relationship was not removed properly from the foreign struct")
	}
	if c.R.ComponentFirmwareSet != nil {
		t.Error("relationship was not removed properly from the foreign struct")
	}
	if d.R.ComponentFirmwareSet != &a {
		t.Error("relationship was not added properly to the foreign struct")
	}
	if e.R.ComponentFirmwareSet != &a {
		t.Error("relationship was not added properly to the foreign struct")
	}

	if a.R.Attributes[0] != &d {
		t.Error("relationship struct slice not set to correct value")
	}
	if a.R.Attributes[1] != &e {
		t.Error("relationship struct slice not set to correct value")
	}
}

func testComponentFirmwareSetToManyRemoveOpAttributes(t *testing.T) {
	var err error

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()

	var a ComponentFirmwareSet
	var b, c, d, e Attribute

	seed := randomize.NewSeed()
	if err = randomize.Struct(seed, &a, componentFirmwareSetDBTypes, false, strmangle.SetComplement(componentFirmwareSetPrimaryKeyColumns, componentFirmwareSetColumnsWithoutDefault)...); err != nil {
		t.Fatal(err)
	}
	foreigners := []*Attribute{&b, &c, &d, &e}
	for _, x := range foreigners {
		if err = randomize.Struct(seed, x, attributeDBTypes, false, strmangle.SetComplement(attributePrimaryKeyColumns, attributeColumnsWithoutDefault)...); err != nil {
			t.Fatal(err)
		}
	}

	if err := a.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Fatal(err)
	}

	err = a.AddAttributes(ctx, tx, true, foreigners...)
	if err != nil {
		t.Fatal(err)
	}

	count, err := a.Attributes().Count(ctx, tx)
	if err != nil {
		t.Fatal(err)
	}
	if count != 4 {
		t.Error("count was wrong:", count)
	}

	err = a.RemoveAttributes(ctx, tx, foreigners[:2]...)
	if err != nil {
		t.Fatal(err)
	}

	count, err = a.Attributes().Count(ctx, tx)
	if err != nil {
		t.Fatal(err)
	}
	if count != 2 {
		t.Error("count was wrong:", count)
	}

	if !queries.IsValuerNil(b.ComponentFirmwareSetID) {
		t.Error("want b's foreign key value to be nil")
	}
	if !queries.IsValuerNil(c.ComponentFirmwareSetID) {
		t.Error("want c's foreign key value to be nil")
	}

	if b.R.ComponentFirmwareSet != nil {
		t.Error("relationship was not removed properly from the foreign struct")
	}
	if c.R.ComponentFirmwareSet != nil {
		t.Error("relationship was not removed properly from the foreign struct")
	}
	if d.R.ComponentFirmwareSet != &a {
		t.Error("relationship to a should have been preserved")
	}
	if e.R.ComponentFirmwareSet != &a {
		t.Error("relationship to a should have been preserved")
	}

	if len(a.R.Attributes) != 2 {
		t.Error("should have preserved two relationships")
	}

	// Removal doesn't do a stable deletion for performance so we have to flip the order
	if a.R.Attributes[1] != &d {
		t.Error("relationship to d should have been preserved")
	}
	if a.R.Attributes[0] != &e {
		t.Error("relationship to e should have been preserved")
	}
}

func testComponentFirmwareSetToManyAddOpFirmwareSetComponentFirmwareSetMaps(t *testing.T) {
	var err error

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()

	var a ComponentFirmwareSet
	var b, c, d, e ComponentFirmwareSetMap

	seed := randomize.NewSeed()
	if err = randomize.Struct(seed, &a, componentFirmwareSetDBTypes, false, strmangle.SetComplement(componentFirmwareSetPrimaryKeyColumns, componentFirmwareSetColumnsWithoutDefault)...); err != nil {
		t.Fatal(err)
	}
	foreigners := []*ComponentFirmwareSetMap{&b, &c, &d, &e}
	for _, x := range foreigners {
		if err = randomize.Struct(seed, x, componentFirmwareSetMapDBTypes, false, strmangle.SetComplement(componentFirmwareSetMapPrimaryKeyColumns, componentFirmwareSetMapColumnsWithoutDefault)...); err != nil {
			t.Fatal(err)
		}
	}

	if err := a.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Fatal(err)
	}
	if err = b.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Fatal(err)
	}
	if err = c.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Fatal(err)
	}

	foreignersSplitByInsertion := [][]*ComponentFirmwareSetMap{
		{&b, &c},
		{&d, &e},
	}

	for i, x := range foreignersSplitByInsertion {
		err = a.AddFirmwareSetComponentFirmwareSetMaps(ctx, tx, i != 0, x...)
		if err != nil {
			t.Fatal(err)
		}

		first := x[0]
		second := x[1]

		if a.ID != first.FirmwareSetID {
			t.Error("foreign key was wrong value", a.ID, first.FirmwareSetID)
		}
		if a.ID != second.FirmwareSetID {
			t.Error("foreign key was wrong value", a.ID, second.FirmwareSetID)
		}

		if first.R.FirmwareSet != &a {
			t.Error("relationship was not added properly to the foreign slice")
		}
		if second.R.FirmwareSet != &a {
			t.Error("relationship was not added properly to the foreign slice")
		}

		if a.R.FirmwareSetComponentFirmwareSetMaps[i*2] != first {
			t.Error("relationship struct slice not set to correct value")
		}
		if a.R.FirmwareSetComponentFirmwareSetMaps[i*2+1] != second {
			t.Error("relationship struct slice not set to correct value")
		}

		count, err := a.FirmwareSetComponentFirmwareSetMaps().Count(ctx, tx)
		if err != nil {
			t.Fatal(err)
		}
		if want := int64((i + 1) * 2); count != want {
			t.Error("want", want, "got", count)
		}
	}
}

func testComponentFirmwareSetsReload(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &ComponentFirmwareSet{}
	if err = randomize.Struct(seed, o, componentFirmwareSetDBTypes, true, componentFirmwareSetColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize ComponentFirmwareSet struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	if err = o.Reload(ctx, tx); err != nil {
		t.Error(err)
	}
}

func testComponentFirmwareSetsReloadAll(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &ComponentFirmwareSet{}
	if err = randomize.Struct(seed, o, componentFirmwareSetDBTypes, true, componentFirmwareSetColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize ComponentFirmwareSet struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	slice := ComponentFirmwareSetSlice{o}

	if err = slice.ReloadAll(ctx, tx); err != nil {
		t.Error(err)
	}
}

func testComponentFirmwareSetsSelect(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &ComponentFirmwareSet{}
	if err = randomize.Struct(seed, o, componentFirmwareSetDBTypes, true, componentFirmwareSetColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize ComponentFirmwareSet struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	slice, err := ComponentFirmwareSets().All(ctx, tx)
	if err != nil {
		t.Error(err)
	}

	if len(slice) != 1 {
		t.Error("want one record, got:", len(slice))
	}
}

var (
	componentFirmwareSetDBTypes = map[string]string{`ID`: `uuid`, `Name`: `string`, `CreatedAt`: `timestamptz`, `UpdatedAt`: `timestamptz`}
	_                           = bytes.MinRead
)

func testComponentFirmwareSetsUpdate(t *testing.T) {
	t.Parallel()

	if 0 == len(componentFirmwareSetPrimaryKeyColumns) {
		t.Skip("Skipping table with no primary key columns")
	}
	if len(componentFirmwareSetAllColumns) == len(componentFirmwareSetPrimaryKeyColumns) {
		t.Skip("Skipping table with only primary key columns")
	}

	seed := randomize.NewSeed()
	var err error
	o := &ComponentFirmwareSet{}
	if err = randomize.Struct(seed, o, componentFirmwareSetDBTypes, true, componentFirmwareSetColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize ComponentFirmwareSet struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	count, err := ComponentFirmwareSets().Count(ctx, tx)
	if err != nil {
		t.Error(err)
	}

	if count != 1 {
		t.Error("want one record, got:", count)
	}

	if err = randomize.Struct(seed, o, componentFirmwareSetDBTypes, true, componentFirmwareSetPrimaryKeyColumns...); err != nil {
		t.Errorf("Unable to randomize ComponentFirmwareSet struct: %s", err)
	}

	if rowsAff, err := o.Update(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	} else if rowsAff != 1 {
		t.Error("should only affect one row but affected", rowsAff)
	}
}

func testComponentFirmwareSetsSliceUpdateAll(t *testing.T) {
	t.Parallel()

	if len(componentFirmwareSetAllColumns) == len(componentFirmwareSetPrimaryKeyColumns) {
		t.Skip("Skipping table with only primary key columns")
	}

	seed := randomize.NewSeed()
	var err error
	o := &ComponentFirmwareSet{}
	if err = randomize.Struct(seed, o, componentFirmwareSetDBTypes, true, componentFirmwareSetColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize ComponentFirmwareSet struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	count, err := ComponentFirmwareSets().Count(ctx, tx)
	if err != nil {
		t.Error(err)
	}

	if count != 1 {
		t.Error("want one record, got:", count)
	}

	if err = randomize.Struct(seed, o, componentFirmwareSetDBTypes, true, componentFirmwareSetPrimaryKeyColumns...); err != nil {
		t.Errorf("Unable to randomize ComponentFirmwareSet struct: %s", err)
	}

	// Remove Primary keys and unique columns from what we plan to update
	var fields []string
	if strmangle.StringSliceMatch(componentFirmwareSetAllColumns, componentFirmwareSetPrimaryKeyColumns) {
		fields = componentFirmwareSetAllColumns
	} else {
		fields = strmangle.SetComplement(
			componentFirmwareSetAllColumns,
			componentFirmwareSetPrimaryKeyColumns,
		)
	}

	value := reflect.Indirect(reflect.ValueOf(o))
	typ := reflect.TypeOf(o).Elem()
	n := typ.NumField()

	updateMap := M{}
	for _, col := range fields {
		for i := 0; i < n; i++ {
			f := typ.Field(i)
			if f.Tag.Get("boil") == col {
				updateMap[col] = value.Field(i).Interface()
			}
		}
	}

	slice := ComponentFirmwareSetSlice{o}
	if rowsAff, err := slice.UpdateAll(ctx, tx, updateMap); err != nil {
		t.Error(err)
	} else if rowsAff != 1 {
		t.Error("wanted one record updated but got", rowsAff)
	}
}
