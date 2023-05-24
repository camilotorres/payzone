package field

import (
	"fmt"
	"testing"

	"github.com/moov-io/iso8583/encoding"
	"github.com/moov-io/iso8583/prefix"
	"github.com/stretchr/testify/require"
)

func Chunks(s string, chunkSize int) []string {
	if len(s) == 0 {
		return nil
	}
	if chunkSize >= len(s) {
		return []string{s}
	}
	var chunks []string = make([]string, 0, (len(s)-1)/chunkSize+1)
	currentLen := 0
	currentStart := 0
	for i := range s {
		if currentLen == chunkSize {
			chunks = append(chunks, s[currentStart:i])
			currentLen = 0
			currentStart = i
		}
		currentLen++
	}
	chunks = append(chunks, s[currentStart:])
	return chunks
}

// ConvertInt convierte el valor de cadena dado de base
// a la base (toBase) definida.
// func ConvertInt(val string, base, toBase int) (string, error) {
// 	i, err := strconv.ParseInt(val, base, 64)
// 	if err != nil {
// 		return "ERROR PUES NADa", err
// 	}
// 	return strconv.FormatInt(i, toBase), nil
// }

func TestBinaryField(t *testing.T) {

	// r := Chunks("0111001000100100010001001000000000101000110000001000000000000000", 4)

	// fmt.Printf("PREPARAMOS EL ARRAY CON LOS BYTES: %s \n", r)

	// var total string = " "
	// for i := 0; i < len(r); i++ {
	// 	println(r[i])
	// 	resultado, _ := ConvertInt(r[i], 2, 16)
	// 	total = total + resultado
	// }

	// println(total)

	// binaryBitmap := "0010"
	// // hexBitmap := "FA5800802C0001000000000000000000"
	// v, errores := ConvertInt(binaryBitmap, 2, 16)
	// fmt.Printf("Convertimos Bits a Hexa: %s, error %s \n", v, errores)

	spec := &Spec{
		Length:      32,
		Description: "Field",
		Enc:         encoding.Binary,
		Pref:        prefix.Binary.Fixed,
	}

	in := []byte("f23884012ee194180000004210000085")

	t.Run("Pack returns binary data", func(t *testing.T) {
		bin := NewBinary(spec)
		bin.SetBytes(in)

		packed, err := bin.Pack()

		require.NoError(t, err)
		require.Equal(t, in, packed)
	})

	t.Run("String returns binary data encoded in HEX", func(t *testing.T) {
		bin := NewBinary(spec)
		bin.value = in

		str, err := bin.String()

		require.NoError(t, err)
		require.Equal(t, "6632333838343031326565313934313830303030303034323130303030303835", str)
	})

	t.Run("Unpack returns binary data", func(t *testing.T) {
		bin := NewBinary(spec)

		n, err := bin.Unpack(in)

		fmt.Println(in)
		fmt.Println(n)
		fmt.Println(bin.value)

		require.NoError(t, err)
		require.Equal(t, len(in), n)
		require.Equal(t, in, bin.value)
	})

	t.Run("SetData sets data to the field", func(t *testing.T) {
		bin := NewBinary(spec)
		bin.SetData(NewBinaryValue(in))

		packed, err := bin.Pack()

		require.NoError(t, err)
		require.Equal(t, in, packed)
	})

	t.Run("Unmarshal gets data from the field", func(t *testing.T) {
		bin := NewBinaryValue([]byte{1, 2, 3})
		val := &Binary{}

		err := bin.Unmarshal(val)

		require.NoError(t, err)
		require.Equal(t, []byte{1, 2, 3}, val.value)
	})

	t.Run("SetBytes sets data to the data field", func(t *testing.T) {
		bin := NewBinary(spec)
		data := &Binary{}
		bin.SetData(data)

		err := bin.SetBytes(in)
		require.NoError(t, err)

		require.Equal(t, in, data.value)
	})

	// SetValue sets data to the data field
	t.Run("SetValue sets data to the data field", func(t *testing.T) {
		bin := NewBinary(spec)
		bin.SetValue(in)

		require.Equal(t, in, bin.Value())
	})

	t.Run("Unpack sets data to data value", func(t *testing.T) {
		bin := NewBinary(spec)
		data := NewBinaryValue([]byte{})
		bin.SetData(data)

		n, err := bin.Unpack(in)

		require.NoError(t, err)
		require.Equal(t, len(in), n)
		require.Equal(t, in, data.value)
	})

	t.Run("UnmarshalJSON unquotes input before handling it", func(t *testing.T) {
		input := []byte(`"500000000000000000000000000000000000000000000000"`)

		bin := NewBinary(spec)
		require.NoError(t, bin.UnmarshalJSON(input))

		str, err := bin.String()
		require.NoError(t, err)

		require.Equal(t, `500000000000000000000000000000000000000000000000`, str)
	})

	t.Run("MarshalJSON returns string hex repr of binary field", func(t *testing.T) {
		bin := NewBinaryValue([]byte{0xAB})
		marshalledJSON, err := bin.MarshalJSON()
		require.NoError(t, err)
		require.Equal(t, `"AB"`, string(marshalledJSON))
	})

	t.Run("returns error for zero value when fixed length and no padding specified", func(t *testing.T) {
		bin := NewBinary(spec)
		_, err := bin.Pack()

		if err != nil {
			fmt.Printf("ERROR EN LA SECCIÓN: %s \n", err.Error())
		}

		require.EqualError(t, err, "failed to encode length: field length: 0 should be fixed: 10")
	})
}

func TestBinaryNil(t *testing.T) {
	var str *Binary = nil

	bs, err := str.Bytes()
	require.NoError(t, err)
	require.Nil(t, bs)

	value, err := str.String()
	require.NoError(t, err)
	require.Equal(t, "", value)

	bs = str.Value()
	require.Nil(t, bs)
}
