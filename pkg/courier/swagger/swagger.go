package swagger

import (
	"bytes"
	"context"
	"fmt"
	"io/ioutil"
	"strings"

	"github.com/profzone/eden-framework/pkg/courier"
	"github.com/profzone/eden-framework/pkg/courier/httpx"
	"github.com/profzone/eden-framework/pkg/courier/transport_http"
)

func getSwaggerJSON() []byte {
	data, err := ioutil.ReadFile("./api/swagger.json")
	if err != nil {
		return data
	}
	return data
}

var SwaggerRouter = courier.NewRouter()

func init() {
	SwaggerRouter.Register(courier.NewRouter(Swagger{}))
	SwaggerRouter.Register(courier.NewRouter(SwaggerDoc{}))
	SwaggerRouter.Register(courier.NewRouter(SwaggerUIBundle{}))
}

type Swagger struct {
	httpx.MethodGet
}

func (s Swagger) Output(c context.Context) (interface{}, error) {
	json := &JSONBytes{}
	json.Write(getSwaggerJSON())
	return json, nil
}

// swagger:strfmt json
type JSONBytes struct {
	bytes.Buffer
}

func (JSONBytes) ContentType() string {
	return httpx.MIMEJSON
}

type SwaggerDoc struct {
	httpx.MethodGet
}

func (SwaggerDoc) Path() string {
	return "/doc"
}

func (s SwaggerDoc) Output(c context.Context) (interface{}, error) {
	html := &httpx.HTML{}
	u := transport_http.GetRequest(c).URL

	html.WriteString(`<!-- HTML for static distribution bundle build -->
<!DOCTYPE html>
<html lang="en">
  <head>
    <meta charset="UTF-8">
    <title>Swagger UI</title>
    <link rel="stylesheet" type="text/css" href="//cdnjs.cloudflare.com/ajax/libs/swagger-ui/3.17.4/swagger-ui.css" >
    <style>
      html
      {
        box-sizing: border-box;
        overflow: -moz-scrollbars-vertical;
        overflow-y: scroll;
      }

      *,
      *:before,
      *:after
      {
        box-sizing: inherit;
      }

      body
      {
        margin:0;
        background: #fafafa;
      }
    </style>
  </head>

  <body>
    <div id="swagger-ui"></div>
    <script src="` + strings.TrimSuffix(u.Path, s.Path()) + "/swagger-ui-bundle.js" + `"> </script>
    <script src="//cdnjs.cloudflare.com/ajax/libs/swagger-ui/3.17.4/swagger-ui-standalone-preset.js"> </script>
    <script>
    window.onload = function() {
      // Build a system
      var ui = SwaggerUIBundle({
        url: "` + strings.TrimSuffix(u.Path, s.Path()) + `",
        dom_id: '#swagger-ui',
        deepLinking: true,
        presets: [
          SwaggerUIBundle.presets.apis,
          SwaggerUIStandalonePreset
        ],
        plugins: [
          SwaggerUIBundle.plugins.DownloadUrl
        ],
        layout: "StandaloneLayout"
      })

      window.ui = ui
    }
  </script>
  </body>
</html>`)

	return html, nil
}

type SwaggerUIBundle struct {
	httpx.MethodGet
}

func (SwaggerUIBundle) Path() string {
	return "/swagger-ui-bundle.js"
}

func (s SwaggerUIBundle) Output(c context.Context) (interface{}, error) {
	html := &httpx.HTML{}

	b, err := Asset("swagger-ui-bundle.js")
	if err != nil {
		return nil, fmt.Errorf("asset: %v", err)
	}

	_, err = html.Write(b)
	if err != nil {
		return nil, fmt.Errorf("html write: %v", err)
	}

	return html, nil
}
