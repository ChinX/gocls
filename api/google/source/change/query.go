package change

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
)

var (
	changesURL = "https://go-review.googlesource.com/changes/"
)

type Query struct {
	ownerID    uint
	startID    uint
	pageNumber uint
	params     params
}

func NewQuery() *Query {
	return &Query{
		ownerID:    881,
		startID:    0,
		pageNumber: 25,
		params:     map[string]string{},
	}
}

func (q *Query) Project(project string) {
	q.params["project"] = project
}

func (q *Query) Page(num uint) {
	if num == 0 {
		q.startID = 0
	} else {
		q.startID = (num - 1) * q.pageNumber
	}
}

func (q *Query) NextPage() {
	q.startID += q.pageNumber
}

func (q *Query) PreviousPage() {
	q.startID -= q.pageNumber
	if q.startID < 0 {
		q.startID = 0
	}
}

func (q *Query) Open() {
	q.params["status"] = "open"
}

func (q *Query) Merged() {
	q.params["status"] = "merged"
}

func (q *Query) Abandoned() {
	q.params["status"] = "abandoned"
}

func (q *Query) Body() io.Reader {
	return nil
}

func (q *Query) Method() string {
	return http.MethodGet
}

func (q *Query) RequestURL() string {
	buf := &bytes.Buffer{}
	buf.WriteString(changesURL)
	buf.WriteByte('?')
	q.writeTo(buf)
	return buf.String()
}

func (q *Query) writeTo(buf *bytes.Buffer) {
	buf.WriteByte('O')
	buf.WriteByte('=')
	buf.WriteString(fmt.Sprintf("%d", q.ownerID))

	buf.WriteByte('&')
	buf.WriteByte('S')
	buf.WriteByte('=')
	buf.WriteString(fmt.Sprintf("%d", q.startID))

	buf.WriteByte('&')
	buf.WriteByte('n')
	buf.WriteByte('=')
	buf.WriteString(fmt.Sprintf("%d", q.pageNumber))

	if len(q.params) != 0 {
		buf.WriteByte('&')
		buf.WriteByte('q')
		buf.WriteByte('=')
		q.params.writeTo(buf)
	}
}

type params map[string]string

func (p params) writeTo(buf *bytes.Buffer) {
	for k, v := range p {
		buf.WriteString(k)
		buf.WriteByte(':')
		buf.WriteString(v)
		buf.WriteString("%20")
	}
}
