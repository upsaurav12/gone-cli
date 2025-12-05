package framework

import "fmt"

type FrameworkConfig struct {
	Imports       string
	Entity        string
	ContextName   string
	ContextType   string
	Bind          string
	JSON          string
	Router        string
	Start         string
	OtherImports  string
	ApiGroup      func(entity string, get string, lowerentity string) string
	Get           string
	FullContext   string
	ToTheClient   string
	Response      string
	ImportRouter  string
	ImportHandler string
}

var FrameworkRegistory = map[string]FrameworkConfig{
	"gin": {
		Imports:      `"github.com/gin-gonic/gin"`,
		ContextName:  "c",
		ContextType:  "*gin.Context",
		Bind:         "c.BindJSON", // func(obj any) error
		JSON:         "c.JSON",     // func(code int, obj any)
		Router:       "*gin.Engine",
		Start:        "gin.Default()",
		OtherImports: `"net/http"`,
		FullContext:  "c *gin.Context",
		ApiGroup: func(entity string, get string, lowerentity string) string {
			apiGroup := fmt.Sprintf(`
				api := r.Group("/api/v1")
				{
					%s := api.Group("/%s")
					{
						%s.%s("", handler.Get%ss)
					}
				}
			`, lowerentity, lowerentity, lowerentity, get, entity)

			return apiGroup
		},
		Get:         "GET",
		ToTheClient: "c.JSON(http.StatusOK, ",
		Response:    "(http.StatusOK,",
		ImportRouter: `
				"net/http"
	"github.com/gin-gonic/gin"
		`,
		ImportHandler: `
		    "github.com/gin-gonic/gin"
    		"net/http"

		`,
	},

	"chi": {
		Imports:     `"github.com/go-chi/chi/v5"`,
		ContextName: "r",
		ContextType: "http.ResponseWriter, *http.Request", // chi passes both
		Bind:        "json.NewDecoder(r.Body).Decode",     // need encoding/json
		JSON:        "render.JSON",                        // from go-chi/render
		Router:      "chi.Router",
		Start:       "chi.NewRouter()",
		OtherImports: `
			"encoding/json"
			"github.com/go-chi/render"
			"net/http"
		`,
		ApiGroup: func(entity string, get string, lowerentity string) string {
			apiGroup := fmt.Sprintf(`
				r.Group(func(r chi.Router) {
					r.%s("/%s", handler.Get%ss)
			})
			`, get, lowerentity, entity)

			return apiGroup
		},
		Get:         "Get",
		FullContext: "w http.ResponseWriter, r *http.Request",
		ToTheClient: "json.NewEncoder(w).Encode(",
		Response:    "(w, r,",
		ImportRouter: `
			"encoding/json"
	"net/http"
	"github.com/go-chi/chi/v5"
		`,
		ImportHandler: `
			"net/http"
			"github.com/go-chi/render"
		`,
	},

	"echo": {
		Imports:      `"github.com/labstack/echo/v4"`,
		ContextName:  "c",
		ContextType:  "*echo.Context",
		Bind:         "c.Bind",
		JSON:         "c.JSON",
		Router:       "*echo.Echo",
		Start:        "echo.New()",
		OtherImports: `"net/http"`,
	},

	"fiber": {
		Imports:      `"github.com/gofiber/fiber/v2"`,
		ContextName:  "c",
		ContextType:  "*fiber.Ctx",
		Bind:         "c.BodyParser", // BodyParser(obj)
		JSON:         "c.JSON",       // JSON(obj)
		Router:       "*fiber.App",
		Start:        "fiber.New()",
		OtherImports: "",
	},

	"mux": {
		Imports:     `"github.com/gorilla/mux"`,
		ContextName: "w, r",
		ContextType: "http.ResponseWriter, *http.Request",
		Bind:        "json.NewDecoder(r.Body).Decode",
		JSON:        `json.NewEncoder(w).Encode`,
		Router:      "*mux.Router",
		Start:       "mux.NewRouter()",
		OtherImports: `
			"encoding/json"
			"net/http"
		`,
	},
}
