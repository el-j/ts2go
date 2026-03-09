class Rectangle {
  private _width: number;
  private _height: number;
  
  constructor(width: number, height: number) {
    this._width = width;
    this._height = height;
  }
  
  // Getter for width
  get width(): number {
    return this._width;
  }
  
  // Setter for width
  set width(value: number) {
    this._width = value;
  }
  
  // Getter for height
  get height(): number {
    return this._height;
  }
  
  // Setter for height
  set height(value: number) {
    this._height = value;
  }
  
  // Computed property (getter only)
  get area(): number {
    return this._width * this._height;
  }
  
  // Method using getters
  public describe(): string {
    return "Rectangle with area";
  }
}

// Test usage - Note: In Go, getters/setters are methods, not properties
// So rect.width becomes rect.GetWidth() and rect.width = 20 becomes rect.SetWidth(20)
const rect = new Rectangle(10, 5);
console.log(rect.describe());
