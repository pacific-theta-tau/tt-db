import React, { useEffect, useState } from 'react';
import { useParams } from 'react-router-dom';
import { EventAttendance, eventAttendanceTableColumns } from "../components/columns"
import { DataTable } from "../components/data-table"


const EventAttendancePage: React.FC = () => {
    const { eventID } = useParams<{ eventID: string }>();
    const [data, setData] = useState<EventAttendance[]>([]);
    const [loading, setLoading] = useState<boolean | null>(true);
    const [error, setError] = useState<string | null>(null);

    useEffect(() => {
        const endpoint = "http://localhost:8080/api/events/" + eventID + "/attendance"
        const fetchData = async () => {
             try {
                setLoading(true)
                const response = await fetch(endpoint, {
                    mode: 'cors',
                    headers: {
                        'Content-Type': 'application/json',
                    }
                });
                if (!response.ok) {
                    throw new Error('Network response was not ok');
                }
                const responseData = await response.json()
                const result: EventAttendance[] = responseData.attendance
                console.log('result:', result)
                setData(result);
            } catch (e) {
                setError((e as Error).message);
                console.log('Error fetching data:', error);
                throw error;
            } finally {
                setLoading(false);
            }
        }
        fetchData()
       }, []);

    if (loading) {
        return <div>Loading...</div>
    }

    if (error) {
        return <div>Error: {error}</div>;
    }

    return (
        <DataTable columns={eventAttendanceTableColumns} data={data} />
    )
};

export default EventAttendancePage 

