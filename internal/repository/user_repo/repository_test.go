package userrepo

import (
	"database/sql"
	"errors"
	"reflect"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/Kaushik1766/LibraryManagement/internal/models"
	"github.com/Kaushik1766/LibraryManagement/internal/models/enums/roles"
	"github.com/google/uuid"
)

func TestNewUserRepository(t *testing.T) {
	db, _, _ := sqlmock.New()
	defer db.Close()

	type args struct {
		db *sql.DB
	}
	tests := []struct {
		name      string
		args      args
		want      *UserRepository
		mockSetup func()
	}{
		{
			name: "valid",
			args: args{
				db: db,
			},
			want: &UserRepository{db: db},
			mockSetup: func() {

			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockSetup()
			if got := NewUserRepository(tt.args.db); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewUserRepository() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUserRepository_AddUser(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()

	type fields struct {
		db *sql.DB
	}
	type args struct {
		name     string
		email    string
		password string
	}
	tests := []struct {
		name      string
		fields    fields
		args      args
		wantErr   bool
		mockSetup func()
	}{
		{
			name: "valid add user",
			fields: fields{
				db: db,
			},
			args: args{
				name:     "kaushik",
				email:    "kaushik@a.com",
				password: "123",
			},
			wantErr: false,
			mockSetup: func() {
				mock.ExpectExec("(?i)insert into users.*").WillReturnResult(sqlmock.NewResult(1, 1))
			},
		},
		{
			name: "database error",
			fields: fields{
				db: db,
			},
			args: args{
				name:     "kaushik",
				email:    "kaushik@a.com",
				password: "123",
			},
			wantErr: true,
			mockSetup: func() {
				mock.ExpectExec("(?i)insert into users.*").WillReturnError(errors.New("database error"))
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := &UserRepository{
				db: tt.fields.db,
			}
			tt.mockSetup()
			if err := u.AddUser(tt.args.name, tt.args.email, tt.args.password); (err != nil) != tt.wantErr {
				t.Errorf("UserRepository.AddUser() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestUserRepository_GetUserByEmail(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()

	user1 := models.User{
		ID:       uuid.New(),
		Name:     "kaushik",
		Email:    "kaushik@a.com",
		Password: "123",
		Role:     roles.Customer,
	}

	type fields struct {
		db *sql.DB
	}
	type args struct {
		email string
	}
	tests := []struct {
		name      string
		fields    fields
		args      args
		want      models.User
		wantErr   bool
		mockSetup func()
	}{
		{
			name: "valid get user",
			fields: fields{
				db: db,
			},
			args: args{
				email: "kaushik@a.com",
			},
			want:    user1,
			wantErr: false,
			mockSetup: func() {
				mock.ExpectQuery("(?i)select .* from users .*").WillReturnRows(sqlmock.NewRows([]string{"id", "name", "email", "password", "role"}).AddRow(user1.ID, user1.Name, user1.Email, user1.Password, user1.Role))
			},
		},
		{
			name: "user not found",
			fields: fields{
				db: db,
			},
			args: args{
				email: "nonexistent@example.com",
			},
			want:    models.User{},
			wantErr: true,
			mockSetup: func() {
				mock.ExpectQuery("(?i)select .* from users .*").WillReturnError(sql.ErrNoRows)
			},
		},
		{
			name: "database error",
			fields: fields{
				db: db,
			},
			args: args{
				email: "error@example.com",
			},
			want:    models.User{},
			wantErr: true,
			mockSetup: func() {
				mock.ExpectQuery("(?i)select .* from users .*").WillReturnError(errors.New("database error"))
			},
		},
		{
			name: "scan error",
			fields: fields{
				db: db,
			},
			args: args{
				email: "scanerror@example.com",
			},
			want:    models.User{},
			wantErr: true,
			mockSetup: func() {
				mock.ExpectQuery("(?i)select .* from users .*").WillReturnRows(sqlmock.NewRows([]string{"id", "name", "email", "password", "role"}).AddRow("invalid-uuid", user1.Name, user1.Email, user1.Password, user1.Role))
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := &UserRepository{
				db: tt.fields.db,
			}
			tt.mockSetup()
			got, err := u.GetUserByEmail(tt.args.email)
			if (err != nil) != tt.wantErr {
				t.Errorf("UserRepository.GetUserByEmail() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("UserRepository.GetUserByEmail() = %v, want %v", got, tt.want)
			}
		})
	}
}
