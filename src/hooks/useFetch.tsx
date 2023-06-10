import { useEffect, useState } from "react";

export function useFetch<T>(call: any) {
    const [data, setData] = useState<T>();
    const [error, setError] = useState();
    const [loading, setLoading] = useState<boolean>(true);

    const reload = async () => {
        setLoading(true);
        await call()
            .then((res) => setData(res.data))
            .catch((error) => setError(error))
            .finally(() => setLoading(false));
    };

    useEffect(() => {
        reload().then();
    }, []);

    return { data, error, loading, reload };
}
