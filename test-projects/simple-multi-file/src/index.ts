import { getUserService, getAllUsers } from './services/userService';

function main(): void {
	console.log('Starting application...');
	
	const userInfo = getUserService();
	console.log(userInfo);
	
	const users = getAllUsers();
	console.log(`Total users: ${users.length}`);
}

main();
