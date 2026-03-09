// Test if/else statement transpilation
function checkAge(age: number): string {
  if (age >= 18) {
    return "adult";
  } else {
    return "minor";
  }
}

function getStatus(score: number): string {
  if (score >= 90) {
    return "A";
  } else if (score >= 80) {
    return "B";
  } else if (score >= 70) {
    return "C";
  } else {
    return "F";
  }
}

// Test the functions
console.log("Age 25:", checkAge(25));
console.log("Age 15:", checkAge(15));
console.log("Score 95:", getStatus(95));
console.log("Score 85:", getStatus(85));
console.log("Score 75:", getStatus(75));
console.log("Score 65:", getStatus(65));
