// Comprehensive switch tests

// Basic switch with break
function testBasicSwitch(): void {
    console.log("Testing basic switch:");
    const day = 3;
    switch (day) {
        case 1:
            console.log("Monday");
            break;
        case 2:
            console.log("Tuesday");
            break;
        case 3:
            console.log("Wednesday");
            break;
        default:
            console.log("Other day");
    }
}

// Switch with string values
function testStringSwitch(): void {
    console.log("Testing string switch:");
    const color = "red";
    switch (color) {
        case "red":
            console.log("Red color");
            break;
        case "blue":
            console.log("Blue color");
            break;
        case "green":
            console.log("Green color");
            break;
        default:
            console.log("Unknown color");
    }
}

// Switch with fall-through (no break)
function testFallThrough(): void {
    console.log("Testing fall-through:");
    const grade = "B";
    switch (grade) {
        case "A":
        case "B":
            console.log("Excellent or Good");
            break;
        case "C":
            console.log("Average");
            break;
        default:
            console.log("Below average");
    }
}

// Switch with multiple statements per case
function testMultipleStatements(): void {
    console.log("Testing multiple statements:");
    const num = 2;
    switch (num) {
        case 1:
            console.log("Case 1: Line 1");
            console.log("Case 1: Line 2");
            break;
        case 2:
            console.log("Case 2: Line 1");
            console.log("Case 2: Line 2");
            console.log("Case 2: Line 3");
            break;
        default:
            console.log("Default: Line 1");
            console.log("Default: Line 2");
    }
}

// Nested switch
function testNestedSwitch(): void {
    console.log("Testing nested switch:");
    const outer = 1;
    const inner = "a";
    
    switch (outer) {
        case 1:
            console.log("Outer case 1");
            switch (inner) {
                case "a":
                    console.log("Inner case a");
                    break;
                case "b":
                    console.log("Inner case b");
                    break;
            }
            break;
        case 2:
            console.log("Outer case 2");
            break;
    }
}

// Main
testBasicSwitch();
testStringSwitch();
testFallThrough();
testMultipleStatements();
testNestedSwitch();
