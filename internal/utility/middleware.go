package utility

// import (
// 	"encoding/json"
// 	"io"
// 	"net/http"

// 	"github.com/DilankaHer/sop-in-go/internal/app"
// 	"github.com/DilankaHer/sop-in-go/internal/logger"
// )

// func Middleware(app *app.App, handlerFunc func(http.ResponseWriter, *http.Request)) func(http.ResponseWriter, *http.Request) {
// 	return func(w http.ResponseWriter, r *http.Request) {
// 		if r.URL.Query() != nil {
// 			p, err := json.Marshal(r.URL.Query())
// 			if err != nil {
// 				app.Logger.Printf("[ACCESS] error marshaling request params: %s", err.Error())
// 			}
// 			access := logger.Access{
// 				Level:         "Access",
// 				Method:        r.Method,
// 				Path:          r.URL.Path,
// 				RequestParams: string(p),
// 			}
// 			r.SetPathValue("params", string(p))
// 			app.Logger.Access(access)
// 		} else {
// 			body, err := io.ReadAll(r.Body)
// 			if err != nil {
// 				app.Logger.Printf("[ACCESS] error reading request body: %s", err.Error())
// 			}
// 			access := logger.Access{
// 				Level:       "Access",
// 				Method:      r.Method,
// 				Path:        r.URL.Path,
// 				RequestBody: string(body),
// 			}
// 			r.SetPathValue("body", string(body))
// 			app.Logger.Access(access)
// 		}
// 		handlerFunc(w, r)
// 	}
// }
