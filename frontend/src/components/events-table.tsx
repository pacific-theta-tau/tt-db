import React, { useState, useEffect } from 'react';
import { Event, eventsTableColumns } from './columns'
import { DataTable } from './data-table'
import { Skeleton } from '@/components/ui/skeleton'

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
                /* uncomment line below to test skeleton during loading */
                // await new Promise(f => setTimeout(f, 3000));
                setLoading(false);
            }
        }

        fetchData();
    }, []);

    if (loading) {
        // Load dummy empty data and skeleton
        const loadingData = Array(5).fill({}) 
        const loadingTableColumns = eventsTableColumns.map((column) => ({
            ...column,
            cell: () => <Skeleton className="h-12"/>,
          }))
        return <DataTable columns={ loadingTableColumns } data={loadingData} />
    }

    if (error) {
        return <div>Error: {error}</div>;
    }

    return (
        <DataTable columns={eventsTableColumns} data={data} />
    )
}

export default EventsTable
