import { useCallback, useEffect, useState } from "react";
import { api } from "../backend/backend";
import { Instance } from "../models/instance";

export default function useInstance(uuid?: string) {
    const [instance, setInstance] = useState<Instance>();

    const reloadInstance = useCallback(() => {
        console.log("Fetching instance", uuid);
        api.instance
            .get(uuid)
            .then((res) => setInstance(res.data))
            .catch(console.error);
    }, [uuid]);

    useEffect(() => {
        reloadInstance();
    }, [uuid]);

    return { instance, setInstance, reloadInstance };
}
