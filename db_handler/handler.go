package db_handler


type DBHandler struct{
	Handler any
}

func NewDBHandler(handler any) *DBHandler {
	return &DBHandler{Handler: handler}
}

