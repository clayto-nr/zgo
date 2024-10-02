package main

import (
    "fmt"
    "net/http"
)

func helloHandler(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintln(w, "Hello, World!")
}

func main() {
    http.HandleFunc("/", helloHandler) // Roteia a requisição para a função helloHandler
    fmt.Println("Servidor escutando na porta 8080...")
    err := http.ListenAndServe(":8080", nil) // Inicia o servidor na porta 8080
    if err != nil {
        fmt.Println("Erro ao iniciar o servidor:", err)
    }
}
