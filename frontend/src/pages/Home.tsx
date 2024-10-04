// Home.tsx (default export)
import React from 'react';
import { Component } from '@/components/graphs/major-pie-chart';
import { Component2 } from '@/components/graphs/pie-test';

const HomePage: React.FC = () => {
    return (
        <div>
            <h1>Home Page</h1>
            <Component />
            <Component2 />
        </div>
       );
};

export default HomePage;
