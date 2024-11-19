import React from 'react';
import BrothersTable from '../components/brothers-table'

//const test: React.ForwardRefExoticComponent<LinkProps & ReactRefAttributes<HTMLAnchorElement>>

const BrothersPage: React.FC = () => {
    return (
        <div>
            <div className="space-y-2 mb-4">
                <h1 className="scroll-m-20 text-3xl font-bold tracking-tight">Brothers List</h1>
                <p className="text-base text-muted-foreground">List of all Brothers in the Lambda Delta Chapter of Theta Tau</p>
            </div>
            <BrothersTable />
        </div>
    );
};

export default BrothersPage;
