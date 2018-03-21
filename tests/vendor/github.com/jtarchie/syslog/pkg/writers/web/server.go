package writers

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/blevesearch/bleve"
	"github.com/jtarchie/syslog/pkg/log"
)

type Server struct {
	port       int
	index      bleve.Index
	httpServer *http.Server
}

type doc struct {
	Hostname string
	Message  string
}

func NewServer(port int) *Server {
	return &Server{
		port: port,
	}
}

func (s *Server) Write(l *syslog.Log) error {
	id := strconv.FormatInt(time.Now().UnixNano(), 10)

	return s.index.Index(id, doc{
		Hostname: string(l.Hostname()),
		Message:  string(l.Message()),
	})
}

func (s *Server) Start() error {
	var (
		err error
	)

	log.Printf("web: starting search index")
	s.index, err = bleve.NewMemOnly(bleve.NewIndexMapping())
	if err != nil {
		return fmt.Errorf("could not start search: %s", err)
	}

	mux := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		html := "<html><body><pre>"

		query := bleve.NewQueryStringQuery(r.URL.Query().Get("q"))
		search := bleve.NewSearchRequest(query)
		search.Highlight = bleve.NewHighlightWithStyle("html")
		search.Highlight.AddField("Message")

		results, err := s.index.Search(search)
		if err != nil {
			html += fmt.Sprintf("Something went wrong: %s", err)
		} else {
			html += fmt.Sprint(results)
		}
		html += "</pre></body></html>"
		w.Write([]byte(html))
	})

	log.Printf("web: starting webserver on port %d", s.port)
	s.httpServer = &http.Server{Addr: fmt.Sprintf("localhost:%d", s.port), Handler: mux}
	return s.httpServer.ListenAndServe()
}

func (s *Server) Addr() string {
	return s.httpServer.Addr
}
