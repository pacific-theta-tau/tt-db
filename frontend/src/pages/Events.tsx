import React from 'react';
import EventsTable from '../components/events-table'


const EventsPage: React.FC = () => {
    return (
        <div>
            <div className="space-y-2 mb-4">
                <h1 className="scroll-m-20 text-3xl font-bold tracking-tight">Events List</h1>
                <p className="text-base text-muted-foreground">List of all recorded events to date</p>
            </div>

            <EventsTable />
        </div>
    );
};

export default EventsPage

