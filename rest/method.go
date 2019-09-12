package rest

type Method int

const (
	GET Method = 1 + iota
	POST
	PUT
	PATCH
	DELETE
)

var method = [...]string{
	"GET",
	"POST",
	"PUT",
	"PATCH",
	"DELETE",
}

func (m Method) String() string { return method[m-1] }

