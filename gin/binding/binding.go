package binding

const (
	MIMEPlain = "text/plain"
)

//是最小的接口，需要被实现才能被用作验证引擎，
type StructValidator interface {
	//可以接收任意类型的参数并且永远不会panic，即使配置错误
	//如果参数不是一个结构体，则会跳过验证并返回nil
	//如果参数是结构体或者结构体指针，则应该验证
	//如果结构体验证事态，则应该返回错误信息，否则返回nil
	ValidateStruct(interface{})error

	//返回底层引擎
	Engine() interface{}
}

var Validator StructValidator = &defaultValidator{}

var (
	JSON = jsonBinding{}
)

func validate(obj interface{})error{
	if Validator == nil{
		return nil
	}

	return Validator.ValidateStruct(obj)
}