package handlers

import (
	"net/http"

	"backend/store"
	"backend/utils"
)

func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		utils.WriteJSON(w, http.StatusMethodNotAllowed, map[string]string{"error": "only POST allowed"}) //[cite: 5]
		return
	}

	var req struct {
		Name    string `json:"name"`
		Email   string `json:"email"`
		IsAdmin bool   `json:"is_admin,omitempty"`
	}
	if err := utils.ReadJSONBody(r, &req); err != nil {
		utils.WriteJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()}) //[cite: 5]
		return
	}
	if req.Name == "" || req.Email == "" {
		utils.WriteJSON(w, http.StatusBadRequest, map[string]string{"error": "name and email required"}) //[cite: 5]
		return
	}

	uid, _ := utils.GenRandomHex(8)     //[cite: 5]
	secret, _ := utils.GenRandomHex(16) //[cite: 5]
	isAdmin := false
	if req.IsAdmin {
		if r.Header.Get("X-Master-Key") == store.MASTER_ADMIN_KEY { //[cite: 5]
			isAdmin = true
		} else {
			utils.WriteJSON(w, http.StatusForbidden, map[string]string{"error": "invalid master key"}) //[cite: 5]
			return
		}
	}

	_, err := store.DB.Exec("INSERT INTO users (id, secret_code, name, email, is_admin) VALUES (?, ?, ?, ?, ?)", uid, secret, req.Name, req.Email, isAdmin)
	if err != nil {
		utils.WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": "failed to create user"})
		return
	}

	user := &store.User{
		ID:         uid,
		SecretCode: secret,
		Name:       req.Name,
		Email:      req.Email,
		Complaints: []string{}, //[cite: 5]
		IsAdmin:    isAdmin,
	}

	utils.WriteJSON(w, http.StatusCreated, user) //[cite: 5]
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Secret string `json:"secret_code"`
	}
	if err := utils.ReadJSONBody(r, &req); err != nil {
		utils.WriteJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()}) //[cite: 5]
		return
	}
	user, err := store.GetUserBySecret(req.Secret)
	if err != nil {
		utils.WriteJSON(w, http.StatusUnauthorized, map[string]string{"error": "invalid secret code"}) //[cite: 5]
		return
	}
	utils.WriteJSON(w, http.StatusOK, user) //[cite: 5]
}

func SubmitComplaintHandler(w http.ResponseWriter, r *http.Request) {
	secret := r.Header.Get("X-Secret-Code") //[cite: 5]
	user, err := store.GetUserBySecret(secret)
	if err != nil {
		utils.WriteJSON(w, http.StatusUnauthorized, map[string]string{"error": "invalid secret code"}) //[cite: 5]
		return
	}

	var req struct {
		Title    string `json:"title"`
		Summary  string `json:"summary"`
		Severity int    `json:"severity"`
	}
	if err := utils.ReadJSONBody(r, &req); err != nil {
		utils.WriteJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()}) //[cite: 5]
		return
	}

	cid, _ := utils.GenRandomHex(8) //[cite: 5]
	_, err = store.DB.Exec("INSERT INTO complaints (id, title, summary, severity, resolved, user_id) VALUES (?, ?, ?, ?, ?, ?)", cid, req.Title, req.Summary, req.Severity, false, user.ID)

	if err != nil {
		utils.WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": "failed to create complaint"})
		return
	}

	c := &store.Complaint{
		ID:       cid,
		Title:    req.Title,
		Summary:  req.Summary,
		Severity: req.Severity,
		UserID:   user.ID,
		Resolved: false, //[cite: 5]
	}

	utils.WriteJSON(w, http.StatusCreated, c) //[cite: 5]
}

func GetAllForUserHandler(w http.ResponseWriter, r *http.Request) {
	secret := r.Header.Get("X-Secret-Code") //[cite: 5]
	user, err := store.GetUserBySecret(secret)
	if err != nil {
		utils.WriteJSON(w, http.StatusUnauthorized, map[string]string{"error": "invalid secret code"}) //[cite: 5]
		return
	}

	type Brief struct {
		ID    string `json:"id"`
		Title string `json:"title"`
	} //[cite: 5]

	rows, err := store.DB.Query("SELECT id, title FROM complaints WHERE user_id = ?", user.ID)
	out := []Brief{} //[cite: 5]
	if err == nil {
		defer rows.Close()
		for rows.Next() {
			var b Brief
			rows.Scan(&b.ID, &b.Title)
			out = append(out, b)
		}
	}
	utils.WriteJSON(w, http.StatusOK, out) //[cite: 5]
}

func GetAllForAdminHandler(w http.ResponseWriter, r *http.Request) {
	secret := r.Header.Get("X-Secret-Code") //[cite: 5]
	admin, err := store.GetUserBySecret(secret)
	if err != nil || !admin.IsAdmin {
		utils.WriteJSON(w, http.StatusForbidden, map[string]string{"error": "admin only"}) //[cite: 5]
		return
	}

	type Entry struct {
		ID       string `json:"id"`
		Title    string `json:"title"`
		UserName string `json:"user_name"`
		Resolved bool   `json:"resolved"`
	} //[cite: 5]

	rows, err := store.DB.Query("SELECT c.id, c.title, u.name, c.resolved FROM complaints c JOIN users u ON c.user_id = u.id")
	out := []Entry{} //[cite: 5]
	if err == nil {
		defer rows.Close()
		for rows.Next() {
			var e Entry
			rows.Scan(&e.ID, &e.Title, &e.UserName, &e.Resolved)
			out = append(out, e)
		}
	}
	utils.WriteJSON(w, http.StatusOK, out) //[cite: 5]
}

func ViewComplaintHandler(w http.ResponseWriter, r *http.Request) {
	secret := r.Header.Get("X-Secret-Code") //[cite: 5]
	user, err := store.GetUserBySecret(secret)
	if err != nil {
		utils.WriteJSON(w, http.StatusUnauthorized, map[string]string{"error": "invalid secret code"}) //[cite: 5]
		return
	}

	id := r.URL.Query().Get("id") //[cite: 5]
	c := &store.Complaint{}

	err = store.DB.QueryRow("SELECT id, title, summary, severity, resolved, user_id FROM complaints WHERE id = ?", id).Scan(&c.ID, &c.Title, &c.Summary, &c.Severity, &c.Resolved, &c.UserID)
	if err != nil {
		utils.WriteJSON(w, http.StatusNotFound, map[string]string{"error": "complaint not found"}) //[cite: 5]
		return
	}

	if user.IsAdmin || c.UserID == user.ID { //[cite: 5]
		utils.WriteJSON(w, http.StatusOK, c) //[cite: 5]
		return
	}
	utils.WriteJSON(w, http.StatusForbidden, map[string]string{"error": "not authorized"}) //[cite: 5]
}

func ResolveComplaintHandler(w http.ResponseWriter, r *http.Request) {
	secret := r.Header.Get("X-Secret-Code") //[cite: 5]
	admin, err := store.GetUserBySecret(secret)
	if err != nil || !admin.IsAdmin {
		utils.WriteJSON(w, http.StatusForbidden, map[string]string{"error": "admin only"}) //[cite: 5]
		return
	}

	var req struct {
		ID string `json:"id"`
	}
	if err := utils.ReadJSONBody(r, &req); err != nil {
		utils.WriteJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()}) //[cite: 5]
		return
	}

	_, err = store.DB.Exec("UPDATE complaints SET resolved = true WHERE id = ?", req.ID)
	if err != nil {
		utils.WriteJSON(w, http.StatusNotFound, map[string]string{"error": "complaint not found"}) //[cite: 5]
		return
	}
	utils.WriteJSON(w, http.StatusOK, map[string]string{"status": "resolved", "id": req.ID}) //[cite: 5]
}
