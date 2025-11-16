// Arrow function tests - expression body only
const add = (a: number, b: number): number => a + b;
const multiply = (x: number): number => x * 2;
const isPositive = (n: number): boolean => n > 0;

// Test the functions
console.log("add(3, 5) =", add(3, 5));
console.log("multiply(7) =", multiply(7));
console.log("isPositive(-5) =", isPositive(-5));
console.log("isPositive(10) =", isPositive(10));
