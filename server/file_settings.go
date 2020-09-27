package server

import (
	"github.com/knoebber/dotfile/db"
	"github.com/knoebber/dotfile/usererror"
	"github.com/pkg/errors"
	"net/http"
)

func updateFile(w http.ResponseWriter, r *http.Request, p *Page) (done bool) {
	currentAlias := p.Vars["alias"]
	alias := r.Form.Get("alias")
	path := r.Form.Get("path")

	record, err := db.File(db.Connection, p.Session.Username, currentAlias)
	if err != nil {
		return p.setError(w, err)
	}

	if err := record.Update(db.Connection, alias, path); err != nil {
		return p.setError(w, err)
	}

	http.Redirect(w, r, "/"+p.Session.Username+"/"+alias, http.StatusSeeOther)
	return true
}

func clearFile(w http.ResponseWriter, r *http.Request, p *Page) (done bool) {
	username := p.Vars["username"]
	alias := p.Vars["alias"]
	tx, err := db.Connection.Begin()
	if err != nil {
		return p.setError(w, errors.Wrap(err, "starting transaction for clear file"))
	}

	if err := db.ClearCommits(tx, username, alias); err != nil {
		return p.setError(w, db.Rollback(tx, err))
	}
	if err := tx.Commit(); err != nil {
		return p.setError(w, errors.Wrap(err, "commiting transaction for clear file"))
	}

	http.Redirect(w, r, "/"+p.Session.Username+"/"+alias+"/commits", http.StatusSeeOther)
	return true
}

func deleteFile(w http.ResponseWriter, r *http.Request, p *Page) (done bool) {
	delete := r.Form.Get("delete")
	username := p.Vars["username"]
	alias := p.Vars["alias"]
	if alias != delete {
		return p.setError(w, usererror.Invalid("Alias does not match"))
	}

	tx, err := db.Connection.Begin()
	if err != nil {
		return p.setError(w, errors.Wrap(err, "starting transaction for delete file"))
	}

	if err := db.DeleteFile(tx, username, alias); err != nil {
		return p.setError(w, db.Rollback(tx, err))
	}
	if err := tx.Commit(); err != nil {
		return p.setError(w, errors.Wrap(err, "commiting transaction for delete file"))
	}

	http.Redirect(w, r, "/"+p.Session.Username, http.StatusSeeOther)
	return true

}

func loadFile(w http.ResponseWriter, r *http.Request, p *Page) (done bool) {
	username := p.Vars["username"]
	alias := p.Vars["alias"]

	file, err := db.File(db.Connection, username, alias)
	if err != nil {
		return p.setError(w, err)
	}

	p.Data["path"] = file.Path
	return
}

func fileSettingsHandler() http.HandlerFunc {
	return createHandler(&pageDescription{
		templateName: "file_settings.tmpl",
		title:        "settings",
		protected:    true,
	})
}

func updateFileHandler() http.HandlerFunc {
	return createHandler(&pageDescription{
		templateName: "update_file.tmpl",
		title:        "update",
		loadData:     loadFile,
		handleForm:   updateFile,
		protected:    true,
	})
}

func clearFileHandler() http.HandlerFunc {
	return createHandler(&pageDescription{
		templateName: "delete_commits.tmpl",
		title:        "clear",
		handleForm:   clearFile,
		protected:    true,
	})
}
func deleteFileHandler() http.HandlerFunc {
	return createHandler(&pageDescription{
		templateName: "delete_file.tmpl",
		title:        "delete",
		handleForm:   deleteFile,
		protected:    true,
	})
}
