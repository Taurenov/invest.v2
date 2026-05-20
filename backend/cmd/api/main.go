package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/fin-helper/backend/internal/alerts"
	"github.com/fin-helper/backend/internal/cache"
	"github.com/fin-helper/backend/internal/config"
	"github.com/fin-helper/backend/internal/domain"
	"github.com/fin-helper/backend/internal/engine"
	"github.com/fin-helper/backend/internal/http/handlers"
	"github.com/fin-helper/backend/internal/http/middleware"
	"github.com/fin-helper/backend/internal/market"
	"github.com/fin-helper/backend/internal/repo"
	"github.com/fin-helper/backend/internal/ws"
	"github.com/google/uuid"
)

func main() {
	cfg := config.Load()
	ctx := context.Background()
	alertHub := alerts.NewHub()

	var (
		users       repo.UserStore
		txStore     repo.TransactionStore
		catStore    repo.CategoryStore
		goalStore   repo.GoalStore
		portStore   repo.PortfolioStore
		watchStore  repo.WatchlistStore
		settings    repo.SettingsStore
		analytics   repo.AnalyticsStore
		summaries   repo.SummaryStore
		tags        repo.TagStore
		budgets     repo.BudgetStore
		recurring   repo.RecurringStore
		pgDB        *repo.Postgres
	)

	if cfg.DatabaseURL != "" {
		pg, err := repo.NewPostgres(ctx, cfg.DatabaseURL)
		if err != nil {
			log.Fatalf("postgres: %v", err)
		}
		pgDB = pg
		users = pg
		txStore = pg
		catStore = pg
		goalStore = pg
		portStore = pg
		watchStore = pg
		settings = pg
		analytics = pg
		summaries = pg
		tags = pg
		budgets = pg
		recurring = pg
		log.Println("using PostgreSQL")
	} else {
		demoUser := uuid.MustParse("00000000-0000-4000-8000-000000000001")
		users = repo.NewMemoryUserStore()
		txStore = repo.NewMemoryStore([]domain.Transaction{
			{ID: uuid.New(), UserID: demoUser, Kind: domain.KindIncome, Amount: 85000, Currency: "RUB", Description: "Зарплата", OccurredAt: time.Now().AddDate(0, 0, -2)},
			{ID: uuid.New(), UserID: demoUser, Kind: domain.KindExpense, Amount: 3200, Currency: "RUB", Description: "Продукты", OccurredAt: time.Now().AddDate(0, 0, -5)},
		})
		catStore = repo.NewMemoryCategoryStore()
		goalStore = repo.NewMemoryGoalStore()
		portStore = repo.NewMemoryPortfolioStore()
		watchStore = repo.NewMemoryWatchlistStore()
		settings = repo.NewMemorySettingsStore()
		analytics = &repo.MemoryAnalyticsStore{Tx: txStore}
		summaries = repo.MemorySummaryStore{}
		tags = repo.NewMemoryTagStore()
		budgets = repo.NewMemoryBudgetStore()
		recurring = repo.NewMemoryRecurringStore()
		log.Println("memory mode — set DATABASE_URL for Postgres")
	}

	eng := engine.NewClient(cfg.EngineHTTP)
	jwtAuth := func(h http.Handler) http.Handler {
		return middleware.JWTOrDevToken(cfg.JWTSecret, cfg.APIToken, h)
	}

	var quoteProvider market.Provider = market.NewMOEX()
	if cfg.RedisURL != "" {
		if rdb, err := cache.New(ctx, cfg.RedisURL); err == nil {
			defer rdb.Close()
			quoteProvider = &market.Cached{Inner: quoteProvider, Cache: rdb, TTL: 45 * time.Second}
		}
	}
	quoteProvider = &fallbackProvider{primary: quoteProvider, fallback: market.Mock{}}

	mux := http.NewServeMux()
	mux.HandleFunc("GET /health", func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"status":"ok"}`))
	})

	authH := &handlers.AuthHandler{Users: users, JWTSecret: cfg.JWTSecret, TokenTTL: 24 * time.Hour * 7}
	mux.HandleFunc("POST /api/v1/auth/register", authH.Register)
	mux.HandleFunc("POST /api/v1/auth/login", authH.Login)
	mux.Handle("GET /api/v1/me", jwtAuth(http.HandlerFunc(authH.Me)))

	meTx := &handlers.MeTransactions{Store: txStore}
	mux.Handle("GET /api/v1/me/transactions", jwtAuth(http.HandlerFunc(meTx.List)))
	mux.Handle("POST /api/v1/me/transactions", jwtAuth(http.HandlerFunc(meTx.Create)))
	mux.Handle("PATCH /api/v1/me/transactions/{txID}", jwtAuth(http.HandlerFunc(meTx.Update)))
	mux.Handle("DELETE /api/v1/me/transactions/{txID}", jwtAuth(http.HandlerFunc(meTx.Delete)))

	meCat := &handlers.MeCategories{Store: catStore}
	mux.Handle("GET /api/v1/me/categories", jwtAuth(http.HandlerFunc(meCat.List)))
	mux.Handle("POST /api/v1/me/categories", jwtAuth(http.HandlerFunc(meCat.Create)))

	meGoals := &handlers.MeGoals{Store: goalStore, Alerts: alertHub}
	mux.Handle("GET /api/v1/me/goals", jwtAuth(http.HandlerFunc(meGoals.List)))
	mux.Handle("POST /api/v1/me/goals", jwtAuth(http.HandlerFunc(meGoals.Create)))
	mux.Handle("POST /api/v1/me/goals/{goalID}/contribute", jwtAuth(http.HandlerFunc(meGoals.Contribute)))

	calcH := &handlers.CalculatorHandler{Engine: eng}
	mux.Handle("POST /api/v1/calculator/roi", jwtAuth(http.HandlerFunc(calcH.ROI)))
	mux.Handle("POST /api/v1/calculator/cagr", jwtAuth(http.HandlerFunc(calcH.CAGR)))
	mux.Handle("POST /api/v1/calculator/savings", jwtAuth(http.HandlerFunc(calcH.Savings)))

	portH := &handlers.PortfolioHandler{Store: portStore, Quotes: quoteProvider}
	mux.Handle("GET /api/v1/me/portfolio", jwtAuth(http.HandlerFunc(portH.Get)))
	mux.Handle("POST /api/v1/me/portfolio/holdings", jwtAuth(http.HandlerFunc(portH.Add)))
	mux.Handle("DELETE /api/v1/me/portfolio/holdings/{holdingID}", jwtAuth(http.HandlerFunc(portH.Remove)))

	watchH := &handlers.WatchlistHandler{Store: watchStore}
	mux.Handle("GET /api/v1/me/watchlist", jwtAuth(http.HandlerFunc(watchH.List)))
	mux.Handle("POST /api/v1/me/watchlist/items", jwtAuth(http.HandlerFunc(watchH.Add)))
	mux.Handle("DELETE /api/v1/me/watchlist/items/{instrumentID}", jwtAuth(http.HandlerFunc(watchH.Remove)))

	setH := &handlers.SettingsHandler{Store: settings}
	mux.Handle("GET /api/v1/me/settings", jwtAuth(http.HandlerFunc(setH.Get)))
	mux.Handle("PATCH /api/v1/me/settings", jwtAuth(http.HandlerFunc(setH.Patch)))

	anH := &handlers.AnalyticsHandler{Store: analytics, Tx: txStore}
	mux.Handle("GET /api/v1/me/analytics", jwtAuth(http.HandlerFunc(anH.Report)))
	mux.Handle("GET /api/v1/me/analytics/export.csv", jwtAuth(http.HandlerFunc(anH.ExportCSV)))

	alertH := &handlers.AlertsHandler{Hub: alertHub}
	mux.Handle("GET /api/v1/me/alerts", jwtAuth(http.HandlerFunc(alertH.List)))
	mux.Handle("POST /api/v1/me/alerts/read", jwtAuth(http.HandlerFunc(alertH.MarkRead)))

	quoteH := &handlers.QuotesHandler{Provider: quoteProvider}
	mux.Handle("GET /api/v1/markets/{symbol}/quote", jwtAuth(http.HandlerFunc(quoteH.Get)))

	sumH := &handlers.SummaryHandler{Store: summaries}
	mux.Handle("GET /api/v1/markets/{symbol}/summary", jwtAuth(http.HandlerFunc(sumH.Get)))

	// Priority 5
	tagsH := &handlers.TagsHandler{Store: tags}
	mux.Handle("GET /api/v1/me/tags", jwtAuth(http.HandlerFunc(tagsH.List)))
	mux.Handle("POST /api/v1/me/tags", jwtAuth(http.HandlerFunc(tagsH.Create)))
	mux.Handle("DELETE /api/v1/me/tags/{tagID}", jwtAuth(http.HandlerFunc(tagsH.Delete)))
	mux.Handle("PUT /api/v1/me/transactions/{txID}/tags", jwtAuth(http.HandlerFunc(tagsH.SetTxTags)))

	budH := &handlers.BudgetsHandler{Store: budgets, Alerts: alertHub}
	mux.Handle("GET /api/v1/me/budgets", jwtAuth(http.HandlerFunc(budH.List)))
	mux.Handle("POST /api/v1/me/budgets", jwtAuth(http.HandlerFunc(budH.Upsert)))
	mux.Handle("DELETE /api/v1/me/budgets/{budgetID}", jwtAuth(http.HandlerFunc(budH.Delete)))

	recH := &handlers.RecurringHandler{Store: recurring}
	mux.Handle("GET /api/v1/me/recurring", jwtAuth(http.HandlerFunc(recH.List)))
	mux.Handle("POST /api/v1/me/recurring", jwtAuth(http.HandlerFunc(recH.Create)))
	mux.Handle("POST /api/v1/me/recurring/{recurringID}/toggle", jwtAuth(http.HandlerFunc(recH.Toggle)))
	mux.Handle("DELETE /api/v1/me/recurring/{recurringID}", jwtAuth(http.HandlerFunc(recH.Delete)))

	if pgDB != nil {
		fcH := &handlers.ForecastHandler{DB: pgDB, Engine: eng}
		mux.Handle("GET /api/v1/markets/{symbol}/forecast", jwtAuth(http.HandlerFunc(fcH.Predict)))
		mux.Handle("GET /api/v1/markets/{symbol}/forecast/history", jwtAuth(http.HandlerFunc(fcH.History)))
		mux.Handle("GET /api/v1/markets/{symbol}/history", jwtAuth(http.HandlerFunc(fcH.PriceHistory)))
	}

	hub := ws.NewHub(quoteProvider)
	mux.Handle("GET /ws/quotes", jwtAuth(hub))

	// recurring runner (simple ticker)
	go func() {
		t := time.NewTicker(1 * time.Minute)
		defer t.Stop()
		for range t.C {
			if recurring != nil {
				n, _ := recurring.RunRecurringDue(context.Background(), time.Now())
				if n > 0 {
					log.Printf("recurring applied: %d", n)
				}
			}
		}
	}()

	go func() {
		ticker := time.NewTicker(30 * time.Second)
		defer ticker.Stop()
		for range ticker.C {
			if watchStore != nil {
				items, err := watchStore.ListWatchlist(context.Background(), uuid.MustParse("00000000-0000-4000-8000-000000000001"))
				if err == nil && len(items) > 0 {
					var syms []struct{ Symbol, Exchange string }
					for _, it := range items {
						syms = append(syms, struct{ Symbol, Exchange string }{it.Symbol, it.Exchange})
					}
					hub.SetSymbols(syms)
				}
			}
		}
	}()

	root := http.Handler(mux)
	if cfg.AllowCORS {
		root = middleware.CORS(root)
	}

	srv := &http.Server{Addr: cfg.APIAddr, Handler: root}
	go hub.RunBroadcast(context.Background(), 3*time.Second)

	go func() {
		log.Printf("api on %s", cfg.APIAddr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal(err)
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	<-stop
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	_ = srv.Shutdown(shutdownCtx)
	if pgDB != nil {
		pgDB.Close()
	}
}

type fallbackProvider struct {
	primary  market.Provider
	fallback market.Provider
}

func (f *fallbackProvider) Quote(ctx context.Context, symbol, exchange string) (*domain.Quote, error) {
	q, err := f.primary.Quote(ctx, symbol, exchange)
	if err == nil {
		return q, nil
	}
	return f.fallback.Quote(ctx, symbol, exchange)
}
