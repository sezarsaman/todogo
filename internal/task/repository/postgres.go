package repository

import (
	"context"

	"task-manager/internal/task/model"

	"github.com/jackc/pgx/v5"
)

type PostgresRepository struct {
	conn *pgx.Conn
}

func NewPostgres(conn *pgx.Conn) *PostgresRepository {
	return &PostgresRepository{conn: conn}
}

func (r *PostgresRepository) Create(ctx context.Context, task *model.Task) error {
	return r.conn.QueryRow(
		ctx,
		`insert into tasks (title, description, status, assignee) values ($1,$2,$3,$4)
		 returning id, created_at, updated_at`,
		task.Title, task.Description, task.Status, task.Assignee,
	).Scan(&task.ID, &task.CreatedAt, &task.UpdatedAt)
}

func (r *PostgresRepository) GetByID(ctx context.Context, id int64) (*model.Task, error) {
	row := r.conn.QueryRow(
		ctx,
		`select id,title,description,status,assignee,created_at,updated_at
		 from tasks where id=$1`, id,
	)

	var t model.Task
	err := row.Scan(
		&t.ID,
		&t.Title,
		&t.Description,
		&t.Status,
		&t.Assignee,
		&t.CreatedAt,
		&t.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	return &t, nil
}

func (r *PostgresRepository) List(ctx context.Context) ([]model.Task, error) {
	rows, err := r.conn.Query(
		ctx,
		`select id,title,description,status,assignee,created_at,updated_at from tasks`,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tasks []model.Task
	for rows.Next() {
		var t model.Task
		if err := rows.Scan(
			&t.ID,
			&t.Title,
			&t.Description,
			&t.Status,
			&t.Assignee,
			&t.CreatedAt,
			&t.UpdatedAt,
		); err != nil {
			return nil, err
		}
		tasks = append(tasks, t)
	}

	return tasks, nil
}

func (r *PostgresRepository) Update(ctx context.Context, task *model.Task) error {
	_, err := r.conn.Exec(
		ctx,
		`update tasks set title=$1,description=$2,status=$3,assignee=$4,updated_at=now()
		 where id=$5`,
		task.Title, task.Description, task.Status, task.Assignee, task.ID,
	)
	return err
}

func (r *PostgresRepository) Delete(ctx context.Context, id int64) error {
	_, err := r.conn.Exec(ctx, `delete from tasks where id=$1`, id)
	return err
}
