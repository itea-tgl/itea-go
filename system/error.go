package system

type DatabaseError struct {
	error string
}

func (de *DatabaseError) Error() string {
	return de.error
}

func (de *DatabaseError) setError(err string) {
	de.error = err
}

type ServerError struct {
	error string
}

func (se *ServerError) Error() string {
	return se.error
}

func (se *ServerError) setError(err string) {
	se.error = err
}

type ParamsError struct {
	error string
}

func (pe *ParamsError) Error() string {
	return pe.error
}

func (pe *ParamsError) setError(err string) {
	pe.error = err
}

func NewDatabaseError(err error) *DatabaseError{
	e := &DatabaseError{}
	e.setError(err.Error())
	return e
}

func NewServerError(err error) *ServerError{
	e := &ServerError{}
	e.setError(err.Error())
	return e
}

func NewParamsError(err error) *ParamsError{
	e := &ParamsError{}
	e.setError(err.Error())
	return e
}
