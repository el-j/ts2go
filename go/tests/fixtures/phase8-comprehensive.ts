// Phase 8 Comprehensive Test: Advanced Logic Support
// Tests: Classes, Inheritance, Access Modifiers, Static Members, Getters/Setters

// Base class with various features
class Animal {
  private _name: string;
  private _age: number;
  
  // Static field
  static totalAnimals: number = 0;
  
  constructor(name: string, age: number) {
    this._name = name;
    this._age = age;
    Animal.totalAnimals = Animal.totalAnimals + 1;
  }
  
  // Public getter
  get name(): string {
    return this._name;
  }
  
  // Public setter
  set name(value: string) {
    this._name = value;
  }
  
  // Private getter
  private get age(): number {
    return this._age;
  }
  
  // Public method
  public speak(): string {
    return "Some sound";
  }
  
  // Private method
  private getAgeInMonths(): number {
    return this._age * 12;
  }
  
  // Public method using private method
  public describe(): string {
    return this._name;
  }
  
  // Static method
  static getTotal(): number {
    return Animal.totalAnimals;
  }
}

// Derived class with inheritance
class Dog extends Animal {
  private _breed: string;
  
  constructor(name: string, age: number, breed: string) {
    super(name, age);
    this._breed = breed;
  }
  
  // Getter for breed
  get breed(): string {
    return this._breed;
  }
  
  // Override parent method
  public speak(): string {
    return "Woof!";
  }
  
  // New method in derived class
  public fetch(): string {
    return "Fetching ball";
  }
}

// Another derived class
class Cat extends Animal {
  private _indoor: boolean;
  
  constructor(name: string, age: number, indoor: boolean) {
    super(name, age);
    this._indoor = indoor;
  }
  
  // Override parent method
  public speak(): string {
    return "Meow!";
  }
  
  // Cat-specific method
  public isIndoor(): boolean {
    return this._indoor;
  }
}

// Test all features
const dog = new Dog("Buddy", 3, "Golden Retriever");
console.log(dog.speak());
console.log(dog.fetch());

const cat = new Cat("Whiskers", 5, true);
console.log(cat.speak());
console.log(cat.isIndoor());

// Test static method
console.log(Animal.getTotal());
