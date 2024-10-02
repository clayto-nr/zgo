package main

import (
    "fmt"
    "net/http"
)

// Função exportada que será chamada pelo Vercel
func Handler(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintln(w, "Hello, World!")
}

func main() {
    http.HandleFunc("/", Handler) // Roteia a requisição para a função Handler
    fmt.Println("Servidor escutando na porta 8080...")
    err := http.ListenAndServe(":8080", nil) // Inicia o servidor na porta 8080
    if err != nil {
        fmt.Println("Erro ao iniciar o servidor:", err)
    }
}
