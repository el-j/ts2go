// Data Processing Script
// Testing: Arrays, objects, loops, functional programming, error handling

interface DataRecord {
    id: number;
    name: string;
    value: number;
    active: boolean;
}

class DataProcessor {
    private data: DataRecord[];

    constructor() {
        this.data = [];
    }

    addRecord(record: DataRecord): void {
        this.data.push(record);
    }

    // Filter active records
    getActiveRecords(): DataRecord[] {
        return this.data.filter(record => record.active);
    }

    // Map to names
    getNames(): string[] {
        return this.data.map(record => record.name);
    }

    // Reduce to sum
    getTotalValue(): number {
        return this.data.reduce((sum, record) => sum + record.value, 0);
    }

    // Find by ID
    findById(id: number): DataRecord | null {
        const record = this.data.find(r => r.id === id);
        return record || null;
    }

    // Process with error handling
    async processRecords(): Promise<void> {
        try {
            for (const record of this.data) {
                await this.processRecord(record);
            }
            console.log('All records processed');
        } catch (error) {
            console.error('Processing failed:', error);
            throw error;
        }
    }

    private async processRecord(record: DataRecord): Promise<void> {
        // Simulated async processing
        if (!record.active) {
            throw new Error(`Record ${record.id} is not active`);
        }
        // Process logic here
    }

    // Destructuring example
    printRecord({ id, name, value }: DataRecord): void {
        console.log(`ID: ${id}, Name: ${name}, Value: ${value}`);
    }

    // Spread operator example
    mergeData(...processors: DataProcessor[]): DataRecord[] {
        const allData: DataRecord[] = [];
        for (const processor of processors) {
            allData.push(...processor.data);
        }
        return allData;
    }
}

// Usage example
const processor = new DataProcessor();
processor.addRecord({ id: 1, name: 'Record 1', value: 100, active: true });
processor.addRecord({ id: 2, name: 'Record 2', value: 200, active: false });
processor.addRecord({ id: 3, name: 'Record 3', value: 300, active: true });

const activeRecords = processor.getActiveRecords();
const totalValue = processor.getTotalValue();

console.log('Active records:', activeRecords.length);
console.log('Total value:', totalValue);
