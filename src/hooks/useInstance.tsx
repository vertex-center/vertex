import { useCallback, useEffect, useState } from "react";
import { api } from "../backend/backend";
import { Instance } from "../models/instance";

export default function useInstance(uuid?: string) {
    const [instance, setInstance] = useState<Instance>();
    const [loading, setLoading] = useState<boolean>(true);

    const reloadInstance = useCallback(() => {
        console.log("Fetching instance", uuid);
        setLoading(true);
        api.instance
            .get(uuid)
            .then((res) => setInstance(res.data))
            .catch(console.error)
            .finally(() => setLoading(false));
    }, [uuid]);

    useEffect(() => {
        reloadInstance();
    }, [uuid]);

    return { instance, setInstance, reloadInstance, loading };
}
