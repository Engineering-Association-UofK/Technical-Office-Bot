package routes

import (
	"log/slog"
	"net/http"

	"github.com/Engineering-Association-UofK/Technical-Office-Bot/internal/config"
	"github.com/Engineering-Association-UofK/Technical-Office-Bot/internal/handler"
	"github.com/Engineering-Association-UofK/Technical-Office-Bot/internal/handler/middleware"
	"github.com/go-chi/chi/v5"
	"github.com/rs/cors"

	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/adapters/humachi"
)

var (
	Fbh *handler.FeedbackHandler
	Hh  *handler.HealthHandler
	Bh  *handler.BackupHandler
)

func HttpStart() {
	router := chi.NewMux()
	router.Use(middleware.Basic)

	conf := huma.DefaultConfig("Technical Office Bot API", "1.0.0")
	conf.Components.SecuritySchemes = map[string]*huma.SecurityScheme{
		"bearerAuth": {
			Type:         "http",
			Scheme:       "bearer",
			BearerFormat: "JWT",
		},
	}
	api := humachi.New(router, conf)

	root := huma.NewGroup(api, "/api/v1")

	////////////////////
	////  Feedback  ////
	////////////////////
	huma.Register(root, huma.Operation{
		OperationID: "post-feedback",
		Method:      http.MethodPost,
		Path:        "/feedback",
		Summary:     "Send feedback",
		Description: "Send feedback to the technical office",
		Tags:        []string{"Account"},
	}, Fbh.PostFeedback)

	adminGroup := huma.NewGroup(root, "/admin")
	adminGroup.UseMiddleware(middleware.HumaAuth(api))

	////////////////////
	////   Health   ////
	////////////////////

	huma.Register(adminGroup, huma.Operation{
		OperationID: "get-overview",
		Method:      http.MethodGet,
		Path:        "/health/overview",
		Summary:     "Get health overview",
		Description: "Get health overview",
		Tags:        []string{"Health"},
		Security:    []map[string][]string{{"bearerAuth": []string{"admin"}}},
	}, Hh.GetOverview)

	huma.Register(adminGroup, huma.Operation{
		OperationID: "get-system-health",
		Method:      http.MethodGet,
		Path:        "/health/system",
		Summary:     "Get system health",
		Description: "Get system health",
		Tags:        []string{"Health"},
		Security:    []map[string][]string{{"bearerAuth": []string{"admin"}}},
	}, Hh.GetSystemHealth)

	huma.Register(adminGroup, huma.Operation{
		OperationID: "get-app-health",
		Method:      http.MethodGet,
		Path:        "/health/app",
		Summary:     "Get app health",
		Description: "Get app health",
		Tags:        []string{"Health"},
		Security:    []map[string][]string{{"bearerAuth": []string{"admin"}}},
	}, Hh.GetAppHealth)

	////////////////////
	////   Backup   ////
	////////////////////

	backupGroup := huma.NewGroup(adminGroup, "")
	backupGroup.UseMiddleware(middleware.HumaRequireRole(api, middleware.RoleTechSupport))

	huma.Register(backupGroup, huma.Operation{
		OperationID: "trigger-backup",
		Method:      http.MethodPost,
		Path:        "/backup/trigger",
		Summary:     "Trigger backup",
		Description: "Manually trigger a system backup",
		Tags:        []string{"Backup"},
		Security:    []map[string][]string{{"bearerAuth": []string{"admin"}}},
	}, Bh.TriggerBackup)

	huma.Register(backupGroup, huma.Operation{
		OperationID: "trigger-restore",
		Method:      http.MethodPost,
		Path:        "/restore/trigger",
		Summary:     "Trigger restore",
		Description: "Manually trigger a system restore from the latest backup",
		Tags:        []string{"Backup"},
		Security:    []map[string][]string{{"bearerAuth": []string{"admin"}}},
	}, Bh.TriggerRestore)

	huma.Register(backupGroup, huma.Operation{
		OperationID: "get-backup-details",
		Method:      http.MethodGet,
		Path:        "/backup/details",
		Summary:     "Get backup details",
		Description: "Get details about the latest backup and schedule",
		Tags:        []string{"Backup"},
		Security:    []map[string][]string{{"bearerAuth": []string{"admin"}}},
	}, Bh.GetBackupDetails)

	huma.Register(backupGroup, huma.Operation{
		OperationID: "change-backup-interval",
		Method:      http.MethodPost,
		Path:        "/backup/interval",
		Summary:     "Change backup interval",
		Description: "Update the frequency of automatic backups",
		Tags:        []string{"Backup"},
		Security:    []map[string][]string{{"bearerAuth": []string{"admin"}}},
	}, Bh.ChangeBackupInterval)

	huma.Register(backupGroup, huma.Operation{
		OperationID: "change-backup-time",
		Method:      http.MethodPost,
		Path:        "/backup/time",
		Summary:     "Change backup time",
		Description: "Update the daily time when automatic backups occur",
		Tags:        []string{"Backup"},
		Security:    []map[string][]string{{"bearerAuth": []string{"admin"}}},
	}, Bh.ChangeBackupTime)

	////////////////////
	////   Config   ////
	////////////////////

	slog.Info("Routes registered successfully")
	slog.Info("Starting server on :" + config.App.Port)

	// CORS
	corsOptions := cors.Options{
		AllowedOrigins: []string{"*"},
		AllowedMethods: []string{"GET", "POST", "OPTIONS"},
		AllowedHeaders: []string{"X-Requested-With", "Content-Type", "Authorization"},
	}

	handler := cors.New(corsOptions).Handler(router)

	if err := http.ListenAndServe(":"+config.App.Port, handler); err != nil {
		panic("Server failed: " + err.Error())
	}
}
