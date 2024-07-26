import { useState, useEffect } from 'react';
import { Event, eventsTableColumns } from './columns'
import { DataTable } from './data-table'

const EventsTable: React.FC = () => {
    const [data, setData] = useState<Event[]>([]);   
    const [loading, setLoading] = useState<boolean | null>(true);
    const [error, setError] = useState<string | null>(null);

    useEffect(() => {
        const endpoint = "http://localhost:8080/api/events"
        const fetchData = async () => {
            try {
                const response = await fetch(endpoint, {
                    mode: 'cors',
                    headers: {
                        'Content-Type': 'application/json',
                    }
                });
                
                if (!response.ok) {
                    throw new Error('Network response was not ok');
                }

                const result: Event[] = await response.json();
                console.log('result:', result)
                setData(result);
            } catch (e) {
                setError((e as Error).message);
                console.log('Error fetching data:', error);
                throw error
            } finally {
                setLoading(false);
            }
        }

        fetchData();
    }, []);

    if (loading) {
        return <div>Loading...</div>
    }

    if (error) {
        return <div>Error: {error}</div>;
    }

    return (
        <DataTable columns={eventsTableColumns} data={data} />
    )
}

export default EventsTable
