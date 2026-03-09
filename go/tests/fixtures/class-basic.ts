// Basic class test
class Person {
  name: string;
  age: number;

  constructor(name: string, age: number) {
    this.name = name;
    this.age = age;
  }

  greet(): string {
    return "Hello, I'm " + this.name;
  }

  getAge(): number {
    return this.age;
  }
}

const person = new Person("Alice", 30);
console.log(person.greet());
console.log(person.getAge());
