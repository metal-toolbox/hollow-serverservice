package db

import (
	"errors"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// Server represents a server in a facility. These are the details of the
// physical server that is located in the facility.
type Server struct {
	ID                  uuid.UUID
	CreatedAt           time.Time
	UpdatedAt           time.Time
	Name                string
	FacilityCode        string
	Attributes          []Attributes
	ServerComponents    []ServerComponent
	VersionedAttributes []VersionedAttributes
}

// BeforeSave ensures that the Server passes validation checks
func (s *Server) BeforeSave(tx *gorm.DB) (err error) {
	if s.ID.String() == uuid.Nil.String() {
		s.ID = uuid.New()
	}

	return nil
}

func serverPreload(db *gorm.DB) *gorm.DB {
	d := db.Preload("VersionedAttributes",
		"(created_at, namespace, server_id) IN (?)",
		db.Table("versioned_attributes").Select("max(created_at), namespace, server_id").Group("namespace").Group("server_id"),
	)

	return d.Preload("ServerComponents.ServerComponentType").Preload("ServerComponents.Attributes").Preload(clause.Associations)
}

// CreateServer will persist a server into the backend datastore
func (s *Store) CreateServer(srv *Server) error {
	return s.db.Create(srv).Error
}

// DeleteServer will delete a server from the datastore.
func (s *Store) DeleteServer(srv *Server) error {
	return s.db.Delete(srv).Error
}

// GetServers will return a list of servers with the requested params, if no
// filter is passed then it will return all servers
func (s *Store) GetServers(filter *ServerFilter, pager *Pagination) ([]Server, int64, error) {
	var (
		srvs  []Server
		count int64
	)

	d := serverPreload(s.db)

	if filter != nil {
		d = filter.apply(d)
	}

	if pager == nil {
		pager = &Pagination{}
	}

	if err := d.Scopes(paginate(*pager)).Find(&srvs).Offset(-1).Limit(-1).Count(&count).Error; err != nil {
		return nil, 0, err
	}

	return srvs, count, nil
}

// FindServerByUUID will return an existing server if one already exists for the
//  given UUID.
func (s *Store) FindServerByUUID(srvUUID uuid.UUID) (*Server, error) {
	var srv Server

	err := serverPreload(s.db).First(&srv, Server{ID: srvUUID}).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrNotFound
		}

		return nil, err
	}

	return &srv, nil
}

// FindOrCreateServerByUUID will return an existing server if one already exists
//  for the given UUID, if one doesn't exist a new one will be created
func (s *Store) FindOrCreateServerByUUID(srvUUID uuid.UUID) (*Server, error) {
	var srv Server

	err := serverPreload(s.db).FirstOrCreate(&srv, Server{ID: srvUUID}).Error
	if err != nil {
		return nil, err
	}

	return &srv, nil
}

// UpdateServer allows you to update the name and facility of a server
func (s *Store) UpdateServer(srvUUID uuid.UUID, newS Server) error {
	srv, err := s.FindServerByUUID(srvUUID)
	if err != nil {
		return err
	}

	values := map[string]interface{}{}

	if newS.Name != "" && newS.Name != srv.Name {
		values["name"] = newS.Name
	}

	if newS.FacilityCode != "" && newS.FacilityCode != srv.FacilityCode {
		values["facility_code"] = newS.FacilityCode
	}

	return s.db.Model(&srv).Updates(values).Error
}
