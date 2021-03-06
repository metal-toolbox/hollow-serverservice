// Code generated by SQLBoiler 4.11.0 (https://github.com/volatiletech/sqlboiler). DO NOT EDIT.
// This file is meant to be re-generated in place and/or deleted at any time.

package models

import (
	"context"
	"database/sql"
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/friendsofgo/errors"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/queries"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
	"github.com/volatiletech/sqlboiler/v4/queries/qmhelper"
	"github.com/volatiletech/strmangle"
)

// ServerSecretType is an object representing the database table.
type ServerSecretType struct {
	ID        string    `boil:"id" json:"id" toml:"id" yaml:"id"`
	Name      string    `boil:"name" json:"name" toml:"name" yaml:"name"`
	Slug      string    `boil:"slug" json:"slug" toml:"slug" yaml:"slug"`
	Builtin   bool      `boil:"builtin" json:"builtin" toml:"builtin" yaml:"builtin"`
	CreatedAt time.Time `boil:"created_at" json:"created_at" toml:"created_at" yaml:"created_at"`
	UpdatedAt time.Time `boil:"updated_at" json:"updated_at" toml:"updated_at" yaml:"updated_at"`

	R *serverSecretTypeR `boil:"-" json:"-" toml:"-" yaml:"-"`
	L serverSecretTypeL  `boil:"-" json:"-" toml:"-" yaml:"-"`
}

var ServerSecretTypeColumns = struct {
	ID        string
	Name      string
	Slug      string
	Builtin   string
	CreatedAt string
	UpdatedAt string
}{
	ID:        "id",
	Name:      "name",
	Slug:      "slug",
	Builtin:   "builtin",
	CreatedAt: "created_at",
	UpdatedAt: "updated_at",
}

var ServerSecretTypeTableColumns = struct {
	ID        string
	Name      string
	Slug      string
	Builtin   string
	CreatedAt string
	UpdatedAt string
}{
	ID:        "server_secret_types.id",
	Name:      "server_secret_types.name",
	Slug:      "server_secret_types.slug",
	Builtin:   "server_secret_types.builtin",
	CreatedAt: "server_secret_types.created_at",
	UpdatedAt: "server_secret_types.updated_at",
}

// Generated where

type whereHelperbool struct{ field string }

func (w whereHelperbool) EQ(x bool) qm.QueryMod  { return qmhelper.Where(w.field, qmhelper.EQ, x) }
func (w whereHelperbool) NEQ(x bool) qm.QueryMod { return qmhelper.Where(w.field, qmhelper.NEQ, x) }
func (w whereHelperbool) LT(x bool) qm.QueryMod  { return qmhelper.Where(w.field, qmhelper.LT, x) }
func (w whereHelperbool) LTE(x bool) qm.QueryMod { return qmhelper.Where(w.field, qmhelper.LTE, x) }
func (w whereHelperbool) GT(x bool) qm.QueryMod  { return qmhelper.Where(w.field, qmhelper.GT, x) }
func (w whereHelperbool) GTE(x bool) qm.QueryMod { return qmhelper.Where(w.field, qmhelper.GTE, x) }

type whereHelpertime_Time struct{ field string }

func (w whereHelpertime_Time) EQ(x time.Time) qm.QueryMod {
	return qmhelper.Where(w.field, qmhelper.EQ, x)
}
func (w whereHelpertime_Time) NEQ(x time.Time) qm.QueryMod {
	return qmhelper.Where(w.field, qmhelper.NEQ, x)
}
func (w whereHelpertime_Time) LT(x time.Time) qm.QueryMod {
	return qmhelper.Where(w.field, qmhelper.LT, x)
}
func (w whereHelpertime_Time) LTE(x time.Time) qm.QueryMod {
	return qmhelper.Where(w.field, qmhelper.LTE, x)
}
func (w whereHelpertime_Time) GT(x time.Time) qm.QueryMod {
	return qmhelper.Where(w.field, qmhelper.GT, x)
}
func (w whereHelpertime_Time) GTE(x time.Time) qm.QueryMod {
	return qmhelper.Where(w.field, qmhelper.GTE, x)
}

var ServerSecretTypeWhere = struct {
	ID        whereHelperstring
	Name      whereHelperstring
	Slug      whereHelperstring
	Builtin   whereHelperbool
	CreatedAt whereHelpertime_Time
	UpdatedAt whereHelpertime_Time
}{
	ID:        whereHelperstring{field: "\"server_secret_types\".\"id\""},
	Name:      whereHelperstring{field: "\"server_secret_types\".\"name\""},
	Slug:      whereHelperstring{field: "\"server_secret_types\".\"slug\""},
	Builtin:   whereHelperbool{field: "\"server_secret_types\".\"builtin\""},
	CreatedAt: whereHelpertime_Time{field: "\"server_secret_types\".\"created_at\""},
	UpdatedAt: whereHelpertime_Time{field: "\"server_secret_types\".\"updated_at\""},
}

// ServerSecretTypeRels is where relationship names are stored.
var ServerSecretTypeRels = struct {
	ServerSecrets string
}{
	ServerSecrets: "ServerSecrets",
}

// serverSecretTypeR is where relationships are stored.
type serverSecretTypeR struct {
	ServerSecrets ServerSecretSlice `boil:"ServerSecrets" json:"ServerSecrets" toml:"ServerSecrets" yaml:"ServerSecrets"`
}

// NewStruct creates a new relationship struct
func (*serverSecretTypeR) NewStruct() *serverSecretTypeR {
	return &serverSecretTypeR{}
}

func (r *serverSecretTypeR) GetServerSecrets() ServerSecretSlice {
	if r == nil {
		return nil
	}
	return r.ServerSecrets
}

// serverSecretTypeL is where Load methods for each relationship are stored.
type serverSecretTypeL struct{}

var (
	serverSecretTypeAllColumns            = []string{"id", "name", "slug", "builtin", "created_at", "updated_at"}
	serverSecretTypeColumnsWithoutDefault = []string{"name", "slug", "created_at", "updated_at"}
	serverSecretTypeColumnsWithDefault    = []string{"id", "builtin"}
	serverSecretTypePrimaryKeyColumns     = []string{"id"}
	serverSecretTypeGeneratedColumns      = []string{}
)

type (
	// ServerSecretTypeSlice is an alias for a slice of pointers to ServerSecretType.
	// This should almost always be used instead of []ServerSecretType.
	ServerSecretTypeSlice []*ServerSecretType
	// ServerSecretTypeHook is the signature for custom ServerSecretType hook methods
	ServerSecretTypeHook func(context.Context, boil.ContextExecutor, *ServerSecretType) error

	serverSecretTypeQuery struct {
		*queries.Query
	}
)

// Cache for insert, update and upsert
var (
	serverSecretTypeType                 = reflect.TypeOf(&ServerSecretType{})
	serverSecretTypeMapping              = queries.MakeStructMapping(serverSecretTypeType)
	serverSecretTypePrimaryKeyMapping, _ = queries.BindMapping(serverSecretTypeType, serverSecretTypeMapping, serverSecretTypePrimaryKeyColumns)
	serverSecretTypeInsertCacheMut       sync.RWMutex
	serverSecretTypeInsertCache          = make(map[string]insertCache)
	serverSecretTypeUpdateCacheMut       sync.RWMutex
	serverSecretTypeUpdateCache          = make(map[string]updateCache)
	serverSecretTypeUpsertCacheMut       sync.RWMutex
	serverSecretTypeUpsertCache          = make(map[string]insertCache)
)

var (
	// Force time package dependency for automated UpdatedAt/CreatedAt.
	_ = time.Second
	// Force qmhelper dependency for where clause generation (which doesn't
	// always happen)
	_ = qmhelper.Where
)

var serverSecretTypeAfterSelectHooks []ServerSecretTypeHook

var serverSecretTypeBeforeInsertHooks []ServerSecretTypeHook
var serverSecretTypeAfterInsertHooks []ServerSecretTypeHook

var serverSecretTypeBeforeUpdateHooks []ServerSecretTypeHook
var serverSecretTypeAfterUpdateHooks []ServerSecretTypeHook

var serverSecretTypeBeforeDeleteHooks []ServerSecretTypeHook
var serverSecretTypeAfterDeleteHooks []ServerSecretTypeHook

var serverSecretTypeBeforeUpsertHooks []ServerSecretTypeHook
var serverSecretTypeAfterUpsertHooks []ServerSecretTypeHook

// doAfterSelectHooks executes all "after Select" hooks.
func (o *ServerSecretType) doAfterSelectHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range serverSecretTypeAfterSelectHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doBeforeInsertHooks executes all "before insert" hooks.
func (o *ServerSecretType) doBeforeInsertHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range serverSecretTypeBeforeInsertHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doAfterInsertHooks executes all "after Insert" hooks.
func (o *ServerSecretType) doAfterInsertHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range serverSecretTypeAfterInsertHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doBeforeUpdateHooks executes all "before Update" hooks.
func (o *ServerSecretType) doBeforeUpdateHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range serverSecretTypeBeforeUpdateHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doAfterUpdateHooks executes all "after Update" hooks.
func (o *ServerSecretType) doAfterUpdateHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range serverSecretTypeAfterUpdateHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doBeforeDeleteHooks executes all "before Delete" hooks.
func (o *ServerSecretType) doBeforeDeleteHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range serverSecretTypeBeforeDeleteHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doAfterDeleteHooks executes all "after Delete" hooks.
func (o *ServerSecretType) doAfterDeleteHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range serverSecretTypeAfterDeleteHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doBeforeUpsertHooks executes all "before Upsert" hooks.
func (o *ServerSecretType) doBeforeUpsertHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range serverSecretTypeBeforeUpsertHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doAfterUpsertHooks executes all "after Upsert" hooks.
func (o *ServerSecretType) doAfterUpsertHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range serverSecretTypeAfterUpsertHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// AddServerSecretTypeHook registers your hook function for all future operations.
func AddServerSecretTypeHook(hookPoint boil.HookPoint, serverSecretTypeHook ServerSecretTypeHook) {
	switch hookPoint {
	case boil.AfterSelectHook:
		serverSecretTypeAfterSelectHooks = append(serverSecretTypeAfterSelectHooks, serverSecretTypeHook)
	case boil.BeforeInsertHook:
		serverSecretTypeBeforeInsertHooks = append(serverSecretTypeBeforeInsertHooks, serverSecretTypeHook)
	case boil.AfterInsertHook:
		serverSecretTypeAfterInsertHooks = append(serverSecretTypeAfterInsertHooks, serverSecretTypeHook)
	case boil.BeforeUpdateHook:
		serverSecretTypeBeforeUpdateHooks = append(serverSecretTypeBeforeUpdateHooks, serverSecretTypeHook)
	case boil.AfterUpdateHook:
		serverSecretTypeAfterUpdateHooks = append(serverSecretTypeAfterUpdateHooks, serverSecretTypeHook)
	case boil.BeforeDeleteHook:
		serverSecretTypeBeforeDeleteHooks = append(serverSecretTypeBeforeDeleteHooks, serverSecretTypeHook)
	case boil.AfterDeleteHook:
		serverSecretTypeAfterDeleteHooks = append(serverSecretTypeAfterDeleteHooks, serverSecretTypeHook)
	case boil.BeforeUpsertHook:
		serverSecretTypeBeforeUpsertHooks = append(serverSecretTypeBeforeUpsertHooks, serverSecretTypeHook)
	case boil.AfterUpsertHook:
		serverSecretTypeAfterUpsertHooks = append(serverSecretTypeAfterUpsertHooks, serverSecretTypeHook)
	}
}

// One returns a single serverSecretType record from the query.
func (q serverSecretTypeQuery) One(ctx context.Context, exec boil.ContextExecutor) (*ServerSecretType, error) {
	o := &ServerSecretType{}

	queries.SetLimit(q.Query, 1)

	err := q.Bind(ctx, exec, o)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, sql.ErrNoRows
		}
		return nil, errors.Wrap(err, "models: failed to execute a one query for server_secret_types")
	}

	if err := o.doAfterSelectHooks(ctx, exec); err != nil {
		return o, err
	}

	return o, nil
}

// All returns all ServerSecretType records from the query.
func (q serverSecretTypeQuery) All(ctx context.Context, exec boil.ContextExecutor) (ServerSecretTypeSlice, error) {
	var o []*ServerSecretType

	err := q.Bind(ctx, exec, &o)
	if err != nil {
		return nil, errors.Wrap(err, "models: failed to assign all query results to ServerSecretType slice")
	}

	if len(serverSecretTypeAfterSelectHooks) != 0 {
		for _, obj := range o {
			if err := obj.doAfterSelectHooks(ctx, exec); err != nil {
				return o, err
			}
		}
	}

	return o, nil
}

// Count returns the count of all ServerSecretType records in the query.
func (q serverSecretTypeQuery) Count(ctx context.Context, exec boil.ContextExecutor) (int64, error) {
	var count int64

	queries.SetSelect(q.Query, nil)
	queries.SetCount(q.Query)

	err := q.Query.QueryRowContext(ctx, exec).Scan(&count)
	if err != nil {
		return 0, errors.Wrap(err, "models: failed to count server_secret_types rows")
	}

	return count, nil
}

// Exists checks if the row exists in the table.
func (q serverSecretTypeQuery) Exists(ctx context.Context, exec boil.ContextExecutor) (bool, error) {
	var count int64

	queries.SetSelect(q.Query, nil)
	queries.SetCount(q.Query)
	queries.SetLimit(q.Query, 1)

	err := q.Query.QueryRowContext(ctx, exec).Scan(&count)
	if err != nil {
		return false, errors.Wrap(err, "models: failed to check if server_secret_types exists")
	}

	return count > 0, nil
}

// ServerSecrets retrieves all the server_secret's ServerSecrets with an executor.
func (o *ServerSecretType) ServerSecrets(mods ...qm.QueryMod) serverSecretQuery {
	var queryMods []qm.QueryMod
	if len(mods) != 0 {
		queryMods = append(queryMods, mods...)
	}

	queryMods = append(queryMods,
		qm.Where("\"server_secrets\".\"server_secret_type_id\"=?", o.ID),
	)

	return ServerSecrets(queryMods...)
}

// LoadServerSecrets allows an eager lookup of values, cached into the
// loaded structs of the objects. This is for a 1-M or N-M relationship.
func (serverSecretTypeL) LoadServerSecrets(ctx context.Context, e boil.ContextExecutor, singular bool, maybeServerSecretType interface{}, mods queries.Applicator) error {
	var slice []*ServerSecretType
	var object *ServerSecretType

	if singular {
		object = maybeServerSecretType.(*ServerSecretType)
	} else {
		slice = *maybeServerSecretType.(*[]*ServerSecretType)
	}

	args := make([]interface{}, 0, 1)
	if singular {
		if object.R == nil {
			object.R = &serverSecretTypeR{}
		}
		args = append(args, object.ID)
	} else {
	Outer:
		for _, obj := range slice {
			if obj.R == nil {
				obj.R = &serverSecretTypeR{}
			}

			for _, a := range args {
				if a == obj.ID {
					continue Outer
				}
			}

			args = append(args, obj.ID)
		}
	}

	if len(args) == 0 {
		return nil
	}

	query := NewQuery(
		qm.From(`server_secrets`),
		qm.WhereIn(`server_secrets.server_secret_type_id in ?`, args...),
	)
	if mods != nil {
		mods.Apply(query)
	}

	results, err := query.QueryContext(ctx, e)
	if err != nil {
		return errors.Wrap(err, "failed to eager load server_secrets")
	}

	var resultSlice []*ServerSecret
	if err = queries.Bind(results, &resultSlice); err != nil {
		return errors.Wrap(err, "failed to bind eager loaded slice server_secrets")
	}

	if err = results.Close(); err != nil {
		return errors.Wrap(err, "failed to close results in eager load on server_secrets")
	}
	if err = results.Err(); err != nil {
		return errors.Wrap(err, "error occurred during iteration of eager loaded relations for server_secrets")
	}

	if len(serverSecretAfterSelectHooks) != 0 {
		for _, obj := range resultSlice {
			if err := obj.doAfterSelectHooks(ctx, e); err != nil {
				return err
			}
		}
	}
	if singular {
		object.R.ServerSecrets = resultSlice
		for _, foreign := range resultSlice {
			if foreign.R == nil {
				foreign.R = &serverSecretR{}
			}
			foreign.R.ServerSecretType = object
		}
		return nil
	}

	for _, foreign := range resultSlice {
		for _, local := range slice {
			if local.ID == foreign.ServerSecretTypeID {
				local.R.ServerSecrets = append(local.R.ServerSecrets, foreign)
				if foreign.R == nil {
					foreign.R = &serverSecretR{}
				}
				foreign.R.ServerSecretType = local
				break
			}
		}
	}

	return nil
}

// AddServerSecrets adds the given related objects to the existing relationships
// of the server_secret_type, optionally inserting them as new records.
// Appends related to o.R.ServerSecrets.
// Sets related.R.ServerSecretType appropriately.
func (o *ServerSecretType) AddServerSecrets(ctx context.Context, exec boil.ContextExecutor, insert bool, related ...*ServerSecret) error {
	var err error
	for _, rel := range related {
		if insert {
			rel.ServerSecretTypeID = o.ID
			if err = rel.Insert(ctx, exec, boil.Infer()); err != nil {
				return errors.Wrap(err, "failed to insert into foreign table")
			}
		} else {
			updateQuery := fmt.Sprintf(
				"UPDATE \"server_secrets\" SET %s WHERE %s",
				strmangle.SetParamNames("\"", "\"", 1, []string{"server_secret_type_id"}),
				strmangle.WhereClause("\"", "\"", 2, serverSecretPrimaryKeyColumns),
			)
			values := []interface{}{o.ID, rel.ID}

			if boil.IsDebug(ctx) {
				writer := boil.DebugWriterFrom(ctx)
				fmt.Fprintln(writer, updateQuery)
				fmt.Fprintln(writer, values)
			}
			if _, err = exec.ExecContext(ctx, updateQuery, values...); err != nil {
				return errors.Wrap(err, "failed to update foreign table")
			}

			rel.ServerSecretTypeID = o.ID
		}
	}

	if o.R == nil {
		o.R = &serverSecretTypeR{
			ServerSecrets: related,
		}
	} else {
		o.R.ServerSecrets = append(o.R.ServerSecrets, related...)
	}

	for _, rel := range related {
		if rel.R == nil {
			rel.R = &serverSecretR{
				ServerSecretType: o,
			}
		} else {
			rel.R.ServerSecretType = o
		}
	}
	return nil
}

// ServerSecretTypes retrieves all the records using an executor.
func ServerSecretTypes(mods ...qm.QueryMod) serverSecretTypeQuery {
	mods = append(mods, qm.From("\"server_secret_types\""))
	q := NewQuery(mods...)
	if len(queries.GetSelect(q)) == 0 {
		queries.SetSelect(q, []string{"\"server_secret_types\".*"})
	}

	return serverSecretTypeQuery{q}
}

// FindServerSecretType retrieves a single record by ID with an executor.
// If selectCols is empty Find will return all columns.
func FindServerSecretType(ctx context.Context, exec boil.ContextExecutor, iD string, selectCols ...string) (*ServerSecretType, error) {
	serverSecretTypeObj := &ServerSecretType{}

	sel := "*"
	if len(selectCols) > 0 {
		sel = strings.Join(strmangle.IdentQuoteSlice(dialect.LQ, dialect.RQ, selectCols), ",")
	}
	query := fmt.Sprintf(
		"select %s from \"server_secret_types\" where \"id\"=$1", sel,
	)

	q := queries.Raw(query, iD)

	err := q.Bind(ctx, exec, serverSecretTypeObj)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, sql.ErrNoRows
		}
		return nil, errors.Wrap(err, "models: unable to select from server_secret_types")
	}

	if err = serverSecretTypeObj.doAfterSelectHooks(ctx, exec); err != nil {
		return serverSecretTypeObj, err
	}

	return serverSecretTypeObj, nil
}

// Insert a single record using an executor.
// See boil.Columns.InsertColumnSet documentation to understand column list inference for inserts.
func (o *ServerSecretType) Insert(ctx context.Context, exec boil.ContextExecutor, columns boil.Columns) error {
	if o == nil {
		return errors.New("models: no server_secret_types provided for insertion")
	}

	var err error
	if !boil.TimestampsAreSkipped(ctx) {
		currTime := time.Now().In(boil.GetLocation())

		if o.CreatedAt.IsZero() {
			o.CreatedAt = currTime
		}
		if o.UpdatedAt.IsZero() {
			o.UpdatedAt = currTime
		}
	}

	if err := o.doBeforeInsertHooks(ctx, exec); err != nil {
		return err
	}

	nzDefaults := queries.NonZeroDefaultSet(serverSecretTypeColumnsWithDefault, o)

	key := makeCacheKey(columns, nzDefaults)
	serverSecretTypeInsertCacheMut.RLock()
	cache, cached := serverSecretTypeInsertCache[key]
	serverSecretTypeInsertCacheMut.RUnlock()

	if !cached {
		wl, returnColumns := columns.InsertColumnSet(
			serverSecretTypeAllColumns,
			serverSecretTypeColumnsWithDefault,
			serverSecretTypeColumnsWithoutDefault,
			nzDefaults,
		)

		cache.valueMapping, err = queries.BindMapping(serverSecretTypeType, serverSecretTypeMapping, wl)
		if err != nil {
			return err
		}
		cache.retMapping, err = queries.BindMapping(serverSecretTypeType, serverSecretTypeMapping, returnColumns)
		if err != nil {
			return err
		}
		if len(wl) != 0 {
			cache.query = fmt.Sprintf("INSERT INTO \"server_secret_types\" (\"%s\") %%sVALUES (%s)%%s", strings.Join(wl, "\",\""), strmangle.Placeholders(dialect.UseIndexPlaceholders, len(wl), 1, 1))
		} else {
			cache.query = "INSERT INTO \"server_secret_types\" %sDEFAULT VALUES%s"
		}

		var queryOutput, queryReturning string

		if len(cache.retMapping) != 0 {
			queryReturning = fmt.Sprintf(" RETURNING \"%s\"", strings.Join(returnColumns, "\",\""))
		}

		cache.query = fmt.Sprintf(cache.query, queryOutput, queryReturning)
	}

	value := reflect.Indirect(reflect.ValueOf(o))
	vals := queries.ValuesFromMapping(value, cache.valueMapping)

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, cache.query)
		fmt.Fprintln(writer, vals)
	}

	if len(cache.retMapping) != 0 {
		err = exec.QueryRowContext(ctx, cache.query, vals...).Scan(queries.PtrsFromMapping(value, cache.retMapping)...)
	} else {
		_, err = exec.ExecContext(ctx, cache.query, vals...)
	}

	if err != nil {
		return errors.Wrap(err, "models: unable to insert into server_secret_types")
	}

	if !cached {
		serverSecretTypeInsertCacheMut.Lock()
		serverSecretTypeInsertCache[key] = cache
		serverSecretTypeInsertCacheMut.Unlock()
	}

	return o.doAfterInsertHooks(ctx, exec)
}

// Update uses an executor to update the ServerSecretType.
// See boil.Columns.UpdateColumnSet documentation to understand column list inference for updates.
// Update does not automatically update the record in case of default values. Use .Reload() to refresh the records.
func (o *ServerSecretType) Update(ctx context.Context, exec boil.ContextExecutor, columns boil.Columns) (int64, error) {
	if !boil.TimestampsAreSkipped(ctx) {
		currTime := time.Now().In(boil.GetLocation())

		o.UpdatedAt = currTime
	}

	var err error
	if err = o.doBeforeUpdateHooks(ctx, exec); err != nil {
		return 0, err
	}
	key := makeCacheKey(columns, nil)
	serverSecretTypeUpdateCacheMut.RLock()
	cache, cached := serverSecretTypeUpdateCache[key]
	serverSecretTypeUpdateCacheMut.RUnlock()

	if !cached {
		wl := columns.UpdateColumnSet(
			serverSecretTypeAllColumns,
			serverSecretTypePrimaryKeyColumns,
		)

		if !columns.IsWhitelist() {
			wl = strmangle.SetComplement(wl, []string{"created_at"})
		}
		if len(wl) == 0 {
			return 0, errors.New("models: unable to update server_secret_types, could not build whitelist")
		}

		cache.query = fmt.Sprintf("UPDATE \"server_secret_types\" SET %s WHERE %s",
			strmangle.SetParamNames("\"", "\"", 1, wl),
			strmangle.WhereClause("\"", "\"", len(wl)+1, serverSecretTypePrimaryKeyColumns),
		)
		cache.valueMapping, err = queries.BindMapping(serverSecretTypeType, serverSecretTypeMapping, append(wl, serverSecretTypePrimaryKeyColumns...))
		if err != nil {
			return 0, err
		}
	}

	values := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(o)), cache.valueMapping)

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, cache.query)
		fmt.Fprintln(writer, values)
	}
	var result sql.Result
	result, err = exec.ExecContext(ctx, cache.query, values...)
	if err != nil {
		return 0, errors.Wrap(err, "models: unable to update server_secret_types row")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "models: failed to get rows affected by update for server_secret_types")
	}

	if !cached {
		serverSecretTypeUpdateCacheMut.Lock()
		serverSecretTypeUpdateCache[key] = cache
		serverSecretTypeUpdateCacheMut.Unlock()
	}

	return rowsAff, o.doAfterUpdateHooks(ctx, exec)
}

// UpdateAll updates all rows with the specified column values.
func (q serverSecretTypeQuery) UpdateAll(ctx context.Context, exec boil.ContextExecutor, cols M) (int64, error) {
	queries.SetUpdate(q.Query, cols)

	result, err := q.Query.ExecContext(ctx, exec)
	if err != nil {
		return 0, errors.Wrap(err, "models: unable to update all for server_secret_types")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "models: unable to retrieve rows affected for server_secret_types")
	}

	return rowsAff, nil
}

// UpdateAll updates all rows with the specified column values, using an executor.
func (o ServerSecretTypeSlice) UpdateAll(ctx context.Context, exec boil.ContextExecutor, cols M) (int64, error) {
	ln := int64(len(o))
	if ln == 0 {
		return 0, nil
	}

	if len(cols) == 0 {
		return 0, errors.New("models: update all requires at least one column argument")
	}

	colNames := make([]string, len(cols))
	args := make([]interface{}, len(cols))

	i := 0
	for name, value := range cols {
		colNames[i] = name
		args[i] = value
		i++
	}

	// Append all of the primary key values for each column
	for _, obj := range o {
		pkeyArgs := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(obj)), serverSecretTypePrimaryKeyMapping)
		args = append(args, pkeyArgs...)
	}

	sql := fmt.Sprintf("UPDATE \"server_secret_types\" SET %s WHERE %s",
		strmangle.SetParamNames("\"", "\"", 1, colNames),
		strmangle.WhereClauseRepeated(string(dialect.LQ), string(dialect.RQ), len(colNames)+1, serverSecretTypePrimaryKeyColumns, len(o)))

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, sql)
		fmt.Fprintln(writer, args...)
	}
	result, err := exec.ExecContext(ctx, sql, args...)
	if err != nil {
		return 0, errors.Wrap(err, "models: unable to update all in serverSecretType slice")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "models: unable to retrieve rows affected all in update all serverSecretType")
	}
	return rowsAff, nil
}

// Delete deletes a single ServerSecretType record with an executor.
// Delete will match against the primary key column to find the record to delete.
func (o *ServerSecretType) Delete(ctx context.Context, exec boil.ContextExecutor) (int64, error) {
	if o == nil {
		return 0, errors.New("models: no ServerSecretType provided for delete")
	}

	if err := o.doBeforeDeleteHooks(ctx, exec); err != nil {
		return 0, err
	}

	args := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(o)), serverSecretTypePrimaryKeyMapping)
	sql := "DELETE FROM \"server_secret_types\" WHERE \"id\"=$1"

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, sql)
		fmt.Fprintln(writer, args...)
	}
	result, err := exec.ExecContext(ctx, sql, args...)
	if err != nil {
		return 0, errors.Wrap(err, "models: unable to delete from server_secret_types")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "models: failed to get rows affected by delete for server_secret_types")
	}

	if err := o.doAfterDeleteHooks(ctx, exec); err != nil {
		return 0, err
	}

	return rowsAff, nil
}

// DeleteAll deletes all matching rows.
func (q serverSecretTypeQuery) DeleteAll(ctx context.Context, exec boil.ContextExecutor) (int64, error) {
	if q.Query == nil {
		return 0, errors.New("models: no serverSecretTypeQuery provided for delete all")
	}

	queries.SetDelete(q.Query)

	result, err := q.Query.ExecContext(ctx, exec)
	if err != nil {
		return 0, errors.Wrap(err, "models: unable to delete all from server_secret_types")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "models: failed to get rows affected by deleteall for server_secret_types")
	}

	return rowsAff, nil
}

// DeleteAll deletes all rows in the slice, using an executor.
func (o ServerSecretTypeSlice) DeleteAll(ctx context.Context, exec boil.ContextExecutor) (int64, error) {
	if len(o) == 0 {
		return 0, nil
	}

	if len(serverSecretTypeBeforeDeleteHooks) != 0 {
		for _, obj := range o {
			if err := obj.doBeforeDeleteHooks(ctx, exec); err != nil {
				return 0, err
			}
		}
	}

	var args []interface{}
	for _, obj := range o {
		pkeyArgs := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(obj)), serverSecretTypePrimaryKeyMapping)
		args = append(args, pkeyArgs...)
	}

	sql := "DELETE FROM \"server_secret_types\" WHERE " +
		strmangle.WhereClauseRepeated(string(dialect.LQ), string(dialect.RQ), 1, serverSecretTypePrimaryKeyColumns, len(o))

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, sql)
		fmt.Fprintln(writer, args)
	}
	result, err := exec.ExecContext(ctx, sql, args...)
	if err != nil {
		return 0, errors.Wrap(err, "models: unable to delete all from serverSecretType slice")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "models: failed to get rows affected by deleteall for server_secret_types")
	}

	if len(serverSecretTypeAfterDeleteHooks) != 0 {
		for _, obj := range o {
			if err := obj.doAfterDeleteHooks(ctx, exec); err != nil {
				return 0, err
			}
		}
	}

	return rowsAff, nil
}

// Reload refetches the object from the database
// using the primary keys with an executor.
func (o *ServerSecretType) Reload(ctx context.Context, exec boil.ContextExecutor) error {
	ret, err := FindServerSecretType(ctx, exec, o.ID)
	if err != nil {
		return err
	}

	*o = *ret
	return nil
}

// ReloadAll refetches every row with matching primary key column values
// and overwrites the original object slice with the newly updated slice.
func (o *ServerSecretTypeSlice) ReloadAll(ctx context.Context, exec boil.ContextExecutor) error {
	if o == nil || len(*o) == 0 {
		return nil
	}

	slice := ServerSecretTypeSlice{}
	var args []interface{}
	for _, obj := range *o {
		pkeyArgs := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(obj)), serverSecretTypePrimaryKeyMapping)
		args = append(args, pkeyArgs...)
	}

	sql := "SELECT \"server_secret_types\".* FROM \"server_secret_types\" WHERE " +
		strmangle.WhereClauseRepeated(string(dialect.LQ), string(dialect.RQ), 1, serverSecretTypePrimaryKeyColumns, len(*o))

	q := queries.Raw(sql, args...)

	err := q.Bind(ctx, exec, &slice)
	if err != nil {
		return errors.Wrap(err, "models: unable to reload all in ServerSecretTypeSlice")
	}

	*o = slice

	return nil
}

// ServerSecretTypeExists checks if the ServerSecretType row exists.
func ServerSecretTypeExists(ctx context.Context, exec boil.ContextExecutor, iD string) (bool, error) {
	var exists bool
	sql := "select exists(select 1 from \"server_secret_types\" where \"id\"=$1 limit 1)"

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, sql)
		fmt.Fprintln(writer, iD)
	}
	row := exec.QueryRowContext(ctx, sql, iD)

	err := row.Scan(&exists)
	if err != nil {
		return false, errors.Wrap(err, "models: unable to check if server_secret_types exists")
	}

	return exists, nil
}

// Upsert attempts an insert using an executor, and does an update or ignore on conflict.
// See boil.Columns documentation for how to properly use updateColumns and insertColumns.
func (o *ServerSecretType) Upsert(ctx context.Context, exec boil.ContextExecutor, updateOnConflict bool, conflictColumns []string, updateColumns, insertColumns boil.Columns) error {
	if o == nil {
		return errors.New("models: no server_secret_types provided for upsert")
	}
	if !boil.TimestampsAreSkipped(ctx) {
		currTime := time.Now().In(boil.GetLocation())

		if o.CreatedAt.IsZero() {
			o.CreatedAt = currTime
		}
		o.UpdatedAt = currTime
	}

	if err := o.doBeforeUpsertHooks(ctx, exec); err != nil {
		return err
	}

	nzDefaults := queries.NonZeroDefaultSet(serverSecretTypeColumnsWithDefault, o)

	// Build cache key in-line uglily - mysql vs psql problems
	buf := strmangle.GetBuffer()
	if updateOnConflict {
		buf.WriteByte('t')
	} else {
		buf.WriteByte('f')
	}
	buf.WriteByte('.')
	for _, c := range conflictColumns {
		buf.WriteString(c)
	}
	buf.WriteByte('.')
	buf.WriteString(strconv.Itoa(updateColumns.Kind))
	for _, c := range updateColumns.Cols {
		buf.WriteString(c)
	}
	buf.WriteByte('.')
	buf.WriteString(strconv.Itoa(insertColumns.Kind))
	for _, c := range insertColumns.Cols {
		buf.WriteString(c)
	}
	buf.WriteByte('.')
	for _, c := range nzDefaults {
		buf.WriteString(c)
	}
	key := buf.String()
	strmangle.PutBuffer(buf)

	serverSecretTypeUpsertCacheMut.RLock()
	cache, cached := serverSecretTypeUpsertCache[key]
	serverSecretTypeUpsertCacheMut.RUnlock()

	var err error

	if !cached {
		insert, ret := insertColumns.InsertColumnSet(
			serverSecretTypeAllColumns,
			serverSecretTypeColumnsWithDefault,
			serverSecretTypeColumnsWithoutDefault,
			nzDefaults,
		)
		update := updateColumns.UpdateColumnSet(
			serverSecretTypeAllColumns,
			serverSecretTypePrimaryKeyColumns,
		)

		if updateOnConflict && len(update) == 0 {
			return errors.New("models: unable to upsert server_secret_types, could not build update column list")
		}

		conflict := conflictColumns
		if len(conflict) == 0 {
			conflict = make([]string, len(serverSecretTypePrimaryKeyColumns))
			copy(conflict, serverSecretTypePrimaryKeyColumns)
		}
		cache.query = buildUpsertQueryCockroachDB(dialect, "\"server_secret_types\"", updateOnConflict, ret, update, conflict, insert)

		cache.valueMapping, err = queries.BindMapping(serverSecretTypeType, serverSecretTypeMapping, insert)
		if err != nil {
			return err
		}
		if len(ret) != 0 {
			cache.retMapping, err = queries.BindMapping(serverSecretTypeType, serverSecretTypeMapping, ret)
			if err != nil {
				return err
			}
		}
	}

	value := reflect.Indirect(reflect.ValueOf(o))
	vals := queries.ValuesFromMapping(value, cache.valueMapping)
	var returns []interface{}
	if len(cache.retMapping) != 0 {
		returns = queries.PtrsFromMapping(value, cache.retMapping)
	}

	if boil.DebugMode {
		_, _ = fmt.Fprintln(boil.DebugWriter, cache.query)
		_, _ = fmt.Fprintln(boil.DebugWriter, vals)
	}

	if len(cache.retMapping) != 0 {
		err = exec.QueryRowContext(ctx, cache.query, vals...).Scan(returns...)
		if err == sql.ErrNoRows {
			err = nil // CockcorachDB doesn't return anything when there's no update
		}
	} else {
		_, err = exec.ExecContext(ctx, cache.query, vals...)
	}
	if err != nil {
		return errors.Wrap(err, "models: unable to upsert server_secret_types")
	}

	if !cached {
		serverSecretTypeUpsertCacheMut.Lock()
		serverSecretTypeUpsertCache[key] = cache
		serverSecretTypeUpsertCacheMut.Unlock()
	}

	return o.doAfterUpsertHooks(ctx, exec)
}
