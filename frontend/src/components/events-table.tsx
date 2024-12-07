import React, { useState, useEffect } from 'react';
import { useQuery } from "@tanstack/react-query";
import { Event, eventsTableColumns } from './columns'
import { DataTable } from './data-table'
import { Skeleton } from '@/components/ui/skeleton'
import SideRowSheet from './sheet/side-row-sheet'
import { EventsForm } from './sheet/forms/events-form'
import { ApiResponse, request } from '@/api/api';


async function fetchTableData() {
    const endpoint = "http://localhost:8080/api/events"
    const result: ApiResponse<Event[]> = await request(endpoint, "GET")
    return result.data
}

export const queryKey = "eventsTableData"

const EventsTable: React.FC = () => {
    const { data, isLoading, isError } = useQuery({ queryKey: [queryKey], queryFn: fetchTableData });

    if (isLoading) {
        // Load dummy empty data and skeleton
        const loadingData = Array(5).fill({}) 
        const loadingTableColumns = eventsTableColumns.map((column) => ({
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
            columns={eventsTableColumns}
            data={data ?? []}
            AddSheet={
                () => <SideRowSheet
                    title="Add new event record"
                    description="Refresh the page once you hit submit to see updated table"
                    FormType={<EventsForm />}
                />
            }
        />
    )
}

export default EventsTable
