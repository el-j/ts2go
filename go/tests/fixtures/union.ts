// Union types example
type StringOrNumber = string | number;
type Status = "active" | "inactive" | "pending";

// Simple function that just returns a string
function getTypeName(): string {
  return "union";
}

console.log(getTypeName());