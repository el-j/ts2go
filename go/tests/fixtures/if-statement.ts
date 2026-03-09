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

function checkAccess(role: string, level: number): boolean {
  if (role === "admin") {
    return true;
  }
  
  if (level > 5) {
    return true;
  }
  
  return false;
}
