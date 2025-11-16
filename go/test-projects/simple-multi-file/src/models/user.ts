// User model
export interface User {
	id: string;
	name: string;
	email: string;
}

export function createUser(name: string, email: string): User {
	return {
		id: generateId(),
		name: name,
		email: email
	};
}

function generateId(): string {
	return Math.random().toString(36).substr(2, 9);
}
