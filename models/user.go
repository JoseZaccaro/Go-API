package models

type Usuario struct {
	ID       int
	Nombre   string
	Apellido string
	Email    string
	Password string
}

type Usuarios []Usuario
