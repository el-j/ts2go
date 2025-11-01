// Comprehensive Phase 7 feature test
type Status = "active" | "inactive" | "pending";

enum Color {
  Red = "red",
  Green = "green",
  Blue = "blue"
}

enum Priority {
  Low,
  Medium,
  High
}

type Coordinate = [number, number];

interface Task {
  name: string;
  status: Status;
  color: Color;
  priority: Priority;
  location: Coordinate;
}

function getDefaultStatus(val: string | null): string {
  return val ?? "unknown";
}

function getPriorityLevel(task: Task): Priority {
  return task.priority;
}

function getTaskName(task: Task): string {
  return task.name;
}

console.log("Phase 7 comprehensive test complete");
