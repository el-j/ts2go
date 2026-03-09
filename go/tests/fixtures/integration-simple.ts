// Simpler integration test

class Animal {
    private name: string;
    
    constructor(name: string) {
        this.name = name;
    }
    
    getName(): string {
        return this.name;
    }
}

function testBasic(): void {
    console.log("Basic test:");
    const x = 10;
    if (x > 5) {
        console.log("  x > 5");
    }
    
    for (let i = 0; i < 3; i++) {
        console.log("  i =", i);
    }
    
    const nums = [1, 2, 3];
    for (const n of nums) {
        console.log("  n =", n);
    }
}

function testArrows(): void {
    console.log("Arrow functions:");
    const add = (a: number, b: number): number => a + b;
    console.log("  add(2, 3) =", add(2, 3));
}

console.log("=== Simple Integration Test ===");
testBasic();
testArrows();
console.log("=== Complete ===");
