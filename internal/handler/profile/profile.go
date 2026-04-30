package profile

import (
	"database/sql"
	"errors"
	"log/slog"
	"mime"
	"net/http"
	"time"

	"github.com/Marcos-Pablo/goth-stack-kickstarter/internal/app"
	"github.com/Marcos-Pablo/goth-stack-kickstarter/internal/db"
	"github.com/Marcos-Pablo/goth-stack-kickstarter/internal/handler"
	"github.com/Marcos-Pablo/goth-stack-kickstarter/internal/handler/auth"
	"github.com/Marcos-Pablo/goth-stack-kickstarter/internal/middleware"
	"github.com/Marcos-Pablo/goth-stack-kickstarter/views"
	slogchi "github.com/samber/slog-chi"
	"modernc.org/sqlite"
	sqlite3 "modernc.org/sqlite/lib"
)

type Handler struct {
	app *app.App
}

func New(app *app.App) *Handler {
	return &Handler{app: app}
}

func (h *Handler) ProfilePage(w http.ResponseWriter, r *http.Request) {
	user, _ := middleware.UserFrom(r.Context())

	avatarURL := ""
	if user.ProfilePictureUrl.Valid {
		avatarURL = h.app.Storage.AvatarURL(user.ProfilePictureUrl.String)
	}
	handler.Render(w, r, http.StatusOK, views.Profile(
		views.User{
			Email:          user.Email,
			Name:           user.Name,
			ProfilePicture: avatarURL,
		}, views.AvatarForm{}, views.ProfileForm{
			Email: user.Email,
			Name:  user.Name,
		}, views.ChangePasswordForm{},
	))
}

func (h *Handler) UpdateProfile(w http.ResponseWriter, r *http.Request) {
	user, _ := middleware.UserFrom(r.Context())
	form := parseAndValidateProfileForm(r)

	if !form.Valid() || form.GeneralError != "" {
		handler.Render(w, r, http.StatusOK, views.ProfilePersonalInfoCard(
			form,
			"",
		))
		return
	}

	dbUser, err := h.app.Queries.UpdatePersonalInfo(r.Context(), db.UpdatePersonalInfoParams{
		ID:        user.ID,
		Email:     form.Email,
		Name:      form.Name,
		UpdatedAt: time.Now(),
	})

	if err != nil {
		if sqliteErr, ok := errors.AsType[*sqlite.Error](err); ok && sqliteErr.Code() == sqlite3.SQLITE_CONSTRAINT_UNIQUE {
			form.AddFieldError("email", "This email is already taken")
			handler.Render(w, r, http.StatusOK, views.ProfilePersonalInfoCard(
				form,
				"",
			))

			return
		}
		slogchi.AddCustomAttributes(r, slog.Any("error", err))
		handler.Render(w, r, http.StatusInternalServerError, views.ServerError())
		return
	}

	handler.Render(w, r, http.StatusOK, views.ProfilePersonalInfoCard(
		views.ProfileForm{
			Email: dbUser.Email,
			Name:  dbUser.Name,
		},
		"Profile updated successfully",
	))
	return
}

func (h *Handler) ChangePassword(w http.ResponseWriter, r *http.Request) {
	user, _ := middleware.UserFrom(r.Context())
	form := parseAndValidateChangePasswordForm(r)

	if !form.Valid() || form.GeneralError != "" {
		handler.Render(w, r, http.StatusOK, views.ProfileSecurityCard(form, ""))
		return
	}

	match, err := auth.CheckPasswordHash(form.CurrentPassword, user.Password)

	if err != nil {
		slogchi.AddCustomAttributes(r, slog.Any("error", err))
		form.GeneralError = "Unexpected error"
		handler.Render(w, r, http.StatusOK, views.ServerError())
		return
	}

	if !match {
		form.AddFieldError("current_password", "Invalid password")
		handler.Render(w, r, http.StatusOK, views.ProfileSecurityCard(
			form,
			"",
		))
		return
	}

	hashed, err := auth.HashPassword(form.NewPassword)

	if err != nil {
		slogchi.AddCustomAttributes(r, slog.Any("error", err))
		form.GeneralError = "Unexpected error"
		handler.Render(w, r, http.StatusInternalServerError, views.ServerError())
		return
	}

	_, err = h.app.Queries.UpdatePassword(r.Context(), db.UpdatePasswordParams{
		ID:       user.ID,
		Password: hashed,
	})

	if err != nil {
		slogchi.AddCustomAttributes(r, slog.Any("error", err))
		handler.Render(w, r, http.StatusInternalServerError, views.ServerError())
		return
	}

	handler.Render(w, r, http.StatusOK, views.ProfileSecurityCard(
		views.ChangePasswordForm{},
		"Password updated successfully"),
	)
}

func (h *Handler) DeleteAccount(w http.ResponseWriter, r *http.Request) {
	user, _ := middleware.UserFrom(r.Context())

	err := h.app.Queries.DeleteUser(r.Context(), user.ID)

	if err != nil {
		slogchi.AddCustomAttributes(r, slog.Any("error", err))
		handler.Render(w, r, http.StatusInternalServerError, views.ServerError())
		return
	}

	http.Redirect(w, r, "/sign-in", http.StatusSeeOther)
}

func (h *Handler) UpdateAvatar(w http.ResponseWriter, r *http.Request) {
	const maxMemory = 5 << 20     // 5 MB
	const maxUploadSize = 5 << 20 // 5 MB

	user, _ := middleware.UserFrom(r.Context())
	http.MaxBytesReader(w, r.Body, maxUploadSize)

	err := r.ParseMultipartForm(maxMemory)

	avatarURL := ""
	if user.ProfilePictureUrl.Valid {
		avatarURL = h.app.Storage.AvatarURL(user.ProfilePictureUrl.String)
	}

	if err != nil {
		slogchi.AddCustomAttributes(r, slog.Any("error", err))
		handler.Render(w, r, http.StatusOK, views.ProfileAvatarCard(views.User{
			Name:           user.Name,
			Email:          user.Email,
			ProfilePicture: avatarURL,
		}, views.AvatarForm{
			GeneralError: "invalid form data",
		}, ""))
		return
	}

	form := views.AvatarForm{}

	file, header, err := r.FormFile("profile_picture")

	if err != nil {
		slogchi.AddCustomAttributes(r, slog.Any("error", err))
		form.GeneralError = "Error uploading file"
		handler.Render(w, r, http.StatusOK, views.ProfileAvatarCard(views.User{
			Name:           user.Name,
			Email:          user.Email,
			ProfilePicture: avatarURL,
		}, form, ""))
		return
	}
	defer file.Close()

	mediaType, _, err := mime.ParseMediaType(header.Header.Get("Content-Type"))

	if err != nil {
		slogchi.AddCustomAttributes(r, slog.Any("error", err))
		form.GeneralError = "Invalid file type"
		handler.Render(w, r, http.StatusOK, views.ProfileAvatarCard(views.User{
			Name:           user.Name,
			Email:          user.Email,
			ProfilePicture: avatarURL,
		}, form, ""))
		return
	}

	if mediaType != "image/jpeg" && mediaType != "image/png" {
		form.GeneralError = "Invalid file type"
		handler.Render(w, r, http.StatusOK, views.ProfileAvatarCard(views.User{
			Name:           user.Name,
			Email:          user.Email,
			ProfilePicture: avatarURL,
		}, form, ""))
		return
	}

	filename := h.app.Storage.NewAvatarFilename(mediaType)
	avatarPath := h.app.Storage.AvatarPath(filename)

	err = h.app.Storage.SaveFile(avatarPath, file)
	if err != nil {
		slogchi.AddCustomAttributes(r, slog.Any("error", err))
		form.GeneralError = "Error uploading file"
		handler.Render(w, r, http.StatusOK, views.ProfileAvatarCard(views.User{
			Name:           user.Name,
			Email:          user.Email,
			ProfilePicture: avatarURL,
		}, form, ""))
		return
	}

	if err := h.app.Storage.DeleteFile(user.ProfilePictureUrl.String); err != nil {
		slogchi.AddCustomAttributes(r, slog.Any("error", err))
	}

	_, err = h.app.Queries.UpdateProfilePicture(r.Context(), db.UpdateProfilePictureParams{
		ID: user.ID,
		ProfilePictureUrl: sql.NullString{
			String: avatarPath,
			Valid:  true,
		},
		UpdatedAt: time.Now(),
	})

	if err != nil {
		slogchi.AddCustomAttributes(r, slog.Any("error", err))
		form.GeneralError = "Error uploading file"
		if cleanupErr := h.app.Storage.DeleteFile(avatarPath); cleanupErr != nil {
			slogchi.AddCustomAttributes(r, slog.Any("error", cleanupErr))
		}

		handler.Render(w, r, http.StatusOK, views.ProfileAvatarCard(views.User{
			Name:           user.Name,
			Email:          user.Email,
			ProfilePicture: avatarURL,
		}, form, ""))
		return
	}

	handler.Render(w, r, http.StatusOK, views.ProfileAvatarCard(views.User{
		Name:           user.Name,
		Email:          user.Email,
		ProfilePicture: h.app.Storage.AvatarURL(avatarPath),
	}, form, "Avatar updated successfully"))
}

func parseAndValidateProfileForm(r *http.Request) views.ProfileForm {
	err := r.ParseForm()

	if err != nil {
		return views.ProfileForm{
			GeneralError: "invalid form data",
		}
	}

	form := views.ProfileForm{
		Name:  r.PostForm.Get("name"),
		Email: r.PostForm.Get("email"),
	}

	form.Validate()
	return form
}

func parseAndValidateChangePasswordForm(r *http.Request) views.ChangePasswordForm {
	err := r.ParseForm()

	if err != nil {
		return views.ChangePasswordForm{
			GeneralError: "invalid form data",
		}
	}

	form := views.ChangePasswordForm{
		CurrentPassword:         r.PostForm.Get("current_password"),
		NewPassword:             r.PostForm.Get("new_password"),
		NewPasswordConfirmation: r.PostForm.Get("new_password_confirmation"),
	}

	form.Validate()
	return form
}
