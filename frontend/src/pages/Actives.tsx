import React, { useState } from 'react';
import { useQuery } from "@tanstack/react-query";
import { brotherStatusTableColumns, BrotherStatus } from "../components/columns"
import { DataTable } from "../components/data-table"
import { Skeleton } from '@/components/ui/skeleton'
import SideRowSheet from '@/components/sheet/side-row-sheet';
import { request, ApiResponse } from '@/api/api'
import { BrotherStatusForm } from '@/components/sheet/forms/brothers-status-form';
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


async function fetchTableData(selectedSemester: string): Promise<BrotherStatus[]> {
    console.log(">Fetching actives from", selectedSemester)
    const endpoint = `http://localhost:8080/api/semesters/${selectedSemester}/statuses`
    const result: ApiResponse<BrotherStatus[]> = await request(endpoint, 'GET')
    console.log(result)

    return result.data
}

async function fetchSemesterData() {
    const endpoint = "http://localhost:8080/api/semesters"
    const responseSemesters: ApiResponse<string[]> = await request(endpoint, 'GET')
    return responseSemesters.data
}

export const activesQueryKey = "activesTableData"

const ActivesPage: React.FC = () => {
    const [selectedSemester, setSelectedSemester] = useState<string>(getSeasonYear)
    console.log('actives page- selectedSemester', selectedSemester)

    // React Query hooks
    const { data: semesterLabels } = useQuery({ queryKey: ["semesters"], queryFn: fetchSemesterData })
    const { data, isLoading, isError } = useQuery({ queryKey: [activesQueryKey, selectedSemester], queryFn: () => fetchTableData(selectedSemester) })

    if (isLoading) {
        // Load dummy empty data and skeleton
        const loadingData = Array(5).fill({}) 
        const loadingTableColumns = brotherStatusTableColumns.map((column) => ({
            ...column,
            cell: () => <Skeleton className="h-12"/>,
          }))
        return <DataTable columns={ loadingTableColumns } data={loadingData} />
    }

    if (isError) {
        return <div className="text-red-500">Error loading table data</div>;
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
                  {semesterLabels && semesterLabels.length > 0 ? (
                      semesterLabels.map((semesterLabel, index) => (
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
                data={data ?? []}
                AddSheet={
                    () =>(
                        <SideRowSheet
                            title="Add new member record"
                            description="Refresh the page once you hit submit"
                            FormType={
                                <BrotherStatusForm selectedSemester={selectedSemester} />
                            }
                        />
                    )
                }
            />
        </div>
   )
}

export default ActivesPage;
