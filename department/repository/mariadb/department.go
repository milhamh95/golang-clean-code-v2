package mariadb

import (
	"context"
	"database/sql"
	"fmt"

	sq "github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"

	employee "github.com/milhamhidayat/golang-clean-code-v2/domain"
)

// DepartmentRepository implement all method from interface
type DepartmentRepository struct {
	DB *sql.DB
}

// NewDepartmentRepository return new department repository
func NewDepartmentRepository(db *sql.DB) DepartmentRepository {
	return DepartmentRepository{
		DB: db,
	}
}

// Create is a repository to insert an article
func (r DepartmentRepository) Create(ctx context.Context, d *employee.Department) (err error) {
	tx, err := r.DB.BeginTx(ctx, nil)
	if err != nil {
		err = errors.Wrap(err, "error starting transaction")
		return
	}

	collectionID := uuid.New().String()
	if d.ID == "" {
		fmt.Println("called")
		d.ID = collectionID
	}

	query, args, err := sq.Insert("departments").
		Columns("id", "name", "description").
		Values(d.ID, d.Name, d.Description).
		ToSql()
	if err != nil {
		err = errors.Wrap(err, "error generating query")
		return
	}

	stmt, err := tx.PrepareContext(ctx, query)
	if err != nil {
		err = errors.Wrap(err, "error prepare context")
		return
	}

	defer func() {
		err = stmt.Close()
		if err != nil {
			log.Warn(err)
		}
	}()

	_, err = r.DB.ExecContext(ctx, query, args...)
	if err != nil {
		err = errors.Wrap(err, "error when inserting department")
		return
	}

	return nil
}

// Fetch is a repository to fetch articles based on parameter
func (r DepartmentRepository) Fetch(ctx context.Context, filter employee.DepartmentFilter) (departments []employee.Department, nextCursor string, err error) {
	return []employee.Department{}, "", nil
}

// Get is a repository to get an article based on parameter
func (r DepartmentRepository) Get(ctx context.Context, departmentID string) (department employee.Department, err error) {
	query, args, err := sq.Select("id", "name", "description", "created_time", "updated_time").
		From("departments").
		Where(sq.Eq{"id": departmentID}).
		ToSql()
	if err != nil {
		err = errors.Wrap(err, "error when building query")
		return
	}

	row := r.DB.QueryRowContext(ctx, query, args...)
	err = row.Scan(
		&department.ID,
		&department.Name,
		&department.Description,
		&department.CreatedTime,
		&department.UpdatedTime,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			err = errors.New("data is not found")
			return
		}
		err = errors.Wrap(err, "error scan the result")
		return
	}

	return
}

// Update is a repository to update an article
func (r DepartmentRepository) Update(ctx context.Context, d employee.Department) (department employee.Department, err error) {
	return employee.Department{}, nil
}

// Delete is a repository to delete an article
func (r DepartmentRepository) Delete(ctx context.Context, departmentID string) (err error) {
	return nil
}
