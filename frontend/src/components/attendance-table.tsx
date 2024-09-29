import React, { useState, useEffect } from 'react';
import { EventAttendance, eventAttendanceTableColumns } from './columns'
import { DataTable } from './data-table'

const AttendanceTable: React.FC = () => {
    const [data, setData] = useState<EventAttendance[]>([]);   
    const [loading, setLoading] = useState<boolean | null>(true);
    const [error, setError] = useState<string | null>(null);

    useEffect(() => {
        const endpoint = "http://localhost:8080/api/events/1/attendance"
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

                const result: EventAttendance[] = await response.json();
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
        <DataTable columns={ eventAttendanceTableColumns } data={data} />
    )
}

export default AttendanceTable 
