package structFlag

type StructMapCheck struct {
	IDPersonFlag map[int64]*BoolStruct
}
type BoolStruct struct {
	CheckFlag bool
	DeletFlag bool
}
