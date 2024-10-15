import React, { useState, useEffect } from 'react';
import { Event, eventsTableColumns } from './columns'
import { DataTable } from './data-table'
import { Skeleton } from '@/components/ui/skeleton'
import AddRowSheet from './sheet/add-row-sheet'
import { EventsForm } from './sheet/forms/events-form'
import { ApiResponse, getData } from '@/api/api';

const EventsTable: React.FC = () => {
    const [data, setData] = useState<Event[]>([]);   
    const [loading, setLoading] = useState<boolean | null>(true);
    const [error, setError] = useState<string | null>(null);

    useEffect(() => {
        const endpoint = "http://localhost:8080/api/events"
        const fetchData = async () => {
            try {
                setLoading(true)
                const result: ApiResponse<Event[]> = await getData(endpoint)
                console.log('result:', result)
                setData(result.data);
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
        <DataTable
            columns={eventsTableColumns}
            data={data}
            AddSheet={
                () => <AddRowSheet
                    title="Add new event record"
                    description="Refresh the page once you hit submit to see updated table"
                    FormType={<EventsForm />}
                />
            }
        />
    )
}

export default EventsTable
