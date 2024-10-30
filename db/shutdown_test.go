package db_test

import (
	"context"
	"database/sql"
	"testing"

	"flamingo.me/flamingo/v3/framework/flamingo"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"

	"github.com/tessig/flamingo-mysql/db"
	"github.com/tessig/flamingo-mysql/db/mocks"
)

func TestShutdownSubscriber_Notify(t *testing.T) {
	t.Parallel()

	type args struct {
		event flamingo.Event
	}

	tests := []struct {
		name      string
		args      args
		wantClose bool
	}{
		{
			name: "Close on shutdown",
			args: args{
				event: &flamingo.ShutdownEvent{},
			},
			wantClose: true,
		},
		{
			name: "No close on other event",
			args: args{
				event: nil,
			},
			wantClose: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			mockDB, mock, err := sqlmock.New()
			if err != nil {
				t.Fatal(err)
			}

			defer func(mockDB *sql.DB) {
				_ = mockDB.Close()
			}(mockDB)

			sqlxDB := sqlx.NewDb(mockDB, "sqlmock")

			if tt.wantClose {
				mock.ExpectClose()
			}

			dbMock := &mocks.DB{}
			dbMock.On("Connection").Return(sqlxDB)

			s := &db.ShutdownSubscriber{}
			s.Inject(dbMock)

			s.Notify(context.Background(), tt.args.event)

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Error(err)
			}
		})
	}
}
