// Express.js Hello World Server
// Testing: HTTP server, middleware, routing, async handlers

import express from 'express';

const app = express();
const PORT = 3000;

// Middleware
app.use(express.json());

// Simple GET route
app.get('/', (req, res) => {
    res.send('Hello, World!');
});

// GET with params
app.get('/user/:id', (req, res) => {
    const userId = req.params.id;
    res.json({ userId: userId, name: 'John Doe' });
});

// POST route
app.post('/api/data', (req, res) => {
    const data = req.body;
    res.json({ success: true, received: data });
});

// Async route handler
app.get('/api/async', async (req, res) => {
    try {
        const result = await fetchData();
        res.json({ data: result });
    } catch (error) {
        res.status(500).json({ error: error.message });
    }
});

async function fetchData(): Promise<string> {
    return "Async data fetched";
}

// Start server
app.listen(PORT, () => {
    console.log(`Server running on port ${PORT}`);
});
