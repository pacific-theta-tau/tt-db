const BASEURL = "http://localhost:8080"

export interface ApiResponse<T> {
    status: string,
    message: string,
    data: T,
}

export const getData = async <T>(endpoint: string, urlParams?: Record<string, string>): Promise<T> => {
    const url = new URL(endpoint, BASEURL);

    // Append URL parameters if provided
    if (urlParams) {
        Object.keys(urlParams).forEach(key => {
            url.searchParams.append(key, urlParams[key])
        })
    }

    const response = await fetch(
        url.toString(),
        {
            mode: 'cors',
            headers: {
                'Content-Type': 'application/json',
            }
        }
    );
    if (!response.ok) {
        throw new Error(`Failed to fetch data from ${url.toString()}: ${response.statusText}`);
    }

    return response.json();
};

// TODO: maybe just have single API call handler and pass method as an argument.
export const requestPOST = async <T>(endpoint: string, body: string): Promise<T> => {
    const response = await fetch(endpoint, {
        method: 'POST',
        body: body,
        mode: 'cors',
        headers: {
            'Content-Type': 'application/json',
        }
    });

    if (!response.ok) {
        throw new Error(`Failed to fetch data: ${response.statusText}`);
    }

    return response.json();
}

export type ApiRequestOptions = {
    method?: 'GET' | 'POST' | 'PUT' | 'DELETE';
    body?: any;
    queryParams?: any;
}

export const request = async <T>(
    endpoint: string,
    method?: 'GET' | 'POST' | 'PUT' | 'DELETE',
    body?: Record<string, unknown>,
    queryParams?: any,
): Promise<T> => {
    console.log("\n---Starting new request")
    console.log("\tmethod:", method)
    console.log("\tqueryParams:", queryParams)
    console.log("\tbody:", body)
    const url = new URL(endpoint, BASEURL)

    if (queryParams) {
        Object.keys(queryParams).forEach(key => {
            url.searchParams.append(key, queryParams[key])
        })
    }

    try {
        const response = await fetch(
            url.toString(),
            {
                method: method? method : 'GET', // Default method = GET
                body: body? JSON.stringify(body) : body,
                mode: 'cors',
                headers: {
                    'Content-Type': 'application/json',
                }
            }
        )

        if (response.status === 200 || response.status === 201) {
            return (await response.json()) as T;
        }

        const errorBody = await response.json().catch(() => null);
        const errorMessage = errorBody?.message || `HTTP error: ${response.status}`;
        throw new Error(errorMessage);
    } catch (error: unknown) {
        console.error('propagate throw');
        throw error; 
    }

}
