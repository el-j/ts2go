// Test: Mixed public and private members
class BankAccount {
  public accountNumber: string;
  private _balance: number;
  private _pin: number;
  
  constructor(accountNumber: string, initialBalance: number, pin: number) {
    this.accountNumber = accountNumber;
    this._balance = initialBalance;
    this._pin = pin;
  }
  
  // Public method using private field
  public getBalance(): number {
    return this._balance;
  }
  
  // Public method modifying private field
  public deposit(amount: number): number {
    this._balance = this._balance + amount;
    return this._balance;
  }
  
  // Public method with private validation
  public withdraw(amount: number, pin: number): boolean {
    return this.validatePin(pin);
  }
  
  // Private validation method
  private validatePin(pin: number): boolean {
    return this._pin === pin;
  }
  
  // Private helper
  private calculateInterest(): number {
    return this._balance * 0.05;
  }
  
  // Public method using private helper
  public addInterest(): number {
    const interest = this.calculateInterest();
    this._balance = this._balance + interest;
    return this._balance;
  }
}

// Test mixed access
const account = new BankAccount("ACC001", 1000, 1234);
console.log(account.accountNumber);
console.log(account.getBalance());
console.log(account.deposit(500));
console.log(account.withdraw(200, 1234));
console.log(account.withdraw(200, 9999));
console.log(account.addInterest());
