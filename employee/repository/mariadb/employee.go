package mariadb

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	sq "github.com/Masterminds/squirrel"
	"github.com/friendsofgo/errors"
	"github.com/go-sql-driver/mysql"
	"github.com/segmentio/ksuid"
	log "github.com/sirupsen/logrus"

	"github.com/milhamhidayat/golang-clean-code-v2/domain"
	"github.com/milhamhidayat/golang-clean-code-v2/pkg/cursor"
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

	query, args, err := sq.Insert("employees").
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

	defer r.closeStatement(stmt)

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
	query, args, err := sq.Select("id", "first_name", "last_name", "birth_place", "date_of_birth", "title", "dept_id", "created_time", "updated_time").
		From("employees").
		Where(sq.Eq{"id": employeeID}).
		ToSql()
	if err != nil {
		return
	}

	lastname := sql.NullString{}
	dateOfBirth := mysql.NullTime{}
	createdTime := mysql.NullTime{}
	updatedTime := mysql.NullTime{}

	row := r.DB.QueryRowContext(ctx, query, args...)
	err = row.Scan(
		&employee.ID,
		&employee.FirstName,
		&lastname,
		&employee.BirthPlace,
		&dateOfBirth,
		&employee.Title,
		&employee.Department.ID,
		&createdTime,
		&updatedTime,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			err = domain.ErrNotFound
			return
		}
		return
	}

	employee.LastName = lastname.String
	employee.SetDateOfBirth(dateOfBirth.Time)
	employee.CreatedTime = createdTime.Time
	employee.UpdatedTime = updatedTime.Time
	return
}

// Fetch is a repository to fetch employees
func (r Repository) Fetch(ctx context.Context, filter domain.EmployeeFilter) (employees []domain.Employee, nextCursor string, err error) {
	employees = make([]domain.Employee, 0)
	qSelect := sq.Select("id", "first_name", "last_name", "birth_place", "date_of_birth", "title", "dept_id", "created_time", "updated_time").
		From("employees")

	if len(filter.IDs) != 0 {
		qSelect = qSelect.Where(sq.Eq{"id": filter.IDs})
		qField := strings.Repeat(",?", len(filter.IDs))
		qOrderBy := fmt.Sprintf("ORDER BY FIELD(id%s)", qField)
		qSelect = qSelect.Suffix(qOrderBy)
	} else if len(filter.DeptIDs) != 0 {
		qSelect = qSelect.Where(sq.Eq{"id": filter.IDs})
		qField := strings.Repeat(",?", len(filter.IDs))
		qOrderBy := fmt.Sprintf("ORDER BY FIELD(id%s)", qField)
		qSelect = qSelect.OrderBy(qOrderBy)
	} else {
		qSelect = qSelect.OrderBy("id desc")

		if filter.Keyword != "" {
			qSelect = qSelect.Where(`first_name LIKE ?`, fmt.Sprint("%", filter.Keyword, "%"))
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
	if err != nil {
		return
	}

	if len(filter.IDs) != 0 || len(filter.DeptIDs) != 0 {
		args = append(args, args...)
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
		lastname := sql.NullString{}
		dateOfBirth := mysql.NullTime{}
		createdTime := mysql.NullTime{}
		updatedTime := mysql.NullTime{}
		e := domain.Employee{}

		err = rows.Scan(
			&e.ID,
			&e.FirstName,
			&lastname,
			&e.BirthPlace,
			&dateOfBirth,
			&e.Title,
			&e.Department.ID,
			&createdTime,
			&updatedTime,
		)
		if err != nil {
			return
		}

		e.LastName = lastname.String
		e.SetDateOfBirth(dateOfBirth.Time)
		e.CreatedTime = createdTime.Time
		e.UpdatedTime = updatedTime.Time
		employees = append(employees, e)
	}

	err = rows.Err()
	if err != nil {
		return
	}

	if len(filter.IDs) != 0 {
		return
	}

	nextCursor = filter.Cursor
	if len(employees) >= 1 {
		id := employees[len(employees)-1].ID
		nextCursor = cursor.EncodeBase64(id)
	}

	return
}

// Update is a repository to update an employee
func (r Repository) Update(ctx context.Context, e domain.Employee) (employee domain.Employee, err error) {
	localTime, err := ntime.GetLocalTime()
	if err != nil {
		return
	}

	lastname := sql.NullString{}

	if e.LastName != "" {
		lastname = sql.NullString{
			Valid:  true,
			String: e.LastName,
		}
	}

	tx, err := r.DB.BeginTx(ctx, nil)
	if err != nil {
		return
	}

	query, args, err := sq.Update("employees").
		SetMap(sq.Eq{
			"first_name":    e.FirstName,
			"last_name":     lastname,
			"birth_place":   e.BirthPlace,
			"date_of_birth": e.DateOfBirth,
			"title":         e.Title,
			"dept_id":       e.Department.ID,
			"updated_time":  localTime,
		}).
		Where(sq.Eq{"id": e.ID}).
		ToSql()
	if err != nil {
		r.rollback(tx, "failed to prepare update employee query")
		return
	}

	stmt, err := tx.PrepareContext(ctx, query)
	if err != nil {
		r.rollback(tx, "failed to prepared update employee statement")
		return
	}

	defer r.closeStatement(stmt)

	res, err := stmt.ExecContext(ctx, args...)
	if err != nil {
		r.rollback(tx, "failed to update employee")
		return
	}

	err = tx.Commit()
	if err != nil {
		r.rollback(tx, "failed to rollback after commit")
		return
	}

	count, err := res.RowsAffected()
	if err != nil {
		return
	}

	if count == 0 {
		err = domain.ErrNotFound
		return
	}

	employee, err = r.Get(ctx, e.ID)
	if err != nil {
		return
	}

	return
}

// Delete is a repository to delete an employee
func (r Repository) Delete(ctx context.Context, employeeID string) (err error) {
	tx, err := r.DB.BeginTx(ctx, nil)
	if err != nil {
		return
	}

	query, args, err := sq.Delete("employees").
		Where(sq.Eq{"id": employeeID}).
		ToSql()
	if err != nil {
		r.rollback(tx, "failed to prepare delete employee query")
		return
	}

	stmt, err := tx.PrepareContext(ctx, query)
	if err != nil {
		r.rollback(tx, "failed to prepare delete employee statement")
	}

	defer r.closeStatement(stmt)

	res, err := stmt.ExecContext(ctx, args...)
	if err != nil {
		r.rollback(tx, "failed to execute delete employee")
		return
	}

	err = tx.Commit()
	if err != nil {
		r.rollback(tx, "failed to commit")
	}

	count, err := res.RowsAffected()
	if err != nil {
		return
	}

	if count == 0 {
		err = domain.ErrNotFound
		return
	}

	return
}

func (r Repository) rollback(tx *sql.Tx, msg string) {
	err := tx.Rollback()
	if err != nil && err != sql.ErrTxDone {
		log.Error(errors.Wrap(err, msg))
	}
}

func (r Repository) closeStatement(stmt *sql.Stmt) {
	err := stmt.Close()
	if err != nil {
		log.Error(err)
	}
}
