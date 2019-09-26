package mariadb

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	sq "github.com/Masterminds/squirrel"
	"github.com/pkg/errors"
	"github.com/segmentio/ksuid"
	log "github.com/sirupsen/logrus"

	"github.com/milhamhidayat/golang-clean-code-v2/domain"
	employee "github.com/milhamhidayat/golang-clean-code-v2/domain"
	"github.com/milhamhidayat/golang-clean-code-v2/pkg/cursor"
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

	collectionID := ksuid.New().String()
	if d.ID == "" {
		d.ID = collectionID
	}

	query, args, err := sq.Insert("departments").
		Columns("id", "name", "description").
		Values(d.ID, d.Name, d.Description).
		ToSql()
	if err != nil {
		defer func() {
			rollback(tx)
		}()
		err = errors.Wrap(err, "error generating query")
		return
	}

	stmt, err := tx.PrepareContext(ctx, query)
	if err != nil {
		defer func() {
			rollback(tx)
		}()
		err = errors.Wrap(err, "error prepare context")
		return
	}

	defer func() {
		err = stmt.Close()
		if err != nil {
			log.Error(err)
		}
	}()

	_, err = stmt.ExecContext(ctx, args...)
	if err != nil {
		defer func() {
			rollback(tx)
		}()
		err = errors.Wrap(err, "error when inserting department")
		return
	}

	err = tx.Commit()
	if err != nil {
		return
	}

	return nil
}

// Fetch is a repository to fetch articles based on parameter
func (r DepartmentRepository) Fetch(ctx context.Context, filter employee.DepartmentFilter) (departments []employee.Department, nextCursor string, err error) {
	qSelect := sq.Select("id", "name", "description", "created_time", "updated_time").
		From("departments")

	if len(filter.IDs) != 0 {
		qSelect = qSelect.Where(sq.Eq{"id": filter.IDs})
		qField := strings.Repeat(",?", len(filter.IDs))
		qOrderBy := fmt.Sprintf("ORDER BY FIELD(id%s)", qField)
		qSelect = qSelect.Suffix(qOrderBy)
	} else {
		qSelect = qSelect.OrderBy(`id desc`)

		if filter.Keyword != "" {
			qSelect = qSelect.Where(`name LIKE ?`, fmt.Sprint("%", filter.Keyword, "%"))
		}

		if filter.Cursor != "" {
			var id string
			id, er := cursor.DecodeBase64(filter.Cursor)
			if er != nil {
				err = er
				return
			}
			qSelect = qSelect.Where(sq.Lt{"id": id})
		}
	}

	query, args, err := qSelect.ToSql()

	if len(filter.IDs) != 0 {
		args = append(args, args...)
	}

	if err != nil {
		return
	}

	rows, err := r.DB.QueryContext(ctx, query, args...)
	if err != nil {
		return
	}

	for rows.Next() {
		d := domain.Department{}

		err = rows.Scan(
			&d.ID,
			&d.Name,
			&d.Description,
			&d.CreatedTime,
			&d.UpdatedTime,
		)
		if err != nil {
			return
		}
		departments = append(departments, d)
	}

	err = rows.Err()

	if len(filter.IDs) != 0 {
		return
	}

	nextCursor = filter.Cursor
	if len(departments) >= 1 {
		fmt.Println("called")
		id := departments[len(departments)-1].ID
		nextCursor = cursor.EncodeBase64(id)
	}

	return
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
	tx, err := r.DB.BeginTx(ctx, nil)
	if err != nil {
		return
	}

	query, args, err := sq.Delete("departments").
		Where(sq.Eq{"id": departmentID}).
		ToSql()
	if err != nil {
		defer func() {
			rollback(tx)
		}()
		return
	}

	stmt, err := tx.PrepareContext(ctx, query)
	if err != nil {
		defer func() {
			rollback(tx)
		}()
		return
	}

	defer func() {
		err = stmt.Close()
		if err != nil {
			log.Error(err)
		}
	}()

	_, err = stmt.ExecContext(ctx, args...)
	if err != nil {
		defer func() {
			rollback(tx)
		}()
	}

	err = tx.Commit()
	if err != nil {
		return
	}

	return
}

func rollback(tx *sql.Tx) {
	err := tx.Rollback()
	if err != nil && err != sql.ErrTxDone {
		log.Error(err)
	}
}
