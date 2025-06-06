package repo_postgres

import (
	"context"
	"database/sql"
	"fmt"
	"log"

	"github.com/informatik-mannheim/cmg-ss2025/services/worker-registry/ports"
	_ "github.com/lib/pq"
)

type Repo struct {
	db *sql.DB
}

var _ ports.Repo = (*Repo)(nil)

func NewRepo(host, port, user, password, dbName string, ctx context.Context) (*Repo, error) {
	connectionString := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbName,
	)

	log.Printf("DEBUG: Connecting to DB with host=%s port=%s db=%s user=%s", host, port, dbName, user) // DEBUG

	db, err := sql.Open("postgres", connectionString)
	if err != nil {
		return nil, err
	}

	if err := db.PingContext(ctx); err != nil {
		return nil, err
	}

	log.Println("DEBUG: Successfully connected to the PostgreSQL database") // DEBUG
	return &Repo{db: db}, nil
}

func (r *Repo) GetWorkers(status ports.WorkerStatus, zone string, ctx context.Context) ([]ports.Worker, error) {
	log.Printf("DEBUG: GetWorkers called with status=%q zone=%q", status, zone) // DEBUG

	query := `SELECT id, status, zone FROM workers WHERE ($1 = '' OR status = $1) AND ($2 = '' OR zone = $2)`
	log.Printf("DEBUG: Executing SQL: %s", query) // DEBUG

	rows, err := r.db.QueryContext(ctx, query, status, zone)
	if err != nil {
		log.Printf("ERROR: SQL query failed: %v", err) // DEBUG
		return nil, err
	}
	defer rows.Close()

	var workers []ports.Worker
	for rows.Next() {
		var w ports.Worker
		if err := rows.Scan(&w.Id, &w.Status, &w.Zone); err != nil {
			log.Printf("ERROR: Failed to scan row: %v", err) // DEBUG
			return nil, err
		}
		log.Printf("DEBUG: Fetched worker: %+v", w) // DEBUG
		workers = append(workers, w)
	}

	if len(workers) == 0 {
		log.Println("DEBUG: No workers found matching the criteria") // DEBUG
	}
	return workers, nil
}

func (r *Repo) GetWorkerById(id string, ctx context.Context) (ports.Worker, error) {
	var w ports.Worker
	query := `SELECT id, status, zone FROM workers WHERE id = $1`
	err := r.db.QueryRowContext(ctx, query, id).Scan(&w.Id, &w.Status, &w.Zone)
	if err == sql.ErrNoRows {
		return ports.Worker{}, ports.NewErrWorkerNotFound(id)
	} else if err != nil {
		return ports.Worker{}, err
	}
	return w, nil
}

func (r *Repo) CreateWorker(worker ports.Worker, ctx context.Context) error {
	if worker.Status == "" || worker.Zone == "" {
		return ports.NewErrCreatingWorkerFailed()
	}
	query := `INSERT INTO workers (id, status, zone) VALUES ($1, $2, $3)`
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
