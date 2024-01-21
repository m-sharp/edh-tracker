import { useEffect, useState } from "react";

interface AsyncResp {
    data: any;
    loading: boolean;
    error: any;
}

export function AsyncComponentHelper(fetcher: Promise<any>): AsyncResp {
    const [data, setData] = useState<any>(null);
    const [loading, setLoading] = useState<boolean>(true);
    const [error, setError] = useState<any>(null);

    useEffect(() => {
        async function fetchData() {
            try {
                const objs = await fetcher;
                setData(objs);
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
