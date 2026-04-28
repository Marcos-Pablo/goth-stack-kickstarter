package auth

import (
	"database/sql"
	"errors"
	"log/slog"
	"net/http"
	"time"

	"github.com/Marcos-Pablo/goth-stack-kickstarter/internal/app"
	"github.com/Marcos-Pablo/goth-stack-kickstarter/internal/db"
	"github.com/Marcos-Pablo/goth-stack-kickstarter/internal/handler"
	"github.com/Marcos-Pablo/goth-stack-kickstarter/views"
	"github.com/google/uuid"
	slogchi "github.com/samber/slog-chi"
	sqlite "modernc.org/sqlite"
	sqlite3 "modernc.org/sqlite/lib"
)

type Handler struct {
	app *app.App
}

func New(app *app.App) *Handler {
	return &Handler{app: app}
}

func (h *Handler) SignInPage(w http.ResponseWriter, r *http.Request) {
	handler.Render(w, r, http.StatusOK, views.SignIn(views.SignInForm{}))
}

func (h *Handler) SignIn(w http.ResponseWriter, r *http.Request) {
	form := parseAndValidateSignInform(r)

	if !form.Valid() || form.GeneralError != "" {
		handler.Render(w, r, http.StatusBadRequest, views.SignIn(form))
		return
	}

	user, err := h.app.Queries.GetUserByEmail(r.Context(), form.Email)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			form.GeneralError = "invalid email or password"
			handler.Render(w, r, http.StatusBadRequest, views.SignIn(form))
			return
		}

		slogchi.AddCustomAttributes(r, slog.Any("error", err))
		handler.Render(w, r, http.StatusInternalServerError, views.ServerError())
		return
	}

	match, err := CheckPasswordHash(form.Password, user.Password)

	if err != nil {
		slogchi.AddCustomAttributes(r, slog.Any("error", err))
		handler.Render(w, r, http.StatusInternalServerError, views.ServerError())
		return
	}

	if !match {
		form.GeneralError = "invalid email or password"
		handler.Render(w, r, http.StatusBadRequest, views.SignIn(form))
		return
	}

	if err := h.app.Sessions.RenewToken(r.Context()); err != nil {
		slogchi.AddCustomAttributes(r, slog.Any("error", err))
		handler.Render(w, r, http.StatusInternalServerError, views.ServerError())
		return
	}

	h.app.Sessions.Put(r.Context(), "user_id", user.ID)
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (h *Handler) SignUpPage(w http.ResponseWriter, r *http.Request) {
	handler.Render(w, r, http.StatusOK, views.SignUp(views.SignUpForm{}))
}

func (h *Handler) SignUp(w http.ResponseWriter, r *http.Request) {
	form := parseAndValidateSignUpform(r)
	if !form.Valid() || form.GeneralError != "" {
		handler.Render(w, r, http.StatusBadRequest, views.SignUp(form))
		return
	}

	hashedPassword, err := HashPassword(form.Password)
	if err != nil {
		slogchi.AddCustomAttributes(r, slog.Any("error", err))
		handler.Render(w, r, http.StatusInternalServerError, views.ServerError())
		return
	}

	user, err := h.app.Queries.CreateUser(r.Context(), db.CreateUserParams{
		ID:        uuid.New().String(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Email:     form.Email,
		Name:      form.Name,
		Password:  hashedPassword,
	})

	if err != nil {
		if sqliteErr, ok := errors.AsType[*sqlite.Error](err); ok && sqliteErr.Code() == sqlite3.SQLITE_CONSTRAINT_UNIQUE {
			form.AddFieldError("email", "This email is already taken")
			handler.Render(w, r, http.StatusInternalServerError, views.SignUp(form))
			return
		}
		slogchi.AddCustomAttributes(r, slog.Any("error", err))
		handler.Render(w, r, http.StatusInternalServerError, views.ServerError())
		return
	}

	if err := h.app.Sessions.RenewToken(r.Context()); err != nil {
		slogchi.AddCustomAttributes(r, slog.Any("error", err))
		handler.Render(w, r, http.StatusInternalServerError, views.ServerError())
		return
	}

	h.app.Sessions.Put(r.Context(), "user_id", user.ID)
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (h *Handler) SignOut(w http.ResponseWriter, r *http.Request) {
	if err := h.app.Sessions.Destroy(r.Context()); err != nil {
		slogchi.AddCustomAttributes(r, slog.Any("error", err))
	}

	if r.Header.Get("HX-Request") == "true" {
		w.Header().Set("HX-Redirect", "/sign-in")
		w.WriteHeader(http.StatusSeeOther)
		return
	}

	http.Redirect(w, r, "/sign-in", http.StatusSeeOther)
}

func parseAndValidateSignUpform(r *http.Request) views.SignUpForm {
	err := r.ParseForm()

	if err != nil {
		return views.SignUpForm{
			GeneralError: "Invalid form data",
		}
	}

	form := views.SignUpForm{
		Email:                r.PostForm.Get("email"),
		Name:                 r.PostForm.Get("name"),
		Password:             r.PostForm.Get("password"),
		PasswordConfirmation: r.PostForm.Get("password_confirmation"),
	}
	form.Validate()

	return form
}

func parseAndValidateSignInform(r *http.Request) views.SignInForm {
	err := r.ParseForm()

	if err != nil {
		return views.SignInForm{
			GeneralError: "Invalid form data",
		}
	}

	form := views.SignInForm{
		Email:    r.PostForm.Get("email"),
		Password: r.PostForm.Get("password"),
	}

	form.Validate()
	return form
}
