// Line chart to show number of current actives throughout the semesters
"use client"

import React, { useState, useEffect } from 'react'
import { TrendingDown } from "lucide-react"
import { CartesianGrid, Line, LineChart, XAxis } from "recharts"

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

export const description = "A linear line chart"

interface activesData {
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
  actives: {
    label: "Actives",
    color: "hsl(var(--chart-1))",
  },
} satisfies ChartConfig

export function LineChartActives() {
    const [chartData, setChartData] = useState<activesData[]>([])
    useEffect(() => {
        setChartData(testChartData)
    }, [])
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
              stroke="var(--color-actives)"
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

