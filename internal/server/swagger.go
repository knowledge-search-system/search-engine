package server

import (
	_ "embed"
	"net/http"
)

//go:embed swaggerdoc/openapi.json
var openapiSpec []byte

const swaggerIndexHTML = `<!DOCTYPE html>
<html>
<head>
  <title>search-engine API docs</title>
  <link rel="stylesheet" href="https://unpkg.com/swagger-ui-dist@5/swagger-ui.css">
</head>
<body>
  <div id="swagger-ui"></div>
  <script src="https://unpkg.com/swagger-ui-dist@5/swagger-ui-bundle.js"></script>
  <script>
    window.onload = () => {
      window.ui = SwaggerUIBundle({
        url: "/docs/openapi.json",
        dom_id: "#swagger-ui",
      });
    };
  </script>
</body>
</html>`

func newSwaggerHandler() http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("/docs", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		_, _ = w.Write([]byte(swaggerIndexHTML))
	})

	mux.HandleFunc("/docs/openapi.json", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write(openapiSpec)
	})

	return mux
}
