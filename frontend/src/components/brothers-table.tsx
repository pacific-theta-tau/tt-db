import React, { useEffect, useState } from 'react';
import { Brother, brothersTableColumns } from "./columns"
import { DataTable } from "./data-table"
import { Skeleton } from "@/components/ui/skeleton"
import { BrotherForm } from './sheet/forms/brothers-form'
import SideRowSheet from './sheet/side-row-sheet';
import { ApiResponse, getData } from '../api/api'


const BrothersTable: React.FC = () => {
    const [data, setData] = useState<Brother[]>([]);
    const [loading, setLoading] = useState<boolean | null>(true);
    const [error, setError] = useState<string | null>(null);

    useEffect(() => {
        const endpoint = "http://localhost:8080/api/brothers"
        const fetchData = async () => {
             try {
                setLoading(true)
                const result: ApiResponse<Brother[]> = await getData(endpoint)
                console.log('result:', result)
                setData(result.data);
            } catch (error: any) {
                setError((error as Error).message);
                console.log('Error fetching data:', error);
                throw error;
            } finally {
                /* uncomment line below to test skeleton during loading */
                // await new Promise(f => setTimeout(f, 3000));
                setLoading(false);
            }
        }
        fetchData()
       }, []);

    if (loading) {
        // Load dummy empty data and skeleton
        const loadingData = Array(5).fill({}) 
        const loadingTableColumns = brothersTableColumns.map((column) => ({
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
            columns={brothersTableColumns}
            data={data}
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
