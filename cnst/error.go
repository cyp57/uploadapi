package cnst

var errMap = make(map[int]string)


func InitErr(){
	// errMap[4001] = "The user was not found or the password was incorrect. The response will be empty."
	// add more error code...
}


func GetErrMsg(errcode int) string{
	msg := ""

	errStr, ok := errMap[errcode]
	if ok {
		msg = errStr
	} 
	return msg
}