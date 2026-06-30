package handlers

import (
	"net/http"

	"backend/store"
	"backend/utils"
)


func WelcomeHandler(w http.ResponseWriter, r *http.Request) {
	
	if r.URL.Path != "/" {
		utils.WriteJSON(w, http.StatusNotFound, map[string]string{"error": "404 page not found"})
		return
	}
	
	
	utils.WriteJSON(w, http.StatusOK, map[string]string{
		"status":  "success",
		"message": "Complaint Portal Backend is running perfectly! 🚀",
	})
}

func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		utils.WriteJSON(w, http.StatusMethodNotAllowed, map[string]string{"error": "only POST allowed"}) 
		return
	}

	var req struct {
		Name    string `json:"name"`
		Email   string `json:"email"`
		IsAdmin bool   `json:"is_admin,omitempty"`
	}
	if err := utils.ReadJSONBody(r, &req); err != nil {
		utils.WriteJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()}) 
		return
	}
	if req.Name == "" || req.Email == "" {
		utils.WriteJSON(w, http.StatusBadRequest, map[string]string{"error": "name and email required"}) 
		return
	}

	uid, _ := utils.GenRandomHex(8)     
	secret, _ := utils.GenRandomHex(16) 
	isAdmin := false
	if req.IsAdmin {
		if r.Header.Get("X-Master-Key") == store.MASTER_ADMIN_KEY { 
			isAdmin = true
		} else {
			utils.WriteJSON(w, http.StatusForbidden, map[string]string{"error": "invalid master key"}) 
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
		Complaints: []string{}, 
		IsAdmin:    isAdmin,
	}

	utils.WriteJSON(w, http.StatusCreated, user) 
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Secret string `json:"secret_code"`
	}
	if err := utils.ReadJSONBody(r, &req); err != nil {
		utils.WriteJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()}) 
		return
	}
	user, err := store.GetUserBySecret(req.Secret)
	if err != nil {
		utils.WriteJSON(w, http.StatusUnauthorized, map[string]string{"error": "invalid secret code"}) 
		return
	}
	utils.WriteJSON(w, http.StatusOK, user) 
}

func SubmitComplaintHandler(w http.ResponseWriter, r *http.Request) {
	secret := r.Header.Get("X-Secret-Code") 
	user, err := store.GetUserBySecret(secret)
	if err != nil {
		utils.WriteJSON(w, http.StatusUnauthorized, map[string]string{"error": "invalid secret code"}) 
		return
	}

	var req struct {
		Title    string `json:"title"`
		Summary  string `json:"summary"`
		Severity int    `json:"severity"`
	}
	if err := utils.ReadJSONBody(r, &req); err != nil {
		utils.WriteJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()}) 
		return
	}

	cid, _ := utils.GenRandomHex(8) 
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
		Resolved: false, 
	}

	utils.WriteJSON(w, http.StatusCreated, c) 
}

func GetAllForUserHandler(w http.ResponseWriter, r *http.Request) {
	secret := r.Header.Get("X-Secret-Code") 
	user, err := store.GetUserBySecret(secret)
	if err != nil {
		utils.WriteJSON(w, http.StatusUnauthorized, map[string]string{"error": "invalid secret code"}) 
		return
	}

	type Brief struct {
		ID    string `json:"id"`
		Title string `json:"title"`
	} 

	rows, err := store.DB.Query("SELECT id, title FROM complaints WHERE user_id = ?", user.ID)
	out := []Brief{} 
	if err == nil {
		defer rows.Close()
		for rows.Next() {
			var b Brief
			rows.Scan(&b.ID, &b.Title)
			out = append(out, b)
		}
	}
	utils.WriteJSON(w, http.StatusOK, out) 
}

func GetAllForAdminHandler(w http.ResponseWriter, r *http.Request) {
	secret := r.Header.Get("X-Secret-Code") 
	admin, err := store.GetUserBySecret(secret)
	if err != nil || !admin.IsAdmin {
		utils.WriteJSON(w, http.StatusForbidden, map[string]string{"error": "admin only"}) 
		return
	}

	type Entry struct {
		ID       string `json:"id"`
		Title    string `json:"title"`
		UserName string `json:"user_name"`
		Resolved bool   `json:"resolved"`
	} 

	rows, err := store.DB.Query("SELECT c.id, c.title, u.name, c.resolved FROM complaints c JOIN users u ON c.user_id = u.id")
	out := []Entry{} 
	if err == nil {
		defer rows.Close()
		for rows.Next() {
			var e Entry
			rows.Scan(&e.ID, &e.Title, &e.UserName, &e.Resolved)
			out = append(out, e)
		}
	}
	utils.WriteJSON(w, http.StatusOK, out) 
}

func ViewComplaintHandler(w http.ResponseWriter, r *http.Request) {
	secret := r.Header.Get("X-Secret-Code") 
	user, err := store.GetUserBySecret(secret)
	if err != nil {
		utils.WriteJSON(w, http.StatusUnauthorized, map[string]string{"error": "invalid secret code"}) 
		return
	}

	id := r.URL.Query().Get("id") 
	c := &store.Complaint{}

	err = store.DB.QueryRow("SELECT id, title, summary, severity, resolved, user_id FROM complaints WHERE id = ?", id).Scan(&c.ID, &c.Title, &c.Summary, &c.Severity, &c.Resolved, &c.UserID)
	if err != nil {
		utils.WriteJSON(w, http.StatusNotFound, map[string]string{"error": "complaint not found"}) 
		return
	}

	if user.IsAdmin || c.UserID == user.ID { 
		utils.WriteJSON(w, http.StatusOK, c) 
		return
	}
	utils.WriteJSON(w, http.StatusForbidden, map[string]string{"error": "not authorized"}) 
}

func ResolveComplaintHandler(w http.ResponseWriter, r *http.Request) {
	secret := r.Header.Get("X-Secret-Code") 
	admin, err := store.GetUserBySecret(secret)
	if err != nil || !admin.IsAdmin {
		utils.WriteJSON(w, http.StatusForbidden, map[string]string{"error": "admin only"}) 
		return
	}

	var req struct {
		ID string `json:"id"`
	}
	if err := utils.ReadJSONBody(r, &req); err != nil {
		utils.WriteJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()}) 
		return
	}

	_, err = store.DB.Exec("UPDATE complaints SET resolved = true WHERE id = ?", req.ID)
	if err != nil {
		utils.WriteJSON(w, http.StatusNotFound, map[string]string{"error": "complaint not found"}) 
		return
	}
	utils.WriteJSON(w, http.StatusOK, map[string]string{"status": "resolved", "id": req.ID}) 
}