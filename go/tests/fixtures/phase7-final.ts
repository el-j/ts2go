// Final Phase 7 Integration Test
// Tests: Union Types, Enums, Tuples, Nullish Coalescing

// Union type with string literals
type Status = "active" | "pending" | "completed";

// String enum
enum LogLevel {
  Debug = "DEBUG",
  Info = "INFO",
  Error = "ERROR"
}

// Numeric enum
enum Priority {
  Low,
  Medium,
  High
}

// Tuple type
type Point = [number, number];

// Interface using all the above
interface Task {
  id: string;
  name: string;
  status: Status;
  level: LogLevel;
  priority: Priority;
  coordinates: Point;
}

// Function using nullish coalescing
function getTaskName(task: Task | null): string {
  return task?.name ?? "Unnamed Task";
}

// Function creating a point (tuple)
function createPoint(x: number, y: number): Point {
  return [x, y];
}

// Function to demonstrate enum usage
function getLogMessage(level: LogLevel): string {
  return "Log level is set";
}

// Function using enum member access
function isHighPriority(priority: Priority): boolean {
  return priority === Priority.High;
}

console.log("Phase 7 final integration test complete!");
