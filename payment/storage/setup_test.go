package storage

import (
	"context"
	"os/exec"
	"testing"
	"time"

	"github.com/testcontainers/testcontainers-go/modules/mssql"
	"gorm.io/driver/sqlserver"
	"gorm.io/gorm"
)

func SetupMSSQL(ctx context.Context, t *testing.T) (*mssql.MSSQLServerContainer, *gorm.DB) {
	container, err := mssql.Run(ctx,
		"mcr.microsoft.com/azure-sql-edge",
		mssql.WithAcceptEULA(),
		mssql.WithPassword("SuperStrong@Passw0rd"),
	)

	time.Sleep(1 * time.Second)

	if err != nil {
		t.Fatalf("Failed to start container: %v", err)
	}

	connStr, err := container.ConnectionString(ctx, "encrypt=disable", "parseTime=True", "time_zone=UTC")
	if err != nil {
		t.Fatalf("Failed to create connection string: %v", err)
	}

	db, err := gorm.Open(sqlserver.Open(connStr), &gorm.Config{})
	if err != nil {
		t.Fatalf("Failed to connect MSSQL: %v", err)
	}
	return container, db
}

func CleanUpMSSQL(container *mssql.MSSQLServerContainer, ctx context.Context, t *testing.T) {
	if err := container.Terminate(ctx); err != nil {
		t.Fatalf("Failed to terminate container: %v", err)
	}
	cmd := exec.Command("docker", "rm", "-f", "$(docker ps -aq)")
	if err := cmd.Run(); err != nil {
		t.Fatalf("Failed to clean up containers: %v", err)
	}
}
