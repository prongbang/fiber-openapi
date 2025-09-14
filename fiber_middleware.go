package fiberopenapi

import (
	"net/http"
	"net/url"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/prongbang/oapigen"
)

func New(cfgs ...oapigen.Config) fiber.Handler {
	m := oapigen.New(cfgs...)
	return func(c *fiber.Ctx) error {
		app := c.App()
		m.SetRoutesProvider(func() []oapigen.Route {
			routes := app.GetRoutes(true)
			out := make([]oapigen.Route, 0, len(routes))
			for _, r := range routes {
				if r.Path == "" || r.Method == "" {
					continue
				}
				out = append(out, oapigen.Route{
					Method: strings.ToUpper(r.Method),
					Path:   r.Path,
				})
			}
			return out
		})

		// serve JSON spec
		if c.Path() == m.Cfg.JSONPath {
			m.EnsureSpecInitialized()
			m.Mu.RLock()
			data := append([]byte(nil), m.SpecJSON...)
			m.Mu.RUnlock()
			c.Set(fiber.HeaderContentType, fiber.MIMEApplicationJSONCharsetUTF8)
			return c.Send(data)
		}

		// serve Scalar
		if c.Path() == m.Cfg.DocsPath {
			html := oapigen.NewScalar(oapigen.DocConfig{Title: m.Cfg.Title, JSONPath: m.Cfg.JSONPath})
			c.Set(fiber.HeaderContentType, fiber.MIMETextHTMLCharsetUTF8)
			return c.SendString(html)
		}

		if m.Cfg.Observe == oapigen.ObsDisable {
			return c.Next()
		}

		method := c.Method()
		reqBody := append([]byte(nil), c.Body()...)
		reqCT := string(c.Request().Header.ContentType())

		queryVals := url.Values{}
		for name, val := range c.Queries() {
			queryVals.Add(name, val)
		}

		reqHeaders := http.Header{}
		c.Request().Header.VisitAll(func(k, v []byte) {
			name := strings.TrimSpace(string(k))
			if name == "" {
				return
			}
			reqHeaders.Add(name, string(v))
		})

		err := c.Next()

		resHeaders := http.Header{}
		c.Response().Header.VisitAll(func(k, v []byte) {
			name := strings.TrimSpace(string(k))
			if name == "" {
				return
			}
			resHeaders.Add(name, string(v))
		})

		status := c.Response().StatusCode()
		resBody := append([]byte(nil), c.Response().Body()...)
		resCT := string(c.Response().Header.ContentType())

		routePattern := c.Path()
		if rt := c.Route(); rt != nil && rt.Path != "" {
			routePattern = rt.Path
		}
		if routePattern == m.Cfg.JSONPath || routePattern == m.Cfg.DocsPath {
			return err
		}

		m.Capture(oapigen.CaptureContext{
			Method:              method,
			Path:                c.Path(),
			RoutePattern:        routePattern,
			QueryParams:         queryVals,
			RequestHeaders:      reqHeaders,
			RequestBody:         reqBody,
			RequestContentType:  reqCT,
			ResponseHeaders:     resHeaders,
			ResponseBody:        resBody,
			ResponseContentType: resCT,
			Status:              status,
		})

		return err
	}
}
