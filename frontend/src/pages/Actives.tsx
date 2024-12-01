import React, { useEffect, useState } from 'react';
import { useParams } from 'react-router-dom';
import { brotherStatusTableColumns, BrotherStatus } from "../components/columns"
import { DataTable } from "../components/data-table"
import { Skeleton } from '@/components/ui/skeleton'
import AddRowSheet from '@/components/sheet/add-row-sheet';
import { getData, ApiResponse } from '@/api/api'
import { BrotherStatusForm } from '@/components/sheet/forms/brothers-status-form';
import { Dropdown } from 'react-day-picker';
import { Select, SelectItem, SelectContent, SelectTrigger, SelectValue } from "@/components/ui/select";


/*
 * Helper function to get today's Semester + Year
 * */
export function getSeasonYear(): string {
  const today = new Date();
  const month = today.getMonth(); // Months are 0-indexed: January is 0, December is 11
  const year = today.getFullYear();

  // Determine the season based on the month
  const season = month < 6 ? 'Spring' : 'Fall';

  return `${season} ${year}`;
}

const ActivesPage: React.FC = () => {
    const [data, setData] = useState<BrotherStatus[]>([]);
    const [loading, setLoading] = useState<boolean | null>(true);
    const [error, setError] = useState<string | null>(null);
    const [semestersDropdownList, setSemestersDropdownList] = useState<string[]>([]);
    // TODO: remove `semester` and use useState instead for setting default
    const [selectedSemester, setSelectedSemester] = useState<string>(getSeasonYear)

    // Fetch/Update table data. Listening to dropdown component selection
    useEffect(() => {
        const endpoint = `http://localhost:8080/api/semesters/${selectedSemester}/statuses`
        const fetchData = async () => {
             try {
                setLoading(true)
                const result: ApiResponse<BrotherStatus[]> = await getData(endpoint)
                setData(result.data);
            } catch (error: any) {
                setError((error as Error).message);
                throw error;
            } finally {
                /* uncomment line below to test skeleton during loading */
                // await new Promise(f => setTimeout(f, 3000));
                setLoading(false);
            }
        }
        fetchData()
       }, [selectedSemester]);

    useEffect(() => {
        const endpoint = "http://localhost:8080/api/semesters"
        const fetchData = async () => {
             try {
                const responseSemesters: ApiResponse<string[]> = await getData(endpoint)
                setSemestersDropdownList(responseSemesters.data)
            } catch (error) {
                console.log('Error fetching data:', error);
                throw error;
            } finally {
            }
        }
        fetchData()
    }, [])

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
        return <div className="text-red-500">Error: {error}</div>;
    }

    return (
        <div>
            <div className="space-y-2 mb-4">
                <h1 className="scroll-m-20 text-3xl font-bold tracking-tight">Actives List</h1>
                <p className="text-base text-muted-foreground">List of active members during {selectedSemester}</p>
            </div>

            {/* Semester Dropdown */}
            <div className="mb-4">
              <Select onValueChange={(value) => setSelectedSemester(value)} defaultValue={selectedSemester}>
                <SelectTrigger className="w-64">
                  <SelectValue placeholder="Select a semester" />
                </SelectTrigger>
                <SelectContent>
                  {semestersDropdownList && semestersDropdownList.length > 0 ? (
                      semestersDropdownList.map((semesterLabel, index) => (
                          <SelectItem key={index.toString()} value={semesterLabel}>{semesterLabel}</SelectItem>
                      ))
                  ) : (
                          <SelectItem value="">Loading...</SelectItem>
                  )}
                </SelectContent>
              </Select>
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
