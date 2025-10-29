# TS2Go Examples

## Example 1: Simple Types and Functions

**Input (simple.ts):**
```typescript
interface Person {
  name: string;
  age: number;
}

type ID = string;

function greet(person: Person): string {
  return "Hello, " + person.name;
}

const message: string = "Test";
console.log(message);
```

**Output (simple.go):**
```go
package main

import (
	"fmt"
)

type Person struct {
	Name string `json:"name"`
	Age float64 `json:"age"`
}

type ID = string

func Greet(person Person) string {
	return "Hello, " + person.Name
}

func main() {
	message := "Test"
	fmt.Println(message)
}
```

## Example 2: Object Literals and Return Types

**Input (advanced.ts):**
```typescript
interface User {
  id: string;
  name: string;
  age: number;
  active: boolean;
}

type UserID = string;

function createUser(name: string, age: number): User {
  return {
    id: "generated-id",
    name: name,
    age: age,
    active: true
  };
}

function getUserName(user: User): string {
  return user.name;
}

const user = createUser("Alice", 30);
const userName = getUserName(user);
console.log(userName);
```

**Output (advanced.go):**
```go
package main

import (
	"fmt"
)

type User struct {
	Id string `json:"id"`
	Name string `json:"name"`
	Age float64 `json:"age"`
	Active bool `json:"active"`
}

type UserID = string

func CreateUser(name string, age float64) User {
	return User{Id: "generated-id", Name: name, Age: age, Active: true}
}

func GetUserName(user User) string {
	return user.Name
}

func main() {
	user := CreateUser("Alice", 30)
	userName := GetUserName(user)
	fmt.Println(userName)
}
```

## Running Examples

```bash
# Transpile
./ts2go --in example.ts --out example.go

# Run the generated Go code
go run example.go
```

## Key Transformations

1. **Interfaces → Structs** with JSON tags
2. **camelCase → PascalCase** for exported functions
3. **Property access** (e.g., `user.name` → `user.Name`)
4. **console.log → fmt.Println**
5. **Object literals** with type inference
6. **Top-level statements** wrapped in `main()`
