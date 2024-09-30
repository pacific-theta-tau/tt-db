// data-table.tsx: component to be used to build Brothers and Events data tables
"use client"

import React from 'react'
import {
    ColumnDef,
    SortingState,
    flexRender,
    getCoreRowModel,
    getPaginationRowModel,
    useReactTable,
    getSortedRowModel,
    VisibilityState,
    getFilteredRowModel,
    Row,
    FilterFn
} from "@tanstack/react-table"

import { Plus } from 'lucide-react'

import {
  DropdownMenu,
  DropdownMenuCheckboxItem,
  DropdownMenuContent,
  DropdownMenuTrigger,
} from "@/components/ui/dropdown-menu"

import {
    Table,
    TableBody,
    TableCell,
    TableHead,
    TableHeader,
    TableRow,
} from "@/components/ui/table"

import { Button } from "@/components/ui/button"

import { Input } from "@/components/ui/input"
import AddRowSheet from './add-row-sheet'


// Custom global filter function to allow filtering for both numbers and strings
const customFilterFn: FilterFn<any> = (row: Row<any>, columnId: string, filterValue: string) => {
    const search = filterValue.toLowerCase();
    let value = row.getValue(columnId) as string;

    if (typeof value === 'number') value = String(value);

    return value?.toLowerCase().includes(search);
};

interface DataTableProps<TData, TValue> {
  columns: ColumnDef<TData, TValue>[]
  data: TData[]
}

export function DataTable<TData, TValue>({
  columns,
  data,
}: DataTableProps<TData, TValue>) {
    const [columnVisibility, setColumnVisibility] = React.useState<VisibilityState>({})
    const [sorting, setSorting] = React.useState<SortingState>([])
    const [globalFilter, setGlobalFilter] = React.useState<any>("")

    const table = useReactTable({
        data,
        columns,
        getCoreRowModel: getCoreRowModel(),
        getPaginationRowModel: getPaginationRowModel(),
        onSortingChange: setSorting,
        getSortedRowModel: getSortedRowModel(),
        onColumnVisibilityChange: setColumnVisibility,
        onGlobalFilterChange: setGlobalFilter,
        getFilteredRowModel: getFilteredRowModel(),
        
        globalFilterFn: customFilterFn,
        state: {
            sorting,
            columnVisibility,
            globalFilter,
        },
    })

  return (
  <div>
    <div className="flex items-center py-4">
        <div>
            <Input
                placeholder="Search..."
                value={ globalFilter }
                onChange={(e) => table.setGlobalFilter(String(e.target.value))}
                className="max-w-sm"
            />
        </div>
        <DropdownMenu>
              <DropdownMenuTrigger asChild>
                <Button variant="outline" className="ml-auto">
                  Columns
                </Button>
              </DropdownMenuTrigger>
              <DropdownMenuContent align="end">
                {table
                  .getAllColumns()
                  .filter(
                    (column) => column.getCanHide()
                  )
                  .map((column) => {
                    return (
                      <DropdownMenuCheckboxItem
                        key={column.id}
                        className="capitalize"
                        checked={column.getIsVisible()}
                        onCheckedChange={(value) =>
                          column.toggleVisibility(!!value)
                        }
                      >
                        {column.id}
                      </DropdownMenuCheckboxItem>
                    )
                  })}
              </DropdownMenuContent>
        </DropdownMenu>

        <AddRowSheet />
    </div>

    <div className="rounded-md border">
      <Table>
        <TableHeader>
          {table.getHeaderGroups().map((headerGroup) => (
            <TableRow key={headerGroup.id}>
              {headerGroup.headers.map((header) => {
                return (
                  <TableHead key={header.id}>
                    {header.isPlaceholder
                      ? null
                      : flexRender(
                          header.column.columnDef.header,
                          header.getContext()
                        )}
                  </TableHead>
                )
              })}
            </TableRow>
          ))}
        </TableHeader>
        <TableBody>
          {table.getRowModel().rows?.length ? (
            table.getRowModel().rows.map((row) => (
              <TableRow
                key={row.id}
                data-state={row.getIsSelected() && "selected"}
              >
                {row.getVisibleCells().map((cell) => (
                  <TableCell key={cell.id}>
                    {flexRender(cell.column.columnDef.cell, cell.getContext())}
                  </TableCell>
                ))}
              </TableRow>
            ))
          ) : (
            <TableRow>
              <TableCell colSpan={columns.length} className="h-24 text-center">
                No results.
              </TableCell>
            </TableRow>
          )}
        </TableBody>
      </Table>
    </div>
        <div>
        <div className="flex items-center justify-end space-x-2 py-4">
        <Button
          variant="outline"
          size="sm"
          onClick={() => table.previousPage()}
          disabled={!table.getCanPreviousPage()}
        >
          Previous
        </Button>
        <Button
          variant="outline"
          size="sm"
          onClick={() => table.nextPage()}
          disabled={!table.getCanNextPage()}
        >
          Next
        </Button>
      </div>
    </div>
    </div>
  )
}

