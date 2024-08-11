package main

import (
    "log"
    "net/http"

    "transaction-service/pkg"
)

func main() {
    r := pkg.NewRouter()
    log.Fatal(http.ListenAndServe(":8080", r))
}
