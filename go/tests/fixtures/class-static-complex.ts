// Test: Static fields and methods with complex interactions
class MathHelper {
  static PI: number = 3.14159;
  static E: number = 2.71828;
  
  // Static method accessing static field
  static circleArea(radius: number): number {
    return MathHelper.PI * radius * radius;
  }
  
  // Static method accessing another static field
  static exponential(power: number): number {
    return MathHelper.E * power;
  }
  
  // Static method accessing multiple static members
  static calculate(): number {
    return MathHelper.PI + MathHelper.E;
  }
}

class Counter {
  private _value: number;
  static instances: number = 0;
  
  constructor(value: number) {
    this._value = value;
    Counter.instances = Counter.instances + 1;
  }
  
  public getValue(): number {
    return this._value;
  }
  
  static getInstanceCount(): number {
    return Counter.instances;
  }
  
  static resetCount(): number {
    Counter.instances = 0;
    return Counter.instances;
  }
}

// Test static methods
console.log(MathHelper.circleArea(5));
console.log(MathHelper.exponential(2));
console.log(MathHelper.calculate());

// Test instance tracking with static
const c1 = new Counter(10);
const c2 = new Counter(20);
const c3 = new Counter(30);
console.log(Counter.getInstanceCount());
console.log(c1.getValue());
console.log(c2.getValue());
console.log(c3.getValue());
console.log(Counter.resetCount());
