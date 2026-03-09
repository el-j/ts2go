class Person {
  // Public field (default)
  public name: string;
  
  // Private field
  private age: number;
  
  // Static field
  static species: string = "Homo sapiens";
  
  constructor(name: string, age: number) {
    this.name = name;
    this.age = age;
  }
  
  // Public method (default)
  public greet(): string {
    return "Hello, I'm " + this.name;
  }
  
  // Private method
  private getAge(): number {
    return this.age;
  }
  
  // Static method
  static getSpecies(): string {
    return Person.species;
  }
  
  // Method using private method
  public isAdult(): boolean {
    return this.getAge() >= 18;
  }
}

// Test usage
const person = new Person("Alice", 30);
console.log(person.name);
console.log(person.greet());
console.log(person.isAdult());
console.log(Person.getSpecies());
