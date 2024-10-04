const baseURL = "http://localhost:8080"

export interface ApiResponse<T> {
    status: string,
    message: string,
    data: T,
}

export const getData = async <T>(endpoint: string, urlParams?: Record<string, string>): Promise<T> => {
    const url = new URL(endpoint, baseURL);

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

