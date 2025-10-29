// Advanced TypeScript example with functions and types
interface User {
  id: string;
  name: string;
  age: number;
  active: boolean;
}

type UserID = string;

function createUser(name: string, age: number): User {
  return {
    id: "generated-id",
    name: name,
    age: age,
    active: true
  };
}

function getUserName(user: User): string {
  return user.name;
}

const user = createUser("Alice", 30);
const userName = getUserName(user);
console.log(userName);
