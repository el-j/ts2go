<script setup lang="ts">
import { ref } from 'vue'
import DataView from 'primevue/dataview'
import Card from 'primevue/card'
import Button from 'primevue/button'
import Tag from 'primevue/tag'

export interface Example {
  id: string
  title: string
  description: string
  category: 'basic' | 'intermediate' | 'advanced'
  typescript: string
  go: string
}

const emit = defineEmits<{
  loadExample: [example: Example]
}>()

const examples = ref<Example[]>([
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
	Name  string \`json:"name"\`
	Age   float64 \`json:"age"\`
	Email string \`json:"email"\`
}

type Address struct {
	Street  string \`json:"street"\`
	City    string \`json:"city"\`
	ZipCode string \`json:"zipCode"\`
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
	name string
}

func NewAnimal(name string) *Animal {
	return &Animal{name: name}
}

func (a *Animal) Speak() string {
	return fmt.Sprintf("%s makes a sound", a.name)
}

type Dog struct {
	Animal
}

func NewDog(name string) *Dog {
	return &Dog{Animal: *NewAnimal(name)}
}

func (d *Dog) Speak() string {
	return fmt.Sprintf("%s - Woof!", d.Animal.Speak())
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
	return a + b
}

func Greet(name string) string {
	return fmt.Sprintf("Hello, %s!", name)
}

func ProcessArray(items []string) float64 {
	return float64(len(items))
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
	ColorRed Color = iota
	ColorGreen
	ColorBlue
)

type Status string

const (
	StatusActive   Status = "ACTIVE"
	StatusInactive Status = "INACTIVE"
	StatusPending  Status = "PENDING"
)

func GetColorName(color Color) string {
	switch color {
	case ColorRed:
		return "Red"
	case ColorGreen:
		return "Green"
	case ColorBlue:
		return "Blue"
	default:
		return ""
	}
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
	resultCh := make(chan interface{}, 1)
	go func() {
		defer func() {
			if r := recover(); r != nil {
				resultCh <- fmt.Errorf("Failed to fetch: %v", r)
			}
		}()
		response := (<-Fetch(url))
		data := (<-response.Text())
		resultCh <- data
	}()
	return resultCh
}

func ProcessData() chan interface{} {
	resultCh := make(chan interface{}, 1)
	go func() {
		data := (<-FetchData("https://api.example.com"))
		fmt.Println(data)
		resultCh <- nil
	}()
	return resultCh
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
	sum := 0.0
	
	for _, num := range numbers {
		if num > 0 {
			sum += num
		} else if num < 0 {
			sum -= num
		}
	}
	
	return sum
}

func Categorize(value float64) string {
	switch {
	case value < 0:
		return "negative"
	case value == 0:
		return "zero"
	case value > 0:
		return "positive"
	default:
		return "unknown"
	}
}`
  }
])

const selectedCategory = ref<string>('all')

const filteredExamples = ref(examples.value)

function filterByCategory(category: string) {
  selectedCategory.value = category
  if (category === 'all') {
    filteredExamples.value = examples.value
  } else {
    filteredExamples.value = examples.value.filter(ex => ex.category === category)
  }
}

function getCategoryColor(category: string): string {
  switch (category) {
    case 'basic': return 'success'
    case 'intermediate': return 'warn'
    case 'advanced': return 'danger'
    default: return 'info'
  }
}
</script>

<template>
  <div class="example-gallery">
    <!-- Category Filter -->
    <div class="mb-4 flex gap-2">
      <Button 
        label="All" 
        :outlined="selectedCategory !== 'all'"
        @click="filterByCategory('all')"
        size="small"
      />
      <Button 
        label="Basic" 
        :outlined="selectedCategory !== 'basic'"
        severity="success"
        @click="filterByCategory('basic')"
        size="small"
      />
      <Button 
        label="Intermediate" 
        :outlined="selectedCategory !== 'intermediate'"
        severity="warn"
        @click="filterByCategory('intermediate')"
        size="small"
      />
      <Button 
        label="Advanced" 
        :outlined="selectedCategory !== 'advanced'"
        severity="danger"
        @click="filterByCategory('advanced')"
        size="small"
      />
    </div>

    <!-- Examples Grid -->
    <DataView :value="filteredExamples" layout="grid">
      <template #grid="slotProps">
        <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
          <Card 
            v-for="example in slotProps.items" 
            :key="example.id"
            class="hover:shadow-lg transition-shadow cursor-pointer"
          >
            <template #title>
              <div class="flex items-center justify-between">
                <span class="text-lg">{{ example.title }}</span>
                <Tag 
                  :value="example.category" 
                  :severity="getCategoryColor(example.category)"
                />
              </div>
            </template>
            <template #content>
              <p class="text-sm text-gray-600 dark:text-gray-400 mb-4">
                {{ example.description }}
              </p>
              <div class="flex gap-2">
                <Button 
                  label="Load Example" 
                  icon="pi pi-code" 
                  size="small"
                  @click="emit('loadExample', example)"
                />
              </div>
            </template>
          </Card>
        </div>
      </template>
    </DataView>
  </div>
</template>

<style scoped>
.example-gallery {
  padding: 1rem;
}
</style>
