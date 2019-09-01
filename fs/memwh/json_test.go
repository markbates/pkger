package memwh

// func Test_File_JSON(t *testing.T) {
// 	r := require.New(t)
//
// 	fs, err := New(here.Info{})
// 	r.NoError(err)
//
// 	f, err := fs.Create("/radio.radio")
// 	r.NoError(err)
// 	_, err = io.Copy(f, strings.NewReader(radio))
// 	r.NoError(err)
// 	r.NoError(f.Close())
//
// 	f, err = fs.Open("/radio.radio")
// 	r.NoError(err)
// 	bi, err := f.Stat()
// 	r.NoError(err)
//
// 	mj, err := json.Marshal(f)
// 	r.NoError(err)
//
// 	f2 := &File{}
//
// 	r.NoError(json.Unmarshal(mj, f2))
//
// 	ai, err := f2.Stat()
// 	r.NoError(err)
//
// 	r.Equal(bi.Size(), ai.Size())
//
// 	r.Equal(radio, string(f2.data))
// }
