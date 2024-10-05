// Home.tsx (default export)
import React from 'react';
import { PieChartMajorsDistribution } from '@/components/graphs/major-pie-chart';
import  { LineChartActives } from '@/components/graphs/actives-line-chart';

const HomePage: React.FC = () => {
    return (
        <div className="flex flex-row">
            <PieChartMajorsDistribution />
            <LineChartActives />
        </div>
       );
};

export default HomePage;
