package internal

import (
    "encoding/json"
    "net/http"
    "strconv"
    "transaction-service/pkg"  // Import only for models, not for router
)

func CreateTransactionHandler(w http.ResponseWriter, r *http.Request) {
    var t pkg.Transaction
    err := json.NewDecoder(r.Body).Decode(&t)
    if err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }

    _, err = GetDB().Exec("INSERT INTO transactions (id, amount, type, parent_id) VALUES ($1, $2, $3, $4)",
        t.ID, t.Amount, t.Type, t.ParentID)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    w.WriteHeader(http.StatusOK)
    json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}

func GetTransactionHandler(w http.ResponseWriter, r *http.Request) {
    id := r.URL.Query().Get("id")
    row := GetDB().QueryRow("SELECT amount, type, parent_id FROM transactions WHERE id = $1", id)

    var t pkg.Transaction
    err := row.Scan(&t.Amount, &t.Type, &t.ParentID)
    if err != nil {
        http.Error(w, err.Error(), http.StatusNotFound)
        return
    }

    w.WriteHeader(http.StatusOK)
    json.NewEncoder(w).Encode(t)
}

func GetTransactionsByTypeHandler(w http.ResponseWriter, r *http.Request) {
    t := r.URL.Query().Get("type")
    rows, err := GetDB().Query("SELECT id FROM transactions WHERE type = $1", t)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    defer rows.Close()

    var ids []int64
    for rows.Next() {
        var id int64
        err := rows.Scan(&id)
        if err != nil {
            http.Error(w, err.Error(), http.StatusInternalServerError)
            return
        }
        ids = append(ids, id)
    }

    w.WriteHeader(http.StatusOK)
    json.NewEncoder(w).Encode(ids)
}

func GetSumHandler(w http.ResponseWriter, r *http.Request) {
    idStr := r.URL.Query().Get("id")
    id, err := strconv.ParseInt(idStr, 10, 64)
    if err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }

    var sum float64
    err = getSumRecursive(id, &sum)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    w.WriteHeader(http.StatusOK)
    json.NewEncoder(w).Encode(map[string]float64{"sum": sum})
}

func getSumRecursive(id int64, sum *float64) error {
    rows, err := GetDB().Query("SELECT amount FROM transactions WHERE parent_id = $1", id)
    if err != nil {
        return err
    }
    defer rows.Close()

    for rows.Next() {
        var amount float64
        err := rows.Scan(&amount)
        if err != nil {
            return err
        }
        *sum += amount
        err = getSumRecursive(id, sum)
        if err != nil {
            return err
        }
    }
    return nil
}
