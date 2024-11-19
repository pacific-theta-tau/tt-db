// Home.tsx (default export)
import React from 'react';
import { PieChartMajorsDistribution } from '@/components/graphs/major-pie-chart';
import  { LineChartActives } from '@/components/graphs/actives-line-chart';

const HomePage: React.FC = () => {
    return (
        <div className="flex">
            <div className="flex-1 p-2">
                <PieChartMajorsDistribution />
            </div>

            <div className="flex-1 p-2">
                <LineChartActives />
            </div>
        </div>
       );
};

export default HomePage;
