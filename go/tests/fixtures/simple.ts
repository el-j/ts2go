// Simple TypeScript example
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
