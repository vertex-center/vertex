import { useEffect, useState } from "react";
import { AxiosError } from "axios";

export function useFetch<T>(call: any) {
    const [data, setData] = useState<T>(undefined);
    const [error, setError] = useState<AxiosError>(undefined);
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
