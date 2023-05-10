package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
)

type Transaction struct {
	ID      int
	Amount  float64
	Success bool
}

func HandleRoot(w http.ResponseWriter, r *http.Request) {

	var wg sync.WaitGroup
	queue := make(chan Transaction)

	// Iniciamos una goroutine para procesar las transacciones autorizadas en orden FIFO
	go func() {
		for transaction := range queue {
			// Aquí puedes realizar cualquier acción adicional con la transacción
			// En este ejemplo, simplemente imprimimos los detalles de la transacción autorizada
			fmt.Printf("Transaction ID: %d, Amount: %.2f, Success: %v\n", transaction.ID, transaction.Amount, transaction.Success)
		}
	}()

	// Simulamos recibir transacciones desde otro servicio
	// Aquí puedes reemplazar esta sección con la lógica de recepción de transacciones en tiempo real
	transactions := []Transaction{
		{ID: 1, Amount: 100.0},
		{ID: 2, Amount: 200.0},
		{ID: 3, Amount: 300.0},
	}

	// Autorizamos cada transacción de forma concurrente
	for _, transaction := range transactions {
		wg.Add(1)
		go authorizeTransaction(transaction, &wg, queue)
	}

	// Esperamos a que todas las transacciones se autoricen
	wg.Wait()

	// Cerramos el canal de transacciones para finalizar la goroutine de procesamiento
	close(queue)
}

func HandleHome(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hola Handle Home")
}

func authorizeTransaction(transaction Transaction, wg *sync.WaitGroup, queue chan Transaction) {
	defer wg.Done()

	// Simulamos la autorización de la transacción
	// Aquí podrías agregar la lógica real de autorización con el autorizador bancario
	// En este ejemplo, simplemente marcamos todas las transacciones como exitosas
	transaction.Success = true

	// Agregamos la transacción autorizada al canal de transacciones
	queue <- transaction
}

func PostRequest(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var metadata MetaData
	err := decoder.Decode(&metadata)
	if err != nil {
		fmt.Fprintf(w, "error: %v", err)
		return
	}

	fmt.Fprintf(w, "Payload %v\n", metadata)

}

func TrxPostRequest(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var trx TRX
	err := decoder.Decode(&trx)
	if err != nil {
		fmt.Fprintf(w, "error: %v", err)
		return
	}

	response, err := trx.ToJson()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(response)

	// fmt.Fprintf(w, "Payload %v\n", trx)
}
