package quic

// import "github.com/lucas-clemente/quic-go"

// // Ref: https://github.com/libp2p/go-libp2p-quic-transport/blob/master/stream.go
// type Stream struct {
// 	Stream quic.Stream
// }

// // quic stream Close currently only closes the write half,
// // we need to close the read half to fully close the stream
// func (s *Stream) Close() error {
// 	s.Stream.CancelRead(0)
// 	return s.Stream.Close()
// }

// func (s *Stream) Read(p []byte) (n int, err error) {
// 	return s.Stream.Read(p)
// }

// func (s *Stream) Write(p []byte) (n int, err error) {
// 	return s.Stream.Write(p)
// }
