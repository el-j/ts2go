// Commander.js CLI Tool
// Testing: CLI parsing, commands, options, arguments

import { Command } from 'commander';

const program = new Command();

program
    .name('mytool')
    .description('A sample CLI tool')
    .version('1.0.0');

// Command with options
program
    .command('greet <name>')
    .description('Greet someone')
    .option('-l, --loud', 'Greet loudly')
    .action((name, options) => {
        const greeting = `Hello, ${name}!`;
        if (options.loud) {
            console.log(greeting.toUpperCase());
        } else {
            console.log(greeting);
        }
    });

// Command with async handler
program
    .command('fetch <url>')
    .description('Fetch data from URL')
    .action(async (url) => {
        try {
            const data = await fetchUrl(url);
            console.log('Data:', data);
        } catch (error) {
            console.error('Error:', error);
        }
    });

async function fetchUrl(url: string): Promise<string> {
    // Simulated async fetch
    return `Data from ${url}`;
}

program.parse(process.argv);
