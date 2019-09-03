package mem

// func Test_File_Read_Memory(t *testing.T) {
// 	r := require.New(t)
//
// 	fs, err := New(here.Info{})
// 	r.NoError(err)
//
// 	f, err := fs.Create("/file_test.go")
// 	r.NoError(err)
// 	_, err = io.Copy(f, bytes.NewReader([]byte("hi!")))
// 	r.NoError(err)
// 	r.NoError(f.Close())
//
// 	f, err = fs.Open("/file_test.go")
// 	r.NoError(err)
// 	fi, err := f.Stat()
// 	r.NoError(err)
// 	r.Equal("/file_test.go", fi.Name())
//
// 	b, err := ioutil.ReadAll(f)
// 	r.NoError(err)
// 	r.Equal(string(b), "hi!")
// }
//
// func Test_File_Write(t *testing.T) {
// 	r := require.New(t)
//
// 	fs, err := New(here.Info{})
// 	r.NoError(err)
//
// 	f, err := fs.Create("/hello.txt")
// 	r.NoError(err)
// 	r.NotNil(f)
//
// 	fi, err := f.Stat()
// 	r.NoError(err)
// 	r.Zero(fi.Size())
//
// 	r.Equal("/hello.txt", fi.Name())
//
// 	mt := fi.ModTime()
// 	r.NotZero(mt)
//
// 	sz, err := io.Copy(f, strings.NewReader(radio))
// 	r.NoError(err)
// 	r.Equal(int64(1381), sz)
//
// 	// because windows can't handle the time precisely
// 	// enough, we have to *force* just a smidge of time
// 	// to ensure the two ModTime's are different.
// 	// i know, i hate it too.
// 	time.Sleep(time.Millisecond)
// 	r.NoError(f.Close())
// 	r.Equal(int64(1381), fi.Size())
// 	r.NotZero(fi.ModTime())
// 	r.NotEqual(mt, fi.ModTime())
// }
