import React, { useEffect, useState } from 'react';
import { useParams } from 'react-router-dom';
import { EventAttendance, eventAttendanceTableColumns } from "../components/columns"
import { DataTable } from "../components/data-table"
import { Skeleton } from '@/components/ui/skeleton'
import AddRowSheet from '@/components/sheet/add-row-sheet';
import { EventAttendanceForm } from '@/components/sheet/forms/event-attendance-form';
import { getData, ApiResponse } from '@/api/api'

type AttendanceData = {
    attendance: EventAttendance[],
    eventCategory: string,
    eventDate: string,
    eventID: number,
    eventLocation: string,
    eventName: string
}

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
                const response: ApiResponse<AttendanceData> = await getData(endpoint)
                const responseData: AttendanceData = response.data
                // const result: AttendanceData[] = response.data !== null ? response.data : []
                console.log('response:', response)
                console.log('attendance:', responseData.attendance)
                setData(responseData.attendance);
            } catch (e) {
                setError((e as Error).message);
                console.log('Error fetching data:', error);
                throw error;
            } finally {
                /* uncomment line below to test skeleton during loading */
                // await new Promise(f => setTimeout(f, 3000))
                setLoading(false);
            }
        }
        fetchData()
       }, []);

    if (loading) {
        // Load dummy empty data and skeleton
        const loadingData = Array(5).fill({}) 
        const loadingTableColumns = eventAttendanceTableColumns.map((column) => ({
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
            columns={eventAttendanceTableColumns}
            data={data}
            AddSheet={
                () => <AddRowSheet
                        title="Add attendance record"
                        description="Refresh page once you hit submit"
                        FormType={<EventAttendanceForm />}
                      />}
        />
    )
};

export default EventAttendancePage 

