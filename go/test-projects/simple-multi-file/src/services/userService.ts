import { User, createUser } from '../models/user';

export function getUserService(): string {
	const user: User = createUser('Alice', 'alice@example.com');
	return `User: ${user.name} (${user.email})`;
}

export function getAllUsers(): User[] {
	return [
		createUser('Alice', 'alice@example.com'),
		createUser('Bob', 'bob@example.com')
	];
}
