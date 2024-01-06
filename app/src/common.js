import { useEffect, useState } from "react";

export function AsyncComponentHelper(fetcher) {
    const [data, setData] = useState(null);
    const [loading, setLoading] = useState(true);
    const [error, setError] = useState(null);

    useEffect(() => {
        async function fetchData() {
            try {
                const games = await fetcher;
                setData(games);
                setLoading(false);
            } catch (error) {
                setError(error);
                setLoading(false);
            }
        }

        fetchData();
    }, []);

    return {
        data: data,
        loading: loading,
        error: error,
    };
}
