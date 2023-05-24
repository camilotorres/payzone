package main

import (
	"crypto/rand"
	"encoding/json"
	"fmt"
	"math/big"
	"net/http"
	"os"
	"sync"

	"github.com/moov-io/iso8583"
	"github.com/moov-io/iso8583/encoding"
	"github.com/moov-io/iso8583/field"
	"github.com/moov-io/iso8583/padding"
	"github.com/moov-io/iso8583/prefix"
	"github.com/moov-io/iso8583/specs"
)

type Transaction struct {
	PAN               string
	SystemTraceNumber int32
	ProcessingCode    string
	ResponseCode      string
	ID                int
	Amount            float64
	Fee               float64
	Success           bool
}

func HandleRoot(w http.ResponseWriter, r *http.Request) {

	isomessage := iso8583.NewMessage(specs.Spec87ASCII)

	isomessage.MTI("0100")
	isomessage.Field(2, "4919108000061104")
	isomessage.Field(3, "003000")
	isomessage.Field(39, "00")

	// bitmapInicial := field.NewBitmap(isomessage.Bitmap().Spec())

	// bitmapInicial.Set(15)

	rawMessage, err := isomessage.Pack()

	if err != nil {
		fmt.Fprintf(w, "error con specificacion ISO: %v", err)
		return
	}

	// bitmap, err := isomessage.Bitmap().String()

	strBitmap, _ := isomessage.Bitmap().String()

	fmt.Printf("rawMessage:%v\n", rawMessage)

	fmt.Printf("BITMAP01:%v\n", strBitmap)

	// data, _ := bitmapInicial.String()

	// fmt.Printf("BITMAP02:%v\n", data)

	// read, _ := bitmapInicial.Unpack([]byte("004000000000000000000000000000000000000000000000"))

	// fmt.Printf("BITMAP03:%v\n", read)

	message1 := iso8583.NewMessage(specs.Spec87ASCII)
	rawMsg := []byte("020042000400000000021612345678901234560609173030123456789ABC1000123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789")
	message1.Unpack([]byte(rawMsg))
	s, err := message1.GetString(2)
	fmt.Printf("BITMAP04:%v\n", s)

	spec := &iso8583.MessageSpec{

		Fields: map[int]field.Field{
			0: field.NewString(&field.Spec{
				Length:      4,
				Description: "Message Type Indicator",
				Enc:         encoding.ASCII,
				Pref:        prefix.ASCII.Fixed,
			}),
			1: field.NewBitmap(&field.Spec{
				Description: "Bitmap",
				Enc:         encoding.ASCII,
				Pref:        prefix.ASCII.Fixed,
			}),

			// Message fields:
			2: field.NewString(&field.Spec{
				Length:      19,
				Description: "Primary Account Number",
				Enc:         encoding.ASCII,
				Pref:        prefix.ASCII.LL,
			}),
			3: field.NewNumeric(&field.Spec{
				Length:      6,
				Description: "Processing Code",
				Enc:         encoding.ASCII,
				Pref:        prefix.ASCII.Fixed,
				Pad:         padding.Left('0'),
			}),
			4: field.NewString(&field.Spec{
				Length:      12,
				Description: "Transaction Amount",
				Enc:         encoding.ASCII,
				Pref:        prefix.ASCII.Fixed,
				Pad:         padding.Left('0'),
			}),
		},
	}

	// create message with defined spec
	message := iso8583.NewMessage(spec)
	// set message type indicator at field 0
	message.MTI("0100")
	// set all message fields you need as strings
	err1 := message.Field(2, "4242424242424242")

	if err1 != nil {
		fmt.Fprintf(w, "error con specificacion ISO: %v", err)
		return
	}

	// handle error
	err = message.Field(3, "123456")

	if err != nil {
		fmt.Fprintf(w, "error con specificacion ISO: %v", err)
		return
	}

	// handle error
	err = message.Field(4, "100")
	// handle error

	if err != nil {
		fmt.Fprintf(w, "error con specificacion ISO: %v", err)
		return
	}

	// generate binary representation of the message into rawMessage
	rawMessage2, err := message.Pack()

	if err != nil {
		fmt.Fprintf(w, "error con specificacion ISO: %v", err)
		return
	}

	bitmap, err := message.Bitmap().String()

	if err != nil {
		fmt.Fprintf(w, "error con el BITMAP: %v", err)
		return
	}

	jsonMessage, err := json.Marshal(message.GetFields())
	if err != nil {
		fmt.Fprintf(w, "error con el JSON: %v", err)
		return
	}

	fmt.Printf("JSON FORMAT %v\n", jsonMessage)
	fmt.Printf("BITMAP: %v\n", bitmap)
	fmt.Println(rawMessage2)
	// now you can send rawMessage over the wire

	iso8583.Describe(message, os.Stdout)

	decoder := json.NewDecoder(r.Body)
	var trx TRX

	err2 := decoder.Decode(&trx)

	if err2 != nil {
		fmt.Fprintf(w, "error: %v", err)
		return
	}

	maximo := big.NewInt(100000000)
	randomico, _ := rand.Int(rand.Reader, maximo)

	trx.PAN = trx.PAN + randomico.String()

	var wg sync.WaitGroup
	queue := make(chan TRX)

	// Iniciamos una goroutine para procesar las transacciones autorizadas en orden FIFO
	go func() {
		for transaction := range queue {
			// Aquí puedes realizar cualquier acción adicional con la transacción
			// En este ejemplo, simplemente imprimimos los detalles de la transacción autorizada
			fmt.Printf("Transaction ID: %d, Amount: %.2f, Success: %v\n", transaction.Id, transaction.Amount, transaction.Success)
		}
	}()

	// Simulamos recibir transacciones desde otro servicio
	// Aquí puedes reemplazar esta sección con la lógica de recepción de transacciones en tiempo real
	// transactions := []Transaction{
	// 	{ID: 1, Amount: 100.0},
	// 	{ID: 2, Amount: 200.0},
	// 	{ID: 3, Amount: 300.0},
	// }

	// Autorizamos cada transacción de forma concurrente
	// for _, transaction := range transactions {
	wg.Add(1)
	go authorizeTransaction(&trx, &wg, queue)
	// }

	// Esperamos a que todas las transacciones se autoricen
	wg.Wait()

	// Cerramos el canal de transacciones para finalizar la goroutine de procesamiento
	close(queue)

	response, err := trx.ToJson()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	w.Write(response)

}

func HandleHome(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hola Handle Home")
}

func authorizeTransaction(transaction *TRX, wg *sync.WaitGroup, queue chan TRX) {
	defer wg.Done()
	fmt.Printf("Transaction ID: %d, Amount: %.2f, Success: %v\n", transaction.Id, transaction.Amount, transaction.Success)

	maximo := big.NewInt(100000000)
	randomico, _ := rand.Int(rand.Reader, maximo)

	transaction.SystemTraceNumber = *randomico

	// Simulamos la autorización de la transacción
	// Aquí podrías agregar la lógica real de autorización con el autorizador bancario
	// En este ejemplo, simplemente marcamos todas las transacciones como exitosas
	transaction.Success = true
	transaction.ProcessingCode = "0210"
	transaction.ResponseCode = "00"

	// Agregamos la transacción autorizada al canal de transacciones
	queue <- *transaction
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

}
