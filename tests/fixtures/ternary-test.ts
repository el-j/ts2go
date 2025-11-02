// Test ternary operator (conditional expression)
function getMax(a: number, b: number): number {
  return a > b ? a : b;
}

function getStatus(age: number): string {
  return age >= 18 ? "adult" : "minor";
}

function getGrade(score: number): string {
  // Nested ternary
  return score >= 90 ? "A" : score >= 80 ? "B" : "F";
}

// Test the functions
console.log("Max of 10 and 20:", getMax(10, 20));
console.log("Max of 30 and 15:", getMax(30, 15));
console.log("Status of 25:", getStatus(25));
console.log("Status of 15:", getStatus(15));
console.log("Grade 95:", getGrade(95));
console.log("Grade 85:", getGrade(85));
console.log("Grade 65:", getGrade(65));
