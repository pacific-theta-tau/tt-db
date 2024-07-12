import React, { useEffect, useState } from 'react';
import { Brother, columns } from "./columns"
import { DataTable } from "./data-table"


const DemoPage: React.FC = () => {
    const [data, setData] = useState<Brother[]>([]);
    const [loading, setLoading] = useState<boolean | null>(true);
    const [error, setError] = useState<string | null>(null);

    useEffect(() => {
        const endpoint = "http://localhost:8080/api/brothers"
        // const endpoint = "http://localhost:3000/data"
        const fetchData = async () => {
             try {
                setLoading(true)
                const response = await fetch(endpoint, {
                    mode: 'cors',
                    headers: {
                        'Content-Type': 'application/json',
                    }
                });
                if (!response.ok) {
                  throw new Error('Network response was not ok');
                }
                const result: Brother[] = await response.json();
                console.log('result:', result)
                setData(result);
            } catch (e) {
                setError((e as Error).message);
                console.log('Error fetching data:', error);
                throw error;
            } finally {
                setLoading(false);
            }
        }
        fetchData()
       }, []);

    if (loading) {
        return <div>Loading...</div>
    }

    if (error) {
        return <div>Error: {error}</div>;
    }

    return (
        <DataTable columns={columns} data={data} />
   )
}

export default DemoPage
