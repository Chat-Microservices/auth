package pg

import (
	"context"
	"fmt"
	"github.com/georgysavva/scany/pgxscan"
	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/semho/chat-microservices/auth/internal/client/db"
	"github.com/semho/chat-microservices/auth/internal/client/db/prettier"
	"log"
)

type pg struct {
	dbc *pgxpool.Pool
}

func NewDB(dbc *pgxpool.Pool) db.DB {
	return &pg{
		dbc: dbc,
	}
}

func (p *pg) ScanAllContext(ctx context.Context, dest any, q db.Query, args ...any) error {
	logQuery(ctx, q, args...)
	rows, err := p.QueryContext(ctx, q, args...)
	if err != nil {
		return err
	}
	return pgxscan.ScanAll(dest, rows)
}

func (p *pg) ScanOneContext(ctx context.Context, dest any, q db.Query, args ...any) error {
	logQuery(ctx, q, args...)
	row, err := p.QueryContext(ctx, q, args...)
	if err != nil {
		return err
	}
	return pgxscan.ScanOne(dest, row)
}

func (p *pg) ExecContext(ctx context.Context, q db.Query, args ...any) (pgconn.CommandTag, error) {
	logQuery(ctx, q, args...)
	return p.dbc.Exec(ctx, q.QueryRow, args...)
}

func (p *pg) QueryContext(ctx context.Context, q db.Query, args ...any) (pgx.Rows, error) {
	logQuery(ctx, q, args...)

	return p.dbc.Query(ctx, q.QueryRow, args...)
}

func (p *pg) QueryRowContext(ctx context.Context, q db.Query, args ...any) pgx.Row {
	logQuery(ctx, q, args...)

	return p.dbc.QueryRow(ctx, q.QueryRow, args...)
}

func (p *pg) Ping(ctx context.Context) error {
	return p.dbc.Ping(ctx)
}

func (p *pg) Close() {
	p.dbc.Close()
}

func logQuery(ctx context.Context, q db.Query, args ...any) {
	prettyQuery := prettier.Pretty(q.QueryRow, prettier.PlaceholderDollar, args...)
	log.Println(ctx, fmt.Sprintf("sql: %s", q.Name), fmt.Sprintf("query: %s", prettyQuery))
}
