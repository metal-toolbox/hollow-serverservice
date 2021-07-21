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
	ID                  uuid.UUID `gorm:"type:uuid;primary_key;default:gen_random_uuid();"`
	CreatedAt           time.Time
	UpdatedAt           time.Time
	Name                string
	FacilityCode        string
	Attributes          []Attributes          `gorm:"polymorphic:Entity;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	ServerComponents    []ServerComponent     `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	VersionedAttributes []VersionedAttributes `gorm:"polymorphic:Entity;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}

func serverPreload(db *gorm.DB) *gorm.DB {
	d := db.Preload("VersionedAttributes",
		"(created_at, namespace, entity_id, entity_type) IN (?)",
		db.Table("versioned_attributes").Select("max(created_at), namespace, entity_id, entity_type").Group("namespace").Group("entity_id").Group("entity_type"),
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
func (s *Store) GetServers(filter *ServerFilter, pager *Pagination) ([]Server, error) {
	var srvs []Server

	d := serverPreload(s.db)

	if filter != nil {
		d = filter.apply(d)
	}

	if pager == nil {
		pager = &Pagination{}
	}

	if err := d.Scopes(paginate(*pager)).Find(&srvs).Error; err != nil {
		return nil, err
	}

	return srvs, nil
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
