// Line chart to show number of current actives throughout the semesters
"use client"

import React, { useState, useEffect } from 'react'
import { TrendingDown } from "lucide-react"
import { CartesianGrid, Line, LineChart, XAxis } from "recharts"
import { getData, ApiResponse } from "@/api/api"

import {
  Card,
  CardContent,
  CardDescription,
  CardFooter,
  CardHeader,
  CardTitle,
} from "@/components/ui/card"
import {
  ChartConfig,
  ChartContainer,
  ChartTooltip,
  ChartTooltipContent,
} from "@/components/ui/chart"


const DEBUG = true
export const description = "A linear line chart"

interface ActivesCount {
    semester: string,
    count: number,
}

const testChartData = [
  { semester: "Fall 2022", count: 31 },
  { semester: "Spring 2023", count: 42 },
  { semester: "Fall 2023", count: 29 },
  { semester: "Spring 2024", count: 34 },
  { semester: "Fall 2024", count: 30 },
]

const chartConfig = {
  count: {
    label: "Actives",
    color: "hsl(var(--chart-1))",
  },
} satisfies ChartConfig

export function LineChartActives() {
    const [chartData, setChartData] = useState<ActivesCount[]>([])
    const [loading, setLoading] = useState<boolean>(true)
    const [error, setError] = useState<string | null>(null)

    useEffect(() => {
        let result
        const fetchData = async () => {
            try {
                // "Get the total count of active members grouped by semester"
                let endpoint = "/api/brothers/statuses/count"
                result = await getData< ApiResponse<ActivesCount[]> >(endpoint, {status: "Active"})
                setChartData(result.data)
            } catch (error: any) {
                setError(error.message)
                console.log(error.message)
            } finally {
                setLoading(false);
            }
        }

        if (DEBUG) {
            console.log("Using test chart data")
            setChartData(testChartData)
            setError(null)
            setLoading(false)
        } else {
            fetchData()
        }
    }, [])


    if (loading) {
        // Load dummy empty data and skeleton
        return <div>Loading...</div>
    }

    if (error) {
        return <div>Error: {error}</div>;
    }

  return (
    <Card>
      <CardHeader>
        <CardTitle>Actives count throughout semesters</CardTitle>
        <CardDescription>Fall 2022 - Fall 2024</CardDescription>
      </CardHeader>
      <CardContent>
        <ChartContainer config={chartConfig}>
          <LineChart
            accessibilityLayer
            data={chartData}
            margin={{
              left: 12,
              right: 12,
            }}
          >
            <CartesianGrid vertical={false} />
            <XAxis
              dataKey="semester"
              tickLine={false}
              axisLine={false}
              tickMargin={8}
              // tickFormatter={(value) => value.slice(0, 4)}
            />
            <ChartTooltip
              cursor={false}
              content={<ChartTooltipContent hideLabel />}
            />
            <Line
              dataKey="count"
              type="linear"
              stroke="var(--color-count)"
              strokeWidth={2}
              dot={false}
            />
          </LineChart>
        </ChartContainer>
      </CardContent>
      <CardFooter className="flex-col items-start gap-2 text-sm">
        <div className="flex gap-2 font-medium leading-none">
          Trending down by 12.8% this semester <TrendingDown className="h-4 w-4" />
        </div>
        <div className="leading-none text-muted-foreground">
          Showing actives from the past 5 semesters
        </div>
      </CardFooter>
    </Card>
  )
}

