package repo_database

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/url"
	"time"

	//"github.com/informatik-mannheim/cmg-ss2025/pkg/logging"
	"github.com/informatik-mannheim/cmg-ss2025/services/job/ports"
	_ "github.com/jackc/pgx/v5/stdlib"
)

type JobStorage struct {
	db *sql.DB
}

var _ ports.JobStorage = (*JobStorage)(nil)

func NewJobStorage(host, port, user, password, dbName, sslMode string, ctx context.Context) (*JobStorage, error) {
	escapedPassword := url.QueryEscape(password)
	connectionString := ""
	if sslMode == "true" {
		connectionString = fmt.Sprintf(
			"postgres://%s:%s@%s:%s/%s?sslmode=require",
			user, escapedPassword, host, port, dbName,
		)
	} else {
		connectionString = fmt.Sprintf(
			"postgres://%s:%s@%s:%s/%s",
			user, escapedPassword, host, port, dbName,
		)
	}

	//logging.Debug("Connecting to DB with Connectionstring:", connectionString)
	fmt.Println("DEBUG: Connecting to DB with Connectionstring:", connectionString)

	db, err := sql.Open("pgx", connectionString)
	if err != nil {
		return nil, err
	}

	if err := db.PingContext(ctx); err != nil {
		return nil, err
	}

	//logging.Debug("Successfully connected to the PostgreSQL database")
	fmt.Println("DEBUG: Successfully connected to the PostgreSQL database")
	return &JobStorage{db: db}, nil
}

func (r *JobStorage) GetJobs(ctx context.Context, status []ports.JobStatus) ([]ports.Job, error) {
	query := `SELECT id, user_id, created_at, updated_at, job_name, image_name, image_version, adjustment_parameters, creation_zone, worker_id, compute_zone, carbon_intensity, carbon_savings, result, error_message, job_status
              FROM jobs`
	var args []interface{}
	if len(status) > 0 {
		var statusStrings []string
		for _, s := range status {
			statusStrings = append(statusStrings, string(s))
		}
		query += " WHERE job_status = ANY($1)"
		args = append(args, statusStrings)
	}
	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var jobs []ports.Job
	for rows.Next() {
		var job ports.Job
		var imageName, imageVersion string
		var paramsJSON []byte
		err := rows.Scan(
			&job.Id, &job.UserID, &job.CreatedAt, &job.UpdatedAt, &job.JobName,
			&imageName, &imageVersion, &paramsJSON, &job.CreationZone,
			&job.WorkerID, &job.ComputeZone, &job.CarbonIntensity, &job.CarbonSaving,
			&job.Result, &job.ErrorMessage, &job.Status,
		)
		if err != nil {
			return nil, err
		}
		if err := json.Unmarshal(paramsJSON, &job.AdjustmentParameters); err != nil {
			return nil, err
		}
		job.Image = ports.ContainerImage{
			Name:    imageName,
			Version: imageVersion,
		}
		jobs = append(jobs, job)
	}
	return jobs, nil
}

func (r *JobStorage) GetJob(ctx context.Context, id string) (ports.Job, error) {
	query := `SELECT id, user_id, created_at, updated_at, job_name, image_name, image_version, adjustment_parameters, creation_zone, worker_id, compute_zone, carbon_intensity, carbon_savings, result, error_message, job_status
              FROM jobs WHERE id = $1`
	var job ports.Job
	var imageName, imageVersion string
	var paramsJSON []byte
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&job.Id, &job.UserID, &job.CreatedAt, &job.UpdatedAt, &job.JobName,
		&imageName, &imageVersion, &paramsJSON, &job.CreationZone,
		&job.WorkerID, &job.ComputeZone, &job.CarbonIntensity, &job.CarbonSaving,
		&job.Result, &job.ErrorMessage, &job.Status,
	)
	if err == sql.ErrNoRows {
		return ports.Job{}, ports.ErrJobNotFound
	} else if err != nil {
		return ports.Job{}, err
	}
	if err := json.Unmarshal(paramsJSON, &job.AdjustmentParameters); err != nil {
		return ports.Job{}, err
	}
	job.Image = ports.ContainerImage{
		Name:    imageName,
		Version: imageVersion,
	}
	return job, nil
}

func (r *JobStorage) CreateJob(ctx context.Context, job ports.Job) error {
	paramsJSON, err := json.Marshal(job.AdjustmentParameters)
	if err != nil {
		return err
	}
	query := `INSERT INTO jobs (id, user_id, created_at, updated_at, job_name, image_name, image_version, adjustment_parameters, creation_zone, worker_id, compute_zone, carbon_intensity, carbon_savings, result, error_message, job_status)
              VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14,$15,$16)`
	_, err = r.db.ExecContext(ctx, query,
		job.Id, job.UserID, job.CreatedAt, job.UpdatedAt, job.JobName,
		job.Image.Name, job.Image.Version, paramsJSON, job.CreationZone,
		job.WorkerID, job.ComputeZone, job.CarbonIntensity, job.CarbonSaving,
		job.Result, job.ErrorMessage, job.Status,
	)
	return err
}

func (r *JobStorage) UpdateJob(ctx context.Context, id string, job ports.Job) (ports.Job, error) {
	paramsJSON, err := json.Marshal(job.AdjustmentParameters)
	if err != nil {
		return ports.Job{}, err
	}
	query := `UPDATE jobs SET
        user_id=$2, updated_at=$3, job_name=$4, image_name=$5, image_version=$6, adjustment_parameters=$7, creation_zone=$8, worker_id=$9, compute_zone=$10, carbon_intensity=$11, carbon_savings=$12, result=$13, error_message=$14, job_status=$15
        WHERE id=$1`
	res, err := r.db.ExecContext(ctx, query,
		id, job.UserID, time.Now(), job.JobName,
		job.Image.Name, job.Image.Version, paramsJSON, job.CreationZone,
		job.WorkerID, job.ComputeZone, job.CarbonIntensity, job.CarbonSaving,
		job.Result, job.ErrorMessage, job.Status,
	)
	if err != nil {
		return ports.Job{}, err
	}
	rowsAffected, _ := res.RowsAffected()
	if rowsAffected == 0 {
		return ports.Job{}, ports.ErrJobNotFound
	}
	return r.GetJob(ctx, id)
}
