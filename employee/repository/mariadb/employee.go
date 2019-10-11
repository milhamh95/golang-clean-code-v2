package mariadb

import (
	"context"
	"database/sql"

	sq "github.com/Masterminds/squirrel"
	"github.com/friendsofgo/errors"
	"github.com/segmentio/ksuid"
	log "github.com/sirupsen/logrus"

	"github.com/milhamhidayat/golang-clean-code-v2/domain"
	ntime "github.com/milhamhidayat/golang-clean-code-v2/pkg/time"
)

// Repository implement all employee repository method from interface
type Repository struct {
	DB *sql.DB
}

// New return new department repository
func New(db *sql.DB) Repository {
	return Repository{
		DB: db,
	}
}

// Create is a repository to insert an employee
func (r Repository) Create(ctx context.Context, e *domain.Employee) (err error) {
	localTime, err := ntime.GetLocalTime()
	if err != nil {
		return
	}

	tx, err := r.DB.BeginTx(ctx, nil)
	if err != nil {
		return
	}

	employeeID := ksuid.New().String()
	if e.ID == "" {
		e.ID = employeeID
	}

	lastname := sql.NullString{}
	if e.LastName != "" {
		lastname = sql.NullString{
			Valid:  true,
			String: e.LastName,
		}
	}

	e.CreatedTime = localTime
	e.UpdatedTime = localTime

	query, args, err := sq.Insert("employee").
		Columns("id", "first_name", "last_name", "birth_place", "date_of_birth", "title", "dept_id", "created_time", "updated_time").
		Values(e.ID, e.FirstName, lastname, e.BirthPlace, e.DateOfBirth, e.Title, e.Department.ID, e.CreatedTime, e.UpdatedTime).
		ToSql()
	if err != nil {
		r.rollback(tx, "failed to generate insert employee query")
		return
	}

	stmt, err := tx.PrepareContext(ctx, query)
	if err != nil {
		r.rollback(tx, "failed to prepared insert employee statement")
		return
	}

	defer func() {
		err = stmt.Close()
		if err != nil {
			log.Error(errors.Wrap(err, "failed to close insert employee statement"))
		}
	}()

	_, err = stmt.ExecContext(ctx, args...)
	if err != nil {
		r.rollback(tx, "failed to execute insert employee statement")
		return
	}

	err = tx.Commit()
	return nil
}

// Get is a repository to get an employee
func (r Repository) Get(ctx context.Context, employeeID string) (employee domain.Employee, err error) {
	query, args, err := sq.Select("id", "first_name", "last_name", "birth_place", "date_of_birth", "title", "dept_id").
		From("employees").
		Where(sq.Eq{"id": employeeID}).
		ToSql()
	if err != nil {
		return
	}

	lastname := sql.NullString{}

	row := r.DB.QueryRowContext(ctx, query, args...)
	err = row.Scan(
		&employee.ID,
		&employee.FirstName,
		&lastname,
		&employee.BirthPlace,
		&employee.DateOfBirth,
		&employee.Title,
		&employee.Department.ID,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			err = errors.Errorf("employee is not found: %s", employeeID)
			return
		}
		return
	}

	return
}

// Fetch is a repository to fetch employees
func (r Repository) Fetch(ctx context.Context, filter domain.EmployeeFilter) (employees []domain.Employee, nextCursor string, err error) {
	return
}

// Update is a repository to update an employee
func (r Repository) Update(ctx context.Context, e domain.Employee) (employee domain.Employee, err error) {
	return
}

// Delete is a repository to delete an employee
func (r Repository) Delete(ctx context.Context, employeeID string) (err error) {
	return
}

func (r Repository) rollback(tx *sql.Tx, msg string) {
	err := tx.Rollback()
	if err != nil && err != sql.ErrTxDone {
		log.Error(errors.Wrap(err, msg))
	}
}
