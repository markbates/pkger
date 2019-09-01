package memfs

// func Test_Walk(t *testing.T) {
// 	r := require.New(t)
//
// 	files := []struct {
// 		name string
// 		body string
// 	}{
// 		{name: "/a/a.txt", body: "A"},
// 		{name: "/a/a.md", body: "Amd"},
// 		{name: "/b/c/d.txt", body: "B"},
// 		{name: "/f.txt", body: "F"},
// 	}
//
// 	sort.Slice(files, func(a, b int) bool {
// 		return files[a].name < files[b].name
// 	})
//
// 	fs, err := New(here.Info{})
// 	r.NoError(err)
//
// 	for _, file := range files {
// 		f, err := fs.Create(file.name)
// 		r.NoError(err)
// 		_, err = io.Copy(f, strings.NewReader(file.body))
// 		r.NoError(err)
// 		r.NoError(f.Close())
// 	}
//
// 	var found []string
// 	err = fs.Walk("/", func(path string, info os.FileInfo, err error) error {
// 		if err != nil {
// 			return err
// 		}
//
// 		found = append(found, path)
// 		return nil
// 	})
// 	r.NoError(err)
//
// 	expected := []string{":/", ":/a", ":/a/a.md", ":/a/a.txt", ":/b", ":/b/c", ":/b/c/d.txt", ":/f.txt"}
// 	r.Equal(expected, found)
//
// 	found = []string{}
// 	err = fs.Walk("/a/", func(path string, info os.FileInfo, err error) error {
// 		if err != nil {
// 			return err
// 		}
//
// 		found = append(found, path)
// 		return nil
// 	})
// 	r.NoError(err)
//
// 	expected = []string{":/a/a.md", ":/a/a.txt"}
// 	r.Equal(expected, found)
// }
