package structFlag

type StructMapCheck struct {
	IDPersonFlag map[int64]*BoolStruct
}
type BoolStruct struct {
	AddNoteFlag   bool
	DeletNoteFlag bool
}
