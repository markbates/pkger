package memwh

// func Test_HTTP_Dir(t *testing.T) {
// 	r := require.New(t)
//
// 	fs := NewFS()
//
// 	r.NoError(Folder.Create(fs))
//
// 	dir, err := fs.Open("/")
// 	r.NoError(err)
// 	ts := httptest.NewServer(http.FileServer(dir))
// 	defer ts.Close()
//
// 	res, err := http.Get(ts.URL + "/")
// 	r.NoError(err)
// 	r.Equal(200, res.StatusCode)
//
// 	b, err := ioutil.ReadAll(res.Body)
// 	r.NoError(err)
// 	r.Contains(string(b), `<a href="/public/images/mark.png">/public/images/mark.png</a>`)
// }
//
// func Test_HTTP_File_Memory(t *testing.T) {
// 	r := require.New(t)
//
// 	fs := NewFS()
// 	r.NoError(Folder.Create(fs))
//
// 	dir, err := fs.Open("/")
// 	r.NoError(err)
// 	ts := httptest.NewServer(http.FileServer(dir))
// 	defer ts.Close()
//
// 	res, err := http.Get(ts.URL + "/public/images/mark.png")
// 	r.NoError(err)
// 	r.Equal(200, res.StatusCode)
//
// 	b, err := ioutil.ReadAll(res.Body)
// 	r.NoError(err)
// 	r.Contains(string(b), `!/public/images/mark.png`)
// }
//
// func Test_HTTP_Dir_Memory_StripPrefix(t *testing.T) {
// 	r := require.New(t)
//
// 	fs := NewFS()
// 	r.NoError(Folder.Create(fs))
//
// 	dir, err := fs.Open("/public")
// 	r.NoError(err)
// 	defer dir.Close()
//
// 	ts := httptest.NewServer(http.StripPrefix("/assets/", http.FileServer(dir)))
// 	defer ts.Close()
//
// 	res, err := http.Get(ts.URL + "/assets/images/mark.png")
// 	r.NoError(err)
// 	r.Equal(200, res.StatusCode)
//
// 	b, _ := ioutil.ReadAll(res.Body)
// 	// r.NoError(err)
// 	r.Contains(string(b), "!/public/images/mark.png")
//
// 	res, err = http.Get(ts.URL + "/assets/images/")
// 	r.NoError(err)
// 	r.Equal(200, res.StatusCode)
//
// 	b, _ = ioutil.ReadAll(res.Body)
// 	// r.NoError(err)
// 	r.Contains(string(b), `<a href="/mark.png">/mark.png</a>`)
// 	r.NotContains(string(b), `/public`)
// 	r.NotContains(string(b), `/images`)
// 	r.NotContains(string(b), `/go.mod`)
// }
