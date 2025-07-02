package repo_postgres

import (
	"context"
	"database/sql"
	"fmt"
	"net/url"

	"github.com/informatik-mannheim/cmg-ss2025/pkg/logging"
	"github.com/informatik-mannheim/cmg-ss2025/services/worker-registry/ports"
	_ "github.com/jackc/pgx/v5/stdlib"
)

type Repo struct {
	db *sql.DB
}

var _ ports.Repo = (*Repo)(nil)

func NewRepo(host, port, user, password, dbName string, sslmode bool, ctx context.Context) (*Repo, error) {
	escapedPassword := url.QueryEscape(password)

	sslModeParam := "disable"
	if sslmode {
		sslModeParam = "require"
	}

	connectionString := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=%s",
		user, escapedPassword, host, port, dbName, sslModeParam,
	)

	logging.Debug("Connecting to DB with Connectionstring:", "conn", connectionString)

	db, err := sql.Open("pgx", connectionString)
	if err != nil {
		return nil, err
	}

	if err := db.PingContext(ctx); err != nil {
		return nil, err
	}

	logging.Debug("Successfully connected to the PostgreSQL database")
	return &Repo{db: db}, nil
}

func (r *Repo) GetWorkers(status ports.WorkerStatus, zone string, ctx context.Context) ([]ports.Worker, error) {
	message := fmt.Sprintf("GetWorkers called with status=%q zone=%q", status, zone)
	logging.Debug(message)

	query := `SELECT id, status, zone FROM workers WHERE ($1 = '' OR status = $1) AND ($2 = '' OR zone = $2)`
	logging.Debug("Executing SQL:", query)

	rows, err := r.db.QueryContext(ctx, query, status, zone)
	if err != nil {
		logging.Warn("SQL query failed:", err)
		return nil, err
	}
	defer rows.Close()

	var workers []ports.Worker
	for rows.Next() {
		var w ports.Worker
		if err := rows.Scan(&w.Id, &w.Status, &w.Zone); err != nil {
			logging.Warn("Failed to scan row:", err)
			return nil, err
		}
		logging.Debug("Fetched worker:", w)
		workers = append(workers, w)
	}

	if len(workers) == 0 {
		logging.Debug("No workers found matching the criteria")
	}
	return workers, nil
}

func (r *Repo) GetWorkerById(id string, ctx context.Context) (ports.Worker, error) {
	var w ports.Worker
	query := `SELECT id, status, zone FROM workers WHERE id = $1`
	logging.Debug("Executing SQL:", query)
	err := r.db.QueryRowContext(ctx, query, id).Scan(&w.Id, &w.Status, &w.Zone)
	if err == sql.ErrNoRows {
		return ports.Worker{}, ports.NewErrWorkerNotFound(id)
	} else if err != nil {
		logging.Debug("Fetched worker:", w)
		return ports.Worker{}, err
	}
	return w, nil
}

func (r *Repo) CreateWorker(worker ports.Worker, ctx context.Context) error {
	if worker.Status == "" || worker.Zone == "" {
		return ports.NewErrCreatingWorkerFailed()
	}
	query := `INSERT INTO workers (id, status, zone) VALUES ($1, $2, $3)`
	logging.Debug("Executing SQL:", query)
	_, err := r.db.ExecContext(ctx, query, worker.Id, worker.Status, worker.Zone)
	if err != nil {
		return ports.NewErrCreatingWorkerFailed()
	}
	return nil
}

func (r *Repo) UpdateWorkerStatus(id string, status ports.WorkerStatus, ctx context.Context) (ports.Worker, error) {
	if !isValidStatus(status) {
		return ports.Worker{}, ports.NewErrUpdatingWorkerFailed(id)
	}
	query := `UPDATE workers SET status = $1 WHERE id = $2`
	logging.Debug("Executing SQL:", query)
	res, err := r.db.ExecContext(ctx, query, status, id)
	if err != nil {
		return ports.Worker{}, ports.NewErrUpdatingWorkerFailed(id)
	}
	rowsAffected, _ := res.RowsAffected()
	if rowsAffected == 0 {
		return ports.Worker{}, ports.NewErrWorkerNotFound(id)
	}
	return r.GetWorkerById(id, ctx)
}

func isValidStatus(status ports.WorkerStatus) bool {
	return status == ports.StatusAvailable || status == ports.StatusRunning
}
