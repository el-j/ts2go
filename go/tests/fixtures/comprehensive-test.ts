// Comprehensive integration test for ts2go transpiler
// Tests all major features implemented

// 1. Type declarations
interface User {
    name: string;
    age: number;
    active: boolean;
}

type Status = "pending" | "active" | "inactive";

enum Color {
    Red,
    Green,
    Blue
}

// 2. Classes with inheritance
class Animal {
    private name: string;
    
    constructor(name: string) {
        this.name = name;
    }
    
    makeSound(): string {
        return "Some sound";
    }
    
    getName(): string {
        return this.name;
    }
}

class Dog extends Animal {
    constructor(name: string) {
        super(name);
    }
    
    makeSound(): string {
        return "Woof!";
    }
}

// 3. Functions with control flow
function testControlFlow(): void {
    // If/else
    const x = 10;
    if (x > 5) {
        console.log("x is greater than 5");
    } else {
        console.log("x is 5 or less");
    }
    
    // For loop
    console.log("For loop:");
    for (let i = 0; i < 3; i++) {
        console.log("  i =", i);
    }
    
    // For...of
    console.log("For...of loop:");
    const fruits = ["apple", "banana"];
    for (const fruit of fruits) {
        console.log("  fruit:", fruit);
    }
    
    // While loop
    console.log("While loop:");
    let count = 0;
    while (count < 2) {
        console.log("  count =", count);
        count++;
    }
    
    // Switch statement
    console.log("Switch:");
    const color = Color.Red;
    switch (color) {
        case Color.Red:
            console.log("  Color is Red");
            break;
        case Color.Green:
            console.log("  Color is Green");
            break;
        default:
            console.log("  Other color");
    }
}

// 4. Arrow functions
function testArrowFunctions(): void {
    const add = (a: number, b: number): number => a + b;
    const multiply = (x: number): number => x * 2;
    
    console.log("Arrow functions:");
    console.log("  add(3, 4) =", add(3, 4));
    console.log("  multiply(5) =", multiply(5));
}

// 5. Classes and objects
function testClassesAndObjects(): void {
    const dog = new Dog("Buddy");
    console.log("Class test:");
    console.log("  Dog name:", dog.getName());
    console.log("  Dog sound:", dog.makeSound());
}

// 6. Arrays and iteration
function testArrays(): void {
    const numbers = [1, 2, 3, 4, 5];
    console.log("Array operations:");
    
    let sum = 0;
    for (const num of numbers) {
        sum = sum + num;
    }
    console.log("  Sum:", sum);
    
    // Find max
    let max = numbers[0];
    for (let i = 1; i < 5; i++) {
        if (numbers[i] > max) {
            max = numbers[i];
        }
    }
    console.log("  Max:", max);
}

// 7. Nested control flow
function testNestedControl(): void {
    console.log("Nested control flow:");
    for (let i = 0; i < 3; i++) {
        for (let j = 0; j < 2; j++) {
            if (i == j) {
                console.log("  i == j at", i);
            }
        }
    }
}

// Main execution
console.log("=== TS2Go Comprehensive Test ===");
console.log("");
testControlFlow();
console.log("");
testArrowFunctions();
console.log("");
testClassesAndObjects();
console.log("");
testArrays();
console.log("");
testNestedControl();
console.log("");
console.log("=== All Tests Complete ===");
