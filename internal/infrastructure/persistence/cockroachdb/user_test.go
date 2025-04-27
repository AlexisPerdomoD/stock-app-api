package cockroachdb_test

import (
	"context"
	"testing"

	"github.com/alexisPerdomoD/stock-app-api/internal/domain"
	"github.com/alexisPerdomoD/stock-app-api/internal/infrastructure/persistence/cockroachdb"
	"github.com/alexisPerdomoD/stock-app-api/internal/infrastructure/service/mock"
	"github.com/alexisPerdomoD/stock-app-api/internal/pkg"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func Test_userRepository_Create(t *testing.T) {
	usr1 := &domain.User{
		UserName: mock.RandomString(12) + "@email.com",
		Password: "1234Abcd",
		Active:   true,
	}

	assert := assert.New(t)
	db := cockroachdb.NewDB()
	tests := []struct {
		name    string
		db      *gorm.DB
		usr     *domain.User
		wantErr bool
	}{
		{
			name:    "should create user",
			db:      db,
			usr:     usr1,
			wantErr: false,
		},
		{
			name:    "should fail to create user with used username",
			db:      db,
			usr:     usr1,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := cockroachdb.NewUserRepository(tt.db)
			gotErr := r.Create(context.Background(), tt.usr)
			if tt.wantErr {
				assert.Error(gotErr)
				return
			}

			assert.Nil(gotErr)
			assert.NotEqual(tt.usr.ID, 0)

		})
	}
}

func Test_userRepository_Get(t *testing.T) {
	usr3 := &domain.User{
		UserName: mock.RandomString(12) + "@email.com",
		Password: "1234Abcd",
		Active:   true,
	}

	db := cockroachdb.NewDB()
	assert := assert.New(t)
	ur := cockroachdb.NewUserRepository(db)

	if err := ur.Create(context.Background(), usr3); err != nil {
		t.Fatal(err)
	}

	tests := []struct {
		name    string
		db      *gorm.DB
		id      uint
		want    *domain.User
		wantErr bool
	}{
		{
			name:    "should return nil if user not found",
			db:      db,
			id:      0,
			want:    nil,
			wantErr: false,
		},
		{
			name:    "should return user if found",
			db:      db,
			id:      usr3.ID,
			want:    usr3,
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := cockroachdb.NewUserRepository(tt.db)
			got, gotErr := r.Get(context.Background(), tt.id)
			if tt.wantErr {
				assert.Error(gotErr)
				return
			}

			assert.Nil(gotErr)
			assert.Equal(tt.want, got)
		})
	}
}

func Test_userRepository_GetByUsername(t *testing.T) {
	usr4 := &domain.User{
		UserName: mock.RandomString(12) + "@email.com",
		Password: "1234Abcd",
		Active:   true,
	}

	db := cockroachdb.NewDB()
	assert := assert.New(t)
	ur := cockroachdb.NewUserRepository(db)

	if err := ur.Create(context.Background(), usr4); err != nil {
		t.Fatal(err)
	}

	tests := []struct {
		name     string
		db       *gorm.DB
		username string
		want     *domain.User
		wantErr  bool
	}{
		{
			name:     "should return nil if user not found",
			db:       db,
			username: "no_registered@email.com",
			want:     nil,
			wantErr:  false,
		},
		{
			name:     "should return user if found",
			db:       db,
			username: usr4.UserName,
			want:     usr4,
			wantErr:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := cockroachdb.NewUserRepository(tt.db)
			got, gotErr := r.GetByUsername(context.Background(), tt.username)
			if tt.wantErr {
				assert.Error(gotErr)
				return
			}

			assert.Nil(gotErr)
			assert.Equal(tt.want, got)

		})
	}
}

func Test_userRepository_RegisterUserStock(t *testing.T) {
	assert := assert.New(t)
	db := cockroachdb.NewDB()
	sr := cockroachdb.NewStockRepository(db)
	sss := mock.MockSourceStockService{}

	usr5 := &domain.User{
		UserName: mock.RandomString(12) + "@email.com",
		Password: "1234Abcd",
		Active:   true,
	}

	ur := cockroachdb.NewUserRepository(db)

	if err := ur.Create(context.Background(), usr5); err != nil {
		t.Fatal(err)
	}

	data, err := sss.Get(context.Background(), nil)
	if err != nil {
		t.Fatal(err)
	}

	if err := sr.Register(context.Background(), data); err != nil {
		t.Fatal(err)
	}

	stocks, err := sr.GetAllPaginated(context.Background(), pkg.PaginationFilter{}, nil)
	if err != nil {
		t.Fatal(err)
	}

	if len(stocks.Items) < 1 {
		t.Fatal("no enough stocks found")
	}

	tests := []struct {
		name    string
		db      *gorm.DB
		userID  uint
		stockID uint
		wantErr bool
	}{
		{
			name:    "should fail with invalid userID",
			db:      db,
			userID:  0,
			stockID: stocks.Items[0].ID,
			wantErr: true,
		}, {
			name:    "should fail with invalid stockID",
			db:      db,
			userID:  usr5.ID,
			stockID: 0,
			wantErr: true,
		}, {
			name:    "should succeed with valid userID and stockID",
			db:      db,
			userID:  usr5.ID,
			stockID: stocks.Items[0].ID,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := cockroachdb.NewUserRepository(tt.db)
			gotErr := r.RegisterUserStock(context.Background(), tt.userID, tt.stockID)
			if tt.wantErr {
				assert.Error(gotErr)
				return
			}

			assert.Nil(gotErr)
		})
	}
}

func Test_userRepository_RemoveUserStock(t *testing.T) {

	assert := assert.New(t)
	db := cockroachdb.NewDB()
	sr := cockroachdb.NewStockRepository(db)

	usr6 := &domain.User{
		UserName: mock.RandomString(12) + "@email.com",
		Password: "1234Abcd",
		Active:   true,
	}

	ur := cockroachdb.NewUserRepository(db)

	if err := ur.Create(context.Background(), usr6); err != nil {
		t.Fatal(err)
	}

	stocks, err := sr.GetAllPaginated(context.Background(), pkg.PaginationFilter{}, nil)
	if err != nil {
		t.Fatal(err)
	}

	if len(stocks.Items) < 1 {
		t.Fatal("no enough stocks found")
	}

	tests := []struct {
		name    string
		db      *gorm.DB
		userID  uint
		stockID uint
		wantErr bool
	}{
		{
			name:    "should fail with invalid userID",
			db:      db,
			userID:  0,
			stockID: usr6.ID,
			wantErr: true,
		},
		{
			name:    "should succeed with valid userID and stockID",
			db:      db,
			userID:  usr6.ID,
			stockID: stocks.Items[0].ID,
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := cockroachdb.NewUserRepository(tt.db)
			gotErr := r.RemoveUserStock(context.Background(), tt.userID, tt.stockID)
			if tt.wantErr {
				assert.Error(gotErr)
				return
			}

			assert.Nil(gotErr)
		})
	}

}
