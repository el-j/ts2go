// Test constructor property assignment
class Animal {
    name: string;
    
    constructor(name: string) {
        this.name = name;
    }
}

const dog = new Animal("Buddy");
console.log(dog.name);
