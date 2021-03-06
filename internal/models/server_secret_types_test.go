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

func testServerSecretTypesUpsert(t *testing.T) {
	t.Parallel()

	if len(serverSecretTypeAllColumns) == len(serverSecretTypePrimaryKeyColumns) {
		t.Skip("Skipping table with only primary key columns")
	}

	seed := randomize.NewSeed()
	var err error
	// Attempt the INSERT side of an UPSERT
	o := ServerSecretType{}
	if err = randomize.Struct(seed, &o, serverSecretTypeDBTypes, true); err != nil {
		t.Errorf("Unable to randomize ServerSecretType struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Upsert(ctx, tx, false, nil, boil.Infer(), boil.Infer()); err != nil {
		t.Errorf("Unable to upsert ServerSecretType: %s", err)
	}

	count, err := ServerSecretTypes().Count(ctx, tx)
	if err != nil {
		t.Error(err)
	}
	if count != 1 {
		t.Error("want one record, got:", count)
	}

	// Attempt the UPDATE side of an UPSERT
	if err = randomize.Struct(seed, &o, serverSecretTypeDBTypes, false, serverSecretTypePrimaryKeyColumns...); err != nil {
		t.Errorf("Unable to randomize ServerSecretType struct: %s", err)
	}

	if err = o.Upsert(ctx, tx, true, nil, boil.Infer(), boil.Infer()); err != nil {
		t.Errorf("Unable to upsert ServerSecretType: %s", err)
	}

	count, err = ServerSecretTypes().Count(ctx, tx)
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

func testServerSecretTypes(t *testing.T) {
	t.Parallel()

	query := ServerSecretTypes()

	if query.Query == nil {
		t.Error("expected a query, got nothing")
	}
}

func testServerSecretTypesDelete(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &ServerSecretType{}
	if err = randomize.Struct(seed, o, serverSecretTypeDBTypes, true, serverSecretTypeColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize ServerSecretType struct: %s", err)
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

	count, err := ServerSecretTypes().Count(ctx, tx)
	if err != nil {
		t.Error(err)
	}

	if count != 0 {
		t.Error("want zero records, got:", count)
	}
}

func testServerSecretTypesQueryDeleteAll(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &ServerSecretType{}
	if err = randomize.Struct(seed, o, serverSecretTypeDBTypes, true, serverSecretTypeColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize ServerSecretType struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	if rowsAff, err := ServerSecretTypes().DeleteAll(ctx, tx); err != nil {
		t.Error(err)
	} else if rowsAff != 1 {
		t.Error("should only have deleted one row, but affected:", rowsAff)
	}

	count, err := ServerSecretTypes().Count(ctx, tx)
	if err != nil {
		t.Error(err)
	}

	if count != 0 {
		t.Error("want zero records, got:", count)
	}
}

func testServerSecretTypesSliceDeleteAll(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &ServerSecretType{}
	if err = randomize.Struct(seed, o, serverSecretTypeDBTypes, true, serverSecretTypeColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize ServerSecretType struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	slice := ServerSecretTypeSlice{o}

	if rowsAff, err := slice.DeleteAll(ctx, tx); err != nil {
		t.Error(err)
	} else if rowsAff != 1 {
		t.Error("should only have deleted one row, but affected:", rowsAff)
	}

	count, err := ServerSecretTypes().Count(ctx, tx)
	if err != nil {
		t.Error(err)
	}

	if count != 0 {
		t.Error("want zero records, got:", count)
	}
}

func testServerSecretTypesExists(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &ServerSecretType{}
	if err = randomize.Struct(seed, o, serverSecretTypeDBTypes, true, serverSecretTypeColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize ServerSecretType struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	e, err := ServerSecretTypeExists(ctx, tx, o.ID)
	if err != nil {
		t.Errorf("Unable to check if ServerSecretType exists: %s", err)
	}
	if !e {
		t.Errorf("Expected ServerSecretTypeExists to return true, but got false.")
	}
}

func testServerSecretTypesFind(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &ServerSecretType{}
	if err = randomize.Struct(seed, o, serverSecretTypeDBTypes, true, serverSecretTypeColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize ServerSecretType struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	serverSecretTypeFound, err := FindServerSecretType(ctx, tx, o.ID)
	if err != nil {
		t.Error(err)
	}

	if serverSecretTypeFound == nil {
		t.Error("want a record, got nil")
	}
}

func testServerSecretTypesBind(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &ServerSecretType{}
	if err = randomize.Struct(seed, o, serverSecretTypeDBTypes, true, serverSecretTypeColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize ServerSecretType struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	if err = ServerSecretTypes().Bind(ctx, tx, o); err != nil {
		t.Error(err)
	}
}

func testServerSecretTypesOne(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &ServerSecretType{}
	if err = randomize.Struct(seed, o, serverSecretTypeDBTypes, true, serverSecretTypeColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize ServerSecretType struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	if x, err := ServerSecretTypes().One(ctx, tx); err != nil {
		t.Error(err)
	} else if x == nil {
		t.Error("expected to get a non nil record")
	}
}

func testServerSecretTypesAll(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	serverSecretTypeOne := &ServerSecretType{}
	serverSecretTypeTwo := &ServerSecretType{}
	if err = randomize.Struct(seed, serverSecretTypeOne, serverSecretTypeDBTypes, false, serverSecretTypeColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize ServerSecretType struct: %s", err)
	}
	if err = randomize.Struct(seed, serverSecretTypeTwo, serverSecretTypeDBTypes, false, serverSecretTypeColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize ServerSecretType struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = serverSecretTypeOne.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}
	if err = serverSecretTypeTwo.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	slice, err := ServerSecretTypes().All(ctx, tx)
	if err != nil {
		t.Error(err)
	}

	if len(slice) != 2 {
		t.Error("want 2 records, got:", len(slice))
	}
}

func testServerSecretTypesCount(t *testing.T) {
	t.Parallel()

	var err error
	seed := randomize.NewSeed()
	serverSecretTypeOne := &ServerSecretType{}
	serverSecretTypeTwo := &ServerSecretType{}
	if err = randomize.Struct(seed, serverSecretTypeOne, serverSecretTypeDBTypes, false, serverSecretTypeColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize ServerSecretType struct: %s", err)
	}
	if err = randomize.Struct(seed, serverSecretTypeTwo, serverSecretTypeDBTypes, false, serverSecretTypeColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize ServerSecretType struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = serverSecretTypeOne.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}
	if err = serverSecretTypeTwo.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	count, err := ServerSecretTypes().Count(ctx, tx)
	if err != nil {
		t.Error(err)
	}

	if count != 2 {
		t.Error("want 2 records, got:", count)
	}
}

func serverSecretTypeBeforeInsertHook(ctx context.Context, e boil.ContextExecutor, o *ServerSecretType) error {
	*o = ServerSecretType{}
	return nil
}

func serverSecretTypeAfterInsertHook(ctx context.Context, e boil.ContextExecutor, o *ServerSecretType) error {
	*o = ServerSecretType{}
	return nil
}

func serverSecretTypeAfterSelectHook(ctx context.Context, e boil.ContextExecutor, o *ServerSecretType) error {
	*o = ServerSecretType{}
	return nil
}

func serverSecretTypeBeforeUpdateHook(ctx context.Context, e boil.ContextExecutor, o *ServerSecretType) error {
	*o = ServerSecretType{}
	return nil
}

func serverSecretTypeAfterUpdateHook(ctx context.Context, e boil.ContextExecutor, o *ServerSecretType) error {
	*o = ServerSecretType{}
	return nil
}

func serverSecretTypeBeforeDeleteHook(ctx context.Context, e boil.ContextExecutor, o *ServerSecretType) error {
	*o = ServerSecretType{}
	return nil
}

func serverSecretTypeAfterDeleteHook(ctx context.Context, e boil.ContextExecutor, o *ServerSecretType) error {
	*o = ServerSecretType{}
	return nil
}

func serverSecretTypeBeforeUpsertHook(ctx context.Context, e boil.ContextExecutor, o *ServerSecretType) error {
	*o = ServerSecretType{}
	return nil
}

func serverSecretTypeAfterUpsertHook(ctx context.Context, e boil.ContextExecutor, o *ServerSecretType) error {
	*o = ServerSecretType{}
	return nil
}

func testServerSecretTypesHooks(t *testing.T) {
	t.Parallel()

	var err error

	ctx := context.Background()
	empty := &ServerSecretType{}
	o := &ServerSecretType{}

	seed := randomize.NewSeed()
	if err = randomize.Struct(seed, o, serverSecretTypeDBTypes, false); err != nil {
		t.Errorf("Unable to randomize ServerSecretType object: %s", err)
	}

	AddServerSecretTypeHook(boil.BeforeInsertHook, serverSecretTypeBeforeInsertHook)
	if err = o.doBeforeInsertHooks(ctx, nil); err != nil {
		t.Errorf("Unable to execute doBeforeInsertHooks: %s", err)
	}
	if !reflect.DeepEqual(o, empty) {
		t.Errorf("Expected BeforeInsertHook function to empty object, but got: %#v", o)
	}
	serverSecretTypeBeforeInsertHooks = []ServerSecretTypeHook{}

	AddServerSecretTypeHook(boil.AfterInsertHook, serverSecretTypeAfterInsertHook)
	if err = o.doAfterInsertHooks(ctx, nil); err != nil {
		t.Errorf("Unable to execute doAfterInsertHooks: %s", err)
	}
	if !reflect.DeepEqual(o, empty) {
		t.Errorf("Expected AfterInsertHook function to empty object, but got: %#v", o)
	}
	serverSecretTypeAfterInsertHooks = []ServerSecretTypeHook{}

	AddServerSecretTypeHook(boil.AfterSelectHook, serverSecretTypeAfterSelectHook)
	if err = o.doAfterSelectHooks(ctx, nil); err != nil {
		t.Errorf("Unable to execute doAfterSelectHooks: %s", err)
	}
	if !reflect.DeepEqual(o, empty) {
		t.Errorf("Expected AfterSelectHook function to empty object, but got: %#v", o)
	}
	serverSecretTypeAfterSelectHooks = []ServerSecretTypeHook{}

	AddServerSecretTypeHook(boil.BeforeUpdateHook, serverSecretTypeBeforeUpdateHook)
	if err = o.doBeforeUpdateHooks(ctx, nil); err != nil {
		t.Errorf("Unable to execute doBeforeUpdateHooks: %s", err)
	}
	if !reflect.DeepEqual(o, empty) {
		t.Errorf("Expected BeforeUpdateHook function to empty object, but got: %#v", o)
	}
	serverSecretTypeBeforeUpdateHooks = []ServerSecretTypeHook{}

	AddServerSecretTypeHook(boil.AfterUpdateHook, serverSecretTypeAfterUpdateHook)
	if err = o.doAfterUpdateHooks(ctx, nil); err != nil {
		t.Errorf("Unable to execute doAfterUpdateHooks: %s", err)
	}
	if !reflect.DeepEqual(o, empty) {
		t.Errorf("Expected AfterUpdateHook function to empty object, but got: %#v", o)
	}
	serverSecretTypeAfterUpdateHooks = []ServerSecretTypeHook{}

	AddServerSecretTypeHook(boil.BeforeDeleteHook, serverSecretTypeBeforeDeleteHook)
	if err = o.doBeforeDeleteHooks(ctx, nil); err != nil {
		t.Errorf("Unable to execute doBeforeDeleteHooks: %s", err)
	}
	if !reflect.DeepEqual(o, empty) {
		t.Errorf("Expected BeforeDeleteHook function to empty object, but got: %#v", o)
	}
	serverSecretTypeBeforeDeleteHooks = []ServerSecretTypeHook{}

	AddServerSecretTypeHook(boil.AfterDeleteHook, serverSecretTypeAfterDeleteHook)
	if err = o.doAfterDeleteHooks(ctx, nil); err != nil {
		t.Errorf("Unable to execute doAfterDeleteHooks: %s", err)
	}
	if !reflect.DeepEqual(o, empty) {
		t.Errorf("Expected AfterDeleteHook function to empty object, but got: %#v", o)
	}
	serverSecretTypeAfterDeleteHooks = []ServerSecretTypeHook{}

	AddServerSecretTypeHook(boil.BeforeUpsertHook, serverSecretTypeBeforeUpsertHook)
	if err = o.doBeforeUpsertHooks(ctx, nil); err != nil {
		t.Errorf("Unable to execute doBeforeUpsertHooks: %s", err)
	}
	if !reflect.DeepEqual(o, empty) {
		t.Errorf("Expected BeforeUpsertHook function to empty object, but got: %#v", o)
	}
	serverSecretTypeBeforeUpsertHooks = []ServerSecretTypeHook{}

	AddServerSecretTypeHook(boil.AfterUpsertHook, serverSecretTypeAfterUpsertHook)
	if err = o.doAfterUpsertHooks(ctx, nil); err != nil {
		t.Errorf("Unable to execute doAfterUpsertHooks: %s", err)
	}
	if !reflect.DeepEqual(o, empty) {
		t.Errorf("Expected AfterUpsertHook function to empty object, but got: %#v", o)
	}
	serverSecretTypeAfterUpsertHooks = []ServerSecretTypeHook{}
}

func testServerSecretTypesInsert(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &ServerSecretType{}
	if err = randomize.Struct(seed, o, serverSecretTypeDBTypes, true, serverSecretTypeColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize ServerSecretType struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	count, err := ServerSecretTypes().Count(ctx, tx)
	if err != nil {
		t.Error(err)
	}

	if count != 1 {
		t.Error("want one record, got:", count)
	}
}

func testServerSecretTypesInsertWhitelist(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &ServerSecretType{}
	if err = randomize.Struct(seed, o, serverSecretTypeDBTypes, true); err != nil {
		t.Errorf("Unable to randomize ServerSecretType struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Whitelist(serverSecretTypeColumnsWithoutDefault...)); err != nil {
		t.Error(err)
	}

	count, err := ServerSecretTypes().Count(ctx, tx)
	if err != nil {
		t.Error(err)
	}

	if count != 1 {
		t.Error("want one record, got:", count)
	}
}

func testServerSecretTypeToManyServerSecrets(t *testing.T) {
	var err error
	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()

	var a ServerSecretType
	var b, c ServerSecret

	seed := randomize.NewSeed()
	if err = randomize.Struct(seed, &a, serverSecretTypeDBTypes, true, serverSecretTypeColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize ServerSecretType struct: %s", err)
	}

	if err := a.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Fatal(err)
	}

	if err = randomize.Struct(seed, &b, serverSecretDBTypes, false, serverSecretColumnsWithDefault...); err != nil {
		t.Fatal(err)
	}
	if err = randomize.Struct(seed, &c, serverSecretDBTypes, false, serverSecretColumnsWithDefault...); err != nil {
		t.Fatal(err)
	}

	b.ServerSecretTypeID = a.ID
	c.ServerSecretTypeID = a.ID

	if err = b.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Fatal(err)
	}
	if err = c.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Fatal(err)
	}

	check, err := a.ServerSecrets().All(ctx, tx)
	if err != nil {
		t.Fatal(err)
	}

	bFound, cFound := false, false
	for _, v := range check {
		if v.ServerSecretTypeID == b.ServerSecretTypeID {
			bFound = true
		}
		if v.ServerSecretTypeID == c.ServerSecretTypeID {
			cFound = true
		}
	}

	if !bFound {
		t.Error("expected to find b")
	}
	if !cFound {
		t.Error("expected to find c")
	}

	slice := ServerSecretTypeSlice{&a}
	if err = a.L.LoadServerSecrets(ctx, tx, false, (*[]*ServerSecretType)(&slice), nil); err != nil {
		t.Fatal(err)
	}
	if got := len(a.R.ServerSecrets); got != 2 {
		t.Error("number of eager loaded records wrong, got:", got)
	}

	a.R.ServerSecrets = nil
	if err = a.L.LoadServerSecrets(ctx, tx, true, &a, nil); err != nil {
		t.Fatal(err)
	}
	if got := len(a.R.ServerSecrets); got != 2 {
		t.Error("number of eager loaded records wrong, got:", got)
	}

	if t.Failed() {
		t.Logf("%#v", check)
	}
}

func testServerSecretTypeToManyAddOpServerSecrets(t *testing.T) {
	var err error

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()

	var a ServerSecretType
	var b, c, d, e ServerSecret

	seed := randomize.NewSeed()
	if err = randomize.Struct(seed, &a, serverSecretTypeDBTypes, false, strmangle.SetComplement(serverSecretTypePrimaryKeyColumns, serverSecretTypeColumnsWithoutDefault)...); err != nil {
		t.Fatal(err)
	}
	foreigners := []*ServerSecret{&b, &c, &d, &e}
	for _, x := range foreigners {
		if err = randomize.Struct(seed, x, serverSecretDBTypes, false, strmangle.SetComplement(serverSecretPrimaryKeyColumns, serverSecretColumnsWithoutDefault)...); err != nil {
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

	foreignersSplitByInsertion := [][]*ServerSecret{
		{&b, &c},
		{&d, &e},
	}

	for i, x := range foreignersSplitByInsertion {
		err = a.AddServerSecrets(ctx, tx, i != 0, x...)
		if err != nil {
			t.Fatal(err)
		}

		first := x[0]
		second := x[1]

		if a.ID != first.ServerSecretTypeID {
			t.Error("foreign key was wrong value", a.ID, first.ServerSecretTypeID)
		}
		if a.ID != second.ServerSecretTypeID {
			t.Error("foreign key was wrong value", a.ID, second.ServerSecretTypeID)
		}

		if first.R.ServerSecretType != &a {
			t.Error("relationship was not added properly to the foreign slice")
		}
		if second.R.ServerSecretType != &a {
			t.Error("relationship was not added properly to the foreign slice")
		}

		if a.R.ServerSecrets[i*2] != first {
			t.Error("relationship struct slice not set to correct value")
		}
		if a.R.ServerSecrets[i*2+1] != second {
			t.Error("relationship struct slice not set to correct value")
		}

		count, err := a.ServerSecrets().Count(ctx, tx)
		if err != nil {
			t.Fatal(err)
		}
		if want := int64((i + 1) * 2); count != want {
			t.Error("want", want, "got", count)
		}
	}
}

func testServerSecretTypesReload(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &ServerSecretType{}
	if err = randomize.Struct(seed, o, serverSecretTypeDBTypes, true, serverSecretTypeColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize ServerSecretType struct: %s", err)
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

func testServerSecretTypesReloadAll(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &ServerSecretType{}
	if err = randomize.Struct(seed, o, serverSecretTypeDBTypes, true, serverSecretTypeColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize ServerSecretType struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	slice := ServerSecretTypeSlice{o}

	if err = slice.ReloadAll(ctx, tx); err != nil {
		t.Error(err)
	}
}

func testServerSecretTypesSelect(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &ServerSecretType{}
	if err = randomize.Struct(seed, o, serverSecretTypeDBTypes, true, serverSecretTypeColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize ServerSecretType struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	slice, err := ServerSecretTypes().All(ctx, tx)
	if err != nil {
		t.Error(err)
	}

	if len(slice) != 1 {
		t.Error("want one record, got:", len(slice))
	}
}

var (
	serverSecretTypeDBTypes = map[string]string{`ID`: `uuid`, `Name`: `string`, `Slug`: `string`, `Builtin`: `bool`, `CreatedAt`: `timestamptz`, `UpdatedAt`: `timestamptz`}
	_                       = bytes.MinRead
)

func testServerSecretTypesUpdate(t *testing.T) {
	t.Parallel()

	if 0 == len(serverSecretTypePrimaryKeyColumns) {
		t.Skip("Skipping table with no primary key columns")
	}
	if len(serverSecretTypeAllColumns) == len(serverSecretTypePrimaryKeyColumns) {
		t.Skip("Skipping table with only primary key columns")
	}

	seed := randomize.NewSeed()
	var err error
	o := &ServerSecretType{}
	if err = randomize.Struct(seed, o, serverSecretTypeDBTypes, true, serverSecretTypeColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize ServerSecretType struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	count, err := ServerSecretTypes().Count(ctx, tx)
	if err != nil {
		t.Error(err)
	}

	if count != 1 {
		t.Error("want one record, got:", count)
	}

	if err = randomize.Struct(seed, o, serverSecretTypeDBTypes, true, serverSecretTypePrimaryKeyColumns...); err != nil {
		t.Errorf("Unable to randomize ServerSecretType struct: %s", err)
	}

	if rowsAff, err := o.Update(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	} else if rowsAff != 1 {
		t.Error("should only affect one row but affected", rowsAff)
	}
}

func testServerSecretTypesSliceUpdateAll(t *testing.T) {
	t.Parallel()

	if len(serverSecretTypeAllColumns) == len(serverSecretTypePrimaryKeyColumns) {
		t.Skip("Skipping table with only primary key columns")
	}

	seed := randomize.NewSeed()
	var err error
	o := &ServerSecretType{}
	if err = randomize.Struct(seed, o, serverSecretTypeDBTypes, true, serverSecretTypeColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize ServerSecretType struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	count, err := ServerSecretTypes().Count(ctx, tx)
	if err != nil {
		t.Error(err)
	}

	if count != 1 {
		t.Error("want one record, got:", count)
	}

	if err = randomize.Struct(seed, o, serverSecretTypeDBTypes, true, serverSecretTypePrimaryKeyColumns...); err != nil {
		t.Errorf("Unable to randomize ServerSecretType struct: %s", err)
	}

	// Remove Primary keys and unique columns from what we plan to update
	var fields []string
	if strmangle.StringSliceMatch(serverSecretTypeAllColumns, serverSecretTypePrimaryKeyColumns) {
		fields = serverSecretTypeAllColumns
	} else {
		fields = strmangle.SetComplement(
			serverSecretTypeAllColumns,
			serverSecretTypePrimaryKeyColumns,
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

	slice := ServerSecretTypeSlice{o}
	if rowsAff, err := slice.UpdateAll(ctx, tx, updateMap); err != nil {
		t.Error(err)
	} else if rowsAff != 1 {
		t.Error("wanted one record updated but got", rowsAff)
	}
}
