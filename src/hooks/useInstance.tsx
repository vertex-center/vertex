import { useCallback, useEffect, useState } from "react";
import { getInstance, Instance } from "../backend/backend";

export default function useInstance(uuid?: string) {
    const [instance, setInstance] = useState<Instance>();

    // const reloadInstance = () => {
    //     console.log("Fetching instance", uuid);
    //     getInstance(uuid)
    //         .then((instance: Instance) => {
    //             setInstance(instance);
    //         })
    //         .catch((err) => {
    //             setInstance(undefined);
    //             console.error(err);
    //         });
    // };

    const reloadInstance = useCallback(() => {
        console.log("Fetching instance", uuid);
        getInstance(uuid).then(setInstance).catch(console.error);
    }, [uuid]);

    useEffect(() => {
        reloadInstance();
    }, [uuid]);

    return { instance, setInstance, reloadInstance };
}
