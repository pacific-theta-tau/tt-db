import React, { useEffect, useState } from 'react';
import { useQuery } from "@tanstack/react-query";
import { Brother, brothersTableColumns } from "./columns"
import { DataTable } from "./data-table"
import { Skeleton } from "@/components/ui/skeleton"
import { BrotherForm } from './sheet/forms/brothers-form'
import SideRowSheet from './sheet/side-row-sheet';
import { ApiResponse, getData } from '../api/api'


async function fetchTableData() {
        const endpoint = "http://localhost:8080/api/brothers"
        try {
            const result: ApiResponse<Brother[]> = await getData(endpoint)
            console.log('result:', result)
            return result.data
        } catch (error: any) {
            console.log('Error fetching data:', error);
            throw error;
        }
}


const BrothersTable: React.FC = () => {
    // const [data, setData] = useState<Brother[]>([]);
    // const [loading, setLoading] = useState<boolean | null>(true);
    // const [error, setError] = useState<string | null>(null);
    const { data, isLoading, isError } = useQuery({queryKey: ["brothersTableData"], queryFn: fetchTableData})

    if (isLoading) {
        // Load dummy empty data and skeleton
        const loadingData = Array(5).fill({}) 
        const loadingTableColumns = brothersTableColumns.map((column) => ({
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
            columns={brothersTableColumns}
            data={data ?? []}
            AddSheet={
                () => <SideRowSheet
                        title="Add new member record"
                        description="Refresh the page once you hit submit"
                        FormType={<BrotherForm />}
                      />}
        />
   )
}

export default BrothersTable 
