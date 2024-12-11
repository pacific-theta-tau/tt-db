import React, { useEffect, useState } from 'react';
import { useQuery } from "@tanstack/react-query";
import { useParams } from 'react-router-dom';
import { EventAttendance, eventAttendanceTableColumns } from "../components/columns"
import { DataTable } from "../components/data-table"
import { Skeleton } from '@/components/ui/skeleton'
import SideRowSheet from '@/components/sheet/side-row-sheet';
import { EventAttendanceForm } from '@/components/sheet/forms/event-attendance-form';
import { request, ApiResponse } from '@/api/api'

type AttendanceData = {
    attendance: EventAttendance[],
    eventCategory: string,
    eventDate: string,
    eventID: number,
    eventLocation: string,
    eventName: string
}


async function fetchTableData(eventID: string): Promise<EventAttendance[]> {
    const endpoint = "http://localhost:8080/api/events/" + eventID + "/attendance"
    const response: ApiResponse<AttendanceData> = await request(endpoint, 'GET')
    const responseData: AttendanceData = response.data
    return responseData.attendance
}

export const attendanceQueryKey = "attendanceQueryData"

const EventAttendancePage: React.FC = () => {
    const { eventID = "" } = useParams<{ eventID: string }>();
    const { data, isLoading, isError } = useQuery({ queryKey: [attendanceQueryKey], queryFn: () => fetchTableData(eventID) })

    if (isLoading) {
        // Load dummy empty data and skeleton
        const loadingData = Array(5).fill({}) 
        const loadingTableColumns = eventAttendanceTableColumns.map((column) => ({
            ...column,
            cell: () => <Skeleton className="h-12"/>,
          }))
        return <DataTable columns={ loadingTableColumns } data={loadingData} />
    }

    if (isError) {
        return <div>Error loading table data</div>;
    }

    return (
        <DataTable
            columns={eventAttendanceTableColumns}
            data={data ?? []}
            AddSheet={
                () => <SideRowSheet
                        title="Add attendance record"
                        description="Refresh page once you hit submit"
                        FormType={<EventAttendanceForm />}
                      />}
        />
    )
};

export default EventAttendancePage 

