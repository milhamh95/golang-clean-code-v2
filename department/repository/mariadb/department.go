package mariadb

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	sq "github.com/Masterminds/squirrel"
	"github.com/friendsofgo/errors"
	"github.com/segmentio/ksuid"
	log "github.com/sirupsen/logrus"

	"github.com/milhamhidayat/golang-clean-code-v2/domain"
	"github.com/milhamhidayat/golang-clean-code-v2/pkg/cursor"
	ntime "github.com/milhamhidayat/golang-clean-code-v2/pkg/time"
)

// Repository implement all department repository method from interface
type Repository struct {
	DB *sql.DB
}

// New return new department repository
func New(db *sql.DB) Repository {
	return Repository{
		DB: db,
	}
}

// Create is a repository to insert a department
func (r Repository) Create(ctx context.Context, d *domain.Department) (err error) {
	localTime, err := ntime.GetLocalTime()
	if err != nil {
		return
	}

	tx, err := r.DB.BeginTx(ctx, nil)
	if err != nil {
		return
	}

	departmentID := ksuid.New().String()
	if d.ID == "" {
		d.ID = departmentID
	}

	d.CreatedTime = localTime
	d.UpdatedTime = localTime

	query, args, err := sq.Insert("departments").
		Columns("id", "name", "description", "created_time", "updated_time").
		Values(d.ID, d.Name, d.Description, d.CreatedTime, d.UpdatedTime).
		ToSql()
	if err != nil {
		r.rollback(tx)
		return
	}

	stmt, err := tx.PrepareContext(ctx, query)
	if err != nil {
		r.rollback(tx)
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
		r.rollback(tx)
		return
	}

	err = tx.Commit()
	if err != nil {
		r.rollback(tx)
		return
	}

	dept, err := r.Get(ctx, d.ID)
	if err != nil {
		return
	}

	d = &dept

	return
}

// Fetch is a repository to fetch department based on parameter
func (r Repository) Fetch(ctx context.Context, filter domain.DepartmentFilter) (departments []domain.Department, nextCursor string, err error) {
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

		if filter.Num > 0 {
			qSelect = qSelect.Limit(uint64(filter.Num))
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

	defer func() {
		err := rows.Close()
		if err != nil {
			log.Error(err)
		}
	}()

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
		id := departments[len(departments)-1].ID
		nextCursor = cursor.EncodeBase64(id)
	}

	return
}

// Get is a repository to get a department based on parameter
func (r Repository) Get(ctx context.Context, departmentID string) (department domain.Department, err error) {
	query, args, err := sq.Select("id", "name", "description", "created_time", "updated_time").
		From("departments").
		Where(sq.Eq{"id": departmentID}).
		ToSql()
	if err != nil {
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
			err = errors.New("department is not found")
			return
		}
		return
	}

	return
}

// Update is a repository to update a department
func (r Repository) Update(ctx context.Context, d domain.Department) (department domain.Department, err error) {
	localTime, err := ntime.GetLocalTime()
	if err != nil {
		return
	}

	tx, err := r.DB.BeginTx(ctx, nil)
	if err != nil {
		return
	}

	query, args, err := sq.Update("departments").
		SetMap(sq.Eq{
			"name":         d.Name,
			"description":  d.Description,
			"updated_time": localTime,
		}).
		Where(sq.Eq{"id": d.ID}).
		ToSql()
	if err != nil {
		r.rollback(tx)
		return
	}

	stmt, err := tx.PrepareContext(ctx, query)
	if err != nil {
		r.rollback(tx)
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
		r.rollback(tx)
		return
	}

	err = tx.Commit()
	if err != nil {
		return
	}

	department, err = r.Get(ctx, d.ID)
	if err != nil {
		return
	}
	return
}

// Delete is a repository to delete a department
func (r Repository) Delete(ctx context.Context, departmentID string) (err error) {
	tx, err := r.DB.BeginTx(ctx, nil)
	if err != nil {
		return
	}

	query, args, err := sq.Delete("departments").
		Where(sq.Eq{"id": departmentID}).
		ToSql()
	if err != nil {
		r.rollback(tx)
		return
	}

	stmt, err := tx.PrepareContext(ctx, query)
	if err != nil {
		r.rollback(tx)
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
		r.rollback(tx)
		return
	}

	err = tx.Commit()
	return
}

func (r Repository) rollback(tx *sql.Tx) {
	err := tx.Rollback()
	if err != nil && err != sql.ErrTxDone {
		log.Error(err)
	}
}
