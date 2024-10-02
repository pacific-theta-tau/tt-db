import React from 'react';
import BrothersTable from '../components/brothers-table'

//const test: React.ForwardRefExoticComponent<LinkProps & ReactRefAttributes<HTMLAnchorElement>>

const BrothersPage: React.FC = () => {
    return (
        <div>
            <BrothersTable />
        </div>
    );
};

export default BrothersPage;
