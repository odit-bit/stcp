package stcp

// func Benchmark_(b *testing.B) {
// 	var log LoginRequest
// 	for i := 0; i < b.N; i++ {
// 		lr := NewLoginRequest("admin", "12345", "22a10", "1")
// 		log = *lr
// 	}
// 	_ = log
// }

// func Benchmark_write(b *testing.B) {
// 	lr := NewLoginRequest("admin", "12345", "22a10", "1")
// 	var dst bytes.Buffer
// 	for i := 0; i < b.N; i++ {
// 		err := binary.Write(&dst, binary.BigEndian, lr)
// 		if err != nil {
// 			b.Log(err)
// 			b.Fail()
// 		}
// 		dst.Reset()
// 	}
// }

// func Benchmark_read(b *testing.B) {
// 	b.StopTimer()
// 	lr := NewLoginRequest("admin", "12345", "22a10", "1")
// 	var temp bytes.Buffer
// 	binary.Write(&temp, binary.BigEndian, lr)

// 	n := temp.Len()
// 	data := make([]byte, n)
// 	copy(data, temp.Bytes())

// 	var src bytes.Reader
// 	src.Reset(data)

// 	b.ResetTimer()
// 	b.StartTimer()
// 	var expect LoginRequest
// 	for i := 0; i < b.N; i++ {
// 		err := binary.Read(&src, binary.BigEndian, &expect)
// 		if err != nil {
// 			b.Log(err)
// 			b.Fail()
// 		}
// 		src.Reset(data)
// 	}
// }

// func Benchmark_sequenced(b *testing.B) {
// 	b.StopTimer()
// 	data := bytes.Repeat([]byte{'X'}, 1000) //[]byte{'A', 'h', 'e', 'l', 'l', 'o'}
// 	s := Sequenced{
// 		Prefix: Prefix{
// 			Length: uint16(1 + len(data)),
// 			Typ:    'S',
// 		},
// 		Payload: data,
// 	}
// 	dst := io.Discard

// 	b.ResetTimer()
// 	b.StartTimer()
// 	for i := 0; i < b.N; i++ {
// 		err := writePrefixBinary(dst, &s.Prefix)
// 		if err != nil {
// 			b.Fatal(err)
// 		}
// 		err = writePayloadBinary(dst, s.Payload)
// 		if err != nil {
// 			b.Fatal(err)
// 		}

// 	}
// }
