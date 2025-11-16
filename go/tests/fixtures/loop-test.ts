// Test all loop types

// Traditional for loop
function testForLoop(): void {
    console.log("Testing traditional for loop:");
    for (let i = 0; i < 5; i++) {
        console.log("i =", i);
    }
}

// For...of loop (array iteration)
function testForOfLoop(): void {
    console.log("Testing for...of loop:");
    const fruits = ["apple", "banana", "cherry"];
    for (const fruit of fruits) {
        console.log("fruit:", fruit);
    }
}

// For...in loop (map key iteration)
function testForInLoop(): void {
    console.log("Testing for...in loop:");
    const numbers = [10, 20, 30];
    for (const index in numbers) {
        console.log("index:", index);
    }
}

// While loop
function testWhileLoop(): void {
    console.log("Testing while loop:");
    let count = 0;
    while (count < 3) {
        console.log("count =", count);
        count++;
    }
}

// Break statement
function testBreak(): void {
    console.log("Testing break:");
    for (let i = 0; i < 10; i++) {
        if (i === 5) {
            console.log("Breaking at", i);
            break;
        }
        console.log("i =", i);
    }
}

// Continue statement
function testContinue(): void {
    console.log("Testing continue:");
    for (let i = 0; i < 5; i++) {
        if (i === 2) {
            console.log("Skipping", i);
            continue;
        }
        console.log("i =", i);
    }
}

// Nested loops
function testNestedLoops(): void {
    console.log("Testing nested loops:");
    for (let i = 0; i < 3; i++) {
        for (let j = 0; j < 2; j++) {
            console.log("i =", i, "j =", j);
        }
    }
}

// Main
testForLoop();
testForOfLoop();
testForInLoop();
testWhileLoop();
testBreak();
testContinue();
testNestedLoops();
