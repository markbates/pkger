package pkger

func Exclude(fl *File, excludes ...string) error {
	fl.excludes = append(fl.excludes, excludes...)
	return nil
}
