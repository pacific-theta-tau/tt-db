import React, { useEffect, useState } from 'react';
import { useParams } from 'react-router-dom';
import { brotherStatusTableColumns, BrotherStatus } from "../components/columns"
import { DataTable } from "../components/data-table"
import { Skeleton } from '@/components/ui/skeleton'
import AddRowSheet from '@/components/sheet/add-row-sheet';
import { getData, ApiResponse } from '@/api/api'
import { BrotherStatusForm } from '@/components/sheet/forms/brothers-status-form';
import { Dropdown } from 'react-day-picker';


const ActivesPage: React.FC = () => {
    const [data, setData] = useState<BrotherStatus[]>([]);
    const [loading, setLoading] = useState<boolean | null>(true);
    const [error, setError] = useState<string | null>(null);
    const { semester = "" } = useParams<{ semester: string }>();

    useEffect(() => {
        const endpoint = `http://localhost:8080/api/semesters/${semester}/statuses`
        const fetchData = async () => {
             try {
                setLoading(true)
                const result: ApiResponse<BrotherStatus[]> = await getData(endpoint)
                console.log('brother status for semester:', result)
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
        const loadingTableColumns = brotherStatusTableColumns.map((column) => ({
            ...column,
            cell: () => <Skeleton className="h-12"/>,
          }))
        return <DataTable columns={ loadingTableColumns } data={loadingData} />
    }

    if (error) {
        return <div>Error: {error}</div>;
    }

    return (
        <div>
            <div className="space-y-2 mb-4">
                <h1 className="scroll-m-20 text-3xl font-bold tracking-tight">{ semester.toString() } Actives</h1>
                <p className="text-base text-muted-foreground">List of active members during {semester.toString()}</p>
            </div>
            <DataTable
                columns={brotherStatusTableColumns}
                data={data}
                AddSheet={
                    () => <AddRowSheet
                            title="Add new member record"
                            description="Refresh the page once you hit submit"
                            FormType={<BrotherStatusForm />}
                          />}
            />
        </div>
   )
}

export default ActivesPage;
