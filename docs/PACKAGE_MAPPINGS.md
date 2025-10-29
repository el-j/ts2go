# Quick Reference: npm to Go Mappings

## Node.js Built-in Modules

| npm Package | Go Equivalent | Status | Notes |
|-------------|---------------|--------|-------|
| `fs` | `github.com/ts2go/runtime/fs` or `os` | ✅ | Runtime wrapper provided |
| `path` | `path/filepath` | ✅ | Direct stdlib mapping |
| `http` / `https` | `net/http` | ✅ | Requires API restructuring |
| `process` | `os` | 🔄 | Partial support |
| `crypto` | `crypto/*` packages | 🔄 | Different API structure |
| `buffer` | `[]byte` | ✅ | Direct type mapping |
| `stream` | `io.Reader/Writer` | 🔄 | Different abstraction |
| `events` | Custom EventEmitter | 🔧 | Runtime implementation needed |
| `url` | `net/url` | ✅ | Direct stdlib mapping |
| `querystring` | `net/url.Values` | ✅ | Direct stdlib mapping |

## Popular npm Packages

### Web Frameworks

| npm Package | Go Equivalent | Difficulty | Notes |
|-------------|---------------|------------|-------|
| `express` | `gin-gonic/gin` | ⚠️ Medium | Context-based, middleware mapping |
| `koa` | `gin-gonic/gin` | ⚠️ Medium | Async middleware requires rework |
| `fastify` | `fiber` | ⚠️ Medium | Performance-focused alternative |
| `nest` | `go-zero` | 🔴 Hard | Complex DI and decorator system |

**Express → Gin Example:**
```typescript
// TypeScript
const express = require('express');
const app = express();

app.get('/users/:id', (req, res) => {
    res.json({ id: req.params.id });
});

app.listen(3000);
```
```go
// Go
import "github.com/gin-gonic/gin"

func main() {
    router := gin.Default()
    
    router.GET("/users/:id", func(c *gin.Context) {
        c.JSON(200, gin.H{"id": c.Param("id")})
    })
    
    router.Run(":3000")
}
```

### HTTP Clients

| npm Package | Go Equivalent | Difficulty | Notes |
|-------------|---------------|------------|-------|
| `axios` | `go-resty/resty` | ✅ Easy | Very similar API |
| `node-fetch` | `net/http` | ✅ Easy | Standard library |
| `got` | `go-resty/resty` | ✅ Easy | Similar features |
| `superagent` | `go-resty/resty` | ✅ Easy | Alternative |

**Axios → Resty Example:**
```typescript
// TypeScript
const response = await axios.get('https://api.com/users', {
    headers: { 'Authorization': 'Bearer token' }
});
```
```go
// Go
client := resty.New()
resp, err := client.R().
    SetHeader("Authorization", "Bearer token").
    Get("https://api.com/users")
```

### Utility Libraries

| npm Package | Go Equivalent | Difficulty | Notes |
|-------------|---------------|------------|-------|
| `lodash` | `samber/lo` | ✅ Easy | Excellent Go alternative |
| `ramda` | `samber/lo` | ✅ Easy | Functional programming |
| `moment` | `time` stdlib | ⚠️ Medium | Different API, layout format |
| `date-fns` | `time` stdlib | ⚠️ Medium | Standard time package |
| `uuid` | `google/uuid` | ✅ Easy | Direct equivalent |
| `nanoid` | `matoous/go-nanoid` | ✅ Easy | Direct equivalent |

**Lodash → Lo Example:**
```typescript
// TypeScript
import _ from 'lodash';

const doubled = _.map([1, 2, 3], x => x * 2);
const filtered = _.filter(users, u => u.active);
```
```go
// Go
import "github.com/samber/lo"

doubled := lo.Map([]int{1, 2, 3}, func(x int, _ int) int {
    return x * 2
})
filtered := lo.Filter(users, func(u User, _ int) bool {
    return u.Active
})
```

### Validation & Parsing

| npm Package | Go Equivalent | Difficulty | Notes |
|-------------|---------------|------------|-------|
| `joi` | `go-playground/validator` | ⚠️ Medium | Struct tags vs schema |
| `yup` | `go-playground/validator` | ⚠️ Medium | Struct tags vs schema |
| `zod` | `go-playground/validator` | ⚠️ Medium | Struct tags vs schema |
| `ajv` | `xeipuuv/gojsonschema` | ⚠️ Medium | JSON Schema validation |
| `validator` | Can transpile | ✅ Easy | Pure TypeScript |

**Joi → Validator Example:**
```typescript
// TypeScript
const schema = Joi.object({
    email: Joi.string().email().required(),
    age: Joi.number().min(18).max(100)
});
```
```go
// Go
type User struct {
    Email string `validate:"required,email"`
    Age   int    `validate:"min=18,max=100"`
}

validate := validator.New()
err := validate.Struct(user)
```

### Database & ORM

| npm Package | Go Equivalent | Difficulty | Notes |
|-------------|---------------|------------|-------|
| `typeorm` | `gorm` | 🔴 Hard | Different patterns |
| `prisma` | `ent` | 🔴 Hard | Schema-first approach |
| `mongoose` | `mongo-go-driver` | ⚠️ Medium | Different API |
| `sequelize` | `gorm` | ⚠️ Medium | Active Record pattern |
| `knex` | `squirrel` | ✅ Easy | Query builder |

### Authentication

| npm Package | Go Equivalent | Difficulty | Notes |
|-------------|---------------|------------|-------|
| `jsonwebtoken` | `golang-jwt/jwt` | ✅ Easy | Direct equivalent |
| `passport` | `markbates/goth` | ⚠️ Medium | OAuth strategies |
| `bcrypt` | `golang.org/x/crypto/bcrypt` | ✅ Easy | Standard library |
| `argon2` | `golang.org/x/crypto/argon2` | ✅ Easy | Standard library |

**JWT Example:**
```typescript
// TypeScript
import jwt from 'jsonwebtoken';

const token = jwt.sign({ userId: 123 }, 'secret');
const decoded = jwt.verify(token, 'secret');
```
```go
// Go
import "github.com/golang-jwt/jwt/v5"

token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
    "userId": 123,
})
tokenString, _ := token.SignedString([]byte("secret"))

parsed, _ := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
    return []byte("secret"), nil
})
```

### Environment & Config

| npm Package | Go Equivalent | Difficulty | Notes |
|-------------|---------------|------------|-------|
| `dotenv` | `joho/godotenv` | ✅ Easy | Direct equivalent |
| `config` | `spf13/viper` | ✅ Easy | More features |
| `convict` | `spf13/viper` | ✅ Easy | Schema validation |

### Testing

| npm Package | Go Equivalent | Difficulty | Notes |
|-------------|---------------|------------|-------|
| `jest` | `testing` stdlib + `stretchr/testify` | ⚠️ Medium | Different structure |
| `mocha` | `testing` stdlib | ⚠️ Medium | Different structure |
| `chai` | `stretchr/testify/assert` | ✅ Easy | Assertions |
| `supertest` | `net/http/httptest` | ✅ Easy | HTTP testing |

### CLI Tools

| npm Package | Go Equivalent | Difficulty | Notes |
|-------------|---------------|------------|-------|
| `commander` | `spf13/cobra` | ✅ Easy | Similar structure |
| `yargs` | `spf13/cobra` | ✅ Easy | Command parsing |
| `inquirer` | `AlecAivazis/survey` | ✅ Easy | Interactive prompts |
| `chalk` | `fatih/color` | ✅ Easy | Terminal colors |
| `ora` | `briandowns/spinner` | ✅ Easy | Loading spinners |

### File Processing

| npm Package | Go Equivalent | Difficulty | Notes |
|-------------|---------------|------------|-------|
| `fs-extra` | `os` + custom helpers | ✅ Easy | Additional file operations |
| `glob` | `doublestar` | ✅ Easy | Pattern matching |
| `chokidar` | `fsnotify` | ✅ Easy | File watching |
| `archiver` | `archive/zip` | ✅ Easy | Standard library |

### Logging

| npm Package | Go Equivalent | Difficulty | Notes |
|-------------|---------------|------------|-------|
| `winston` | `sirupsen/logrus` | ✅ Easy | Structured logging |
| `pino` | `rs/zerolog` | ✅ Easy | Fast JSON logger |
| `bunyan` | `rs/zerolog` | ✅ Easy | Structured logging |
| `debug` | `log` stdlib | ✅ Easy | Simple debugging |

### Data Processing

| npm Package | Go Equivalent | Difficulty | Notes |
|-------------|---------------|------------|-------|
| `csv-parse` | `encoding/csv` | ✅ Easy | Standard library |
| `xml2js` | `encoding/xml` | ✅ Easy | Standard library |
| `yaml` | `gopkg.in/yaml.v3` | ✅ Easy | Direct equivalent |
| `jsdom` | `PuerkitoBio/goquery` | ⚠️ Medium | HTML parsing |

## Unsupported Categories

### Frontend Frameworks (Not Applicable)
- ❌ `react`, `vue`, `angular`, `svelte`
- **Alternative:** Keep frontend separate or use Go templates (`html/template`, `templ`)

### Build Tools (Not Needed in Go)
- ❌ `webpack`, `rollup`, `vite`, `parcel`
- **Alternative:** Go's built-in build system

### Browser APIs (Not Applicable)
- ❌ `dom`, `window`, `localStorage`
- **Alternative:** Server-side alternatives

### Electron/Desktop (Different Ecosystem)
- ❌ `electron`
- **Alternative:** `wails`, `fyne` for Go desktop apps

## Auto-Transpilable Packages

These packages are pure TypeScript and can be automatically transpiled:

✅ `validator` - String validation  
✅ `@types/*` - Type definitions (converted to Go types)  
✅ Most utility libraries without native dependencies  
✅ Pure algorithm implementations  

## Legend

- ✅ **Easy** - Direct mapping, minimal changes
- ⚠️ **Medium** - API restructuring required
- 🔴 **Hard** - Significant architectural changes
- 🔧 **Runtime** - Custom runtime implementation needed
- 🔄 **Partial** - Some features supported
- ❌ **Unsupported** - Not applicable or no equivalent

## How to Use This Guide

1. **Check your `package.json`** dependencies
2. **Look up each package** in this guide
3. **Note the difficulty level** and Go equivalent
4. **Plan migration strategy** based on complexity
5. **Use `ts2go analyze`** to get automated analysis

## Contributing

Found a good npm→Go mapping? Add it to:
- `internal/deps/mappings.go` (code)
- This guide (documentation)
- `docs/DEPENDENCY_GUIDE.md` (detailed implementation)

Submit a PR with examples and tests!
