# my_inventory
Example API created using go lang and MUX router for learning purpose

## Code Structure

```go
type App struct {
    Router *mux.Router
    Db     *sql.DB
}

func (app *App) Initialize() error
func (app *App) Run(address string)
func (app *App) createProduct(w http.ResponseWriter, r *http.Request)
func (app *App) deleteProduct(w http.ResponseWriter, r *http.Request)
func (app *App) getProduct(w http.ResponseWriter, r *http.Request)
func (app *App) getProducts(w http.ResponseWriter, r *http.Request)
func (app *App) handleRoutes()
func (app *App) updateProduct(w http.ResponseWriter, r *http.Request)
```

```go
type product struct { // size=40 (0x28)
    ID       int    `json:"id"`
    Name     string `json:"name"`
    Quantity int    `json:"quantity"`
    Price    int    `json:"price"`
}

func (p *product) createProduct(db *sql.DB) error
func (p *product) deleteProduct(db *sql.DB) error
func (p *product) getProductByID(db *sql.DB) error
func (p *product) updateProduct(db *sql.DB) error
```