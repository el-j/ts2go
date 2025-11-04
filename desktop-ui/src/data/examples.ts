export interface Example {
  id: string
  title: string
  description: string
  category: 'basic' | 'intermediate' | 'advanced'
  typescript: string
  go: string
}

export const EXAMPLES: Example[] = [
  {
    id: 'interface',
    title: 'Interface Definition',
    description: 'Basic interface transpilation with multiple field types',
    category: 'basic',
    typescript: `interface Person {
  name: string;
  age: number;
  email: string;
}

interface Address {
  street: string;
  city: string;
  zipCode: string;
}`,
    go: `type Person struct {
\tName  string \`json:"name"\`
\tAge   float64 \`json:"age"\`
\tEmail string \`json:"email"\`
}

type Address struct {
\tStreet  string \`json:"street"\`
\tCity    string \`json:"city"\`
\tZipCode string \`json:"zipCode"\`
}`
  },
  {
    id: 'class',
    title: 'Class with Inheritance',
    description: 'Class transpilation with constructor, methods, and inheritance',
    category: 'intermediate',
    typescript: `class Animal {
  private name: string;
  
  constructor(name: string) {
    this.name = name;
  }
  
  speak(): string {
    return \`\${this.name} makes a sound\`;
  }
}

class Dog extends Animal {
  constructor(name: string) {
    super(name);
  }
  
  speak(): string {
    return \`\${super.speak()} - Woof!\`;
  }
}`,
    go: `type Animal struct {
\tname string
}

func NewAnimal(name string) *Animal {
\treturn &Animal{name: name}
}

func (a *Animal) Speak() string {
\treturn fmt.Sprintf("%s makes a sound", a.name)
}

type Dog struct {
\tAnimal
}

func NewDog(name string) *Dog {
\treturn &Dog{Animal: *NewAnimal(name)}
}

func (d *Dog) Speak() string {
\treturn fmt.Sprintf("%s - Woof!", d.Animal.Speak())
}`
  },
  {
    id: 'function',
    title: 'Function Declarations',
    description: 'Various function types with typed parameters and return values',
    category: 'basic',
    typescript: `function add(a: number, b: number): number {
  return a + b;
}

function greet(name: string): string {
  return \`Hello, \${name}!\`;
}

function processArray(items: string[]): number {
  return items.length;
}`,
    go: `func Add(a float64, b float64) float64 {
\treturn a + b
}

func Greet(name string) string {
\treturn fmt.Sprintf("Hello, %s!", name)
}

func ProcessArray(items []string) float64 {
\treturn float64(len(items))
}`
  },
  {
    id: 'enum',
    title: 'Enum Definitions',
    description: 'Both numeric and string enums with usage',
    category: 'basic',
    typescript: `enum Color {
  Red,
  Green,
  Blue
}

enum Status {
  Active = "ACTIVE",
  Inactive = "INACTIVE",
  Pending = "PENDING"
}

function getColorName(color: Color): string {
  switch (color) {
    case Color.Red: return "Red";
    case Color.Green: return "Green";
    case Color.Blue: return "Blue";
  }
}`,
    go: `type Color int

const (
\tColorRed Color = iota
\tColorGreen
\tColorBlue
)

type Status string

const (
\tStatusActive   Status = "ACTIVE"
\tStatusInactive Status = "INACTIVE"
\tStatusPending  Status = "PENDING"
)

func GetColorName(color Color) string {
\tswitch color {
\tcase ColorRed:
\t\treturn "Red"
\tcase ColorGreen:
\t\treturn "Green"
\tcase ColorBlue:
\t\treturn "Blue"
\tdefault:
\t\treturn ""
\t}
}`
  },
  {
    id: 'async',
    title: 'Async/Await with Promises',
    description: 'Asynchronous functions with error handling',
    category: 'advanced',
    typescript: `async function fetchData(url: string): Promise<string> {
  try {
    const response = await fetch(url);
    const data = await response.text();
    return data;
  } catch (error) {
    throw new Error(\`Failed to fetch: \${error}\`);
  }
}

async function processData(): Promise<void> {
  const data = await fetchData("https://api.example.com");
  console.log(data);
}`,
    go: `func FetchData(url string) chan interface{} {
\tresultCh := make(chan interface{}, 1)
\tgo func() {
\t\tdefer func() {
\t\t\tif r := recover(); r != nil {
\t\t\t\tresultCh <- fmt.Errorf("Failed to fetch: %v", r)
\t\t\t}
\t\t}()
\t\tresponse := (<-Fetch(url))
\t\tdata := (<-response.Text())
\t\tresultCh <- data
\t}()
\treturn resultCh
}

func ProcessData() chan interface{} {
\tresultCh := make(chan interface{}, 1)
\tgo func() {
\t\tdata := (<-FetchData("https://api.example.com"))
\t\tfmt.Println(data)
\t\tresultCh <- nil
\t}()
\treturn resultCh
}`
  },
  {
    id: 'control-flow',
    title: 'Control Flow Structures',
    description: 'If/else, loops, and switch statements',
    category: 'basic',
    typescript: `function processNumbers(numbers: number[]): number {
  let sum = 0;
  
  for (const num of numbers) {
    if (num > 0) {
      sum += num;
    } else if (num < 0) {
      sum -= num;
    }
  }
  
  return sum;
}

function categorize(value: number): string {
  switch (true) {
    case value < 0:
      return "negative";
    case value === 0:
      return "zero";
    case value > 0:
      return "positive";
    default:
      return "unknown";
  }
}`,
    go: `func ProcessNumbers(numbers []float64) float64 {
\tsum := 0.0
\t
\tfor _, num := range numbers {
\t\tif num > 0 {
\t\t\tsum += num
\t\t} else if num < 0 {
\t\t\tsum -= num
\t\t}
\t}
\t
\treturn sum
}

func Categorize(value float64) string {
\tswitch {
\tcase value < 0:
\t\treturn "negative"
\tcase value == 0:
\t\treturn "zero"
\tcase value > 0:
\t\treturn "positive"
\tdefault:
\t\treturn "unknown"
\t}
}`
  }
]
