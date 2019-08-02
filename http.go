package pkger

// type HTTPFile struct {
// 	io.Reader
// 	StatFn    func() (os.FileInfo, error)
// 	ReaddirFn func(int) ([]os.FileInfo, error)
// }
//
// func (h HTTPFile) Readdir(n int) ([]os.FileInfo, error) {
// 	if h.ReaddirFn != nil {
// 		return h.ReaddirFn(n)
// 	}
// 	return nil, nil
// }
//
// func (h HTTPFile) Stat() (os.FileInfo, error) {
// 	if h.StatFn != nil {
// 		return h.StatFn()
// 	}
// 	return nil, nil
// }
//
// func (h HTTPFile) Close() error {
// 	if c, ok := h.Reader.(io.Closer); ok {
// 		return c.Close()
// 	}
// 	return nil
// }
//
// func (h HTTPFile) Seek(offset int64, whence int) (int64, error) {
// 	if sk, ok := h.Reader.(io.Seeker); ok {
// 		return sk.Seek(offset, whence)
// 	}
// 	return 0, nil
// }
