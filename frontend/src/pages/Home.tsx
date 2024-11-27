// Home.tsx (default export)
import React from 'react';
import { PieChartMajorsDistribution } from '@/components/graphs/major-pie-chart';
import  { LineChartActives } from '@/components/graphs/actives-line-chart';
import { Separator } from '@radix-ui/react-separator';

const HomePage: React.FC = () => {
    return (
        <main>
            <section className="flex flex-col items-start gap-2 border-b border-border/40 pt-2 pb-4 dark:border-border">
                <h1 className="text-3xl font-bold leading-tight tracking-tighter md:text-4xl lg:leading-[1.1]">
                    Pacific Theta Tau Database
                </h1>
                <p className="max-w-2xl text-lg font-light text-foreground">
                    Store, read, and visualize historic data from the Lambda Delta Chapter of Theta Tau
                </p>
            </section>
            <Separator className="my-4" />
            <section className="pb-2">
                <h1 className="text-2xl font-normal leading-tight tracking-tighter pl-4 md:text-4xl lg:leading-[1.1]">
                    Analytics
                </h1>
            </section>
            <div className="flex">
                <div className="flex-1 p-2">
                    <PieChartMajorsDistribution />
                </div>

                <div className="flex-1 p-2">
                    <LineChartActives />
                </div>
            </div>
        </main>
       );
};

export default HomePage;
