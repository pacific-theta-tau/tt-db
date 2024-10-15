"use client"

import React, { useState, useEffect } from "react"
import { TrendingUp } from "lucide-react"
import { Label, Pie, PieChart } from "recharts"

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
  ChartLegend,
  ChartLegendContent
} from "@/components/ui/chart"
import { getData, ApiResponse } from "@/api/api"


const DEBUG = true 
export const description = "A donut chart showing members by major"

interface MajorCount {
    major: string,
    count: number,
}

interface Count {
    count: number
}

const chartConfig = {
  count: {
    label: "Brothers",
  },
  bioengineering: {
    label: "Bioengineering",
    color: "hsl(var(--chart-1))",
  },
  civilengineering: {
    label: "Civil Engineering",
    color: "hsl(var(--chart-2))",
  },
  computerengineering: {
    label: "Computer Engineering",
    color: "hsl(var(--chart-3))",
  },
  computerscience: {
    label: "Computer Science",
    color: "hsl(var(--chart-4))",
  },
  electricalengineering: {
    label: "Electrical Engineering",
    color: "hsl(var(--chart-5))",
  },
  engineeringphysics: {
    label: "Engineering Physics",
    color: "hsl(var(--chart-6))",
  },
  mechanicalengineering: {
    label: "Mechanical Engineering",
    color: "hsl(var(--chart-7))",
  },
} satisfies ChartConfig

// mapping of major to chart color to inject into ChartData
const majorColors: { [key: string]: string } = {
  bioengineering: "var(--color-bioengineering)",
  civilengineering: "var(--color-civilengineering)",
  computerengineering: "var(--color-computerengineering)",
  computerscience: "var(--color-computerscience)",
  electricalengineering: "var(--color-electricalengineering)",
  engineeringphysics: "var(--color-engineeringphysics)",
  mechanicalengineering: "var(--color-mechanicalengineering)",
}

// Dummy data used for debugging/testing. set `DEBUG=true` to use
const testChartData = [
    { major: "computerscience", count: 60, fill: "var(--color-computerscience)" },
    { major: "computerengineering", count: 40, fill: "var(--color-computerengineering)" },
    { major: "bioengineering", count: 50, fill: "var(--color-bioengineering)" },
    { major: "civilengineering", count: 30, fill: "var(--color-civilengineering)" },
    { major: "electricalengineering", count: 20, fill: "var(--color-electricalengineering)" },
]

export function PieChartMajorsDistribution() {
    const [chartData, setChartData] = React.useState<MajorCount[]>([]);
    const [totalCount, setTotalCount] = React.useState(0);
    const [loading, setLoading] = React.useState<boolean>(true);
    const [error, setError] = React.useState<string | null>(null);

    useEffect(() => {
        let result
        const fetchData = async () => {
            try {
                let endpoint = "/api/brothers/majors/count" 
                result = await getData< ApiResponse<MajorCount[]> >(endpoint)
                console.log(result)
                const processedData = result.data.map(item => ({
                    ...item,
                    // sadly the value of 'major' needs to match with their respective chartConfig keys (no spaces) so this is a solution for now
                    major: item.major.replace(/\s+/g, '').toLowerCase(),
                    fill: majorColors[item.major.replace(/\s+/g, '').toLowerCase()] || "hsl(var(--chart-8))"
                }))
                console.log("Processed Data:", processedData)
                setChartData(processedData)

                endpoint = "/api/brothers/count" 
                result = await getData< ApiResponse<Count> >(endpoint)
                console.log(result)
                setTotalCount(result.data.count)
            } catch (error: any) {
                setError(error.message)
            } finally {
                setLoading(false);
            }
        }
        if (DEBUG) {
            setChartData(testChartData)
            setTotalCount(200)
            setLoading(false)
        } else {
            fetchData()
        }
    }, []);

    if (loading) {
        // Load dummy empty data and skeleton
        return <div>Loading...</div>
    }

    if (error) {
        return <div>Error: {error}</div>;
    }

  return (
    <Card className="flex flex-col">
      <CardHeader className="items-center pb-0">
        <CardTitle>Major distribution across members</CardTitle>
        <CardDescription>All Members</CardDescription>
      </CardHeader>
      <CardContent className="flex-1 pb-0">
        <ChartContainer
          config={chartConfig}
          className="mx-auto min-h-[200px] max-h-[450px]"
        >
          <PieChart>
            <ChartTooltip
              cursor={false}
              content={<ChartTooltipContent hideLabel />}
            />
            <Pie
              data={chartData}
              dataKey="count"
              nameKey="major"
              innerRadius={60}
              strokeWidth={5}
            >
              <Label
                content={({ viewBox }) => {
                  if (viewBox && "cx" in viewBox && "cy" in viewBox) {
                    return (
                      <text
                        x={viewBox.cx}
                        y={viewBox.cy}
                        textAnchor="middle"
                        dominantBaseline="middle"
                      >
                        <tspan
                          x={viewBox.cx}
                          y={viewBox.cy}
                          className="fill-foreground text-3xl font-bold"
                        >
                          {totalCount.toLocaleString()}
                        </tspan>
                        <tspan
                          x={viewBox.cx}
                          y={(viewBox.cy || 0) + 24}
                          className="fill-muted-foreground"
                        >
                          Brothers
                        </tspan>
                      </text>
                    )
                  }
                }}
              />
            </Pie>
            <ChartLegend
              content={<ChartLegendContent nameKey="major" />}
              className="-translate-y-2 flex-wrap gap-2 [&>*]:basis-1/4 [&>*]:justify-center text-sm"
            />
          </PieChart>
        </ChartContainer>
      </CardContent>
      <CardFooter className="flex-col gap-2 text-sm">
        <div className="leading-none text-muted-foreground">
            Showing total members across all engineering majors
        </div>
      </CardFooter>
    </Card>
  )
}

