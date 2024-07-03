package routes

import (
	"api/usuarios/config"
	"api/usuarios/models"
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"

	"golang.org/x/crypto/bcrypt"
)

func RegistrarUsuario(w http.ResponseWriter, r *http.Request) {
	//TODO: Get JSON data
	var usuario models.Usuario
	err := json.NewDecoder(r.Body).Decode(&usuario)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	//TODO: Validate JSON data
	if usuario.Nombre == "" || usuario.Apellido == "" || usuario.Email == "" || usuario.Password == "" {
		http.Error(w, "Campos obligatorios", http.StatusBadRequest)
		return
	}

	config.ConnectToDB()
	defer config.CloseConnection()

	//TODO: Check if user exists
	checking := "SELECT * FROM usuarios WHERE email = '" + usuario.Email + "'"
	resultChecking, errSel := config.MySqlDatabase.Query(checking)
	defer resultChecking.Close()
	if errSel != nil {
		http.Error(w, errSel.Error(), http.StatusInternalServerError)
		return
	}
	if resultChecking.Next() {
		http.Error(w, "El usuario ya existe", http.StatusBadRequest)
		return
	}

	// Check if password is strong
	if len(usuario.Password) < 8 {
		http.Error(w, "La contraseña debe tener al menos 8 caracteres", http.StatusBadRequest)
		return
	}

	// Check if email is valid
	if !validateEmail(usuario.Email) {
		http.Error(w, "El correo electrónico no es válido", http.StatusBadRequest)
		return
	}

	//TODO hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(usuario.Password), 14)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	usuario.Password = string(hashedPassword)

	//TODO Save JSON data
	sql := "INSERT INTO usuarios (nombre, apellido, email, password) VALUES ('" + usuario.Nombre + "', '" + usuario.Apellido + "', '" + usuario.Email + "', '" + usuario.Password + "')"

	result, err := config.MySqlDatabase.Exec(sql)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// Get the last inserted ID
	id, err := result.LastInsertId()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Update usuario.ID with the last inserted ID
	usuario.ID = int(id)

	// Build and return JSON response
	response := map[string]interface{}{
		"status": http.StatusText(http.StatusCreated),
		"id":     usuario.ID,
	}

	jsonData, err := json.Marshal(response)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Set response headers
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Location", "/users")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.WriteHeader(http.StatusCreated)
	_, err = w.Write(jsonData)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

}

func LoguearUsuario(w http.ResponseWriter, r *http.Request) {

	//TODO: Get JSON data
	var usuario models.Usuario
	err := json.NewDecoder(r.Body).Decode(&usuario)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	//TODO: Validate JSON data
	if usuario.Email == "" || usuario.Password == "" {
		http.Error(w, "Campos obligatorios", http.StatusBadRequest)
		return
	}

	config.ConnectToDB()
	defer config.CloseConnection()

	//TODO: Check if user exists
	sql := "SELECT * FROM usuarios WHERE email = '" + usuario.Email + "'"
	result, errSel := config.MySqlDatabase.Query(sql)
	if errSel != nil {
		http.Error(w, errSel.Error(), http.StatusInternalServerError)
		return
	}
	defer result.Close()

	//TODO: Check if password is correct
	var hashedPassword string
	var userDB models.Usuario
	for result.Next() {
		err := result.Scan(&userDB.ID, &userDB.Nombre, &userDB.Apellido, &userDB.Email, &userDB.Password)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		hashedPassword = userDB.Password
		fmt.Printf("Password: %s\n", hashedPassword)
		fmt.Printf("Password: %s\n", usuario.Password)
		err = bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(usuario.Password))
		if err != nil {
			http.Error(w, "Contraseña incorrecta", http.StatusBadRequest)
			return
		}
		usuario = userDB
	}

	// Build and return JSON response
	response := map[string]interface{}{
		"status":   http.StatusText(http.StatusOK),
		"id":       usuario.ID,
		"nombre":   usuario.Nombre,
		"apellido": usuario.Apellido,
		"email":    usuario.Email,
		"token":    "token",
	}

	jsonData, err := json.Marshal(response)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Set response headers
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Location", "/users")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.WriteHeader(http.StatusOK)
	_, err = w.Write(jsonData)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

}

// Helper function to validate email
func validateEmail(email string) bool {
	emailRegex := `^[a-zA-Z0-9.!#$%&'*+/=?^_` + "`" + `{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$`
	emailRegexCompiled := regexp.MustCompile(emailRegex)
	return emailRegexCompiled.MatchString(email)
}
