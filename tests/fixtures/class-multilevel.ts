// Test: Multiple levels of inheritance
class GrandParent {
  public name: string;
  
  constructor(name: string) {
    this.name = name;
  }
  
  public greet(): string {
    return "Hello from GrandParent";
  }
}

class Parent extends GrandParent {
  public age: number;
  
  constructor(name: string, age: number) {
    super(name);
    this.age = age;
  }
  
  public greet(): string {
    return "Hello from Parent";
  }
  
  public introduce(): string {
    return this.name;
  }
}

class Child extends Parent {
  public school: string;
  
  constructor(name: string, age: number, school: string) {
    super(name, age);
    this.school = school;
  }
  
  public greet(): string {
    return "Hello from Child";
  }
  
  public info(): string {
    return this.school;
  }
}

// Test multi-level inheritance
const child = new Child("Alice", 10, "Elementary");
console.log(child.greet());
console.log(child.introduce());
console.log(child.info());
console.log(child.name);
