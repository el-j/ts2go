// Test optional chaining operator
interface User {
  name: string;
  address?: {
    street: string;
    city: string;
  };
}

function getCity(user: User): string {
  return user.address?.city ?? "Unknown";
}

function getName(user: User): string {
  return user?.name ?? "Anonymous";
}

console.log("Optional chaining test complete");
