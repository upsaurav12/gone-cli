package framework

import "fmt"

type FrameworkConfig struct {
	Name          string
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
	Returnable    string
	ReturnKeyword string
	HTTPHandler   string
	Entities      []string
}

var FrameworkRegistory = map[string]FrameworkConfig{
	"gin": {
		Name:         "gin",
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
				{
					%s := api.Group("/%s")
					{
						%s.%s("", %sHandler.Get%ss)
					}
				}
			`, lowerentity, lowerentity, lowerentity, get, lowerentity, entity)

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
		Returnable:    "",
		ReturnKeyword: "",
		HTTPHandler:   "http.Handler",
	},

	"chi": {
		Name:        "chi",
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
		Returnable:    "",
		ReturnKeyword: "",
		HTTPHandler:   "http.Handler",
	},

	"echo": {
		Name:         "echo",
		Imports:      `"github.com/labstack/echo/v4"`,
		ContextName:  "c",
		ContextType:  "echo.Context",
		Bind:         "c.Bind", // func(obj interface{}) error
		JSON:         "c.JSON", // func(int, interface{}) error
		Router:       "*echo.Echo",
		Start:        "echo.New()",
		OtherImports: `"net/http"`,

		ApiGroup: func(entity, get, lowerentity string) string {
			return fmt.Sprintf(`
            api := r.Group("/api/v1")
            {
                %s := api.Group("/%s")
                {
                    %s.%s("", handler.Get%ss)
                }
            }
        `, lowerentity, lowerentity, lowerentity, get, entity)
		},

		Get: "GET",

		FullContext: "c echo.Context",

		ToTheClient: "c.JSON(http.StatusOK, ",

		Response: "(http.StatusOK,",

		ImportRouter: `
        "net/http"
        "github.com/labstack/echo/v4"
    `,

		ImportHandler: `
        "net/http"
        "github.com/labstack/echo/v4"
    `,
		Returnable:    "error",
		ReturnKeyword: "return",
		HTTPHandler:   "http.Handler",
	},

	"fiber": {
		Name:         "fiber",
		Imports:      `"github.com/gofiber/fiber/v2"`,
		ContextName:  "c",
		ContextType:  "*fiber.Ctx",
		Bind:         "c.BodyParser", // func(obj interface{}) error
		JSON:         "c.JSON",       // func(obj interface{}) error
		Router:       "*fiber.App",
		Start:        "fiber.New()",
		OtherImports: `"net/http"`,

		ApiGroup: func(entity, get, lowerentity string) string {
			return fmt.Sprintf(`
            api := r.Group("/api/v1")
            {
                %s := api.Group("/%s")
                {
                    %s.%s("", handler.Get%ss)
                }
            }
        `, lowerentity, lowerentity, lowerentity, get, entity)
		},

		Get: "Get",

		FullContext: "c *fiber.Ctx",

		ToTheClient: "c.JSON(",

		Response: "(fiber.StatusOK,",

		ImportRouter: `
        "github.com/gofiber/fiber/v2"
    `,

		ImportHandler: `
        "github.com/gofiber/fiber/v2"
    `,

		Returnable:    "error",
		ReturnKeyword: "return",
		HTTPHandler:   "*fiber.App",
	},

	"mux": {
		Name:        "mux",
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
